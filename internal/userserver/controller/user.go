// Package controller 实现了用户服务的 HTTP 请求处理器。
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/userserver/cache"
)

// UserController 处理用户相关的 HTTP 请求，依赖缓存工厂实现业务逻辑。
type UserController struct {
	cache cache.Factory
}

// NewUserController 创建并返回一个新的 UserController 实例。
func NewUserController(c cache.Factory) *UserController {
	return &UserController{cache: c}
}

// Hello 是测试接口，向缓存写入并读取用户名，用于验证缓存连通性。
func (uc *UserController) Hello(ctx *gin.Context) {
	u := uc.cache.User()
	if err := u.SetUser(ctx, "username", "hanjian"); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	username, err := u.GetUser(ctx, "username")
	if err != nil {
		ctx.JSON(http.StatusNotFound, "")
		return
	}
	ctx.JSON(http.StatusOK, username)
}
