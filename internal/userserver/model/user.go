package model

import (
	"fmt"
	"time"

	"github.com/hanzhuoxian/mall/pkg/auth"
	userv1 "github.com/hanzhuoxian/mall/proto/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SysUser struct {
	ObjectMeta `json:"metadata,omitempty"`
	Username   string     `json:"username" gorm:"column:username;type:varchar(64);comment:用户名" validate:"omitempty"`
	Email      string     `json:"email" gorm:"column:email;type:varchar(100);comment:邮箱" validate:"required,email,min=1,max=100"`
	Phone      string     `json:"phone" gorm:"column:phone;type:varchar(20);comment:手机号" validate:"omitempty"`
	Nickname   string     `json:"nickname" gorm:"column:nickname;type:varchar(30);comment:昵称" validate:"required,min=1,max=30"`
	Password   string     `json:"password,omitempty" gorm:"column:password;type:varchar(255);comment:密码" validate:"required"`
	Status     int        `json:"status" gorm:"column:status;comment:状态 1启用 0禁用" validate:"omitempty"`
	LoginedAt  *time.Time `json:"loginedAt,omitempty" gorm:"column:logined_at;comment:最后登录时间"`
}

func (u *SysUser) TableName() string {
	return "sys_user"
}

type SysUserList struct {
	ListMeta `json:",inline"`

	Items []*SysUser `json:"items"`
}

func (u *SysUser) Compare(password string) error {
	if err := auth.ComparePassword(u.Password, password); err != nil {
		return fmt.Errorf("failed to compile password: %w", err)
	}
	return nil
}

func (u *SysUser) ToProto() *userv1.User {
	return &userv1.User{
		InstanceId: u.InstanceID,
		Name:       u.Name,
		Email:      u.Email,
		Phone:      u.Phone,
		Username:   u.Username,
		Nickname:   u.Nickname,
		Status:     int32(u.Status),
		LoginedAt:  timeToProto(u.LoginedAt),
		CreatedAt:  timestamppb.New(u.CreatedAt),
		UpdatedAt:  timestamppb.New(u.UpdatedAt),
	}
}

func timeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
