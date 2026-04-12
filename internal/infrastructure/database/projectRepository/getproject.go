package projectRepository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
	domainRepo "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/repository"
)

// GetByID retrieves a project by its ID or returns nil if not found.
func (r *projectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)

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

// GetByUserID retrieves projects accessible to a user with pagination and indicates if more results exist.
func (r *projectRepository) GetByUserID(
	ctx context.Context,
	userID string,
	pagination domainRepo.Pagination,
) ([]*models.Project, bool, error) {

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

	limit := pagination.Limit
	offset := pagination.Offset

	rows, err := r.db.Query(ctx, query, userID, limit+1, offset)
	if err != nil {
		return nil, false, err
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
			return nil, false, err
		}

		projects = append(projects, &p)
	}

	hasNext := false
	if len(projects) > limit {
		hasNext = true
		projects = projects[:limit]
	}

	return projects, hasNext, nil
}

// IsPartOfProject checks whether a user is associated with a project as owner or assignee.
func (r *projectRepository) IsPartOfProject(
	ctx context.Context,
	projectID, userID string,
) (bool, error) {

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
					SELECT 1
					FROM tasks t
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
