package mysql

import (
	"context"

	"gorm.io/gorm"

	"github.com/hanzhuoxian/mall/internal/userserver/model"
)

type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{ds.db}
}

func (user *users) Get(ctx context.Context, username string, opts model.GetOptions) (*model.User, error) {
	u := &model.User{}
	err := user.db.Where("username = ?", username).First(u).Error
	return u, err
}

func (user *users) GetByEmail(ctx context.Context, email string, opts model.GetOptions) (*model.User, error) {
	u := &model.User{}
	err := user.db.Where("email = ?", email).First(u).Error
	return u, err
}

func (user *users) GetByPhone(ctx context.Context, phone string, opts model.GetOptions) (*model.User, error) {
	u := &model.User{}
	err := user.db.Where("phone = ?", phone).First(u).Error
	return u, err
}

func (user *users) GetByInstanceID(ctx context.Context, instanceID string, opts model.GetOptions) (*model.User, error) {
	u := &model.User{}
	err := user.db.Where("instance_id = ?", instanceID).First(u).Error
	return u, err
}

func (user *users) Create(ctx context.Context, u *model.User, opts model.CreateOptions) error {
	return user.db.Create(u).Error
}

func (user *users) Update(ctx context.Context, u *model.User, opts model.UpdateOptions) error {
	return user.db.Save(u).Error
}

func (user *users) Delete(ctx context.Context, instanceID string, opts model.DeleteOptions) error {
	if opts.Unscoped {
		user.db = user.db.Unscoped()
	}
	return user.db.Where("instance_id = ?", instanceID).Delete(&model.User{}).Error
}

func (user *users) DeleteCollection(ctx context.Context, instanceIDs []string, opts model.DeleteOptions) error {
	if opts.Unscoped {
		user.db = user.db.Unscoped()
	}
	return user.db.Where("instance_id IN ?", instanceIDs).Delete(&model.User{}).Error
}

func (user *users) List(ctx context.Context, opts model.ListOptions) (*model.UserList, error) {
	userList := &model.UserList{}
	return userList, nil
}
