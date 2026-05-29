package mysql

import (
	"errors"
	"sync"

	"github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/userserver/store"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Close() error {
	db := ds.db
	if db == nil {
		return nil
	}
	d, err := db.DB()

	if err != nil {
		return err
	}
	return d.Close()
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

func GetMySQLFactoryOr(opts *options.MySQLOptions) (store.Factory, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, errors.New("failed to get mysql factory")
	}
	var err error
	once.Do(func() {
		var db *gorm.DB
		db, err = opts.NewClient()
		if err != nil {
			return
		}
		mysqlFactory = &datastore{db}
	})

	if err != nil {
		return nil, err
	}
	return mysqlFactory, nil
}
