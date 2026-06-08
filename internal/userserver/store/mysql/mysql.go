package mysql

import (
	"fmt"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/userserver/store"
)

// ProviderSet is used by Wire.
var ProviderSet = wire.NewSet(NewDatastore)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Close() error {
	if ds.db == nil {
		return nil
	}
	d, err := ds.db.DB()
	if err != nil {
		return err
	}
	return d.Close()
}

// NewDatastore creates a MySQL-backed store.Factory.
func NewDatastore(opts *options.MySQLOptions) (store.Factory, error) {
	db, err := opts.NewClient()
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}
	return &datastore{db: db}, nil
}
