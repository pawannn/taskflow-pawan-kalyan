package projectRepository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

func (pR *ProjectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func (pR *ProjectRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
			SELECT DISTINCT p.id, p.name, p.description, p.owner_id, p.created_at, p.updated_at
			FROM projects p
			LEFT JOIN tasks t ON t.project_id = p.id
			WHERE p.owner_id = $1 OR t.assignee_id = $1
		`

	rows, err := pR.db.Query(ctx, query, userID)
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
