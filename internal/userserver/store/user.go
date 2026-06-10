package store

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/types"
)

type UserStore interface {
	Get(ctx context.Context, username string, opts types.GetOptions) (*types.User, error)
	GetByEmail(ctx context.Context, email string, opts types.GetOptions) (*types.User, error)
	GetByPhone(ctx context.Context, phone string, opts types.GetOptions) (*types.User, error)
	GetByInstanceID(ctx context.Context, instanceID string, opts types.GetOptions) (*types.User, error)
	Create(ctx context.Context, user *types.User, opts types.CreateOptions) error
	Update(ctx context.Context, user *types.User, opts types.UpdateOptions) error
	Delete(ctx context.Context, instanceID string, opts types.DeleteOptions) error
	DeleteCollection(ctx context.Context, instanceIDs []string, opts types.DeleteOptions) error
	List(ctx context.Context, opts types.ListOptions) (*types.UserList, error)
}
