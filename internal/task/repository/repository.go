package repository

import (
	"context"
	"database/sql"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type repository struct {
	db         *sql.DB
	idProvider domain.IDProvider
}

// New create a new task repository.
func New(db *sql.DB, idProvider domain.IDProvider) domain.TaskRepository {
	return &repository{db: db, idProvider: idProvider}
}

// Store save a new task.
func (r *repository) Store(ctx context.Context, t *entity.Task) (entity.TaskID, error) {
	id := r.idProvider.Generate()
	q := `INSERT INTO tasks (id, user_id, content, description, due_date) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, q, id, t.UserID, t.Content, t.Description, t.DueDate)
	if err != nil {
		return "", err
	}
	return entity.TaskID(id), nil
}
