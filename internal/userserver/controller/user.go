package controller

import (
	"context"

	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

// UserServiceServer 实现了 userv1.UserServiceServer 接口。
type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer
}

// NewUserServiceServer 创建 UserServiceServer 实例。
func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{}
}

// GetUser 根据 user_id 返回用户信息。
func (s *UserServiceServer) GetUser(_ context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	return &userv1.GetUserResponse{
		UserId:   req.UserId,
		Username: "demo",
		Email:    "demo@example.com",
	}, nil
}

// ListUsers 返回分页用户列表。
func (s *UserServiceServer) ListUsers(_ context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	return &userv1.ListUsersResponse{
		Users: []*userv1.GetUserResponse{
			{UserId: "1", Username: "demo", Email: "demo@example.com"},
		},
		Total: 1,
	}, nil
}
