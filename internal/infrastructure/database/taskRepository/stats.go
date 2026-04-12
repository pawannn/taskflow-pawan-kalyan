package taskRepository

import (
	"context"
	"time"

	models "github.com/pawannn/taskflow-pawan-kalyan/backend/internal/domain/models"
)

// GetProjectStats retrieves aggregated task statistics for a given project.
func (r *taskRepository) GetProjectStats(ctx context.Context, projectID string) (*models.ProjectStats, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	statusQuery := `
		SELECT
			COUNT(*) FILTER (WHERE status = 'todo') AS todo,
			COUNT(*) FILTER (WHERE status = 'in_progress') AS in_progress,
			COUNT(*) FILTER (WHERE status = 'done') AS done,
			COUNT(*) AS total
		FROM tasks
		WHERE project_id = $1
	`

	var stats models.ProjectStats
	err := r.db.QueryRow(ctx, statusQuery, projectID).Scan(
		&stats.StatusCounts.Todo,
		&stats.StatusCounts.InProgress,
		&stats.StatusCounts.Done,
		&stats.Total,
	)
	if err != nil {
		return nil, err
	}

	assigneeQuery := `
		SELECT assignee_id, COUNT(*) AS count
		FROM tasks
		WHERE project_id = $1
		GROUP BY assignee_id
		ORDER BY count DESC
	`

	rows, err := r.db.Query(ctx, assigneeQuery, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ac models.AssigneeCount
		if err := rows.Scan(&ac.AssigneeID, &ac.Count); err != nil {
			return nil, err
		}
		stats.AssigneeCounts = append(stats.AssigneeCounts, ac)
	}

	return &stats, nil
}
