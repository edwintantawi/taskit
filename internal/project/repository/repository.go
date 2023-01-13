package repository

import (
	"context"
	"database/sql"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type Repository struct {
	db         *sql.DB
	idProvider domain.IDProvider
}

// New create a new project repository.
func New(db *sql.DB, idProvider domain.IDProvider) Repository {
	return Repository{db: db, idProvider: idProvider}
}

// Store save a new project to database.
func (r *Repository) Store(ctx context.Context, p *entity.Project) (entity.ProjectID, error) {
	id := entity.ProjectID(r.idProvider.Generate())
	q := `INSERT INTO projects (id, user_id, title) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, q, id, p.UserID, p.Title)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) FindAllByUserID(ctx context.Context, userID entity.UserID) ([]entity.Project, error) {
	q := `SELECT id, user_id, title, created_at, updated_at FROM projects WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]entity.Project, 0)
	for rows.Next() {
		var project entity.Project
		err := rows.Scan(&project.ID, &project.UserID, &project.Title, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
