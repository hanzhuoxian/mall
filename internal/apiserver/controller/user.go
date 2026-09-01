package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/apiserver/grpcclient"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
	"github.com/hanzhuoxian/mall/internal/pkg/response"
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

type createUserRequest struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Phone    string `json:"phone"`
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateUser 创建新用户。
func (h *UserController) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.userClient.CreateUser(c.Request.Context(), &userv1.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Username: req.Username,
		Nickname: req.Nickname,
		Password: req.Password,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, resp)
}

type updateUserRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Status   int32  `json:"status"`
}

// UpdateUser 更新指定用户的信息。
func (h *UserController) UpdateUser(c *gin.Context) {
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.userClient.UpdateUser(c.Request.Context(), &userv1.UpdateUserRequest{
		InstanceId: c.Param("id"),
		Email:      req.Email,
		Phone:      req.Phone,
		Nickname:   req.Nickname,
		Password:   req.Password,
		Status:     req.Status,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, resp)
}

type deleteUserQuery struct {
	Unscoped bool `form:"unscoped"`
}

// DeleteUser 删除指定用户，unscoped=true 时硬删除。
func (h *UserController) DeleteUser(c *gin.Context) {
	var q deleteUserQuery
	_ = c.ShouldBindQuery(&q)
	resp, err := h.userClient.DeleteUser(c.Request.Context(), &userv1.DeleteUserRequest{
		InstanceId: c.Param("id"),
		Unscoped:   q.Unscoped,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, resp)
}

type deleteCollectionRequest struct {
	InstanceIds []string `json:"instanceIds" binding:"required,min=1"`
	Unscoped    bool     `json:"unscoped"`
}

// DeleteCollection 批量删除用户，unscoped=true 时硬删除。
func (h *UserController) DeleteCollection(c *gin.Context) {
	var req deleteCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.userClient.DeleteCollection(c.Request.Context(), &userv1.DeleteCollectionRequest{
		InstanceIds: req.InstanceIds,
		Unscoped:    req.Unscoped,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, resp)
}

// GetMe 获取当前登录用户的信息。
func (h *UserController) GetMe(c *gin.Context) {
	instanceID := c.GetString(middleware.UserIdentifier)
	resp, err := h.userClient.GetUser(c.Request.Context(), &userv1.GetUserRequest{
		InstanceId: instanceID,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, resp)
}

// GetUser 根据路径参数 id 获取单个用户信息。
func (h *UserController) GetUser(c *gin.Context) {
	resp, err := h.userClient.GetUser(c.Request.Context(), &userv1.GetUserRequest{
		InstanceId: c.Param("id"),
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, resp)
}

type listUsersQuery struct {
	Page     int32 `form:"page"`
	PageSize int32 `form:"pageSize"`
}

// ListUsers 获取分页用户列表。
func (h *UserController) ListUsers(c *gin.Context) {
	var q listUsersQuery
	_ = c.ShouldBindQuery(&q)
	resp, err := h.userClient.ListUsers(c.Request.Context(), &userv1.ListUsersRequest{
		Page:     q.Page,
		PageSize: q.PageSize,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, resp)
}
