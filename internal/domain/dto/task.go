package dto

import (
	"time"

	"github.com/edwintantawi/taskit/internal/domain/entity"
)

// TaskCreateIn represents the input of task creation.
type TaskCreateIn struct {
	UserID      entity.UserID   `json:"-"`
	Content     string          `json:"content"`
	Description string          `json:"description"`
	DueDate     entity.NullTime `json:"due_date"`
}

func (t *TaskCreateIn) Validate() error {
	switch {
	case t.Content == "":
		return ErrContentEmpty
	}
	return nil
}

// TaskCreateOut represents the output of task creation.
type TaskCreateOut struct {
	ID entity.TaskID `json:"id"`
}

// TaskGetAllIn represents the input of task retrieval.
type TaskGetAllIn struct {
	UserID entity.UserID `json:"-"`
}

// TaskGetAllOut represents the output of task retrieval.
type TaskGetAllOut struct {
	ID          entity.TaskID   `json:"id"`
	Content     string          `json:"content"`
	Description string          `json:"description"`
	IsDone      bool            `json:"is_done"`
	DueDate     entity.NullTime `json:"due_date"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// TaskRemoveIn represents the input of task removal.
type TaskRemoveIn struct {
	TaskID entity.TaskID `json:"-"`
	UserID entity.UserID `json:"-"`
}

// TaskGetByIDIn represents the input of task retrieval.
type TaskGetByIDIn struct {
	TaskID entity.TaskID `json:"-"`
	UserID entity.UserID `json:"-"`
}

// TaskGetByIDOut represents the output of task retrieval.
type TaskGetByIDOut struct {
	ID          entity.TaskID   `json:"id"`
	Content     string          `json:"content"`
	Description string          `json:"description"`
	IsDone      bool            `json:"is_done"`
	DueDate     entity.NullTime `json:"due_date"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// TaskUpdateIn represents the input of task update
type TaskUpdateIn struct {
	TaskID      entity.TaskID   `json:"-"`
	UserID      entity.UserID   `json:"-"`
	Content     string          `json:"content"`
	Description string          `json:"description"`
	IsDone      bool            `json:"is_done"`
	DueDate     entity.NullTime `json:"due_date"`
}

func (t *TaskUpdateIn) Validate() error {
	switch {
	case t.Content == "":
		return ErrContentEmpty
	}
	return nil
}

// TaskUpdateOut represents the output of task update
type TaskUpdateOut struct {
	ID entity.TaskID `json:"id"`
}
