package mysql

import (
	"context"

	"gorm.io/gorm/clause"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

func (s *datastore) Address(ctx context.Context, id, userID string) (model.Address, error) {
	var address model.Address
	err := s.db.WithContext(ctx).First(&address, "id = ? AND user_id = ?", id, userID).Error
	return address, normalizeError(err)
}

func (s *datastore) ListAddresses(
	ctx context.Context,
	userID string,
	opts model.ListOptions,
) (model.Page[model.Address], error) {
	return list[model.Address](s.db.WithContext(ctx).Model(&model.Address{}).Where("user_id = ?", userID), opts)
}

func (s *datastore) CreateAddress(ctx context.Context, address *model.Address) error {
	return normalizeError(s.db.WithContext(ctx).Create(address).Error)
}

func (s *datastore) DeleteAddress(ctx context.Context, id, userID string) error {
	return normalizeError(
		s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Address{}).Error,
	)
}

func (s *datastore) ListCart(
	ctx context.Context,
	userID string,
	opts model.ListOptions,
) (model.Page[model.CartItem], error) {
	return list[model.CartItem](
		s.db.WithContext(ctx).Model(&model.CartItem{}).Preload("Product").Where("user_id = ?", userID),
		opts,
	)
}

func (s *datastore) CartByIDs(ctx context.Context, userID string, ids []string) ([]model.CartItem, error) {
	cart := []model.CartItem{}
	err := lock(
		s.db.WithContext(ctx),
		true,
	).Where("user_id = ? AND id IN ?", userID, ids).
		Order("product_id").
		Find(&cart).
		Error
	return cart, normalizeError(err)
}

func (s *datastore) PutCart(ctx context.Context, item *model.CartItem) error {
	return normalizeError(s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity", "updated_at"}),
	}).Create(item).Error)
}

func (s *datastore) DeleteCart(ctx context.Context, userID string, ids []string) error {
	return normalizeError(
		s.db.WithContext(ctx).Where("user_id = ? AND id IN ?", userID, ids).Delete(&model.CartItem{}).Error,
	)
}

func (s *datastore) ListFavorites(ctx context.Context, userID string, opts model.ListOptions) (model.Page[model.Favorite], error) {
	db := s.db.WithContext(ctx)
	query := db.Model(&model.Favorite{}).Preload("Product").Where("user_id = ? AND product_id IN (?)", userID,
		db.Model(&model.Product{}).Select("id").Where("status = ? AND shop_id IN (?)", model.ProductStatusApproved,
			db.Model(&model.Shop{}).Select("id").Where("status = ?", model.ShopStatusApproved)))
	return list[model.Favorite](query, opts)
}

func (s *datastore) PutFavorite(ctx context.Context, favorite *model.Favorite) error {
	return normalizeError(s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "product_id"}}, DoNothing: true,
	}).Create(favorite).Error)
}

func (s *datastore) DeleteFavorite(ctx context.Context, userID, productID string) error {
	return normalizeError(s.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).Delete(&model.Favorite{}).Error)
}
