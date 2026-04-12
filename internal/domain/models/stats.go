package models

type ProjectStats struct {
	Total          int             `json:"total"`
	StatusCounts   StatusCounts    `json:"status_counts"`
	AssigneeCounts []AssigneeCount `json:"assignee_counts"`
}

type StatusCounts struct {
	Todo       int `json:"todo"`
	InProgress int `json:"in_progress"`
	Done       int `json:"done"`
}

type AssigneeCount struct {
	AssigneeID *string `json:"assignee_id"`
	Count      int     `json:"count"`
}
