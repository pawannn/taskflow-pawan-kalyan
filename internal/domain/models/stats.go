package models

// ProjectStats represents aggregated task statistics for a project.
type ProjectStats struct {
	Total          int             `json:"total"`
	StatusCounts   StatusCounts    `json:"status_counts"`
	AssigneeCounts []AssigneeCount `json:"assignee_counts"`
}

// StatusCounts represents task counts grouped by status.
type StatusCounts struct {
	Todo       int `json:"todo"`
	InProgress int `json:"in_progress"`
	Done       int `json:"done"`
}

// AssigneeCount represents task count for a specific assignee (nil if unassigned).
type AssigneeCount struct {
	AssigneeID *string `json:"assignee_id"`
	Count      int     `json:"count"`
}
