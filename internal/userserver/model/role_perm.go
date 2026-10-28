package model

// SysRolePermission associates a role with a permission.
type SysRolePermission struct {
	ObjectMeta   `       json:"metadata,omitempty"`
	PermissionID uint64 `json:"permissionId"       gorm:"uniqueIndex:uk_role_perm;comment:权限ID"`
	RoleID       uint64 `json:"roleId"             gorm:"uniqueIndex:uk_role_perm;comment:角色ID"`
}

// TableName returns the database table name for SysRolePermission.
func (p *SysRolePermission) TableName() string {
	return "sys_role_permission"
}
