package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/apiserver/grpcclient"
	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

// UserController 处理用户相关的 HTTP 请求，通过 gRPC 调用 userserver。
type UserController struct {
	userClient *grpcclient.UserClient
}

// NewUserController 创建 UserController 实例。
func NewUserController(userClient *grpcclient.UserClient) *UserController {
	return &UserController{userClient: userClient}
}

// GetUser 根据路径参数 id 获取单个用户信息。
func (h *UserController) GetUser(c *gin.Context) {
	resp, err := h.userClient.GetUser(c.Request.Context(), &userv1.GetUserRequest{
		UserId: c.Param("id"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ListUsers 获取用户列表。
func (h *UserController) ListUsers(c *gin.Context) {
	resp, err := h.userClient.ListUsers(c.Request.Context(), &userv1.ListUsersRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
