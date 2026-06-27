package mallctl

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/hanzhuoxian/mall/internal/pkg/admin"
	"github.com/hanzhuoxian/mall/internal/pkg/options"
	"github.com/hanzhuoxian/mall/internal/types"
	"github.com/hanzhuoxian/mall/pkg/app"
	"github.com/hanzhuoxian/mall/pkg/auth"
	"github.com/hanzhuoxian/mall/pkg/nflag"
)

type initAdminOptions struct {
	MySQL    *options.MySQLOptions
	Username string
	Password string
	Email    string
	Nickname string
}

func newInitAdminOptions() *initAdminOptions {
	return &initAdminOptions{
		MySQL:    options.NewMySQLOptions(),
		Username: admin.AdminUserName,
		Nickname: "Admin",
	}
}

func (o *initAdminOptions) Flags() nflag.NamedFlagSets {
	var fss nflag.NamedFlagSets
	o.MySQL.AddFlags(fss.FlagSet("mysql"))

	fs := fss.FlagSet("admin")
	fs.StringVar(&o.Username, "username", o.Username, "Admin account username")
	fs.StringVar(&o.Password, "password", o.Password, "Admin account password (required)")
	fs.StringVar(&o.Email, "email", o.Email, "Admin account email")
	fs.StringVar(&o.Nickname, "nickname", o.Nickname, "Admin account nickname")
	return fss
}

func (o *initAdminOptions) Validate() []error {
	var errs []error
	errs = append(errs, o.MySQL.Validate()...)
	if o.Password == "" {
		errs = append(errs, fmt.Errorf("--password is required"))
	}
	if o.Username == "" {
		errs = append(errs, fmt.Errorf("--username is required"))
	}
	return errs
}

func newInitAdminCmd() *app.Command {
	opts := newInitAdminOptions()
	return app.NewCommand(
		"init-admin",
		"Initialize the admin account in the database",
		app.WithCommandOptions(opts),
		app.WithCommandRunFunc(func(args []string) error {
			return runInitAdmin(opts)
		}),
	)
}

func runInitAdmin(opts *initAdminOptions) error {
	o := *opts.MySQL
	o.Database = "mall_user"
	db, err := o.NewClient()
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	ctx := context.Background()

	existing := &types.User{}
	err = db.WithContext(ctx).Where("username = ?", opts.Username).First(existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("query user: %w", err)
	}
	if err == nil {
		fmt.Printf("Admin account %q already exists (instance_id: %s), skipping.\n", opts.Username, existing.InstanceID)
		return nil
	}

	hashed, err := auth.GeneratePassword(opts.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user := &types.User{
		ObjectMeta: types.ObjectMeta{
			InstanceID: uuid.New().String(),
			Name:       opts.Username,
		},
		Username: opts.Username,
		Nickname: opts.Nickname,
		Email:    opts.Email,
		Password: hashed,
		Status:   1,
	}

	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("create admin: %w", err)
	}

	fmt.Printf("Admin account %q created successfully (instance_id: %s).\n", opts.Username, user.InstanceID)
	return nil
}
