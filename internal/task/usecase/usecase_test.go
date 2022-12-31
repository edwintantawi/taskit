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

func (s *TaskUsecaseTestSuite) TestGetAll() {
	type args struct {
		ctx     context.Context
		payload *domain.GetAllTaskIn
	}
	type expected struct {
		output []domain.GetAllTaskOut
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
				payload: &domain.GetAllTaskIn{UserID: "user-xxxxx"},
			},
			expected: expected{
				output: nil,
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindAllByUserID", context.Background(), entity.UserID("user-xxxxx")).
					Return(nil, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error nil and tasks when success",
			args: args{
				ctx:     context.Background(),
				payload: &domain.GetAllTaskIn{UserID: "user-xxxxx"},
			},
			expected: expected{
				output: []domain.GetAllTaskOut{
					{ID: "task-xxxxx", Content: "task_xxxxx_content", Description: "task_xxxxx_description", IsCompleted: false, DueDate: nil, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
					{ID: "task-yyyyy", Content: "task_yyyyy_content", Description: "task_yyyyy_description", IsCompleted: true, DueDate: &test.TimeAfterNow, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
				},
			},
			setup: func(d *dependency) {
				tasks := []entity.Task{
					{ID: "task-xxxxx", Content: "task_xxxxx_content", Description: "task_xxxxx_description", IsCompleted: false, DueDate: nil, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
					{ID: "task-yyyyy", Content: "task_yyyyy_content", Description: "task_yyyyy_description", IsCompleted: true, DueDate: &test.TimeAfterNow, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
				}

				d.taskRepository.On("FindAllByUserID", context.Background(), entity.UserID("user-xxxxx")).
					Return(tasks, nil)
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
			output, err := usecase.GetAll(t.args.ctx, t.args.payload)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.output, output)
		})
	}
}

func (s *TaskUsecaseTestSuite) TestRemove() {
	type args struct {
		ctx     context.Context
		payload *domain.RemoveTaskIn
	}
	type expected struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when task repository FindByID return unexpected error",
			args: args{
				ctx: context.Background(),
				payload: &domain.RemoveTaskIn{
					TaskID: "task-xxxxx",
					UserID: "user-xxxxx",
				},
			},
			expected: expected{
				err: test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(entity.Task{}, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error when task repository FindByID return unexpected error",
			args: args{
				ctx: context.Background(),
				payload: &domain.RemoveTaskIn{
					TaskID: "task-xxxxx",
					UserID: "user-xxxxx",
				},
			},
			expected: expected{
				err: domain.ErrTaskAuthorization,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(entity.Task{UserID: "user-yyyyy"}, nil)
			},
		},
		{
			name: "it should return error when task repository DeleteByID return unexpected error",
			args: args{
				ctx: context.Background(),
				payload: &domain.RemoveTaskIn{
					TaskID: "task-xxxxx",
					UserID: "user-xxxxx",
				},
			},
			expected: expected{
				err: test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(entity.Task{UserID: "user-xxxxx"}, nil)

				d.taskRepository.On("DeleteByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(test.ErrUnexpected)
			},
		},
		{
			name: "it should return error nil when success delete task",
			args: args{
				ctx: context.Background(),
				payload: &domain.RemoveTaskIn{
					TaskID: "task-xxxxx",
					UserID: "user-xxxxx",
				},
			},
			expected: expected{
				err: nil,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(entity.Task{UserID: "user-xxxxx"}, nil)

				d.taskRepository.On("DeleteByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(nil)
			},
		},
	}

	for _, t := range tests {
		deps := &dependency{
			taskRepository: &mocks.TaskRepository{},
		}
		t.setup(deps)

		usecase := New(deps.taskRepository)
		err := usecase.Remove(t.args.ctx, t.args.payload)

		s.Equal(t.expected.err, err)
	}
}
