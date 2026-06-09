package auth

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
)

const AuthBasicName = "Basic"

type BasicStrategy struct {
	compare func(identifier, password string) bool
}

var _ middleware.AuthStrategy = &BasicStrategy{}

func NewBasicStrategy(compare func(identifier, password string) bool) BasicStrategy {
	return BasicStrategy{compare: compare}
}

func (b *BasicStrategy) AuthFunc() gin.HandlerFunc {
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
		if len(pair) != 2 || b.compare(pair[0], pair[1]) {
			c.Abort()
		}
		c.Next()
	}
}
