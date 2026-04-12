package projectHandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (pH *projectHandler) stats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := pH.engine.ParseContext(ctx)

	projectID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(projectID) {
		pH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, domain.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	stats, err := pH.projectService.GetProjectStats(ctx, projectID, meta.UserID)
	if !err.IsEmpty() {
		if err.Data != nil {
			pH.engine.Log.Error(ctx, "fetch project stats failed", "error", err.Data)
		}
		pH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	pH.engine.SendResponse(w, meta.ReqID, http.StatusOK, "stats fetched", stats)
}
