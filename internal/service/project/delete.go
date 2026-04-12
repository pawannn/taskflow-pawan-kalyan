package projectService

import (
	"context"
	"net/http"

	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
)

func (s *ProjectService) DeleteProject(ctx context.Context, projectID, userID string) apperr.Err {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	if project == nil {
		return apperr.NewErr(http.StatusNotFound, apperr.ErrNotFound, nil)
	}

	if project.OwnerID != userID {
		return apperr.NewErr(http.StatusForbidden, apperr.ErrForbidden, nil)
	}

	err = s.projectRepo.Delete(ctx, projectID)
	if err != nil {
		return apperr.NewErr(http.StatusInternalServerError, apperr.ErrInternalError, err)
	}

	return apperr.NoErr
}
