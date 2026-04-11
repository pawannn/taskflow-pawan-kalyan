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

	if meta.UserID == "" {
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusForbidden, domain.ErrForbidded, nil)
		return
	}

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	fields := utils.ValidateRequired(map[string]string{
		"name": req.Name,
	})

	if len(fields) > 0 {
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, "validation failed", fields)
		return
	}

	project, err := pH.projectService.Create(ctx, req.Name, req.Description, meta.UserID)
	if !err.IsEmpty() {
		pH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	resp := ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID,
		CreatedAt:   project.CreatedAt,
	}

	pH.engine.SendResponse(w, meta.ReqID, http.StatusCreated, "project created", resp)
}
