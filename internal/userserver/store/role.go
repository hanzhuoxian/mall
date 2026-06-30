package store

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/userserver/model"
)

type RoleStore interface {
	GetByInstanceID(ctx context.Context, instanceID string, opts model.GetOptions) (*model.SysRole, error)
	Create(ctx context.Context, u *model.SysRole, opts model.CreateOptions) error
	Update(ctx context.Context, u *model.SysRole, opts model.UpdateOptions) error
	Delete(ctx context.Context, instanceID string, opts model.DeleteOptions) error
	DeleteCollection(ctx context.Context, instanceIDs []string, opts model.DeleteOptions) error
	List(ctx context.Context, opts model.ListRoleOptions) (*model.SysRoleList, error)
}
