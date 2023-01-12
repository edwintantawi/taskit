package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
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
	type args struct {
		ctx     context.Context
		payload *dto.TaskCreateIn
	}
	type expected struct {
		output dto.TaskCreateOut
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
				ctx: context.Background(),
				payload: &dto.TaskCreateIn{
					UserID:      "user-xxxxx",
					Content:     "task_content",
					Description: "content_description",
					DueDate:     entity.NullTime{NullTime: sql.NullTime{Valid: false}},
				},
			},
			expected: expected{
				output: dto.TaskCreateOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("Store", context.Background(), &entity.Task{
					UserID:      "user-xxxxx",
					Content:     "task_content",
					Description: "content_description",
					DueDate:     entity.NullTime{NullTime: sql.NullTime{Valid: false}},
				}).Return(entity.TaskID(""), test.ErrUnexpected)
			},
		},
		{
			name: "it should return error nil and output when task respository return nil error",
			args: args{
				ctx: context.Background(),
				payload: &dto.TaskCreateIn{
					UserID:      "user-xxxxx",
					Content:     "task_content",
					Description: "content_description",
					DueDate:     entity.NullTime{NullTime: sql.NullTime{Valid: false}},
				},
			},
			expected: expected{
				output: dto.TaskCreateOut{ID: "task-xxxxx"},
				err:    nil,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("Store", context.Background(), &entity.Task{
					UserID:      "user-xxxxx",
					Content:     "task_content",
					Description: "content_description",
					DueDate:     entity.NullTime{NullTime: sql.NullTime{Valid: false}},
				}).Return(entity.TaskID("task-xxxxx"), nil)
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
		payload *dto.TaskGetAllIn
	}
	type expected struct {
		output []dto.TaskGetAllOut
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
				payload: &dto.TaskGetAllIn{UserID: "user-xxxxx"},
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
				payload: &dto.TaskGetAllIn{UserID: "user-xxxxx"},
			},
			expected: expected{
				output: []dto.TaskGetAllOut{
					{ID: "task-xxxxx", Content: "task_xxxxx_content", Description: "task_xxxxx_description", IsDone: false, DueDate: entity.NullTime{NullTime: sql.NullTime{Valid: false}}, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
					{ID: "task-yyyyy", Content: "task_yyyyy_content", Description: "task_yyyyy_description", IsDone: true, DueDate: entity.NullTime{NullTime: sql.NullTime{Time: test.TimeAfterNow, Valid: true}}, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
				},
			},
			setup: func(d *dependency) {
				tasks := []entity.Task{
					{ID: "task-xxxxx", Content: "task_xxxxx_content", Description: "task_xxxxx_description", IsDone: false, DueDate: entity.NullTime{NullTime: sql.NullTime{Valid: false}}, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
					{ID: "task-yyyyy", Content: "task_yyyyy_content", Description: "task_yyyyy_description", IsDone: true, DueDate: entity.NullTime{NullTime: sql.NullTime{Time: test.TimeAfterNow, Valid: true}}, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
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
		payload *dto.TaskRemoveIn
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
				payload: &dto.TaskRemoveIn{
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
			name: "it should return error ErrTaskAuthorization when task not own by the user",
			args: args{
				ctx: context.Background(),
				payload: &dto.TaskRemoveIn{
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
				payload: &dto.TaskRemoveIn{
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
				payload: &dto.TaskRemoveIn{
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
		d := &dependency{
			taskRepository: &mocks.TaskRepository{},
		}
		t.setup(d)

		usecase := New(d.taskRepository)
		err := usecase.Remove(t.args.ctx, t.args.payload)

		s.Equal(t.expected.err, err)
	}
}

func (s *TaskUsecaseTestSuite) TestGetByID() {
	type args struct {
		ctx     context.Context
		payload *dto.TaskGetByIDIn
	}
	type expected struct {
		output dto.TaskGetByIDOut
		err    error
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
				ctx:     context.Background(),
				payload: &dto.TaskGetByIDIn{},
			},
			expected: expected{
				output: dto.TaskGetByIDOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindByID", context.Background(), entity.TaskID("")).
					Return(entity.Task{}, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error ErrTaskAuthorization when task not own by the user",
			args: args{
				ctx: context.Background(),
				payload: &dto.TaskGetByIDIn{
					TaskID: "task-xxxxx",
					UserID: "user-xxxxx",
				},
			},
			expected: expected{
				output: dto.TaskGetByIDOut{},
				err:    domain.ErrTaskAuthorization,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(entity.Task{UserID: "user-yyyyy"}, nil)
			},
		},
		{
			name: "it should return error nil when success get task",
			args: args{
				ctx: context.Background(),
				payload: &dto.TaskGetByIDIn{
					TaskID: "task-xxxxx",
					UserID: "user-xxxxx",
				},
			},
			expected: expected{
				output: dto.TaskGetByIDOut{
					ID:          "task-xxxxx",
					Content:     "task_content",
					Description: "task_description",
					IsDone:      true,
					DueDate:     entity.NullTime{NullTime: sql.NullTime{Time: test.TimeAfterNow, Valid: true}},
					CreatedAt:   test.TimeBeforeNow,
					UpdatedAt:   test.TimeBeforeNow,
				},
				err: nil,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(entity.Task{
						ID:          "task-xxxxx",
						UserID:      "user-xxxxx",
						Content:     "task_content",
						Description: "task_description",
						IsDone:      true,
						DueDate:     entity.NullTime{NullTime: sql.NullTime{Time: test.TimeAfterNow, Valid: true}},
						CreatedAt:   test.TimeBeforeNow,
						UpdatedAt:   test.TimeBeforeNow,
					}, nil)
			},
		},
	}

	for _, t := range tests {
		d := &dependency{
			taskRepository: &mocks.TaskRepository{},
		}
		t.setup(d)

		usecase := New(d.taskRepository)
		output, err := usecase.GetByID(t.args.ctx, t.args.payload)

		s.Equal(t.expected.err, err)
		s.Equal(t.expected.output, output)
	}
}

func (s *TaskUsecaseTestSuite) TestUpdate() {
	type args struct {
		ctx     context.Context
		payload *dto.TaskUpdateIn
	}
	type expected struct {
		output dto.TaskUpdateOut
		err    error
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
				ctx:     context.Background(),
				payload: &dto.TaskUpdateIn{TaskID: "task-xxxxx"},
			},
			expected: expected{
				output: dto.TaskUpdateOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(entity.Task{}, test.ErrUnexpected)
			},
		},
		{
			name: "it should return error ErrTaskAuthorization when task is not own by the user",
			args: args{
				ctx:     context.Background(),
				payload: &dto.TaskUpdateIn{TaskID: "task-xxxxx", UserID: "user-xxxxx"},
			},
			expected: expected{
				output: dto.TaskUpdateOut{},
				err:    domain.ErrTaskAuthorization,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(entity.Task{UserID: "user-yyyyy"}, nil)
			},
		},
		{
			name: "it should return error when task repository Update return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &dto.TaskUpdateIn{TaskID: "task-xxxxx", UserID: "user-xxxxx"},
			},
			expected: expected{
				output: dto.TaskUpdateOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(entity.Task{UserID: "user-xxxxx"}, nil)

				d.taskRepository.On("Update", context.Background(), &entity.Task{UserID: "user-xxxxx"}).
					Return(entity.TaskID(""), test.ErrUnexpected)
			},
		},
		{
			name: "it should return error nil when success update",
			args: args{
				ctx: context.Background(),
				payload: &dto.TaskUpdateIn{
					TaskID:      "task-xxxxx",
					UserID:      "user-xxxxx",
					Content:     "new_content",
					Description: "new_description",
					IsDone:      true,
					DueDate:     entity.NullTime{NullTime: sql.NullTime{Time: test.TimeAfterNow, Valid: true}},
				},
			},
			expected: expected{
				output: dto.TaskUpdateOut{
					ID: "task-xxxxx",
				},
				err: nil,
			},
			setup: func(d *dependency) {
				d.taskRepository.On("FindByID", context.Background(), entity.TaskID("task-xxxxx")).
					Return(entity.Task{
						ID:          "task-xxxxx",
						UserID:      "user-xxxxx",
						Content:     "task_content",
						Description: "task_description",
						IsDone:      false,
						DueDate:     entity.NullTime{NullTime: sql.NullTime{Valid: false}},
						CreatedAt:   test.TimeBeforeNow,
						UpdatedAt:   test.TimeBeforeNow,
					}, nil)

				d.taskRepository.On("Update", context.Background(), &entity.Task{
					ID:          "task-xxxxx",
					UserID:      "user-xxxxx",
					Content:     "new_content",
					Description: "new_description",
					IsDone:      true,
					DueDate:     entity.NullTime{NullTime: sql.NullTime{Time: test.TimeAfterNow, Valid: true}},
					CreatedAt:   test.TimeBeforeNow,
					UpdatedAt:   test.TimeBeforeNow,
				}).Return(entity.TaskID("task-xxxxx"), nil)
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
			output, err := usecase.Update(t.args.ctx, t.args.payload)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.output, output)
		})
	}
}
