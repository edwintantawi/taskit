package entity

import "time"

type TaskID string

// Task represents a task in the system.
type Task struct {
	ID          TaskID
	UserID      UserID
	ProjectID   NullString
	Content     string
	Description string
	IsCompleted bool
	DueDate     NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
