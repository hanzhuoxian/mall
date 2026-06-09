package service

import "github.com/appleboy/gin-jwt/v3/store"

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
