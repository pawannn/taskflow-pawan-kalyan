package repository

import "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"

type ProjectRepository interface {
	Create(project *models.Project) error
	GetByID(id string) (*models.Project, error)
	GetByUserID(userID string) ([]*models.Project, error)
	Update(project *models.Project) error
	Delete(id string) error
}
