package projectHandler

import (
	"encoding/json"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
)

func (pH *projectHandler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	meta := pH.engine.ParseContext(ctx)
	if meta.UserID == "" {
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusForbidden, domain.ErrForbidded, nil)
		return
	}

	query := r.URL.Query()

	projectID := query.Get("id")
	if projectID == "" {
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "validation failed", map[string]string{
			"id": "is required",
		})
		return
	}

	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pH.engine.SendResponse(w, meta.ReqID, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	project, err := pH.projectService.UpdateProject(ctx, projectID, meta.UserID, req.Name, req.Description)
	if !err.IsEmpty() {
		pH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	pH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "project updated", project)
}
