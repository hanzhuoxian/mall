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

func (user *users) Get(ctx context.Context, username string, opts model.GetOptions) (*model.SysUser, error) {
	u := &model.SysUser{}
	err := user.db.WithContext(ctx).Where("username = ?", username).First(u).Error
	return u, err
}

func (user *users) GetByEmail(ctx context.Context, email string, opts model.GetOptions) (*model.SysUser, error) {
	u := &model.SysUser{}
	err := user.db.WithContext(ctx).Where("email = ?", email).First(u).Error
	return u, err
}

func (user *users) GetByPhone(ctx context.Context, phone string, opts model.GetOptions) (*model.SysUser, error) {
	u := &model.SysUser{}
	err := user.db.WithContext(ctx).Where("phone = ?", phone).First(u).Error
	return u, err
}

func (user *users) GetByInstanceID(ctx context.Context, instanceID string, opts model.GetOptions) (*model.SysUser, error) {
	u := &model.SysUser{}
	err := user.db.WithContext(ctx).Where("instance_id = ?", instanceID).First(u).Error
	return u, err
}

func (user *users) Create(ctx context.Context, u *model.SysUser, opts model.CreateOptions) error {
	return user.db.WithContext(ctx).Create(u).Error
}

func (user *users) Update(ctx context.Context, u *model.SysUser, opts model.UpdateOptions) error {
	return user.db.WithContext(ctx).Save(u).Error
}

func (user *users) Delete(ctx context.Context, instanceID string, opts model.DeleteOptions) error {
	db := user.db.WithContext(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	return db.Where("instance_id = ?", instanceID).Delete(&model.SysUser{}).Error
}

func (user *users) DeleteCollection(ctx context.Context, instanceIDs []string, opts model.DeleteOptions) error {
	db := user.db.WithContext(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	return db.Where("instance_id IN ?", instanceIDs).Delete(&model.SysUser{}).Error
}

func (user *users) List(ctx context.Context, opts model.ListUserOptions) (*model.SysUserList, error) {
	userList := &model.SysUserList{}
	tx := user.db.WithContext(ctx).Model(&model.SysUser{})
	if opts.Keyword != "" {
		like := "%" + opts.Keyword + "%"
		tx = tx.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", like, like, like)
	}
	if opts.Status != nil {
		tx = tx.Where("user_status = ?", *opts.Status)
	}
	if err := tx.Count(&userList.TotalCount).Error; err != nil {
		return nil, err
	}
	tx = tx.Offset(int(*opts.Offset)).Limit(int(*opts.Limit))
	if err := tx.Find(&userList.Items).Error; err != nil {
		return nil, err
	}
	return userList, nil
}
