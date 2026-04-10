package models

import "time"

// -------- STATUS --------

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

func (p TaskPriority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

type Task struct {
	ID          string
	Title       string
	Description *string
	Status      TaskStatus
	Priority    TaskPriority
	ProjectID   string
	AssigneeID  *string
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
