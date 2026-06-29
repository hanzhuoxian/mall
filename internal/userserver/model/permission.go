package model

type SysPermission struct {
	ObjectMeta `json:"metadata,omitempty"`
	ParentID   uint64 `json:"parent_id" gorm:"column:parent_id;comment:父级权限ID 0为顶级"`
	PermCode   string `json:"perm_code" gorm:"column:perm_code;type:varchar(64);comment:权限编码"`
	PermType   string `json:"perm_type" gorm:"column:perm_type;type:varchar(20);comment:权限类型 menu/button/api"`
	Path       string `json:"path" gorm:"column:path;type:varchar(255);comment:路由路径"`
	Icon       string `json:"icon" gorm:"column:icon;type:varchar(64);comment:图标"`
	Sort       int    `json:"sort" gorm:"column:sort;comment:排序"`
	PermStatus int8   `json:"perm_status" gorm:"column:perm_status;comment:状态 1启用 0禁用"`
}

func (p *SysPermission) TableName() string {
	return "sys_permission"
}

type SysPermissionList struct {
	ListMeta `json:",inline"`
	Items    []*SysPermission `json:"items"`
}
