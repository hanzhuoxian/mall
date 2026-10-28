package model

// SysUserRole associates a user with a role.
type SysUserRole struct {
	ObjectMeta `       json:"metadata,omitempty"`
	UserID     uint64 `json:"userId"             gorm:"uniqueIndex:uk_user_role;comment:用户ID"`
	RoleID     uint64 `json:"roleId"             gorm:"uniqueIndex:uk_user_role;comment:角色ID"`
}

// TableName returns the database table name for SysUserRole.
func (p *SysUserRole) TableName() string {
	return "sys_user_role"
}
