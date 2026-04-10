package repository

import "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"

type TaskRepository interface {
	Create(task *models.Task) error
	GetByID(id string) (*models.Task, error)
	GetByProjectID(projectID string, status *string, assigneeID *string) ([]*models.Task, error)
	Update(task *models.Task) error
	Delete(id string) error
}
