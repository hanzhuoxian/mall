package controller

import (
	"github.com/google/wire"

	marketcontroller "github.com/hanzhuoxian/mall/internal/marketplace/controller"
)

// ProviderSet is used by Wire.
var ProviderSet = wire.NewSet(
	NewUserController,
	NewCaptchaController,
	NewRegisterController,
	NewControllers,
	NewMarketplaceService,
	NewMarketplaceController,
)

// Controllers 汇聚所有控制器实例，便于在路由注册时统一传递。
type Controllers struct {
	Marketplace *marketcontroller.MarketplaceController
	User        *UserController
	Captcha     *CaptchaController
	Register    *RegisterController
}

// NewControllers 创建并返回包含所有控制器的 Controllers 实例。
func NewControllers(
	uc *UserController,
	cc *CaptchaController,
	mc *marketcontroller.MarketplaceController,
	rc *RegisterController,
) *Controllers {
	return &Controllers{
		User:        uc,
		Marketplace: mc,
		Captcha:     cc,
		Register:    rc,
	}
}
