package model

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hanzhuoxian/mall/pkg/auth"
	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
)

// SysUser represents a user account.
type SysUser struct {
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty" gorm:"column:last_login_at;comment:最后登录时间"`
	Username    string     `json:"username"              gorm:"column:username;type:varchar(64);uniqueIndex;comment:用户名" validate:"omitempty"`
	Email       string     `json:"email"                 gorm:"column:email;type:varchar(100);uniqueIndex;comment:邮箱"    validate:"required,email,min=1,max=100"`
	Phone       string     `json:"phone"                 gorm:"column:phone;type:varchar(20);comment:手机号"                validate:"omitempty"`
	Nickname    string     `json:"nickname"              gorm:"column:nickname;type:varchar(30);comment:昵称"              validate:"required,min=1,max=30"`
	Password    string     `json:"password,omitempty"    gorm:"column:password;type:varchar(255);comment:密码"             validate:"required"`
	ObjectMeta  `json:"metadata,omitempty"`
	UserStatus  int `json:"userStatus"            gorm:"column:user_status;comment:状态 1启用 0禁用"                    validate:"omitempty"`
}

// TableName returns the database table name for SysUser.
func (u *SysUser) TableName() string {
	return "sys_user"
}

// SysUserList contains a page of users.
type SysUserList struct {
	Items    []*SysUser `json:"items"`
	ListMeta `           json:",inline"`
}

// ListUserOptions controls user list queries.
type ListUserOptions struct {
	Status *int `json:"status"  form:"status"`
	ListOptions
	Keyword string `json:"keyword" form:"keyword"`
}

// Compare verifies password against the stored password hash.
func (u *SysUser) Compare(password string) error {
	if err := auth.ComparePassword(u.Password, password); err != nil {
		return fmt.Errorf("failed to compile password: %w", err)
	}
	return nil
}

// ToProto converts the user model to its protobuf representation.
func (u *SysUser) ToProto() *userv1.User {
	return &userv1.User{
		InstanceId:  u.InstanceID,
		Name:        u.Name,
		Email:       u.Email,
		Phone:       u.Phone,
		Username:    u.Username,
		Nickname:    u.Nickname,
		Status:      int32(u.UserStatus),
		LastLoginAt: timeToProto(u.LastLoginAt),
		CreatedAt:   timestamppb.New(u.CreatedAt),
		UpdatedAt:   timestamppb.New(u.UpdatedAt),
	}
}

func timeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
