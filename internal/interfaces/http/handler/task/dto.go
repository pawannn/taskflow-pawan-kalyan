package taskHandler

import (
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	engine "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
)

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	AssigneeID  *string `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
	AssigneeID  *string `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
}

type TasksResponse struct {
	Tasks          []*models.Task        `json:"tasks"`
	PaginationInfo engine.PaginationInfo `json:"pagination"`
}
