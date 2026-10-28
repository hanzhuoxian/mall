package mysql

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

func (s *datastore) CreateCheckout(ctx context.Context, checkout *model.Checkout) error {
	return normalizeError(s.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(checkout).Error)
}

func (s *datastore) CheckoutByKey(ctx context.Context, userID, key string) (model.Checkout, error) {
	var checkout model.Checkout
	err := lock(
		s.db.WithContext(ctx),
		true,
	).Where("user_id = ? AND request_key = ?", userID, key).
		First(&checkout).
		Error
	return checkout, normalizeError(err)
}

func preloadItems(q *gorm.DB) *gorm.DB {
	return q.Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("id") })
}

func (s *datastore) OrdersByCheckout(ctx context.Context, checkoutID string) ([]model.Order, error) {
	orders := []model.Order{}
	err := preloadItems(s.db.WithContext(ctx)).Where("checkout_id = ?", checkoutID).Order("id").Find(&orders).Error
	return orders, normalizeError(err)
}

func (s *datastore) Order(ctx context.Context, id string, locked bool) (model.Order, error) {
	var order model.Order
	err := lock(s.db.WithContext(ctx), locked).First(&order, "id = ?", id).Error
	return order, normalizeError(err)
}

func (s *datastore) orderScope(ctx context.Context, filter model.OrderFilter) *gorm.DB {
	db := s.db.WithContext(ctx)
	q := db.Model(&model.Order{})
	if filter.DenyAll {
		q = q.Where("1 = 0")
	}
	if filter.OwnerID != "" {
		q = q.Where("shop_id IN (?)", db.Model(&model.Shop{}).Select("id").Where("owner_id = ?", filter.OwnerID))
	}
	if filter.UserID != "" {
		q = q.Where("user_id = ?", filter.UserID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	return q
}

func (s *datastore) ListOrders(
	ctx context.Context,
	filter model.OrderFilter,
	opts model.ListOptions,
) (model.Page[model.Order], error) {
	return list[model.Order](preloadItems(s.orderScope(ctx, filter)), opts)
}

func (s *datastore) CreateOrder(ctx context.Context, order *model.Order) error {
	return normalizeError(s.db.WithContext(ctx).Create(order).Error)
}

func (s *datastore) SaveOrder(ctx context.Context, order *model.Order) error {
	return normalizeError(s.db.WithContext(ctx).Save(order).Error)
}

func (s *datastore) SetOrderStatus(ctx context.Context, id, status string) error {
	return normalizeError(
		s.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error,
	)
}

func (s *datastore) OrderItems(ctx context.Context, orderID string) ([]model.OrderItem, error) {
	items := []model.OrderItem{}
	err := s.db.WithContext(ctx).Where("order_id = ?", orderID).Order("product_id").Find(&items).Error
	return items, normalizeError(err)
}

func (s *datastore) ExpiredOrderIDs(ctx context.Context, now time.Time, limit int) ([]string, error) {
	ids := []string{}
	err := s.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("status = ? AND expires_at <= ?", "pending", now).
		Limit(limit).
		Pluck("id", &ids).
		Error
	return ids, normalizeError(err)
}

func (s *datastore) Summary(ctx context.Context, filter model.OrderFilter) (model.Summary, error) {
	var summary model.Summary
	if err := s.orderScope(ctx, filter).Count(&summary.Orders).Error; err != nil {
		return summary, normalizeError(err)
	}
	var sums struct {
		GMV        int64
		Commission int64
	}
	err := s.orderScope(ctx, filter).
		Where("status IN ?", []string{"paid", "shipped", "completed"}).
		Select("COALESCE(SUM(amount),0) AS gmv, COALESCE(SUM(commission),0) AS commission").
		Scan(&sums).
		Error
	if err != nil {
		return summary, normalizeError(err)
	}
	summary.GMV = sums.GMV
	summary.Commission = sums.Commission
	err = s.orderScope(ctx, filter).
		Where("status = ?", "completed").
		Select("COALESCE(SUM(amount - commission),0)").
		Scan(&summary.Payable).
		Error
	return summary, normalizeError(err)
}
