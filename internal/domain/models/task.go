package models

import "time"

// TaskStatus represents the status of a task.
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

// IsValid checks if the task status is valid.
func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

// TaskPriority represents the priority level of a task.
type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

// IsValid checks if the task priority is valid.
func (p TaskPriority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

// Task represents a task with its metadata and relationships.
type Task struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description *string       `json:"description"`
	Status      TaskStatus    `json:"status"`
	Priority    *TaskPriority `json:"priority"`
	ProjectID   string        `json:"project_id"`
	AssigneeID  *string       `json:"assignee_id"`
	CreatorID   string        `json:"creator_id"`
	DueDate     *time.Time    `json:"due_date"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
