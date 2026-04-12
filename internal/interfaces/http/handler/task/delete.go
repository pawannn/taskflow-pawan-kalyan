package taskHandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	domain "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain"
	Error "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/pkg/taskflowErr"
	utils "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/utils"
)

func (tH *taskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meta := tH.engine.ParseContext(ctx)

	taskID := chi.URLParam(r, "id")
	if !utils.IsValidUUID(taskID) {
		tH.engine.Log.Warn(ctx, domain.ErrValidationFailed, "fields", "id")
		tH.engine.SendErrorResponse(w, meta.ReqID, http.StatusBadRequest, domain.ErrValidationFailed, map[string]string{
			"id": "is invalid",
		})
		return
	}

	err := tH.taskService.DeleteTask(ctx, taskID, meta.UserID)
	if err != Error.NoErr {
		if err.Data != nil {
			tH.engine.Log.Error(ctx, "Task deletion failed", "error", err.Data)
		}

		tH.engine.SendErrorResponse(w, meta.ReqID, err.Code, err.Message, nil)
		return
	}

	tH.engine.SendNoContent(w)
}
