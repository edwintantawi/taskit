package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/test"
)

type TaskRepositoryTestSuite struct {
	suite.Suite
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}

type dependency struct {
	mockDB     sqlmock.Sqlmock
	idProvider *mocks.IDProvider
}

func (s *TaskRepositoryTestSuite) TestStore() {
	type args struct {
		ctx  context.Context
		task *entity.Task
	}
	type expected struct {
		taskID entity.TaskID
		err    error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when database fail to store",
			args: args{
				ctx:  context.Background(),
				task: &entity.Task{UserID: "user-xxxxx", Content: "task_content", Description: "task_description", DueDate: &test.TimeAfterNow},
			},
			expected: expected{
				taskID: "",
				err:    test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.idProvider.On("Generate").Return("task-xxxxx")
				d.mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO tasks (id, user_id, content, description, due_date)`)).
					WithArgs("task-xxxxx", "user-xxxxx", "task_content", "task_description", &test.TimeAfterNow).
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error nil and task id when successfully store",
			args: args{
				ctx:  context.Background(),
				task: &entity.Task{UserID: "user-xxxxx", Content: "task_content", Description: "task_description", DueDate: &test.TimeAfterNow},
			},
			expected: expected{
				taskID: "task-xxxxx",
				err:    nil,
			},
			setup: func(d *dependency) {
				d.idProvider.On("Generate").Return("task-xxxxx")
				d.mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO tasks (id, user_id, content, description, due_date)`)).
					WithArgs("task-xxxxx", "user-xxxxx", "task_content", "task_description", &test.TimeAfterNow).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				s.FailNow("an error '%s' was not expected when opening a database mock connection", err)
			}

			d := &dependency{
				mockDB:     mockDB,
				idProvider: &mocks.IDProvider{},
			}
			t.setup(d)

			repository := New(db, d.idProvider)
			taskID, err := repository.Store(t.args.ctx, t.args.task)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.taskID, taskID)
		})
	}
}

func (s *TaskRepositoryTestSuite) TestFindAllByUserID() {
	type args struct {
		ctx    context.Context
		userID entity.UserID
	}
	type expected struct {
		tasks         []entity.Task
		allowAnyError bool
		err           error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when database fail to query",
			args: args{
				ctx:    context.Background(),
				userID: "user-xxxxx",
			},
			expected: expected{
				tasks: nil,
				err:   test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, content, description, is_completed, due_date, created_at, updated_at FROM tasks WHERE user_id = $1`)).
					WithArgs("user-xxxxx").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error when database rows fail to scan",
			args: args{
				ctx:    context.Background(),
				userID: "user-xxxxx",
			},
			expected: expected{
				tasks:         nil,
				allowAnyError: true,
				err:           errors.New("anything"),
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "content", "description", "is_completed", "due_date", "created_at", "updated_at"}).
					AddRow(nil, "task_xxxxx_content", "task_yyyyy_description", false, nil, test.TimeBeforeNow, test.TimeBeforeNow)

				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, content, description, is_completed, due_date, created_at, updated_at FROM tasks WHERE user_id = $1`)).
					WithArgs("user-xxxxx").
					WillReturnRows(mockRow)
			},
		},
		{
			name: "it should return error when database rows error",
			args: args{
				ctx:    context.Background(),
				userID: "user-xxxxx",
			},
			expected: expected{
				tasks: nil,
				err:   test.ErrRows,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "content", "description", "is_completed", "due_date", "created_at", "updated_at"}).
					AddRow("task-xxxxx", "task_xxxxx_content", "task_yyyyy_description", false, nil, test.TimeBeforeNow, test.TimeBeforeNow).
					AddRow("task-yyyyy", "task_yyyyy_content", "task_yyyyy_description", true, test.TimeAfterNow, test.TimeBeforeNow, test.TimeBeforeNow).
					RowError(1, test.ErrRows)

				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, content, description, is_completed, due_date, created_at, updated_at FROM tasks WHERE user_id = $1`)).
					WithArgs("user-xxxxx").
					WillReturnRows(mockRow)
			},
		},
		{
			name: "it should return error nil and empty slice task when successfully query with no tasks",
			args: args{
				ctx:    context.Background(),
				userID: "user-xxxxx",
			},
			expected: expected{
				tasks: []entity.Task{},
				err:   nil,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "content", "description", "is_completed", "due_date", "created_at", "updated_at"})
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, content, description, is_completed, due_date, created_at, updated_at FROM tasks WHERE user_id = $1`)).
					WithArgs("user-xxxxx").
					WillReturnRows(mockRow)
			},
		},
		{
			name: "it should return error nil and all task when successfully query",
			args: args{
				ctx:    context.Background(),
				userID: "user-xxxxx",
			},
			expected: expected{
				tasks: []entity.Task{
					{ID: "task-xxxxx", Content: "task_xxxxx_content", Description: "task_yyyyy_description", IsCompleted: false, DueDate: nil, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
					{ID: "task-yyyyy", Content: "task_yyyyy_content", Description: "task_yyyyy_description", IsCompleted: true, DueDate: &test.TimeAfterNow, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
				},
				err: nil,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "content", "description", "is_completed", "due_date", "created_at", "updated_at"}).
					AddRow("task-xxxxx", "task_xxxxx_content", "task_yyyyy_description", false, nil, test.TimeBeforeNow, test.TimeBeforeNow).
					AddRow("task-yyyyy", "task_yyyyy_content", "task_yyyyy_description", true, test.TimeAfterNow, test.TimeBeforeNow, test.TimeBeforeNow)

				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, content, description, is_completed, due_date, created_at, updated_at FROM tasks WHERE user_id = $1`)).
					WithArgs("user-xxxxx").
					WillReturnRows(mockRow)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			db, mockDB, err := sqlmock.New()
			if err != nil {
				s.FailNow("an error '%s' was not expected when opening a database mock connection", err)
			}

			d := &dependency{
				mockDB: mockDB,
			}
			t.setup(d)

			repository := New(db, d.idProvider)
			tasks, err := repository.FindAllByUserID(t.args.ctx, t.args.userID)

			if t.expected.allowAnyError {
				s.Error(err)
			} else {
				s.Equal(t.expected.err, err)
			}
			s.Equal(t.expected.tasks, tasks)
		})
	}
}
