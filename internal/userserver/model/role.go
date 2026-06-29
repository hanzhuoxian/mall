package model

type SysRole struct {
	ObjectMeta  `json:"metadata,omitempty"`
	RoleCode    string `json:"role_code" gorm:"uniqueIndex:uk_role_code;comment:角色编码"`
	Description string `json:"description" gorm:"comment:角色描述"`
	RoleStatus  int    `json:"role_status" gorm:"comment:状态 1启用 0禁用"`
}

func (r *SysRole) TableName() string {
	return "sys_role"
}

type SysRoleList struct {
	ListMeta `json:",inline"`
	Item     []SysRole `json:"item"`
}
