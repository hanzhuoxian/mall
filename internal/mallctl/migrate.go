package mallctl

import (
	"fmt"

	"github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/types"
	"github.com/hanzhuoxian/mall/pkg/app"
	"github.com/hanzhuoxian/mall/pkg/nflag"
)

type migrateOptions struct {
	MySQL *options.MySQLOptions
}

func newMigrateOptions() *migrateOptions {
	return &migrateOptions{
		MySQL: options.NewMySQLOptions(),
	}
}

func (o *migrateOptions) Flags() nflag.NamedFlagSets {
	var fss nflag.NamedFlagSets
	o.MySQL.AddFlags(fss.FlagSet("mysql"))
	return fss
}

func (o *migrateOptions) Validate() []error {
	return o.MySQL.Validate()
}

func newMigrateCmd() *app.Command {
	opts := newMigrateOptions()
	return app.NewCommand(
		"migrate",
		"Migrate database tables for the mall platform",
		app.WithCommandOptions(opts),
		app.WithCommandRunFunc(func(args []string) error {
			return runMigrate(opts)
		}),
	)
}

func runMigrate(opts *migrateOptions) error {
	userModels := []any{
		&types.User{},
		&types.SysRole{},
	}
	if err := migrateDB(opts.MySQL, "mall_user", userModels...); err != nil {
		return err
	}

	return nil
}

func migrateDB(base *options.MySQLOptions, database string, models ...any) error {
	o := *base
	o.Database = database
	db, err := o.NewClient()
	if err != nil {
		return fmt.Errorf("connect %s: %w", database, err)
	}

	return db.AutoMigrate(models...)
}
