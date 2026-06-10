package middleware

import "github.com/gin-gonic/gin"

// UserIdentifier 是上下文中用于标识当前用户的键名。
const UserIdentifier = "identifier"
const UserInstanceID = "instanceID"

// AuthStrategy 定义认证策略接口，任何认证实现只需返回一个 gin.HandlerFunc。
// 返回的 HandlerFunc 应处理认证逻辑并在认证通过时将必要的信息写入 Context。
type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

// AuthOperator 是认证策略的包装器，用于在运行时注入不同的认证实现。
type AuthOperator struct {
	strategy AuthStrategy
}

// SetStrategy 设置当前使用的认证策略实现。
func (a *AuthOperator) SetStrategy(authStrategy AuthStrategy) {
	a.strategy = authStrategy
}

// AuthFunc 返回当前策略对应的 gin.HandlerFunc，满足 AuthStrategy 接口。
func (a *AuthOperator) AuthFunc() gin.HandlerFunc {
	return a.strategy.AuthFunc()
}
