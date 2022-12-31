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

// Remove remove a task.
func (u *usecase) Remove(ctx context.Context, payload *domain.RemoveTaskIn) error {
	task, err := u.taskRepository.FindByID(ctx, payload.TaskID)
	if err != nil {
		return err
	}
	if task.UserID != payload.UserID {
		return domain.ErrTaskAuthorization
	}
	if err := u.taskRepository.DeleteByID(ctx, payload.TaskID); err != nil {
		return err
	}
	return nil
}

// GetByID get task by id.
func (u *usecase) GetByID(ctx context.Context, payload *domain.GetTaskByIDIn) (domain.GetTaskByIDOut, error) {
	task, err := u.taskRepository.FindByID(ctx, payload.TaskID)
	if err != nil {
		return domain.GetTaskByIDOut{}, err
	}
	if task.UserID != payload.UserID {
		return domain.GetTaskByIDOut{}, domain.ErrTaskAuthorization
	}

	output := domain.GetTaskByIDOut{
		ID:          task.ID,
		Content:     task.Content,
		Description: task.Description,
		IsCompleted: task.IsCompleted,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	return output, nil
}

func (u *usecase) Update(ctx context.Context, payload *domain.UpdateTaskIn) (domain.UpdateTaskOut, error) {
	task, err := u.taskRepository.FindByID(ctx, payload.TaskID)
	if err != nil {
		return domain.UpdateTaskOut{}, err
	}
	if task.UserID != payload.UserID {
		return domain.UpdateTaskOut{}, domain.ErrTaskAuthorization
	}

	task.Content = payload.Content
	task.Description = payload.Description
	task.IsCompleted = payload.IsCompleted
	task.DueDate = payload.DueDate

	taskID, err := u.taskRepository.Update(ctx, &task)
	if err != nil {
		return domain.UpdateTaskOut{}, err
	}
	return domain.UpdateTaskOut{ID: taskID}, nil
}
