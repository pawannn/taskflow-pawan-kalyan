package taskHandler

import (
	"net/http"

	engine "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	middlewares "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/middlewares"
	taskService "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/task"
)

type taskHandler struct {
	engine      *engine.HttpEngine
	middleware  *middlewares.MiddlewareHandler
	taskService *taskService.TaskService
}

func NewTaskHandler(
	engine *engine.HttpEngine,
	middleware *middlewares.MiddlewareHandler,
	taskService *taskService.TaskService,
) *taskHandler {
	return &taskHandler{
		engine:      engine,
		middleware:  middleware,
		taskService: taskService,
	}
}

func (h *taskHandler) AddRoutes() {
	h.engine.AddRoutes([]engine.Route{
		{
			Method:      http.MethodGet,
			Endpoint:    "/projects/{id}/tasks",
			Description: "List tasks — support `?status=` and `?assignee=` filters",
			Controller:  h.getByProject,
			Middleware: []engine.Middleware{
				h.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodPost,
			Endpoint:    "/projects/{id}/tasks",
			Description: "Create a task",
			Controller:  h.create,
			Middleware: []engine.Middleware{
				h.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodPatch,
			Endpoint:    "/tasks/{id}",
			Description: "Update title, description, status, priority, assignee, due_date",
			Controller:  h.Update,
			Middleware: []engine.Middleware{
				h.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodDelete,
			Endpoint:    "/tasks/{id}",
			Description: "Delete task (project owner or task creator only)",
			Controller:  h.Delete,
			Middleware: []engine.Middleware{
				h.middleware.ValidateAuthToken,
			},
		},
	})
}
