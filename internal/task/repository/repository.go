package repository

import (
	"context"
	"database/sql"
	"errors"

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
func (r *repository) FindAllByUserID(ctx context.Context, userID entity.UserID) ([]entity.Task, error) {
	q := `SELECT id, content, description, is_completed, due_date, created_at, updated_at FROM tasks WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]entity.Task, 0)
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(&task.ID, &task.Content, &task.Description, &task.IsCompleted, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
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

// VerifyAvailableByID check if a task is available by id.
func (r *repository) VerifyAvailableByID(ctx context.Context, taskID entity.TaskID) error {
	var id string
	q := `SELECT id FROM tasks WHERE id = $1`
	row := r.db.QueryRowContext(ctx, q, taskID)
	err := row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrTaskNotFound
	} else if err != nil {
		return err
	}
	return nil
}

// DeleteByID delete a task by id.
func (r *repository) DeleteByID(ctx context.Context, taskID entity.TaskID) error {
	q := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, taskID)
	if err != nil {
		return err
	}
	return nil
}
