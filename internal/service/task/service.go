package taskService

import (
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
)

type TaskService struct {
	taskRepo    domainRepo.TaskRepository
	projectRepo domainRepo.ProjectRepository
	userRepo    domainRepo.UserRepository
}

func NewTaskService(
	taskRepo domainRepo.TaskRepository,
	projectRepo domainRepo.ProjectRepository,
	userRepo domainRepo.UserRepository,
) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}
