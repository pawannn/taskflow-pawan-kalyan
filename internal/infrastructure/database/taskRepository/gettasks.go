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
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT id, title, description, status, priority,
		       project_id, assignee_id, creator_id, due_date,
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
		&t.CreatorID,
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

func (r *taskRepository) GetByProjectID(
	ctx context.Context,
	projectID string,
	filter *domainRepo.TaskFilter,
	pagination *domainRepo.Pagination,
) ([]*models.Task, bool, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT id, title, description, status, priority,
		       project_id, assignee_id, creator_id, due_date,
		       created_at, updated_at
		FROM tasks
		WHERE project_id = $1
	`

	args := []interface{}{projectID}
	argIndex := 2

	if filter != nil {
		if filter.Status != nil {
			query += " AND status = $" + fmt.Sprint(argIndex)
			args = append(args, *filter.Status)
			argIndex++
		}

		if filter.AssigneeID != nil {
			query += " AND assignee_id = $" + fmt.Sprint(argIndex)
			args = append(args, *filter.AssigneeID)
			argIndex++
		}
	}

	query += " ORDER BY created_at DESC"

	limit := 10
	offset := 0

	if pagination != nil {
		limit = pagination.Limit
		offset = pagination.Offset
	}

	query += " LIMIT $" + fmt.Sprint(argIndex)
	args = append(args, limit+1)
	argIndex++

	query += " OFFSET $" + fmt.Sprint(argIndex)
	args = append(args, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, false, err
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
			&t.CreatorID,
			&t.DueDate,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, false, err
		}

		tasks = append(tasks, &t)
	}

	hasNext := false

	if len(tasks) > limit {
		hasNext = true
		tasks = tasks[:limit]
	}

	return tasks, hasNext, nil
}

func (r *taskRepository) CanUpdateTask(ctx context.Context, taskID, userID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM tasks t
			JOIN projects p ON p.id = t.project_id
			WHERE t.id = $1
			AND (
				p.owner_id = $2
				OR t.assignee_id = $2
			)
		)
	`

	var allowed bool

	err := r.db.QueryRow(ctx, query, taskID, userID).Scan(&allowed)
	if err != nil {
		return false, err
	}

	return allowed, nil
}
