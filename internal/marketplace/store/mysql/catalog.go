package mysql

import (
	"context"

	"gorm.io/gorm"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

func (s *datastore) Category(ctx context.Context, id string) (model.Category, error) {
	var category model.Category
	err := s.db.WithContext(ctx).First(&category, "id = ?", id).Error
	return category, normalizeError(err)
}

func (s *datastore) ListCategories(ctx context.Context, opts model.ListOptions) (model.Page[model.Category], error) {
	return list[model.Category](s.db.WithContext(ctx).Model(&model.Category{}), opts)
}

func (s *datastore) CreateCategory(ctx context.Context, category *model.Category) error {
	return normalizeError(s.db.WithContext(ctx).Create(category).Error)
}

func (s *datastore) Product(ctx context.Context, id string, locked bool) (model.Product, error) {
	var product model.Product
	err := lock(s.db.WithContext(ctx), locked).First(&product, "id = ?", id).Error
	return product, normalizeError(err)
}

func (s *datastore) CatalogProduct(ctx context.Context, id string) (model.Product, error) {
	var product model.Product
	db := s.db.WithContext(ctx)
	err := db.Model(&model.Product{}).
		Where("id = ? AND status = ? AND shop_id IN (?)", id, "approved",
			db.Model(&model.Shop{}).Select("id").Where("status = ?", "approved")).
		First(&product).Error
	return product, normalizeError(err)
}

func (s *datastore) ProductsByIDs(ctx context.Context, ids []string) ([]model.Product, error) {
	products := []model.Product{}
	err := s.db.WithContext(ctx).Where("id IN ?", ids).Find(&products).Error
	return products, normalizeError(err)
}

func (s *datastore) ListProducts(
	ctx context.Context,
	filter model.ProductFilter,
	opts model.ListOptions,
) (model.Page[model.Product], error) {
	db := s.db.WithContext(ctx)
	q := db.Model(&model.Product{})
	if filter.Catalog {
		q = q.Where(
			"status = ? AND shop_id IN (?)",
			"approved",
			db.Model(&model.Shop{}).Select("id").Where("status = ?", "approved"),
		)
	}
	if filter.OwnerID != "" {
		q = q.Where("shop_id IN (?)", db.Model(&model.Shop{}).Select("id").Where("owner_id = ?", filter.OwnerID))
	}
	if filter.Name != "" {
		q = q.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.CategoryID != "" {
		q = q.Where("category_id = ?", filter.CategoryID)
	}
	return list[model.Product](q, opts)
}

func (s *datastore) SaveProduct(ctx context.Context, product *model.Product) error {
	return normalizeError(s.db.WithContext(ctx).Save(product).Error)
}

func (s *datastore) AdjustStock(ctx context.Context, id string, delta int64) error {
	return normalizeError(
		s.db.WithContext(ctx).
			Model(&model.Product{}).
			Where("id = ?", id).
			Update("stock", gorm.Expr("stock + ?", delta)).
			Error,
	)
}
