package auth

import (
	ginjwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
)

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
