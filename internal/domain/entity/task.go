package entity

import "time"

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
