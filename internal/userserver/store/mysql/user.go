package mysql

import (
	"context"

	"gorm.io/gorm"

	"github.com/hanzhuoxian/mall/internal/types"
)

type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{ds.db}
}

func (user *users) Get(ctx context.Context, username string, opts types.GetOptions) (*types.User, error) {
	u := &types.User{}
	err := user.db.First(u).Error
	return u, err
}

func (user *users) Create(ctx context.Context, u *types.User, opts types.CreateOptions) error {
	return user.db.Create(u).Error
}

func (user *users) Update(ctx context.Context, u *types.User, opts types.UpdateOptions) error {
	return user.db.Save(u).Error
}

func (user *users) Delete(ctx context.Context, instanceID string, opts types.DeleteOptions) error {
	if opts.Unscoped {
		user.db = user.db.Unscoped()
	}
	return user.db.Where("instance_id = ?", instanceID).Delete(&types.User{}).Error
}

func (user *users) DeleteCollection(ctx context.Context, instanceIDs []string, opts types.DeleteOptions) error {
	if opts.Unscoped {
		user.db = user.db.Unscoped()
	}
	return user.db.Where("instance_id IN ?", instanceIDs).Delete(&types.User{}).Error
}

func (user *users) List(ctx context.Context, opts types.ListOptions) (*types.UserList, error) {
	userList := &types.UserList{}
	return userList, nil
}
