// Package service 实现商城业务校验、权限控制及事务编排。
package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
)

const (
	maxProductDescriptionLength = 10_000
	maxProductImageLength       = 1_000
	maxPageNumber               = 100_000
	defaultPageSize             = 20
	checkoutOrderLifetime       = 30 * time.Minute
	defaultActorID              = "system"
)

var (
	// ErrInvalid 表示请求参数或资源状态不允许当前操作。
	ErrInvalid = errors.New("参数无效或当前状态不允许此操作")
	// ErrForbidden 表示当前操作人无权访问资源。
	ErrForbidden = errors.New("没有权限操作此资源")
	// ErrStock 表示商品无法购买或库存不足。
	ErrStock = errors.New("商品已下架或库存不足")
	// ErrPaymentDisabled 表示沙箱付款及退款未启用。
	ErrPaymentDisabled = errors.New("未配置真实支付，沙箱支付当前关闭")
)

// Service 通过存储接口执行业务操作，不持有 HTTP 或数据库实现。
type Service struct {
	store          store.Store
	search         SearchIndex
	sandboxPayment bool
}

// SearchIndex indexes and queries public catalog documents.
type SearchIndex interface {
	Search(context.Context, string, model.ListOptions) (model.Page[model.Product], error)
	Index(context.Context, model.Product) error
	Delete(context.Context, string) error
}

// NewService 使用指定存储创建商城服务。
func NewService(ds store.Store, sandboxPayment bool) *Service {
	return &Service{store: ds, sandboxPayment: sandboxPayment}
}

// SetSearchIndex enables Elasticsearch backed catalog search.
func (s *Service) SetSearchIndex(index SearchIndex) { s.search = index }

// Close 释放商城存储资源。
func (s *Service) Close() error { return s.store.Close() }

func base() model.Base { return model.Base{ID: uuid.NewString()} }

func audit(ctx context.Context, tx store.Store, actor model.Actor, action, target, detail string) error {
	return tx.CreateAudit(ctx, &model.Audit{
		Base:     base(),
		ActorID:  actor.ID,
		Action:   action,
		TargetID: target,
		Detail:   detail,
	})
}

func searchEvent(ctx context.Context, tx store.Store, productID, operation string) error {
	return tx.CreateSearchOutbox(ctx, &model.SearchOutbox{Base: base(), ProductID: productID, Operation: operation})
}

func ownedShop(ctx context.Context, tx store.Store, actor model.Actor, id string) (model.Shop, error) {
	shop, err := tx.Shop(ctx, id, true)
	if err != nil {
		return shop, err
	}
	if shop.OwnerID != actor.ID && !actor.Admin {
		return shop, ErrForbidden
	}
	return shop, nil
}

func validText(value string, maxLength int) bool {
	return strings.TrimSpace(value) != "" && len([]rune(value)) <= maxLength
}

func validProduct(product model.Product) bool {
	return validText(product.Name, model.MaxProductNameLength) &&
		len([]rune(product.Description)) <= maxProductDescriptionLength &&
		len(product.Image) <= maxProductImageLength &&
		(product.Image == "" || strings.HasPrefix(product.Image, "https://")) &&
		product.Price > 0 && product.Price <= model.MaxProductPriceCents &&
		product.Stock >= 0 && product.Stock <= model.MaxProductStock
}

func normalizePage(opts model.ListOptions) model.ListOptions {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Page > maxPageNumber {
		opts.Page = maxPageNumber
	}
	if opts.PageSize < 1 || opts.PageSize > model.MaxPageSize {
		opts.PageSize = defaultPageSize
	}
	return opts
}
