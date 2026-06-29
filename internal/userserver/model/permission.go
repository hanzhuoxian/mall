package model

type Permission struct {
	ObjectMeta `json:"metadata,omitempty"`
	ParentID   uint64 `json:"parent_id" gorm:"comment:父级权限ID 0为顶级"`
	PermCode   string `json:"perm_code" gorm:"comment:权限编码"`
	PermType   string `json:"perm_type" gorm:"comment:权限类型 menu/button/api"`
	Path       string `json:"path" gorm:"comment:路由路径"`
	Icon       string `json:"icon" gorm:"comment:图标"`
	Sort       int    `json:"sort" gorm:"comment:排序"`
	PermStatus int8   `json:"perm_status" gorm:"comment:状态 1启用 0禁用"`
}

func (p *Permission) TableName() string {
	return "permission"
}

type PermissionList struct {
	ListMeta `json:",inline"`
	Item     []Permission `json:"item"`
}
