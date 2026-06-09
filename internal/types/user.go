package types

import "time"

type User struct {
	ObjectMeta `json:"metadata,omitempty"`
	Email      string    `json:"email" gorm:"column:email" validate:"required,email,min=1,max=100"`
	Phone      string    `json:"phone" gorm:"column:phone" validate:"omitempty"`
	Username   string    `json:"username" gorm:"column:phone" validate:"omitempty"`
	Nickname   string    `json:"nickname" gorm:"column:nickname" validate:"required,min=1,max=30"`
	Password   string    `json:"password,omitempty" gorm:"column:password" validate:"required"`
	Status     int       `json:"status" gorm:"column:status" validate:"omitempty"`
	LoginedAt  time.Time `json:"loginedAt,omitempty" gorm:"column:logined_at"`
}

type UserList struct {
	ListMeta `json:",inline"`

	Items []*User `json:"items"`
}

func (u *User) TableName() string {
	return "user"
}
