package dto

import (
	"time"

	"github.com/edwintantawi/taskit/internal/domain/entity"
)

// CreateTaskIn represents the input of task creation.
type CreateTaskIn struct {
	UserID      entity.UserID `json:"-"`
	Content     string        `json:"content"`
	Description string        `json:"description"`
	DueDate     *time.Time    `json:"due_date"`
}

// CreateTaskOut represents the output of task creation.
type CreateTaskOut struct {
	ID entity.TaskID `json:"id"`
}

// GetAllTaskIn represents the input of task retrieval.
type GetAllTaskIn struct {
	UserID entity.UserID `json:"-"`
}

// GetAllTaskOut represents the output of task retrieval.
type GetAllTaskOut struct {
	ID          entity.TaskID `json:"id"`
	Content     string        `json:"content"`
	Description string        `json:"description"`
	IsCompleted bool          `json:"is_completed"`
	DueDate     *time.Time    `json:"due_date"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// RemoveTaskIn represents the input of task removal.
type RemoveTaskIn struct {
	TaskID entity.TaskID `json:"-"`
	UserID entity.UserID `json:"-"`
}

// GetTaskByIDIn represents the input of task retrieval.
type GetTaskByIDIn struct {
	TaskID entity.TaskID `json:"-"`
	UserID entity.UserID `json:"-"`
}

// GetTaskByIDOut represents the output of task retrieval.
type GetTaskByIDOut struct {
	ID          entity.TaskID `json:"id"`
	Content     string        `json:"content"`
	Description string        `json:"description"`
	IsCompleted bool          `json:"is_completed"`
	DueDate     *time.Time    `json:"due_date"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// UpdateTaskIn represents the input of task update
type UpdateTaskIn struct {
	TaskID      entity.TaskID `json:"-"`
	UserID      entity.UserID `json:"-"`
	Content     string        `json:"content"`
	Description string        `json:"description"`
	IsCompleted bool          `json:"is_completed"`
	DueDate     *time.Time    `json:"due_date"`
}

// UpdateTaskOut represents the output of task update
type UpdateTaskOut struct {
	ID entity.TaskID `json:"id"`
}
