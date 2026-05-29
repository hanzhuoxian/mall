package controller

import "github.com/google/wire"

// ProviderSet is used by Wire.
var ProviderSet = wire.NewSet(NewUserController, NewControllers)

type Controllers struct {
	User *UserController
}

func NewControllers(uc *UserController) *Controllers {
	return &Controllers{User: uc}
}
