package projectHandler

import (
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
)

func (pH *projectHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meta := pH.engine.ParseContext(ctx)
	if meta.UserID == "" {
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusForbidden, domain.ErrForbidded, nil)
		return
	}

	query := r.URL.Query()
	projectID := query.Get("id")

	err := pH.projectService.DeleteProject(ctx, projectID, meta.UserID)
	if !err.IsEmpty() {
		pH.engine.SendResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	pH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "project deleted", nil)
}
