package service

import (
	"context"
	"errors"
	"time"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
)

func transition(o model.Order, action string, buyer, seller, platform, sandbox bool) (string, error) {
	switch action {
	case model.OrderActionPay:
		if !sandbox {
			return "", errors.New("真实支付尚未配置；仅显式启用沙箱支付后可演示")
		}
		if buyer && o.Status == model.OrderStatusPending && time.Now().Before(o.ExpiresAt) {
			return model.OrderStatusPaid, nil
		}
	case model.OrderActionCancel:
		if buyer && o.Status == model.OrderStatusPending {
			return model.OrderCanceled, nil
		}
	case model.OrderActionShip:
		if seller && o.Status == model.OrderStatusPaid {
			return model.OrderStatusShipped, nil
		}
	case model.OrderActionReceive:
		if buyer && o.Status == model.OrderStatusShipped {
			return model.OrderStatusCompleted, nil
		}
	case model.OrderActionRefund:
		if buyer && (o.Status == model.OrderStatusPaid || o.Status == model.OrderStatusShipped || o.Status == model.OrderStatusCompleted) {
			return model.OrderStatusRefunding, nil
		}
	case model.OrderActionApproveRefund:
		if sandbox && (seller || platform) && o.Status == model.OrderStatusRefunding {
			return model.OrderStatusRefunded, nil
		}
	case model.OrderActionRejectRefund:
		if (seller || platform) && o.Status == model.OrderStatusRefunding &&
			(o.RefundFrom == model.OrderStatusPaid || o.RefundFrom == model.OrderStatusShipped || o.RefundFrom == model.OrderStatusCompleted) {
			return o.RefundFrom, nil
		}
	}
	return "", ErrInvalid
}
func restoreStock(ctx context.Context, tx store.Store, o model.Order) error {
	items, err := tx.OrderItems(ctx, o.ID)
	if err != nil {
		return err
	}
	for _, i := range items {
		if err := tx.AdjustStock(ctx, i.ProductID, i.Quantity); err != nil {
			return err
		}
	}
	return nil
}

// ActOrder 根据操作人的身份及当前订单状态执行订单操作。
func (s *Service) ActOrder(ctx context.Context, a model.Actor, id string, in model.OrderAction) error {
	if (in.Action == model.OrderActionPay || in.Action == model.OrderActionApproveRefund) && !s.sandboxPayment {
		return ErrPaymentDisabled
	}
	return s.store.Transaction(ctx, func(tx store.Store) error {
		o, err := tx.Order(ctx, id, true)
		if err != nil {
			return err
		}
		shop, err := tx.Shop(ctx, o.ShopID, false)
		if err != nil {
			return err
		}
		buyer, seller := o.UserID == a.ID, shop.OwnerID == a.ID && shop.Status == model.ShopStatusApproved
		if !buyer && !seller && !a.Admin {
			return ErrForbidden
		}
		next, err := transition(o, in.Action, buyer, seller, a.Admin, s.sandboxPayment)
		if err != nil {
			return err
		}
		switch in.Action {
		case model.OrderActionShip:
			if !validText(in.Carrier, model.MaxCarrierLength) || !validText(in.Tracking, model.MaxTrackingLength) {
				return ErrInvalid
			}
			o.Carrier = in.Carrier
			o.Tracking = in.Tracking
		case model.OrderActionRefund:
			if !validText(in.Reason, model.MaxReasonLength) {
				return ErrInvalid
			}
			o.RefundFrom = o.Status
			o.RefundReason = in.Reason
		case model.OrderActionRejectRefund:
			if !validText(in.Reason, model.MaxReasonLength) {
				return ErrInvalid
			}
		}
		if next == model.OrderCanceled || (next == model.OrderStatusRefunded && o.RefundFrom == model.OrderStatusPaid) {
			if err = restoreStock(ctx, tx, o); err != nil {
				return err
			}
		}
		o.Status = next
		if err = tx.SaveOrder(ctx, &o); err != nil {
			return err
		}
		return audit(ctx, tx, a, "order."+in.Action, id, in.Reason)
	})
}

// ExpireOrders is safe to invoke concurrently: each order is checked under a row lock.
func (s *Service) ExpireOrders(ctx context.Context) error {
	ids, err := s.store.ExpiredOrderIDs(ctx, time.Now(), model.MaxPageSize)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if err := s.store.Transaction(ctx, func(tx store.Store) error {
			o, err := tx.Order(ctx, id, true)
			if err != nil {
				return err
			}
			if o.Status != model.OrderStatusPending {
				return nil
			}
			if err := restoreStock(ctx, tx, o); err != nil {
				return err
			}
			if err := tx.SetOrderStatus(ctx, o.ID, model.OrderCanceled); err != nil {
				return err
			}
			return audit(ctx, tx, model.Actor{ID: defaultActorID}, "order.expire", id, "")
		}); err != nil {
			return err
		}
	}
	return nil
}
