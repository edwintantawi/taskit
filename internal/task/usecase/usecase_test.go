package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/test"
)

type TaskUsecaseTestSuite struct {
	suite.Suite
}

func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseTestSuite))
}

type dependency struct {
	taskRepository *mocks.TaskRepository
}

func (s *TaskUsecaseTestSuite) TestCreate() {
	s.Run("it should return error when validation fail", func() {
		ctx := context.Background()
		payload := &domain.CreateTaskIn{}
		usecase := New(nil)

		output, err := usecase.Create(ctx, payload)

		s.Error(err)
		s.Empty(output)
	})

	type args struct {
		ctx     context.Context
		payload *domain.CreateTaskIn
	}
	type expected struct {
		output domain.CreateTaskOut
		err    error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when task respository return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.CreateTaskIn{UserID: "user-xxxxx", Content: "task_content", Description: "content_description", DueDate: nil},
			},
			expected: expected{
				output: domain.CreateTaskOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("Store", context.Background(), &entity.Task{UserID: "user-xxxxx", Content: "task_content", Description: "content_description", DueDate: nil}).
					Return(entity.TaskID(""), test.ErrUnexpected)
			},
		},
		{
			name: "it should return error nil and output when task respository return nil error",
			args: args{
				ctx:     context.Background(),
				payload: &domain.CreateTaskIn{UserID: "user-xxxxx", Content: "task_content", Description: "content_description", DueDate: nil},
			},
			expected: expected{
				output: domain.CreateTaskOut{ID: "task-xxxxx"},
				err:    nil,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("Store", context.Background(), &entity.Task{UserID: "user-xxxxx", Content: "task_content", Description: "content_description", DueDate: nil}).
					Return(entity.TaskID("task-xxxxx"), nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			d := &dependency{
				taskRepository: &mocks.TaskRepository{},
			}
			t.setup(d)

			usecase := New(d.taskRepository)
			output, err := usecase.Create(t.args.ctx, t.args.payload)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.output, output)
		})
	}
}
