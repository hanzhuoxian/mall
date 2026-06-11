package apiserver

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/appleboy/gin-jwt/v3/core"
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/apiserver/config"
	"github.com/hanzhuoxian/mall/internal/apiserver/grpcclient"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware/auth"
	"github.com/hanzhuoxian/mall/internal/pkg/response"
	"github.com/hanzhuoxian/mall/pkg/log"
	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

type loginInfo struct {
	Identifier string `form:"identifier" json:"identifier" binding:"required,identifier"`
	Password   string `form:"password" json:"password" binding:"required,password"`
}

// NewAutoAuth 创建 Auto 认证策略，根据 Authorization 头自动选择 Basic 或 JWT。
// cfg 和 userClient 均由 Wire 注入。
func NewAutoAuth(cfg *config.Config, userClient *grpcclient.UserClient) middleware.AuthStrategy {
	return auth.NewAutoStrategy(NewBasicAuth(userClient), NewJWTAuth(cfg, userClient))
}

func NewBasicAuth(userClient *grpcclient.UserClient) auth.BasicStrategy {
	return auth.NewBasicStrategy(func(identifier, password string) (string, bool) {
		resp, err := userClient.AuthenticateUser(context.Background(), &userv1.AuthenticateUserRequest{
			Identifier: identifier,
			Password:   password,
		})
		if err != nil {
			return "", false
		}
		return resp.GetInstanceId(), true
	})
}

func NewJWTAuth(cfg *config.Config, userClient *grpcclient.UserClient) auth.JWTStrategy {
	ginjwt, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "mall",
		SigningAlgorithm: "HS256",
		Key:              []byte(cfg.ServerRunOptions.JWTSecret),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour * 24,
		Authenticator:    authenticator(userClient),
		LoginResponse: func(c *gin.Context, token *core.Token) {
			response.Write(c, nil, nil)
		},
		LogoutResponse: func(c *gin.Context) {
			c.JSON(http.StatusOK, nil)
		},
		RefreshResponse: func(c *gin.Context, token *core.Token) {
			c.JSON(http.StatusOK, token)
		},
		IdentityHandler: func(c *gin.Context) any {
			claims := jwt.ExtractClaims(c)
			return claims[jwt.IdentityKey]
		},
		IdentityKey: middleware.UserIdentifier,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": message})
		},
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	return auth.NewJWTStrategy(*ginjwt)
}

func authenticator(userClient *grpcclient.UserClient) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var login loginInfo
		var err error

		if c.Request.Header.Get("Authorization") != "" {
			login, err = parseWithHeader(c)
		} else {
			login, err = parseWithBody(c)
		}

		if err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		resp, err := userClient.AuthenticateUser(context.Background(), &userv1.AuthenticateUserRequest{
			Identifier: login.Identifier,
			Password:   login.Password,
		})
		return resp.InstanceId, nil
	}
}

func parseWithHeader(c *gin.Context) (loginInfo, error) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		log.Errorf("get basic string from Authorization header failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	payload, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		log.Errorf("decode basic string: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		log.Errorf("parse payload failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	return loginInfo{
		Identifier: pair[0],
		Password:   pair[1],
	}, nil
}

func parseWithBody(c *gin.Context) (loginInfo, error) {
	var login loginInfo
	if err := c.ShouldBindJSON(&login); err != nil {
		log.Errorf("parse login parameters: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	return login, nil
}
