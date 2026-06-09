package service

import "github.com/appleboy/gin-jwt/v3/store"

type UserSrv interface {
}

type userSrv struct {
	store store.Factory
}

func newUsers(srv *service) UserSrv {
	return &userSrv{
		store: srv.store,
	}
}
