package apiserver

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/appleboy/gin-jwt/v3/core"
	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hanzhuoxian/mall/internal/apiserver/config"
	"github.com/hanzhuoxian/mall/internal/apiserver/grpcclient"
	"github.com/hanzhuoxian/mall/internal/pkg/coder"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware/auth"
	"github.com/hanzhuoxian/mall/internal/pkg/response"
	"github.com/hanzhuoxian/mall/pkg/errors"
	"github.com/hanzhuoxian/mall/pkg/logger"
	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

const (
	loginFailKeyPrefix = "login:fail:"
	loginFailThreshold = 2 // 连续失败次数达到此值后要求验证码
	loginFailExpiry    = 15 * time.Minute
)

type loginInfo struct {
	Identifier  string `form:"identifier"   json:"identifier"   binding:"required"`
	Password    string `form:"password"     json:"password"     binding:"required"`
	CaptchaID   string `form:"captcha_id"   json:"captcha_id"`
	CaptchaCode string `form:"captcha_code" json:"captcha_code"`
}

// NewAutoAuth 创建 Auto 认证策略，根据 Authorization 头自动选择 Basic 或 JWT。
// cfg、userClient、captcha 和 rdb 均由 Wire 注入。
func NewAutoAuth(cfg *config.Config, userClient *grpcclient.UserClient, captcha *base64Captcha.Captcha, rdb redis.UniversalClient) middleware.AuthStrategy {
	return auth.NewAutoStrategy(NewBasicAuth(userClient), NewJWTAuth(cfg, userClient, captcha, rdb))
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

func NewJWTAuth(cfg *config.Config, userClient *grpcclient.UserClient, captcha *base64Captcha.Captcha, rdb redis.UniversalClient) auth.JWTStrategy {
	ginjwt, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "mall",
		SigningAlgorithm: "HS256",
		Key:              []byte(cfg.ServerRunOptions.JWTSecret),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour * 24,
		Authenticator:    authenticator(userClient, captcha, rdb),
		PayloadFunc:      payload(),
		LoginResponse: func(c *gin.Context, token *core.Token) {
			response.Success(c, token)
		},
		LogoutResponse: func(c *gin.Context) {
			response.Success(c, nil)
		},
		RefreshResponse: func(c *gin.Context, token *core.Token) {
			response.Success(c, token)
		},
		IdentityKey: middleware.UserIdentifier,
		Unauthorized: func(c *gin.Context, code int, message string) {
			response.Fail(c, errors.New(coder.ErrTokenInvalid, "Token invalid"))
		},
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	return auth.NewJWTStrategy(*ginjwt)
}

func authenticator(userClient *grpcclient.UserClient, captcha *base64Captcha.Captcha, rdb redis.UniversalClient) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var login loginInfo
		var err error

		isBasic := strings.HasPrefix(c.Request.Header.Get("Authorization"), "Basic")
		if isBasic {
			login, err = parseWithHeader(c)
		} else {
			login, err = parseWithBody(c)
		}
		if err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		if !isBasic {
			ctx := context.Background()
			failKey := loginFailKeyPrefix + login.Identifier
			failCount, _ := rdb.Get(ctx, failKey).Int()
			if failCount >= loginFailThreshold {
				if !captcha.Verify(login.CaptchaID, login.CaptchaCode, true) {
					return "", errors.New(coder.ErrCaptchaInvalid, "captcha invalid or expired")
				}
			}

			resp, err := userClient.AuthenticateUser(ctx, &userv1.AuthenticateUserRequest{
				Identifier: login.Identifier,
				Password:   login.Password,
			})
			if err != nil {
				// 失败计数：第一次失败时设置过期时间
				count, _ := rdb.Incr(ctx, failKey).Result()
				if count == 1 {
					rdb.Expire(ctx, failKey, loginFailExpiry)
				}
				logAuthError(login.Identifier, err)
				return "", errors.New(coder.ErrPasswordIncorrect, "password incorrect")
			}
			// 登录成功，清除失败计数
			rdb.Del(ctx, failKey)
			return resp.InstanceId, nil
		}

		resp, err := userClient.AuthenticateUser(context.Background(), &userv1.AuthenticateUserRequest{
			Identifier: login.Identifier,
			Password:   login.Password,
		})
		if err != nil {
			logAuthError(login.Identifier, err)
			return "", errors.New(coder.ErrPasswordIncorrect, "password incorrect")
		}
		return resp.InstanceId, nil
	}
}

func payload() func(data any) jwtv5.MapClaims {
	return func(data any) jwtv5.MapClaims {
		if instanceID, ok := data.(string); ok {
			return jwtv5.MapClaims{middleware.UserIdentifier: instanceID}
		}
		return jwtv5.MapClaims{}
	}
}

func parseWithHeader(c *gin.Context) (loginInfo, error) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		logger.Errorf("get basic string from Authorization header failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	payload, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		logger.Errorf("decode basic string: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		logger.Errorf("parse payload failed")

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
		logger.Errorf("parse login parameters: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	return login, nil
}

// logAuthError logs auth failures at WARN for expected cases (user not found,
// wrong password) and ERROR for unexpected internal failures.
func logAuthError(identifier string, err error) {
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.NotFound, codes.Unauthenticated:
		logger.Warnf("authenticate user %q failed: %s", identifier, st.Message())
	default:
		logger.Errorf("authenticate user %q failed: %v", identifier, err)
	}
}
