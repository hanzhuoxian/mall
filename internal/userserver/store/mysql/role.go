package mysql

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/userserver/model"
	"gorm.io/gorm"
)

type roles struct {
	db *gorm.DB
}

func newRoles(ds *datastore) *roles {
	return &roles{ds.db}
}

func (role *roles) GetByInstanceID(ctx context.Context, instanceID string, opts model.GetOptions) (*model.SysRole, error) {
	r := &model.SysRole{}
	err := role.db.WithContext(ctx).Where("instance_id = ?", instanceID).First(r).Error
	return r, err
}

func (role *roles) Create(ctx context.Context, u *model.SysRole, opts model.CreateOptions) error {
	return role.db.WithContext(ctx).Create(u).Error
}

func (role *roles) Update(ctx context.Context, u *model.SysRole, opts model.UpdateOptions) error {
	return role.db.WithContext(ctx).Save(u).Error
}

func (role *roles) Delete(ctx context.Context, instanceID string, opts model.DeleteOptions) error {
	db := role.db.WithContext(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	return db.Where("instance_id = ?", instanceID).Delete(&model.SysRole{}).Error
}

func (role *roles) DeleteCollection(ctx context.Context, instanceIDs []string, opts model.DeleteOptions) error {
	db := role.db.WithContext(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	return db.Where("instance_id IN ?", instanceIDs).Delete(&model.SysRole{}).Error
}

func (role *roles) List(ctx context.Context, opts model.ListRoleOptions) (*model.SysRoleList, error) {
	roleList := &model.SysRoleList{}
	tx := role.db.WithContext(ctx).Model(&model.SysRole{})
	if opts.Keyword != "" {
		like := "%" + opts.Keyword + "%"
		tx = tx.Where("role_code LIKE ? OR name LIKE ?", like, like)
	}
	if opts.RoleStatus != nil {
		tx = tx.Where("role_status = ?", *opts.RoleStatus)
	}
	if err := tx.Count(&roleList.TotalCount).Error; err != nil {
		return nil, err
	}
	if err := tx.Offset(int(*opts.Offset)).Limit(int(*opts.Limit)).Find(&roleList.Items).Error; err != nil {
		return nil, err
	}
	return roleList, nil
}
