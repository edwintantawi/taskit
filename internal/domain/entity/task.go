package entity

import (
	"errors"
	"time"
)

var (
	ErrContentEmpty = errors.New("task.entity.content_empty")
)

type TaskID string
type TaskDueDate *time.Time

// Task represents a task in the system.
type Task struct {
	ID          TaskID
	UserID      UserID
	Content     string
	Description string
	IsCompleted bool
	DueDate     TaskDueDate
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate task fields.
func (t *Task) Validate() error {
	switch {
	case t.Content == "":
		return ErrContentEmpty
	}
	return nil
}
