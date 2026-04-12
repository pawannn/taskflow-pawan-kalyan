package integration_test

import (
	"context"

	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
)

type mockProjectRepo struct {
	projects map[string]*models.Project
}

func newMockProjectRepo() *mockProjectRepo {
	return &mockProjectRepo{projects: make(map[string]*models.Project)}
}

func (m *mockProjectRepo) Create(_ context.Context, p *models.Project) error {
	m.projects[p.ID] = p
	return nil
}

func (m *mockProjectRepo) GetByID(_ context.Context, id string) (*models.Project, error) {
	p, ok := m.projects[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockProjectRepo) GetByUserID(_ context.Context, userID string, _ domainRepo.Pagination) ([]*models.Project, bool, error) {
	var result []*models.Project
	for _, p := range m.projects {
		if p.OwnerID == userID {
			result = append(result, p)
		}
	}
	return result, false, nil
}

func (m *mockProjectRepo) Update(_ context.Context, p *models.Project) error {
	m.projects[p.ID] = p
	return nil
}

func (m *mockProjectRepo) Delete(_ context.Context, id string) error {
	delete(m.projects, id)
	return nil
}

// IsPartOfProject returns true if the user owns any project with this ID
// or has any task assigned in it. For test purposes: owner check is enough.
func (m *mockProjectRepo) IsPartOfProject(_ context.Context, projectID, userID string) (bool, error) {
	p, ok := m.projects[projectID]
	if !ok {
		return false, nil
	}
	return p.OwnerID == userID, nil
}

type mockTaskRepo struct {
	tasks map[string]*models.Task
}

func newMockTaskRepo() *mockTaskRepo {
	return &mockTaskRepo{tasks: make(map[string]*models.Task)}
}

func (m *mockTaskRepo) Create(_ context.Context, t *models.Task) error {
	m.tasks[t.ID] = t
	return nil
}

func (m *mockTaskRepo) GetByID(_ context.Context, id string) (*models.Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (m *mockTaskRepo) Update(_ context.Context, t *models.Task) error {
	m.tasks[t.ID] = t
	return nil
}

func (m *mockTaskRepo) Delete(_ context.Context, id string) error {
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepo) CanUpdateTask(_ context.Context, taskID, userID string) (bool, error) {
	t, ok := m.tasks[taskID]
	if !ok {
		return false, nil
	}
	return t.CreatorID == userID || (t.AssigneeID != nil && *t.AssigneeID == userID), nil
}

func (m *mockTaskRepo) GetByProjectID(_ context.Context, projectID string, _ *domainRepo.TaskFilter, _ *domainRepo.Pagination) ([]*models.Task, bool, error) {
	var result []*models.Task
	for _, t := range m.tasks {
		if t.ProjectID == projectID {
			result = append(result, t)
		}
	}
	return result, false, nil
}

func (m *mockTaskRepo) GetProjectStats(_ context.Context, projectID string) (*models.ProjectStats, error) {
	var stats models.ProjectStats
	for _, t := range m.tasks {
		if t.ProjectID != projectID {
			continue
		}
		stats.Total++
		switch t.Status {
		case models.StatusTodo:
			stats.StatusCounts.Todo++
		case models.StatusInProgress:
			stats.StatusCounts.InProgress++
		case models.StatusDone:
			stats.StatusCounts.Done++
		}
	}
	return &stats, nil
}
