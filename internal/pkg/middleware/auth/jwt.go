package auth

import (
	ginjwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
)

// AuthBasicName 表示 HTTP Basic 认证在 Authorization 头中使用的方案名称。
const AuthJWTName = "Bearer"

// JWTStrategy
type JWTStrategy struct {
	ginjwt.GinJWTMiddleware
}

var _ middleware.AuthStrategy = &JWTStrategy{}

func NewJWTStrategy(g ginjwt.GinJWTMiddleware) JWTStrategy {
	return JWTStrategy{g}
}

func (j JWTStrategy) AuthFunc() gin.HandlerFunc {
	return j.MiddlewareFunc()
}
