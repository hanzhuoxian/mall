package model

type GetOptions struct {
}

type CreateOptions struct {
	DryRun []string `json:"dryRun,omitempty"`
}

type UpdateOptions struct {
	DryRun []string `json:"dryRun,omitempty"`
}

type PatchOptions struct {
	DryRun []string `json:"dryRun,omitempty"`
	Force  bool     `json:"force,omitempty"`
}

type DeleteOptions struct {
	DryRun   []string `json:"dryRun,omitempty"`
	Unscoped bool     `json:"unscoped"`
}

type ListOptions struct {
	LabelSelector  string `json:"labelSelector,omitempty" form:"labelSelector"`
	FieldSelector  string `json:"fieldSelector,omitempty" form:"fieldSelector"`
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`
	Offset         *int64 `json:"offset,omitempty" form:"offset"`
	Limit          *int64 `json:"limit,omitempty" form:"limit"`
	Export         bool   `json:"export" form:"export"`
}
