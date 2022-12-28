package usecase

import (
	"context"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type usecase struct {
	taskRepository domain.TaskRepository
}

// New create a new usecase.
func New(taskRepository domain.TaskRepository) domain.TaskUsecase {
	return &usecase{taskRepository: taskRepository}
}

// Create create a new task.
func (u *usecase) Create(ctx context.Context, payload *domain.CreateTaskIn) (domain.CreateTaskOut, error) {
	task := &entity.Task{UserID: payload.UserID, Content: payload.Content, Description: payload.Description, DueDate: payload.DueDate}
	if err := task.Validate(); err != nil {
		return domain.CreateTaskOut{}, err
	}
	taskID, err := u.taskRepository.Store(ctx, task)
	if err != nil {
		return domain.CreateTaskOut{}, err
	}
	return domain.CreateTaskOut{ID: taskID}, nil
}
