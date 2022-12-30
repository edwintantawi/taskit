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

// FindAllByUserID get all tasks owned by a user by user id.
func (r *repository) FindAllByUserID(ctx context.Context, t *entity.Task) ([]entity.Task, error) {
	q := `SELECT id, content, description, due_date, created_at, updated_at FROM tasks WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, q, t.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]entity.Task, 0)
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(&task.ID, &task.Content, &task.Description, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
