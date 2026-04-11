package projectHandler

import (
	"encoding/json"
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (pH *projectHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := pH.engine.ParseContext(ctx)

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pH.engine.Log.Warn(ctx, domain.ErrInvalidReqBody, "error", err)
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	fields := utils.ValidateRequired(map[string]string{
		"name": req.Name,
	})

	if len(fields) > 0 {
		pH.engine.Log.Warn(ctx, "validation failed", "fields", fields)
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "validation failed", fields)
		return
	}

	project, err := pH.projectService.Create(ctx, req.Name, req.Description, meta.UserID)
	if !err.IsEmpty() {
		if err.Data != nil {
			pH.engine.Log.Error(ctx, "project creation", "error", err.Data)
		}
		pH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	pH.engine.Log.Info(ctx, "project created", "project_id", project.ID)

	pH.engine.SendResponse(w, meta.ReqID, http.StatusCreated, "project created", project)
}
