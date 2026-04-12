package projectHandler

import (
	"encoding/json"
	"net/http"

	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (h *projectHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := h.engine.ParseContext(ctx)

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.engine.Log.Warn(ctx, apperr.ErrInvalidReqBody, "error", err)
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrInvalidReqBody, nil)
		return
	}

	fields := utils.ValidateRequired(map[string]string{
		"name": req.Name,
	})

	if len(fields) > 0 {
		h.engine.Log.Warn(ctx, apperr.ErrValidationFailed, "fields", fields)
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrValidationFailed, fields)
		return
	}

	project, err := h.projectService.Create(ctx, req.Name, req.Description, meta.UserID)
	if !err.IsEmpty() {
		if err.Data != nil {
			h.engine.Log.Error(ctx, "create project", "error", err.Data)
		}
		h.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	h.engine.Log.Info(ctx, "project created", "project_id", project.ID)

	h.engine.SendResponse(w, meta.ReqID, http.StatusCreated, "project created", project)
}
