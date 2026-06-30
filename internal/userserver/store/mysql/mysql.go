package mysql

import (
	"fmt"

	"github.com/google/wire"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"

	"github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/userserver/store"
)

const DBNameUser = "mall_user"

// ProviderSet is used by Wire.
var ProviderSet = wire.NewSet(NewDatastore)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *datastore) Roles() store.RoleStore {
	return newRoles(ds)
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
	// otelgorm 插件为每条 SQL 生成子 span（含语句、表名等属性），归属调用方传入的 ctx 链路。
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, fmt.Errorf("install gorm otel plugin: %w", err)
	}
	return &datastore{db: db}, nil
}
