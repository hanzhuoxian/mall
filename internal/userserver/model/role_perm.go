package model

type SysRolePermission struct {
	ObjectMeta   `json:"metadata,omitempty"`
	PermissionID uint64 `json:"permission_id" gorm:"uniqueIndex:uk_role_perm;comment:权限ID"`
	RoleID       uint64 `json:"role_id" gorm:"uniqueIndex:uk_role_perm;comment:角色ID"`
}

func (p *SysRolePermission) TableName() string {
	return "sys_role_permission"
}
