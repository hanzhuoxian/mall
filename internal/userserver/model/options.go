package model

// GetOptions controls a single-record lookup.
type GetOptions struct{}

// CreateOptions controls record creation.
type CreateOptions struct {
	DryRun []string `json:"dryRun,omitempty"`
}

// UpdateOptions controls record updates.
type UpdateOptions struct {
	DryRun []string `json:"dryRun,omitempty"`
}

// PatchOptions controls partial record updates.
type PatchOptions struct {
	DryRun []string `json:"dryRun,omitempty"`
	Force  bool     `json:"force,omitempty"`
}

// DeleteOptions controls record deletion.
type DeleteOptions struct {
	DryRun   []string `json:"dryRun,omitempty"`
	Unscoped bool     `json:"unscoped"`
}

// ListOptions controls pagination for list queries.
type ListOptions struct {
	Offset *int64 `json:"offset,omitempty"`
	Limit  *int64 `json:"limit,omitempty"`
}
