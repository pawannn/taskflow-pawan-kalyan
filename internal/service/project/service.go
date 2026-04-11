package projectService

import domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"

type ProjectService struct {
	projectRepo domainRepo.ProjectRepository
	taskRepo    domainRepo.TaskRepository
}

func NewProjectRepository(
	projectRepository domainRepo.ProjectRepository,
	taskRepo domainRepo.TaskRepository,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepository,
		taskRepo:    taskRepo,
	}
}
