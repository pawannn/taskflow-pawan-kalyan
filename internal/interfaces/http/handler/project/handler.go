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

func (pH *projectHandler) AddRoutes() {
	pH.engine.AddRoutes([]engine.Route{
		{
			Method:      http.MethodGet,
			Endpoint:    "/projects",
			Description: "List projects the current user owns or has tasks in",
			Controller:  pH.getAll,
			Middleware: []engine.Middleware{
				pH.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodPost,
			Endpoint:    "/projects",
			Description: "Create a project (owner = current user)",
			Controller:  pH.create,
			Middleware: []engine.Middleware{
				pH.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodGet,
			Endpoint:    "/projects/:id",
			Description: "Get project details + its tasks",
			Controller:  pH.getByID,
			Middleware: []engine.Middleware{
				pH.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodPatch,
			Endpoint:    "/projects/:id",
			Description: "Update name/description (owner only)",
			Controller:  pH.update,
			Middleware: []engine.Middleware{
				pH.middleware.ValidateAuthToken,
			},
		},
		{
			Method:      http.MethodDelete,
			Endpoint:    "/projects/:id",
			Description: "Delete project and all its tasks (owner only)",
			Controller:  pH.delete,
			Middleware: []engine.Middleware{
				pH.middleware.ValidateAuthToken,
			},
		},
	})
}
