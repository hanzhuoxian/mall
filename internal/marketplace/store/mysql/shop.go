package mysql

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

func (s *datastore) Shop(ctx context.Context, id string, locked bool) (model.Shop, error) {
	var shop model.Shop
	err := lock(s.db.WithContext(ctx), locked).First(&shop, "id = ?", id).Error
	return shop, normalizeError(err)
}

func (s *datastore) ShopByOwner(ctx context.Context, ownerID string, locked bool) (model.Shop, error) {
	var shop model.Shop
	err := lock(s.db.WithContext(ctx), locked).Where("owner_id = ?", ownerID).First(&shop).Error
	return shop, normalizeError(err)
}

func (s *datastore) ShopsByOwner(ctx context.Context, ownerID string) ([]model.Shop, error) {
	shops := []model.Shop{}
	err := s.db.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&shops).Error
	return shops, normalizeError(err)
}

func (s *datastore) ListShops(
	ctx context.Context,
	ownerID string,
	opts model.ListOptions,
) (model.Page[model.Shop], error) {
	q := s.db.WithContext(ctx).Model(&model.Shop{})
	if ownerID != "" {
		q = q.Where("owner_id = ?", ownerID)
	}
	return list[model.Shop](q, opts)
}

func (s *datastore) SaveShop(ctx context.Context, shop *model.Shop) error {
	return normalizeError(s.db.WithContext(ctx).Save(shop).Error)
}
