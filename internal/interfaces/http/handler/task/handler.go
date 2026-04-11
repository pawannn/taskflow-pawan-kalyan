package taskHandler

import (
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/middlewares"
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

func (tH *taskHandler) AddRoutes() {
	tH.engine.AddRoutes([]engine.Route{
		{
			Method:      http.MethodPost,
			Endpoint:    "/projects/:id/tasks",
			Description: "Create a task",
			Controller:  tH.create,
			Middleware: []engine.Middleware{
				tH.middleware.ValidateAuthToken,
			},
		},
	})
}
