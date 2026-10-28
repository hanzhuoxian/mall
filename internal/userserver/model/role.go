package model

// SysRole represents a system role.
type SysRole struct {
	RoleCode    string `json:"roleCode"           gorm:"column:role_code;type:varchar(64);uniqueIndex:uk_role_code;comment:角色编码"`
	Description string `json:"description"        gorm:"column:description;type:varchar(255);comment:角色描述"`
	ObjectMeta  `       json:"metadata,omitempty"`
	RoleStatus  int `json:"roleStatus"         gorm:"column:role_status;comment:状态 1启用 0禁用"`
}

// TableName returns the database table name for SysRole.
func (r *SysRole) TableName() string {
	return "sys_role"
}

// SysRoleList contains a page of roles.
type SysRoleList struct {
	Items    []*SysRole `json:"items"`
	ListMeta `           json:",inline"`
}

// ListRoleOptions controls role list queries.
type ListRoleOptions struct {
	RoleStatus *int `json:"roleStatus" form:"roleStatus"`
	ListOptions
	Keyword string `json:"keyword"    form:"keyword"`
}
