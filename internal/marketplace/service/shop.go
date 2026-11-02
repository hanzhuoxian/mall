package service

import (
	"context"
	"errors"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
)

// ApplyShop 创建入驻申请，或重新提交已经驳回的申请。
func (s *Service) ApplyShop(ctx context.Context, actor model.Actor, in model.ShopInput) (model.Shop, error) {
	if !validText(in.Name, model.MaxEntityNameLength) || !validText(in.Contact, model.MaxContactLength) || !validText(in.License, model.MaxLicenseLength) {
		return model.Shop{}, ErrInvalid
	}
	shop := model.Shop{
		Base:          base(),
		OwnerID:       actor.ID,
		Name:          in.Name,
		Contact:       in.Contact,
		License:       in.License,
		Status:        model.ShopStatusPending,
		CommissionBPS: model.DefaultCommissionBPS,
	}
	err := s.store.Transaction(ctx, func(tx store.Store) error {
		old, err := tx.ShopByOwner(ctx, actor.ID, true)
		if err == nil {
			if old.Status != model.ShopStatusRejected {
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
	if (in.Status != model.ShopStatusApproved && in.Status != model.ShopStatusRejected && in.Status != model.ShopStatusSuspended) ||
		in.CommissionBPS < 0 || in.CommissionBPS > model.CommissionBasisPoints || len([]rune(in.Reason)) > model.MaxReasonLength ||
		((in.Status == model.ShopStatusRejected || in.Status == model.ShopStatusSuspended) && !validText(in.Reason, model.MaxReasonLength)) {
		return ErrInvalid
	}
	err := s.store.Transaction(ctx, func(tx store.Store) error {
		shop, err := tx.Shop(ctx, id, true)
		if err != nil {
			return err
		}
		shop.Status, shop.Reason, shop.CommissionBPS = in.Status, in.Reason, in.CommissionBPS
		if err = tx.SaveShop(ctx, &shop); err != nil {
			return err
		}
		products, err := tx.ProductsByShop(ctx, shop.ID, model.ListOptions{Page: 1, PageSize: model.MaxPageSize})
		if err != nil {
			return err
		}
		operation := model.SearchOperationDelete
		if shop.Status == model.ShopStatusApproved {
			operation = model.SearchOperationIndex
		}
		for _, product := range products.Items {
			if err = searchEvent(ctx, tx, product.ID, operation); err != nil {
				return err
			}
		}
		return audit(ctx, tx, actor, "shop."+in.Status, shop.ID, in.Reason)
	})
	return err
}

// CreateCategory 由平台创建分类并记录审计。
func (s *Service) CreateCategory(ctx context.Context, actor model.Actor, name string) (model.Category, error) {
	if !actor.Admin {
		return model.Category{}, ErrForbidden
	}
	if !validText(name, model.MaxEntityNameLength) {
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
