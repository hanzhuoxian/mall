package auth

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
)

// AuthBasicName 表示 HTTP Basic 认证在 Authorization 头中使用的方案名称。
const AuthBasicName = "Basic"

// BasicStrategy 实现了一个基于 HTTP Basic 的简单认证策略。
// compare 函数用于校验给定的 identifier 与 password 是否匹配。
type BasicStrategy struct {
	authenticate func(identifier, password string) (instanceID string, ok bool)
}

var _ middleware.AuthStrategy = &BasicStrategy{}

// NewBasicStrategy 使用给定的比较函数构造并返回一个 BasicStrategy。
func NewBasicStrategy(authenticate func(identifier, password string) (instanceID string, ok bool)) BasicStrategy {
	return BasicStrategy{authenticate: authenticate}
}

// AuthFunc 返回一个 gin.HandlerFunc，用于解析 Authorization 头并调用 authenticate
// 进行认证。认证失败时会调用 c.Abort() 终止请求链路。
func (b BasicStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != AuthBasicName {
			c.Abort()
		}
		payload, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			c.Abort()
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			c.Abort()
		}
		if instanceID, ok := b.authenticate(pair[0], pair[1]); ok {
			c.Set(middleware.UserIdentifier, pair[0])
			c.Set(middleware.UserInstanceID, instanceID)
		}
		c.Next()
	}
}
