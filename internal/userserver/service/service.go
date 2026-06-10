package service

import (
	"github.com/google/wire"
	"github.com/hanzhuoxian/mall/internal/userserver/store"
)

var ProviderSet = wire.NewSet(NewService)

type Service interface {
	Users() UserSrv
}

type service struct {
	store store.Factory
}

func NewService(store store.Factory) Service {
	return &service{store: store}
}

func (s *service) Users() UserSrv {
	return newUsers(s)
}
