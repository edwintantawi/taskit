package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
)

type TaskUsecaseTestSuite struct {
	suite.Suite
}

func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseTestSuite))
}

func (s *TaskUsecaseTestSuite) TestCreate() {
	s.Run("it should return error when task validation fail", func() {
		ctx := context.Background()
		payload := &domain.CreateTaskIn{}
		usecase := New(nil)

		r, err := usecase.Create(ctx, payload)

		s.Error(err)
		s.Empty(r)
	})

	s.Run("it should return error when fail to store task", func() {
		ctx := context.Background()
		payload := domain.CreateTaskIn{Content: "Task 1", Description: "Task 1 Description"}
		output := domain.CreateTaskOut{}

		mockRepo := &mocks.TaskRepository{}
		mockRepo.On("Store", ctx, mock.Anything).Return(output.ID, errors.New("create task fail"))

		usecase := New(mockRepo)
		r, err := usecase.Create(ctx, &payload)

		s.Equal(errors.New("create task fail"), err)
		s.Empty(r)
	})

	s.Run("it should return return error nil when create task successfully", func() {
		ctx := context.Background()
		payload := domain.CreateTaskIn{Content: "Task 1", Description: "Task 1 Description"}
		task := entity.Task{Content: payload.Content, Description: payload.Description, DueDate: payload.DueDate}
		output := domain.CreateTaskOut{ID: entity.TaskID("xxxxx")}

		mockRepo := &mocks.TaskRepository{}
		mockRepo.On("Store", ctx, &task).Return(output.ID, nil)

		usecase := New(mockRepo)
		r, err := usecase.Create(ctx, &payload)

		s.NoError(err)
		s.Equal(output, r)
	})
}
