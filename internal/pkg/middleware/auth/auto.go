package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
)

const authHeaderCount = 2

type AutoStrategy struct {
	basic middleware.AuthStrategy
	jwt   middleware.AuthStrategy
}

func NewAutoStrategy(basic, jwt middleware.AuthStrategy) middleware.AuthStrategy {
	return &AutoStrategy{
		basic: basic,
		jwt:   jwt,
	}
}

func (auto *AutoStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		operator := middleware.AuthOperator{}
		authHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", authHeaderCount)
		if len(authHeader) != authHeaderCount {
			c.Abort()
			return
		}

		switch authHeader[0] {
		case AuthBasicName:
			operator.SetStrategy(auto.basic)
		case AuthJWTName:
			operator.SetStrategy(auto.jwt)
		default:
			c.Abort()
		}

		operator.AuthFunc()(c)
		c.Next()
	}
}
