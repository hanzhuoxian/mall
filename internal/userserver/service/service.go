package service

import (
	"github.com/google/wire"

	"github.com/hanzhuoxian/mall/internal/userserver/store"
)

// ProviderSet contains the user service dependency providers.
var ProviderSet = wire.NewSet(NewService)

// Service exposes user-related business operations.
type Service interface {
	Users() UserSrv
}

type service struct {
	store store.Factory
}

// NewService creates the user service backed by the provided store.
func NewService(store store.Factory) Service {
	return &service{store: store}
}

func (s *service) Users() UserSrv {
	return newUsers(s)
}
