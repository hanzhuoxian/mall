package auth

import (
	ginjwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
)

// AuthJWTName is the Authorization scheme name used for JWT authentication.
const AuthJWTName = "Bearer"

// JWTStrategy adapts the gin-jwt middleware to the application's auth strategy.
type JWTStrategy struct {
	ginjwt.GinJWTMiddleware
}

var _ middleware.AuthStrategy = &JWTStrategy{}

// NewJWTStrategy wraps the configured gin-jwt middleware.
func NewJWTStrategy(g ginjwt.GinJWTMiddleware) JWTStrategy {
	return JWTStrategy{g}
}

// AuthFunc returns the JWT authentication handler.
func (j JWTStrategy) AuthFunc() gin.HandlerFunc {
	return j.MiddlewareFunc()
}
