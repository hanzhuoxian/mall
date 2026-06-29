package store

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/userserver/model"
)

type UserStore interface {
	Get(ctx context.Context, username string, opts model.GetOptions) (*model.SysUser, error)
	GetByEmail(ctx context.Context, email string, opts model.GetOptions) (*model.SysUser, error)
	GetByPhone(ctx context.Context, phone string, opts model.GetOptions) (*model.SysUser, error)
	GetByInstanceID(ctx context.Context, instanceID string, opts model.GetOptions) (*model.SysUser, error)
	Create(ctx context.Context, user *model.SysUser, opts model.CreateOptions) error
	Update(ctx context.Context, user *model.SysUser, opts model.UpdateOptions) error
	Delete(ctx context.Context, instanceID string, opts model.DeleteOptions) error
	DeleteCollection(ctx context.Context, instanceIDs []string, opts model.DeleteOptions) error
	List(ctx context.Context, opts model.ListOptions) (*model.SysUserList, error)
}
