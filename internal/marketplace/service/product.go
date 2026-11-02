package service

import (
	"context"
	"strconv"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
)

// SaveProduct 保存商品并重新进入审核，不通过编辑覆盖已预留库存。
func (s *Service) SaveProduct(ctx context.Context, a model.Actor, p model.Product) (model.Product, error) {
	err := s.store.Transaction(ctx, func(tx store.Store) error {
		shop, err := ownedShop(ctx, tx, a, p.ShopID)
		if err != nil {
			return err
		}
		if shop.Status != model.ShopStatusApproved {
			return ErrForbidden
		}
		if !validProduct(p) {
			return ErrInvalid
		}
		if _, err = tx.Category(ctx, p.CategoryID); err != nil {
			return err
		}
		if p.ID != "" {
			var old model.Product
			old, err = tx.Product(ctx, p.ID, true)
			if err != nil {
				return err
			}
			if old.ShopID != p.ShopID {
				return store.ErrNotFound
			}
			p.Base = old.Base
			// Stock is changed only through the inventory endpoint; editing cannot overwrite reservations.
			p.Stock = old.Stock
		} else {
			p.Base = base()
		}
		p.Status = model.ProductStatusPending
		p.Reason = ""
		if err = tx.SaveProduct(ctx, &p); err != nil {
			return err
		}
		if err = searchEvent(ctx, tx, p.ID, model.SearchOperationDelete); err != nil {
			return err
		}
		return audit(ctx, tx, a, "product.submit", p.ID, p.Name)
	})
	return p, err
}

// ReviewProduct 审核或下架商品，普通商家仅能下架自己的商品。
func (s *Service) ReviewProduct(ctx context.Context, actor model.Actor, id string, in model.ProductReview) error {
	if in.Status != model.ProductStatusApproved && in.Status != model.ProductStatusRejected && in.Status != model.ProductStatusOffline {
		return ErrInvalid
	}
	if !actor.Admin && in.Status != model.ProductStatusOffline {
		return ErrForbidden
	}
	if len([]rune(in.Reason)) > model.MaxReasonLength || (in.Status == model.ProductStatusRejected && !validText(in.Reason, model.MaxReasonLength)) {
		return ErrInvalid
	}
	err := s.store.Transaction(ctx, func(tx store.Store) error {
		product, err := tx.Product(ctx, id, false)
		if err != nil {
			return err
		}
		shop, err := ownedShop(ctx, tx, actor, product.ShopID)
		if err != nil {
			return err
		}
		product, err = tx.Product(ctx, product.ID, true)
		if err != nil {
			return err
		}
		if in.Status == model.ProductStatusApproved && (shop.Status != model.ShopStatusApproved || product.Status != model.ProductStatusPending) {
			return ErrInvalid
		}
		product.Status, product.Reason = in.Status, in.Reason
		if err = tx.SaveProduct(ctx, &product); err != nil {
			return err
		}
		if in.Status == model.ProductStatusApproved {
			err = searchEvent(ctx, tx, product.ID, model.SearchOperationIndex)
		} else {
			err = searchEvent(ctx, tx, product.ID, model.SearchOperationDelete)
		}
		if err != nil {
			return err
		}
		return audit(ctx, tx, actor, "product."+in.Status, product.ID, in.Reason)
	})
	return err
}

// AdjustStock 在锁定商家及商品后调整库存。
func (s *Service) AdjustStock(ctx context.Context, actor model.Actor, id string, delta int64) error {
	if delta == 0 || delta > model.MaxStockAdjustment || delta < -model.MaxStockAdjustment {
		return ErrInvalid
	}
	return s.store.Transaction(ctx, func(tx store.Store) error {
		product, err := tx.Product(ctx, id, false)
		if err != nil {
			return err
		}
		if _, err = ownedShop(ctx, tx, actor, product.ShopID); err != nil {
			return err
		}
		product, err = tx.Product(ctx, product.ID, true)
		if err != nil {
			return err
		}
		if product.Stock+delta < 0 || product.Stock+delta > model.MaxProductStock {
			return ErrInvalid
		}
		if err = tx.AdjustStock(ctx, product.ID, delta); err != nil {
			return err
		}
		return audit(ctx, tx, actor, "stock.adjust", product.ID, strconv.FormatInt(delta, 10))
	})
}
