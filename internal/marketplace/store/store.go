// Package store 定义商城的持久化接口。
package store

import (
	"context"
	"errors"
	"time"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

// ErrNotFound 表示指定资源不存在，与具体数据库实现无关。
var ErrNotFound = errors.New("record not found")

// Store 封装商城存储操作；事务回调内的所有操作使用同一数据库事务。
type Store interface {
	Transaction(context.Context, func(Store) error) error
	Close() error

	Shop(context.Context, string, bool) (model.Shop, error)
	ShopByOwner(context.Context, string, bool) (model.Shop, error)
	ShopsByOwner(context.Context, string) ([]model.Shop, error)
	ListShops(context.Context, string, model.ListOptions) (model.Page[model.Shop], error)
	SaveShop(context.Context, *model.Shop) error

	Category(context.Context, string) (model.Category, error)
	ListCategories(context.Context, model.ListOptions) (model.Page[model.Category], error)
	CreateCategory(context.Context, *model.Category) error

	Product(context.Context, string, bool) (model.Product, error)
	CatalogProduct(context.Context, string) (model.Product, error)
	ProductsByIDs(context.Context, []string) ([]model.Product, error)
	ListProducts(context.Context, model.ProductFilter, model.ListOptions) (model.Page[model.Product], error)
	SaveProduct(context.Context, *model.Product) error
	AdjustStock(context.Context, string, int64) error

	Address(context.Context, string, string) (model.Address, error)
	ListAddresses(context.Context, string, model.ListOptions) (model.Page[model.Address], error)
	CreateAddress(context.Context, *model.Address) error
	DeleteAddress(context.Context, string, string) error

	ListCart(context.Context, string, model.ListOptions) (model.Page[model.CartItem], error)
	CartByIDs(context.Context, string, []string) ([]model.CartItem, error)
	PutCart(context.Context, *model.CartItem) error
	DeleteCart(context.Context, string, []string) error

	CreateCheckout(context.Context, *model.Checkout) error
	CheckoutByKey(context.Context, string, string) (model.Checkout, error)
	OrdersByCheckout(context.Context, string) ([]model.Order, error)
	Order(context.Context, string, bool) (model.Order, error)
	ListOrders(context.Context, model.OrderFilter, model.ListOptions) (model.Page[model.Order], error)
	CreateOrder(context.Context, *model.Order) error
	SaveOrder(context.Context, *model.Order) error
	SetOrderStatus(context.Context, string, string) error
	OrderItems(context.Context, string) ([]model.OrderItem, error)
	ExpiredOrderIDs(context.Context, time.Time, int) ([]string, error)
	Summary(context.Context, model.OrderFilter) (model.Summary, error)

	CreateAudit(context.Context, *model.Audit) error
	ListAudits(context.Context, model.ListOptions) (model.Page[model.Audit], error)
}
