package types

type SysRole struct {
	ObjectMeta  `json:"metadata,omitempty"`
	RoleCode    string `json:"role_code"`
	Description string `json:"description"`
	RoleStatus  int    `json:"role_status"`
}

func (r *SysRole) TableName() string {
	return "sys_role"
}

type SysRoleList struct {
	ListMeta `json:",inline"`
	Item     []SysRole `json:"item"`
}
