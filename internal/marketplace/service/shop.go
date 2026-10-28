package service

import (
	"context"
	"errors"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
)

// ApplyShop 创建入驻申请，或重新提交已经驳回的申请。
func (s *Service) ApplyShop(ctx context.Context, actor model.Actor, in model.ShopInput) (model.Shop, error) {
	if !validText(in.Name, 100) || !validText(in.Contact, 100) || !validText(in.License, 200) {
		return model.Shop{}, ErrInvalid
	}
	shop := model.Shop{
		Base:          base(),
		OwnerID:       actor.ID,
		Name:          in.Name,
		Contact:       in.Contact,
		License:       in.License,
		Status:        "pending",
		CommissionBPS: 500,
	}
	err := s.store.Transaction(ctx, func(tx store.Store) error {
		old, err := tx.ShopByOwner(ctx, actor.ID, true)
		if err == nil {
			if old.Status != "rejected" {
				return ErrInvalid
			}
			shop.Base = old.Base
		} else if !errors.Is(err, store.ErrNotFound) {
			return err
		}
		if err = tx.SaveShop(ctx, &shop); err != nil {
			return err
		}
		return audit(ctx, tx, actor, "shop.apply", shop.ID, shop.Name)
	})
	return shop, err
}

// ReviewShop 由平台审核商家状态并设置佣金比例。
func (s *Service) ReviewShop(ctx context.Context, actor model.Actor, id string, in model.ShopReview) error {
	if !actor.Admin {
		return ErrForbidden
	}
	if (in.Status != "approved" && in.Status != "rejected" && in.Status != "suspended") ||
		in.CommissionBPS < 0 || in.CommissionBPS > 10000 || len([]rune(in.Reason)) > 500 ||
		((in.Status == "rejected" || in.Status == "suspended") && !validText(in.Reason, 500)) {
		return ErrInvalid
	}
	return s.store.Transaction(ctx, func(tx store.Store) error {
		shop, err := tx.Shop(ctx, id, true)
		if err != nil {
			return err
		}
		shop.Status, shop.Reason, shop.CommissionBPS = in.Status, in.Reason, in.CommissionBPS
		if err = tx.SaveShop(ctx, &shop); err != nil {
			return err
		}
		return audit(ctx, tx, actor, "shop."+in.Status, shop.ID, in.Reason)
	})
}

// CreateCategory 由平台创建分类并记录审计。
func (s *Service) CreateCategory(ctx context.Context, actor model.Actor, name string) (model.Category, error) {
	if !actor.Admin {
		return model.Category{}, ErrForbidden
	}
	if !validText(name, 100) {
		return model.Category{}, ErrInvalid
	}
	category := model.Category{Base: base(), Name: name}
	err := s.store.Transaction(ctx, func(tx store.Store) error {
		if err := tx.CreateCategory(ctx, &category); err != nil {
			return err
		}
		return audit(ctx, tx, actor, "category.create", category.ID, category.Name)
	})
	return category, err
}
