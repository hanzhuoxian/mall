package service

import (
	"context"

	"github.com/hanzhuoxian/mall/internal/types"
	"github.com/hanzhuoxian/mall/internal/userserver/store"
	"github.com/hanzhuoxian/mall/pkg/regex"
	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

type identifierType int

const (
	IdentifierUsername identifierType = iota
	IdentifierEmail
	IdentifierPhone
)

type UserSrv interface {
	AuthenticateUser(ctx context.Context, req *userv1.AuthenticateUserRequest) (*userv1.AuthenticateUserResponse, error)
	Get(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error)
}

type userSrv struct {
	store store.Factory
}

func newUsers(srv *service) UserSrv {
	return &userSrv{
		store: srv.store,
	}
}

func (s *userSrv) Get(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := s.store.Users().GetByInstanceID(ctx, req.InstanceId, types.GetOptions{})
	return &userv1.GetUserResponse{
		User: user.ToProto(),
	}, err
}

func (s *userSrv) AuthenticateUser(ctx context.Context, req *userv1.AuthenticateUserRequest) (*userv1.AuthenticateUserResponse, error) {
	var user *types.User
	var err error
	switch DetectIdentifierType(req.Identifier) {
	case IdentifierEmail:
		user, err = s.store.Users().GetByEmail(ctx, req.Identifier, types.GetOptions{})
	case IdentifierPhone:
		user, err = s.store.Users().GetByPhone(ctx, req.Identifier, types.GetOptions{})
	default:
		user, err = s.store.Users().GetByPhone(ctx, req.Identifier, types.GetOptions{})
	}

	if err != nil {
		return nil, err
	}
	return &userv1.AuthenticateUserResponse{InstanceId: user.InstanceID}, nil
}

// DetectIdentifierType
func DetectIdentifierType(id string) identifierType {
	if regex.Email.MatchString(id) {
		return IdentifierEmail
	}
	if regex.Phone.MatchString(id) {
		return IdentifierPhone
	}
	return IdentifierUsername
}
