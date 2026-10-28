package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/marketplace/service"
	"github.com/hanzhuoxian/mall/internal/pkg/response"
)

func (h *MarketplaceController) authenticated(c *gin.Context) {
	id := c.GetString(h.identityKey)
	if id == "" {
		response.Fail(c, &response.HTTPError{Status: http.StatusUnauthorized, Message: "请先登录"})
		c.Abort()
		return
	}
	current, err := h.identity(c.Request.Context(), id)
	if err != nil {
		fail(c, service.ErrForbidden)
		return
	}
	c.Set(actorKey, current)
	c.Next()
}

// AdminOnly 为用户管理等外部路由校验平台管理员身份。
func (h *MarketplaceController) AdminOnly(c *gin.Context) {
	current, err := h.identity(c.Request.Context(), c.GetString(h.identityKey))
	if err != nil || !current.Admin {
		fail(c, service.ErrForbidden)
		return
	}
	c.Next()
}

func requireAdmin(c *gin.Context) {
	if !actor(c).Admin {
		fail(c, service.ErrForbidden)
		return
	}
	c.Next()
}
