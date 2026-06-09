package middleware

import "github.com/gin-gonic/gin"

const UserIdentifier = "identifier"

// AuthStrategy 定义策略接口
type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

// AuthOperator 定义策略结构体
type AuthOperator struct {
	strategy AuthStrategy
}

// SetStrategy 设置策略
func (a *AuthOperator) SetStrategy(authStrategy AuthStrategy) {
	a.strategy = authStrategy
}

// AuthFunc 实现 AuthStrategy 接口
func (a *AuthOperator) AuthFunc() gin.HandlerFunc {
	return a.strategy.AuthFunc()
}
