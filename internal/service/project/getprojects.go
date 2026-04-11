package projectService

import (
	"context"
	"net/http"

	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
	TaskFlowErr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
)

func (pS *ProjectService) GetProjects(ctx context.Context, userID string, page int, limit int) ([]*models.Project, TaskFlowErr.Err) {
	offset := (page - 1) * limit
	pagination := domainRepo.Pagination{
		Offset: offset,
		Limit:  limit,
	}

	projects, err := pS.projectRepo.GetByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrFetchProject, err)
	}

	return projects, TaskFlowErr.NoErr
}

func (s *ProjectService) GetProjectByID(ctx context.Context, projectID string) (*models.Project, []*models.Task, TaskFlowErr.Err) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrFetchProject, err)
	}

	if project == nil {
		return nil, nil, TaskFlowErr.NewErr(http.StatusNotFound, domain.ErrNotFound, nil)
	}

	tasks, err := s.taskRepo.GetByProjectID(ctx, projectID, nil, nil)
	if err != nil {
		return nil, nil, TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrFetchProject, err)
	}

	return project, tasks, TaskFlowErr.NoErr
}
