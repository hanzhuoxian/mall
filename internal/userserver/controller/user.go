package controller

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/userserver/service"
	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

// UserServiceServer 实现了 userv1.UserServiceServer 接口。
type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer
	svc service.Service
}

// NewUserServiceServer 创建 UserServiceServer 实例。
func NewUserServiceServer(svc service.Service) *UserServiceServer {
	return &UserServiceServer{svc: svc}
}

// GetUser 根据 user_id 返回用户信息。
func (s *UserServiceServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	return s.svc.Users().Get(ctx, req)
}

// ListUsers 返回分页用户列表。
func (s *UserServiceServer) ListUsers(_ context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	return &userv1.ListUsersResponse{
		Users: []*userv1.User{},
		Total: 1,
	}, nil
}

func (s *UserServiceServer) AuthenticateUser(ctx context.Context, req *userv1.AuthenticateUserRequest) (*userv1.AuthenticateUserResponse, error) {
	return s.svc.Users().AuthenticateUser(ctx, req)
}
