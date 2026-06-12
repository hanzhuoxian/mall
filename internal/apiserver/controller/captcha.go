package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"

	"github.com/hanzhuoxian/mall/internal/pkg/response"
)

// CaptchaController 处理图形验证码的生成请求。
type CaptchaController struct {
	captcha *base64Captcha.Captcha
}

// NewCaptchaController 创建 CaptchaController 实例。
func NewCaptchaController(captcha *base64Captcha.Captcha) *CaptchaController {
	return &CaptchaController{captcha: captcha}
}

type generateCaptchaResponse struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

// Generate 生成图形验证码并返回 id 和 base64 图片。
func (h *CaptchaController) Generate(c *gin.Context) {
	id, b64s, _, err := h.captcha.Generate()
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, generateCaptchaResponse{ID: id, Image: b64s})
}
