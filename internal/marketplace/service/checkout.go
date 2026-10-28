package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
)

func fingerprint(in model.CheckoutInput) string {
	ids := append([]string(nil), in.CartIDs...)
	sort.Strings(ids)
	data, _ := json.Marshal(struct {
		Address string
		IDs     []string
	}{in.AddressID, ids})
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Checkout 按商家拆单，在同一事务内预留库存并处理幂等重试。
func (s *Service) Checkout(ctx context.Context, a model.Actor, in model.CheckoutInput) ([]model.Order, error) {
	orders := []model.Order{}
	if !validText(in.RequestKey, 100) || len(in.CartIDs) == 0 || len(in.CartIDs) > 100 {
		return nil, ErrInvalid
	}
	seen := map[string]bool{}
	for _, id := range in.CartIDs {
		if seen[id] {
			return nil, ErrInvalid
		}
		seen[id] = true
	}
	err := s.store.Transaction(ctx, func(tx store.Store) error {
		batch := model.Checkout{
			Base:        base(),
			UserID:      a.ID,
			RequestKey:  in.RequestKey,
			Fingerprint: fingerprint(in),
		}
		if err := tx.CreateCheckout(ctx, &batch); err != nil {
			return err
		}
		stored, err := tx.CheckoutByKey(ctx, a.ID, in.RequestKey)
		if err != nil {
			return err
		}
		if stored.Fingerprint != batch.Fingerprint {
			return ErrInvalid
		}
		if stored.ID != batch.ID {
			orders, err = tx.OrdersByCheckout(ctx, stored.ID)
			return err
		}
		address, err := tx.Address(ctx, in.AddressID, a.ID)
		if err != nil {
			return err
		}
		cart, err := tx.CartByIDs(ctx, a.ID, in.CartIDs)
		if err != nil {
			return err
		}
		if len(cart) != len(in.CartIDs) {
			return ErrInvalid
		}
		// 所有商品与结算流程都先锁定商家，再锁定商品。
		ids := []string{}
		for _, c := range cart {
			ids = append(ids, c.ProductID)
		}
		products, err := tx.ProductsByIDs(ctx, ids)
		if err != nil {
			return err
		}
		shopIDs := []string{}
		shops := map[string]model.Shop{}
		for _, p := range products {
			if _, ok := shops[p.ShopID]; !ok {
				shops[p.ShopID] = model.Shop{}
				shopIDs = append(shopIDs, p.ShopID)
			}
		}
		sort.Strings(shopIDs)
		for _, id := range shopIDs {
			shop, err := tx.Shop(ctx, id, true)
			if err != nil {
				return err
			}
			if shop.Status != "approved" {
				return ErrStock
			}
			shops[id] = shop
		}
		groups := map[string]*model.Order{}
		for _, c := range cart {
			p, err := tx.Product(ctx, c.ProductID, true)
			if err != nil {
				return err
			}
			if p.Status != "approved" || p.Stock < c.Quantity || c.Quantity < 1 || c.Quantity > 999 {
				return ErrStock
			}
			o := groups[p.ShopID]
			if o == nil {
				o = &model.Order{
					Base:       base(),
					CheckoutID: batch.ID,
					UserID:     a.ID,
					ShopID:     p.ShopID,
					Status:     "pending",
					Recipient:  address.Recipient,
					Phone:      address.Phone,
					Address:    address.Detail,
					ExpiresAt:  time.Now().Add(30 * time.Minute),
				}
				groups[p.ShopID] = o
			}
			o.Amount += p.Price * c.Quantity
			o.Items = append(o.Items, model.OrderItem{
				Base:      base(),
				OrderID:   o.ID,
				ProductID: p.ID,
				Name:      p.Name,
				Price:     p.Price,
				Quantity:  c.Quantity,
			})
			if err := tx.AdjustStock(ctx, p.ID, -c.Quantity); err != nil {
				return err
			}
		}
		for _, id := range shopIDs {
			o := groups[id]
			o.Commission = o.Amount * shops[id].CommissionBPS / 10000
			if err := tx.CreateOrder(ctx, o); err != nil {
				return err
			}
			orders = append(orders, *o)
		}
		if err := tx.DeleteCart(ctx, a.ID, in.CartIDs); err != nil {
			return err
		}
		return audit(ctx, tx, a, "checkout.create", batch.ID, fmt.Sprintf("%d orders", len(orders)))
	})
	return orders, err
}
