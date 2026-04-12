package projectHandler

import (
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	taskHandler "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/handler/task"
)

type ProjectTasksResponse struct {
	Project *models.Project `json:"project"`
	// Tasks   []*models.Task  `json:"tasks"`
	Tasks taskHandler.TasksResponse `json:"tasks"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type userProjectsResponse struct {
	Projects       []*models.Project     `json:"projects"`
	PaginationInfo engine.PaginationInfo `json:"pagination"`
}
