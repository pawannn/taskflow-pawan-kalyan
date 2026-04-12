package taskHandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	apperr "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/apperror"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (h *taskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := h.engine.ParseContext(ctx)

	taskID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(taskID) {
		h.engine.Log.Warn(ctx, apperr.ErrValidationFailed, "fields", "id")
		h.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, apperr.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	err := h.taskService.DeleteTask(ctx, taskID, meta.UserID)
	if err != apperr.NoErr {
		if err.Data != nil {
			h.engine.Log.Error(ctx, "Task deletion failed", "error", err.Data)
		}

		h.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	h.engine.SendNoContent(w)
}
