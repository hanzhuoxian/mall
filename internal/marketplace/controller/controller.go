// Package controller 将商城 HTTP 请求转换为业务服务调用。
package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/service"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
	"github.com/hanzhuoxian/mall/internal/pkg/response"
)

// Identity 根据认证标识查询当前用户及平台权限。
type Identity func(context.Context, string) (model.Actor, error)

// MarketplaceController 处理商城 HTTP 请求，数据库访问由业务服务承担。
type MarketplaceController struct {
	service     *service.Service
	identity    Identity
	identityKey string
}

// NewMarketplaceController 创建商城控制器。
func NewMarketplaceController(svc *service.Service, identity Identity, identityKey string) *MarketplaceController {
	return &MarketplaceController{service: svc, identity: identity, identityKey: identityKey}
}

const actorKey = "marketplace.actor"

func actor(c *gin.Context) model.Actor { return c.MustGet(actorKey).(model.Actor) }

func fail(c *gin.Context, err error) {
	status, message := http.StatusInternalServerError, "操作失败，请稍后重试"
	switch {
	case errors.Is(err, service.ErrInvalid), errors.Is(err, service.ErrStock):
		status, message = http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrForbidden):
		status, message = http.StatusForbidden, err.Error()
	case errors.Is(err, store.ErrNotFound):
		status, message = http.StatusNotFound, err.Error()
	case errors.Is(err, service.ErrPaymentDisabled):
		status, message = http.StatusConflict, err.Error()
	}
	response.Fail(c, &response.HTTPError{Status: status, Message: message, Cause: err})
	c.Abort()
}

func respond(c *gin.Context, data any, err error) {
	if err != nil {
		fail(c, err)
		return
	}
	// 无数据的写操作保留商城原有的 data: null 响应。
	if data == nil {
		data = json.RawMessage("null")
	}
	response.Success(c, data)
}

func bind(c *gin.Context, input any) bool {
	if err := c.ShouldBindJSON(input); err != nil {
		fail(c, service.ErrInvalid)
		return false
	}
	return true
}

func pagination(c *gin.Context) model.ListOptions {
	page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	return model.ListOptions{Page: page, PageSize: size}
}
