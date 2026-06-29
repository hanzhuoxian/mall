package model

type SysUserRole struct {
	ObjectMeta `json:"metadata,omitempty"`
	UserID     uint64 `json:"user_id" gorm:"uniqueIndex:uk_user_role;comment:用户ID"`
	RoleID     uint64 `json:"role_id" gorm:"uniqueIndex:uk_user_role;comment:角色ID"`
}

func (p *SysUserRole) TableName() string {
	return "sys_user_role"
}
