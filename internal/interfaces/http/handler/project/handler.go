package projectHandler

import (
	"net/http"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/middlewares"
	projectService "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/project"
)

type projectHandler struct {
	engine         *engine.HttpEngine
	middleware     *middlewares.MiddlewareHandler
	projectService *projectService.ProjectService
}

func NewProjectHandler(
	engine *engine.HttpEngine,
	projectService *projectService.ProjectService,
	middleware *middlewares.MiddlewareHandler,
) *projectHandler {
	return &projectHandler{
		engine:         engine,
		projectService: projectService,
		middleware:     middleware,
	}
}

func (h *projectHandler) AddRoutes() {
	h.engine.AddRoutes([]engine.Route{
		{
			Method:      http.MethodGet,
			Endpoint:    "/projects",
			Description: "List projects the current user owns or has tasks in",
			Controller:  h.userProjects,
			Middleware: []engine.Middleware{
				h.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodPost,
			Endpoint:    "/projects",
			Description: "Create a project (owner = current user)",
			Controller:  h.create,
			Middleware: []engine.Middleware{
				h.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodGet,
			Endpoint:    "/projects/{id}",
			Description: "Get project details + its tasks",
			Controller:  h.projectByID,
			Middleware: []engine.Middleware{
				h.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodPatch,
			Endpoint:    "/projects/{id}",
			Description: "Update name/description (owner only)",
			Controller:  h.update,
			Middleware: []engine.Middleware{
				h.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodDelete,
			Endpoint:    "/projects/{id}",
			Description: "Delete project and all its tasks (owner only)",
			Controller:  h.delete,
			Middleware: []engine.Middleware{
				h.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodGet,
			Endpoint:    "/projects/{id}/stats",
			Description: "Get task counts by status and by assignee",
			Controller:  h.stats,
			Middleware: []engine.Middleware{
				h.middleware.ValidateAuthToken,
			},
		},
	})
}
