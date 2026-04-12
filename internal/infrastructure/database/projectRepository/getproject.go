package projectRepository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
)

func (pR *projectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE id = $1`

	row := pR.db.QueryRow(ctx, query, id)

	var project models.Project

	err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.OwnerID,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &project, nil
}

func (pR *projectRepository) GetByUserID(ctx context.Context, userID string, pagination domainRepo.Pagination) ([]*models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT p.id, p.name, p.description, p.owner_id, p.created_at, p.updated_at
		FROM projects p
		WHERE
			p.owner_id = $1
			OR EXISTS (
				SELECT 1 FROM tasks t
				WHERE t.project_id = p.id AND t.assignee_id = $1
			)
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := pR.db.Query(ctx, query, userID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project

	for rows.Next() {
		var p models.Project

		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.OwnerID,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}

		projects = append(projects, &p)
	}

	return projects, nil
}

func (r *projectRepository) IsPartOfProject(ctx context.Context, projectID, userID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM projects p
			WHERE p.id = $1
			AND (
				p.owner_id = $2
				OR EXISTS (
					SELECT 1 FROM tasks t
					WHERE t.project_id = p.id
					AND t.assignee_id = $2
				)
			)
		)
	`

	var isAuthorized bool

	err := r.db.QueryRow(ctx, query, projectID, userID).Scan(&isAuthorized)
	if err != nil {
		return false, err
	}

	return isAuthorized, nil
}
