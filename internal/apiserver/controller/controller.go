package controller

import "github.com/google/wire"

// ProviderSet is used by Wire.
var ProviderSet = wire.NewSet(
	NewUserController,
	NewCaptchaController,
	NewControllers,
)

// Controllers 汇聚所有控制器实例，便于在路由注册时统一传递。
type Controllers struct {
	User    *UserController
	Captcha *CaptchaController
}

// NewControllers 创建并返回包含所有控制器的 Controllers 实例。
func NewControllers(uc *UserController, cc *CaptchaController) *Controllers {
	return &Controllers{
		User:    uc,
		Captcha: cc,
	}
}
