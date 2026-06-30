package model

type SysRole struct {
	ObjectMeta  `json:"metadata,omitempty"`
	RoleCode    string `json:"role_code" gorm:"column:role_code;type:varchar(64);uniqueIndex:uk_role_code;comment:角色编码"`
	Description string `json:"description" gorm:"column:description;type:varchar(255);comment:角色描述"`
	RoleStatus  int    `json:"role_status" gorm:"column:role_status;comment:状态 1启用 0禁用"`
}

func (r *SysRole) TableName() string {
	return "sys_role"
}

type SysRoleList struct {
	ListMeta `json:",inline"`
	Items    []*SysRole `json:"items"`
}

type ListRoleOptions struct {
	ListOptions
	Keyword    string `json:"keyword" form:"keyword"`
	RoleStatus *int   `json:"roleStatus" form:"roleStatus"`
}
