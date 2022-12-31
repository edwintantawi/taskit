package domain

import (
	"time"

	"github.com/edwintantawi/taskit/internal/domain/entity"
)

// CreateUserIn represents the input of user creation.
type CreateUserIn struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUserOut represents the output of user creation.
type CreateUserOut struct {
	ID    entity.UserID `json:"id"`
	Email string        `json:"email"`
}

// LoginAuthIn represent login input.
type LoginAuthIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginAuthOut represent login output.
type LoginAuthOut struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// LogoutAuthIn represent logout input.
type LogoutAuthIn struct {
	RefreshToken string `json:"refresh_token"`
}

// GetProfileAuthIn represent get profile input.
type GetProfileAuthIn struct {
	UserID entity.UserID `json:"user_id"`
}

// GetProfileAuthOut represent get profile output.
type GetProfileAuthOut struct {
	ID    entity.UserID `json:"id"`
	Name  string        `json:"name"`
	Email string        `json:"email"`
}

// RefreshAuthIn represent refresh input.
type RefreshAuthIn struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshAuthOut represent refresh output.
type RefreshAuthOut struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

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
