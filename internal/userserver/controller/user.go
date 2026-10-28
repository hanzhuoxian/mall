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

// CreateUser creates a user through the user service.
func (s *UserServiceServer) CreateUser(
	ctx context.Context,
	req *userv1.CreateUserRequest,
) (*userv1.CreateUserResponse, error) {
	return s.svc.Users().Create(ctx, req)
}

// UpdateUser updates a user through the user service.
func (s *UserServiceServer) UpdateUser(
	ctx context.Context,
	req *userv1.UpdateUserRequest,
) (*userv1.UpdateUserResponse, error) {
	return s.svc.Users().Update(ctx, req)
}

// DeleteUser deletes a user through the user service.
func (s *UserServiceServer) DeleteUser(
	ctx context.Context,
	req *userv1.DeleteUserRequest,
) (*userv1.DeleteUserResponse, error) {
	return s.svc.Users().Delete(ctx, req)
}

// DeleteCollection deletes a collection of users through the user service.
func (s *UserServiceServer) DeleteCollection(
	ctx context.Context,
	req *userv1.DeleteCollectionRequest,
) (*userv1.DeleteCollectionResponse, error) {
	return s.svc.Users().DeleteCollection(ctx, req)
}

// GetUser retrieves a user through the user service.
func (s *UserServiceServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	return s.svc.Users().Get(ctx, req)
}

// ListUsers retrieves users through the user service.
func (s *UserServiceServer) ListUsers(
	ctx context.Context,
	req *userv1.ListUsersRequest,
) (*userv1.ListUsersResponse, error) {
	return s.svc.Users().List(ctx, req)
}

// AuthenticateUser validates user credentials through the user service.
func (s *UserServiceServer) AuthenticateUser(
	ctx context.Context,
	req *userv1.AuthenticateUserRequest,
) (*userv1.AuthenticateUserResponse, error) {
	return s.svc.Users().AuthenticateUser(ctx, req)
}
