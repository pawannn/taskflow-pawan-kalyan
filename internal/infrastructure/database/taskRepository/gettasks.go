package taskRepository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
)

func (r *taskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, description, status, priority,
		       project_id, assignee_id, due_date,
		       created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var t models.Task

	err := row.Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Priority,
		&t.ProjectID,
		&t.AssigneeID,
		&t.DueDate,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (r *taskRepository) GetByProjectID(ctx context.Context, projectID string, filter *domainRepo.TaskFilter, pagination *domainRepo.Pagination) ([]*models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, description, status, priority,
		       project_id, assignee_id, due_date,
		       created_at, updated_at
		FROM tasks
		WHERE project_id = $1
	`

	args := []interface{}{projectID}
	argIndex := 2

	if filter.Status != nil {
		query += " AND status = $" + fmt.Sprint(argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.Status != nil {
		query += " AND assignee_id = $" + fmt.Sprint(argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task

	for rows.Next() {
		var t models.Task

		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.ProjectID,
			&t.AssigneeID,
			&t.DueDate,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, &t)
	}

	return tasks, nil
}
