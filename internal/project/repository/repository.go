package repository

import (
	"context"
	"database/sql"
	"errors"

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

// FindAllByUserID get all project by user id
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

// FindByID get project by project id.
func (r *Repository) FindByID(ctx context.Context, projectID entity.ProjectID) (entity.Project, error) {
	var project entity.Project
	q := `SELECT id, user_id, title, created_at, updated_at FROM projects WHERE id = $1`
	row := r.db.QueryRowContext(ctx, q, projectID)
	err := row.Scan(&project.ID, &project.UserID, &project.Title, &project.CreatedAt, &project.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Project{}, domain.ErrProjectNotFound
	} else if err != nil {
		return entity.Project{}, err
	}
	return project, nil
}
