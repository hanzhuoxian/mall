package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/hanzhuoxian/mall/internal/userserver/model"
	"github.com/hanzhuoxian/mall/internal/userserver/store"
	"github.com/hanzhuoxian/mall/pkg/auth"
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
	Create(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error)
	Update(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error)
	Delete(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error)
	DeleteCollection(ctx context.Context, req *userv1.DeleteCollectionRequest) (*userv1.DeleteCollectionResponse, error)
	Get(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error)
	List(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error)
	AuthenticateUser(ctx context.Context, req *userv1.AuthenticateUserRequest) (*userv1.AuthenticateUserResponse, error)
}

type userSrv struct {
	store store.Factory
}

func newUsers(srv *service) UserSrv {
	return &userSrv{
		store: srv.store,
	}
}

func (s *userSrv) Create(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	hashedPassword, err := auth.GeneratePassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &model.SysUser{
		ObjectMeta: model.ObjectMeta{
			InstanceID: uuid.New().String(),
			Name:       req.Name,
		},
		Email:    req.Email,
		Phone:    req.Phone,
		Username: req.Username,
		Nickname: req.Nickname,
		Password: hashedPassword,
	}
	if err := s.store.Users().Create(ctx, user, model.CreateOptions{}); err != nil {
		return nil, err
	}
	return &userv1.CreateUserResponse{User: user.ToProto()}, nil
}

func (s *userSrv) Update(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	user, err := s.store.Users().GetByInstanceID(ctx, req.InstanceId, model.GetOptions{})
	if err != nil {
		return nil, err
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Password != "" {
		hashed, err := auth.GeneratePassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	}
	if req.Status != 0 {
		user.UserStatus = int(req.Status)
	}
	if err := s.store.Users().Update(ctx, user, model.UpdateOptions{}); err != nil {
		return nil, err
	}
	return &userv1.UpdateUserResponse{User: user.ToProto()}, nil
}

func (s *userSrv) Delete(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	err := s.store.Users().Delete(ctx, req.InstanceId, model.DeleteOptions{Unscoped: req.Unscoped})
	return &userv1.DeleteUserResponse{}, err
}

func (s *userSrv) DeleteCollection(ctx context.Context, req *userv1.DeleteCollectionRequest) (*userv1.DeleteCollectionResponse, error) {
	err := s.store.Users().DeleteCollection(ctx, req.InstanceIds, model.DeleteOptions{Unscoped: req.Unscoped})
	return &userv1.DeleteCollectionResponse{}, err
}

func (s *userSrv) Get(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := s.store.Users().GetByInstanceID(ctx, req.InstanceId, model.GetOptions{})
	return &userv1.GetUserResponse{User: user.ToProto()}, err
}

func (s *userSrv) List(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	pageSize := int64(req.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	page := int64(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	userList, err := s.store.Users().List(ctx, model.ListUserOptions{
		ListOptions: model.ListOptions{Offset: &offset, Limit: &pageSize},
		Keyword:     req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	users := make([]*userv1.User, 0, len(userList.Items))
	for _, u := range userList.Items {
		users = append(users, u.ToProto())
	}
	return &userv1.ListUsersResponse{Users: users, Total: userList.TotalCount}, nil
}

func (s *userSrv) AuthenticateUser(ctx context.Context, req *userv1.AuthenticateUserRequest) (*userv1.AuthenticateUserResponse, error) {
	var user *model.SysUser
	var err error
	switch DetectIdentifierType(req.Identifier) {
	case IdentifierEmail:
		user, err = s.store.Users().GetByEmail(ctx, req.Identifier, model.GetOptions{})
	case IdentifierPhone:
		user, err = s.store.Users().GetByPhone(ctx, req.Identifier, model.GetOptions{})
	default:
		user, err = s.store.Users().Get(ctx, req.Identifier, model.GetOptions{})
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, err
	}
	if err := auth.ComparePassword(user.Password, req.Password); err != nil {
		return nil, status.Error(codes.Unauthenticated, "password incorrect")
	}
	return &userv1.AuthenticateUserResponse{InstanceId: user.InstanceID}, nil
}

func DetectIdentifierType(id string) identifierType {
	if regex.Email.MatchString(id) {
		return IdentifierEmail
	}
	if regex.Phone.MatchString(id) {
		return IdentifierPhone
	}
	return IdentifierUsername
}
