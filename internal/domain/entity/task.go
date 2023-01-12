package entity

import "time"

type TaskID string

// Task represents a task in the system.
type Task struct {
	ID          TaskID
	UserID      UserID
	Content     string
	Description string
	IsDone      bool
	DueDate     NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
