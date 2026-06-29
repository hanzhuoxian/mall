package types

type Permission struct {
	ObjectMeta `json:"metadata,omitempty"`
	ParentID   uint64 `json:"parent_id"`
	PermCode   string `json:"perm_code"`
	PermType   string `json:"perm_type"`
	Path       string `json:"path"`
	Icon       string `json:"icon"`
}
