package controller

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"

	"github.com/hanzhuoxian/mall/internal/apiserver/grpcclient"
	"github.com/hanzhuoxian/mall/internal/pkg/response"
	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

var usernamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{2,31}$`)

// RegisterController 处理公开注册，通过用户服务创建普通账号。
type RegisterController struct {
	userClient *grpcclient.UserClient
	captcha    *base64Captcha.Captcha
}

// NewRegisterController 创建注册控制器。
func NewRegisterController(users *grpcclient.UserClient, captcha *base64Captcha.Captcha) *RegisterController {
	return &RegisterController{userClient: users, captcha: captcha}
}

type registerRequest struct {
	Username    string `json:"username"    binding:"required"`
	Email       string `json:"email"       binding:"required,email,max=100"`
	Password    string `json:"password"    binding:"required,min=12,max=72"`
	CaptchaID   string `json:"captchaId"   binding:"required"`
	CaptchaCode string `json:"captchaCode" binding:"required"`
}

// Register 消费验证码并创建普通启用账号，不接受请求指定商家或平台权限。
func (h *RegisterController) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, &response.HTTPError{Status: http.StatusBadRequest, Message: "请检查用户名、邮箱及密码（12–72 字节）"})
		return
	}
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if !usernamePattern.MatchString(req.Username) || req.Username == "admin" || len(req.Password) > 72 {
		response.Fail(
			c,
			&response.HTTPError{Status: http.StatusBadRequest, Message: "用户名须为 3–32 位字母、数字或下划线，以字母开头，不可使用保留名称"},
		)
		return
	}
	if !h.captcha.Verify(req.CaptchaID, req.CaptchaCode, true) {
		response.Fail(c, &response.HTTPError{Status: http.StatusBadRequest, Message: "验证码错误或已过期"})
		return
	}
	_, err := h.userClient.CreateUser(c.Request.Context(), &userv1.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Username,
		Name:     req.Username,
	})
	if err != nil {
		response.Fail(c, &response.HTTPError{Status: http.StatusConflict, Message: "注册失败，用户名或邮箱可能已被使用", Cause: err})
		return
	}
	response.Success(c, gin.H{"registered": true})
}
