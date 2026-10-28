package service

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

// Context 返回用户的商家身份及沙箱支付配置。
func (s *Service) Context(ctx context.Context, actor model.Actor) (model.Context, error) {
	shops, err := s.store.ShopsByOwner(ctx, actor.ID)
	return model.Context{UserID: actor.ID, Admin: actor.Admin, Shops: shops, SandboxPayment: s.sandboxPayment}, err
}

// Catalog 查询已经审核且所属商家可营业的商品。
func (s *Service) Catalog(
	ctx context.Context,
	name, categoryID string,
	opts model.ListOptions,
) (model.Page[model.Product], error) {
	return s.store.ListProducts(
		ctx,
		model.ProductFilter{Name: name, CategoryID: categoryID, Catalog: true},
		normalizePage(opts),
	)
}

// CatalogProduct 返回公开目录中可见的商品。
func (s *Service) CatalogProduct(ctx context.Context, id string) (model.Product, error) {
	return s.store.CatalogProduct(ctx, id)
}

// Shops 返回平台或当前商家的店铺列表。
func (s *Service) Shops(
	ctx context.Context,
	actor model.Actor,
	opts model.ListOptions,
) (model.Page[model.Shop], error) {
	ownerID := actor.ID
	if actor.Admin {
		ownerID = ""
	}
	return s.store.ListShops(ctx, ownerID, normalizePage(opts))
}

// Categories 返回商品分类列表。
func (s *Service) Categories(ctx context.Context, opts model.ListOptions) (model.Page[model.Category], error) {
	return s.store.ListCategories(ctx, normalizePage(opts))
}

// Products 按名称及审核状态查询当前操作人可管理的商品。
func (s *Service) Products(
	ctx context.Context,
	actor model.Actor,
	name, status string,
	opts model.ListOptions,
) (model.Page[model.Product], error) {
	filter := model.ProductFilter{Name: name, Status: status}
	if !actor.Admin {
		filter.OwnerID = actor.ID
	}
	return s.store.ListProducts(ctx, filter, normalizePage(opts))
}

// Addresses 返回当前用户的收货地址。
func (s *Service) Addresses(
	ctx context.Context,
	actor model.Actor,
	opts model.ListOptions,
) (model.Page[model.Address], error) {
	return s.store.ListAddresses(ctx, actor.ID, normalizePage(opts))
}

// Cart 返回当前用户的购物车。
func (s *Service) Cart(
	ctx context.Context,
	actor model.Actor,
	opts model.ListOptions,
) (model.Page[model.CartItem], error) {
	return s.store.ListCart(ctx, actor.ID, normalizePage(opts))
}

func orderFilter(actor model.Actor, scope string) model.OrderFilter {
	switch scope {
	case "platform":
		return model.OrderFilter{DenyAll: !actor.Admin}
	case "merchant":
		return model.OrderFilter{OwnerID: actor.ID}
	default:
		return model.OrderFilter{UserID: actor.ID}
	}
}

// Orders 按买家、商家或平台范围查询订单。
func (s *Service) Orders(
	ctx context.Context,
	actor model.Actor,
	scope, status string,
	opts model.ListOptions,
) (model.Page[model.Order], error) {
	filter := orderFilter(actor, scope)
	filter.Status = status
	return s.store.ListOrders(ctx, filter, normalizePage(opts))
}

// Dashboard 返回当前操作人可见的经营汇总。
func (s *Service) Dashboard(ctx context.Context, actor model.Actor, scope string) (model.Summary, error) {
	return s.store.Summary(ctx, orderFilter(actor, scope))
}

// Audits 仅允许平台管理员读取审计记录。
func (s *Service) Audits(
	ctx context.Context,
	actor model.Actor,
	opts model.ListOptions,
) (model.Page[model.Audit], error) {
	if !actor.Admin {
		return model.Page[model.Audit]{}, ErrForbidden
	}
	return s.store.ListAudits(ctx, normalizePage(opts))
}
