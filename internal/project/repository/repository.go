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

func New(db *sql.DB, idProvider domain.IDProvider) Repository {
	return Repository{db: db, idProvider: idProvider}
}

func (r *Repository) Store(ctx context.Context, p *entity.Project) (entity.ProjectID, error) {
	id := entity.ProjectID(r.idProvider.Generate())
	q := `INSERT INTO projects (id, user_id, title) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, q, id, p.UserID, p.Title)
	if err != nil {
		return "", err
	}
	return id, nil
}
