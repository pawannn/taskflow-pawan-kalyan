package taskService

import (
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
)

type TaskService struct {
	taskRepo    domainRepo.TaskRepository
	projectRepo domainRepo.ProjectRepository
}

func NewTaskService(
	taskRepo domainRepo.TaskRepository,
	projectRepo domainRepo.ProjectRepository,
) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}
