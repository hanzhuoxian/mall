package service

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
)

// CreateAddress 创建属于当前用户的收货地址。
func (s *Service) CreateAddress(ctx context.Context, actor model.Actor, in model.Address) (model.Address, error) {
	if !validText(in.Recipient, model.MaxEntityNameLength) || !validText(in.Phone, model.MaxPhoneLength) || !validText(in.Detail, model.MaxAddressDetailLength) {
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
	if in.Quantity < 1 || in.Quantity > model.MaxCartQuantity {
		return ErrInvalid
	}
	product, err := s.store.Product(ctx, in.ProductID, false)
	if err != nil {
		return err
	}
	if product.Status != model.ProductStatusApproved {
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

// PutFavorite adds an approved product to the current user's favorites.
func (s *Service) PutFavorite(ctx context.Context, actor model.Actor, productID string) error {
	if _, err := s.store.CatalogProduct(ctx, productID); err != nil {
		return err
	}
	return s.store.PutFavorite(ctx, &model.Favorite{Base: base(), UserID: actor.ID, ProductID: productID})
}

// DeleteFavorite removes a product only from the current user's favorites.
func (s *Service) DeleteFavorite(ctx context.Context, actor model.Actor, productID string) error {
	return s.store.DeleteFavorite(ctx, actor.ID, productID)
}
