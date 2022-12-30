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

// GetAll get all tasks.
func (u *usecase) GetAll(ctx context.Context, payload *domain.GetAllTaskIn) ([]domain.GetAllTaskOut, error) {
	tasks, err := u.taskRepository.FindAllByUserID(ctx, payload.UserID)
	if err != nil {
		return nil, err
	}

	output := make([]domain.GetAllTaskOut, len(tasks))
	for i, task := range tasks {
		output[i] = domain.GetAllTaskOut{
			ID:          task.ID,
			Content:     task.Content,
			Description: task.Description,
			IsCompleted: task.IsCompleted,
			DueDate:     task.DueDate,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		}
	}
	return output, nil
}
