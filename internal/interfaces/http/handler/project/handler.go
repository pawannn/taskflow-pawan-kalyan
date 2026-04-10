package projectHandler

import (
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/interfaces/http/engine"
	projectService "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/service/project"
)

type ProjectHandler struct {
	engine         *engine.HttpEngine
	projectService *projectService.ProjectService
}

func NewProjectHandler(engine *engine.HttpEngine, projectService *projectService.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		engine:         engine,
		projectService: projectService,
	}
}
