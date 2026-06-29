package model

type UserRole struct {
	ObjectMeta `json:"metadata,omitempty"`
	UserID     uint64 `json:"user_id" gorm:"uniqueIndex:uk_user_role;comment:用户ID"`
	RoleID     uint64 `json:"role_id" gorm:"uniqueIndex:uk_user_role;comment:角色ID"`
}

func (p *UserRole) TableName() string {
	return "user_role"
}
