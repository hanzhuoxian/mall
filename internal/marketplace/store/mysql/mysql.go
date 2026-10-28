// Package mysql implements marketplace persistence with GORM.
package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
)

type datastore struct{ db *gorm.DB }

// NewStore creates a store backed by the supplied database.
func NewStore(db *gorm.DB) store.Store { return &datastore{db: db} }

func (s *datastore) Transaction(ctx context.Context, fn func(store.Store) error) error {
	return normalizeError(s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&datastore{db: tx})
	}))
}

func (s *datastore) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func normalizeError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return store.ErrNotFound
	}
	return err
}

func lock(q *gorm.DB, enabled bool) *gorm.DB {
	if enabled {
		return q.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	return q
}

func list[T any](q *gorm.DB, opts model.ListOptions) (model.Page[T], error) {
	page := model.Page[T]{Items: []T{}}
	if err := q.Count(&page.Total).Error; err != nil {
		return page, normalizeError(err)
	}
	err := q.Order("created_at DESC, id").
		Offset((opts.Page - 1) * opts.PageSize).
		Limit(opts.PageSize).
		Find(&page.Items).
		Error
	return page, normalizeError(err)
}
