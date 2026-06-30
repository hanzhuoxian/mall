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
	Offset *int64 `json:"offset,omitempty"`
	Limit  *int64 `json:"limit,omitempty"`
}
