package projectHandler

import (
	"time"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

type ProjectTasksResponse struct {
	Project *models.Project `json:"project"`
	Tasks   []*models.Task  `json:"tasks"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}
