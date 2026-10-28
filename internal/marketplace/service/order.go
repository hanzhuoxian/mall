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
	case "pay":
		if !sandbox {
			return "", errors.New("真实支付尚未配置；仅显式启用沙箱支付后可演示")
		}
		if buyer && o.Status == "pending" && time.Now().Before(o.ExpiresAt) {
			return "paid", nil
		}
	case "cancel":
		if buyer && o.Status == "pending" {
			return model.OrderCanceled, nil
		}
	case "ship":
		if seller && o.Status == "paid" {
			return "shipped", nil
		}
	case "receive":
		if buyer && o.Status == "shipped" {
			return "completed", nil
		}
	case "refund":
		if buyer && (o.Status == "paid" || o.Status == "shipped" || o.Status == "completed") {
			return "refunding", nil
		}
	case "approve-refund":
		if sandbox && (seller || platform) && o.Status == "refunding" {
			return "refunded", nil
		}
	case "reject-refund":
		if (seller || platform) && o.Status == "refunding" &&
			(o.RefundFrom == "paid" || o.RefundFrom == "shipped" || o.RefundFrom == "completed") {
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
	if (in.Action == "pay" || in.Action == "approve-refund") && !s.sandboxPayment {
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
		buyer, seller := o.UserID == a.ID, shop.OwnerID == a.ID && shop.Status == "approved"
		if !buyer && !seller && !a.Admin {
			return ErrForbidden
		}
		next, err := transition(o, in.Action, buyer, seller, a.Admin, s.sandboxPayment)
		if err != nil {
			return err
		}
		switch in.Action {
		case "ship":
			if !validText(in.Carrier, 100) || !validText(in.Tracking, 100) {
				return ErrInvalid
			}
			o.Carrier = in.Carrier
			o.Tracking = in.Tracking
		case "refund":
			if !validText(in.Reason, 500) {
				return ErrInvalid
			}
			o.RefundFrom = o.Status
			o.RefundReason = in.Reason
		case "reject-refund":
			if !validText(in.Reason, 500) {
				return ErrInvalid
			}
		}
		if next == model.OrderCanceled || (next == "refunded" && o.RefundFrom == "paid") {
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
	ids, err := s.store.ExpiredOrderIDs(ctx, time.Now(), 100)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if err := s.store.Transaction(ctx, func(tx store.Store) error {
			o, err := tx.Order(ctx, id, true)
			if err != nil {
				return err
			}
			if o.Status != "pending" {
				return nil
			}
			if err := restoreStock(ctx, tx, o); err != nil {
				return err
			}
			if err := tx.SetOrderStatus(ctx, o.ID, model.OrderCanceled); err != nil {
				return err
			}
			return audit(ctx, tx, model.Actor{ID: "system"}, "order.expire", id, "")
		}); err != nil {
			return err
		}
	}
	return nil
}
