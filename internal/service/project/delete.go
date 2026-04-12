package projectService

import (
	"context"
	"net/http"

	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	TaskFlowErr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
)

func (s *ProjectService) DeleteProject(ctx context.Context, projectID, userID string) TaskFlowErr.Err {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	if project == nil {
		return TaskFlowErr.NewErr(http.StatusNotFound, domain.ErrNotFound, nil)
	}

	if project.OwnerID != userID {
		return TaskFlowErr.NewErr(http.StatusForbidden, domain.ErrForbidden, nil)
	}

	err = s.projectRepo.Delete(ctx, projectID)
	if err != nil {
		return TaskFlowErr.NewErr(http.StatusInternalServerError, domain.ErrInternalError, err)
	}

	return TaskFlowErr.NoErr
}
