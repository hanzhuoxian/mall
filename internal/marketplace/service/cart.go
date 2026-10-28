package service

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
)

// CreateAddress 创建属于当前用户的收货地址。
func (s *Service) CreateAddress(ctx context.Context, actor model.Actor, in model.Address) (model.Address, error) {
	if !validText(in.Recipient, 100) || !validText(in.Phone, 30) || !validText(in.Detail, 500) {
		return model.Address{}, ErrInvalid
	}
	in.Base, in.UserID = base(), actor.ID
	err := s.store.CreateAddress(ctx, &in)
	return in, err
}

// DeleteAddress 仅删除当前用户自己的地址。
func (s *Service) DeleteAddress(ctx context.Context, actor model.Actor, id string) error {
	return s.store.DeleteAddress(ctx, id, actor.ID)
}

// PutCart 校验商品及数量后更新当前用户购物车。
func (s *Service) PutCart(ctx context.Context, actor model.Actor, in model.CartItem) error {
	if in.Quantity < 1 || in.Quantity > 999 {
		return ErrInvalid
	}
	product, err := s.store.Product(ctx, in.ProductID, false)
	if err != nil {
		return err
	}
	if product.Status != "approved" {
		return store.ErrNotFound
	}
	if product.Stock < in.Quantity {
		return ErrStock
	}
	in.Base, in.UserID, in.Product = base(), actor.ID, nil
	return s.store.PutCart(ctx, &in)
}

// DeleteCart 仅删除当前用户自己的购物车条目。
func (s *Service) DeleteCart(ctx context.Context, actor model.Actor, id string) error {
	return s.store.DeleteCart(ctx, actor.ID, []string{id})
}
