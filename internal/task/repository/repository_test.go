package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
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

			deps := &dependency{
				mockDB:     mockDB,
				idProvider: &mocks.IDProvider{},
			}
			t.setup(deps)

			repository := New(db, deps.idProvider)
			taskID, err := repository.Store(t.args.ctx, t.args.task)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.taskID, taskID)
		})
	}
}

func (s *TaskRepositoryTestSuite) TestFindByID() {
	type args struct {
		ctx    context.Context
		taskID entity.TaskID
	}
	type expected struct {
		task entity.Task
		err  error
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
				taskID: "task-xxxxx",
			},
			expected: expected{
				task: entity.Task{},
				err:  test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, content, description, is_completed, due_date, created_at, updated_at FROM tasks WHERE id = $1")).
					WithArgs("task-xxxxx").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error ErrTaskNotFound when task not found",
			args: args{
				ctx:    context.Background(),
				taskID: "task-xxxxx",
			},
			expected: expected{
				task: entity.Task{},
				err:  domain.ErrTaskNotFound,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, content, description, is_completed, due_date, created_at, updated_at FROM tasks WHERE id = $1")).
					WithArgs("task-xxxxx").
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "it should return error when database scan fail",
			args: args{
				ctx:    context.Background(),
				taskID: "task-xxxxx",
			},
			expected: expected{
				task: entity.Task{},
				err:  test.ErrRowScan,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, content, description, is_completed, due_date, created_at, updated_at FROM tasks WHERE id = $1")).
					WithArgs("task-xxxxx").
					WillReturnError(test.ErrRowScan)
			},
		},
		{
			name: "it should return error nil and task when success",
			args: args{
				ctx:    context.Background(),
				taskID: "task-xxxxx",
			},
			expected: expected{
				task: entity.Task{
					ID:          "task-xxxxx",
					UserID:      "user-xxxxx",
					Content:     "task_content",
					Description: "task_description",
					IsCompleted: true,
					DueDate:     &test.TimeAfterNow,
					CreatedAt:   test.TimeBeforeNow,
					UpdatedAt:   test.TimeBeforeNow,
				},
				err: nil,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "user_id", "content", "description", "is_completed", "due_date", "created_at", "updated_at"}).
					AddRow("task-xxxxx", "user-xxxxx", "task_content", "task_description", true, test.TimeAfterNow, test.TimeBeforeNow, test.TimeBeforeNow)

				d.mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, content, description, is_completed, due_date, created_at, updated_at FROM tasks WHERE id = $1")).
					WithArgs("task-xxxxx").
					WillReturnRows(mockRow)
			},
		},
	}

	for _, t := range tests {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a database mock connection", err)
		}

		deps := &dependency{
			mockDB: mockDB,
		}
		t.setup(deps)

		repository := New(db, nil)
		task, err := repository.FindByID(t.args.ctx, t.args.taskID)

		s.Equal(t.expected.err, err)
		s.Equal(t.expected.task, task)
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
					{ID: "task-xxxxx", Content: "task_xxxxx_content", Description: "task_xxxxx_description", IsCompleted: false, DueDate: nil, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
					{ID: "task-yyyyy", Content: "task_yyyyy_content", Description: "task_yyyyy_description", IsCompleted: true, DueDate: &test.TimeAfterNow, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
				},
				err: nil,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "content", "description", "is_completed", "due_date", "created_at", "updated_at"}).
					AddRow("task-xxxxx", "task_xxxxx_content", "task_xxxxx_description", false, nil, test.TimeBeforeNow, test.TimeBeforeNow).
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

			deps := &dependency{
				mockDB: mockDB,
			}
			t.setup(deps)

			repository := New(db, deps.idProvider)
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

func (s *TaskRepositoryTestSuite) TestVerifyAvailableByID() {
	type args struct {
		ctx    context.Context
		taskID entity.TaskID
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
			name: "it should return error when database fail to query",
			args: args{
				ctx:    context.Background(),
				taskID: "task-xxxxx",
			},
			expected: expected{
				err: test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM tasks WHERE id = $1`)).
					WithArgs("task-xxxxx").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error ErrTaskNotFound when task not found",
			args: args{
				ctx:    context.Background(),
				taskID: "task-xxxxx",
			},
			expected: expected{
				err: domain.ErrTaskNotFound,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM tasks WHERE id = $1`)).
					WithArgs("task-xxxxx").
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "it should return error when fail to scan row",
			args: args{
				ctx:    context.Background(),
				taskID: "task-xxxxx",
			},
			expected: expected{
				err: test.ErrRowScan,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM tasks WHERE id = $1`)).
					WithArgs("task-xxxxx").
					WillReturnError(test.ErrRowScan)
			},
		},
		{
			name: "it should return error nil when task found",
			args: args{
				ctx:    context.Background(),
				taskID: "task-xxxxx",
			},
			expected: expected{
				err: nil,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id"}).AddRow("task-xxxxx")
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM tasks WHERE id = $1`)).
					WithArgs("task-xxxxx").
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

			deps := &dependency{
				mockDB: mockDB,
			}
			t.setup(deps)

			repository := New(db, nil)
			err = repository.VerifyAvailableByID(t.args.ctx, t.args.taskID)

			s.Equal(t.expected.err, err)
		})
	}
}

func (s *TaskRepositoryTestSuite) TestDeleteByID() {
	type args struct {
		ctx    context.Context
		taskID entity.TaskID
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
			name: "it should return error when database fail to query",
			args: args{
				ctx:    context.Background(),
				taskID: "task-xxxxx",
			},
			expected: expected{
				err: test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM tasks WHERE id = $1`)).
					WithArgs("task-xxxxx").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return nil when success to delete",
			args: args{
				ctx:    context.Background(),
				taskID: "task-xxxxx",
			},
			expected: expected{
				err: nil,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM tasks WHERE id = $1`)).
					WithArgs("task-xxxxx").
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

			deps := &dependency{
				mockDB: mockDB,
			}
			t.setup(deps)

			repository := New(db, nil)
			err = repository.DeleteByID(t.args.ctx, t.args.taskID)

			s.Equal(t.expected.err, err)
		})
	}
}

func (s *TaskRepositoryTestSuite) TestUpdate() {
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
			name: "it should return error when database fail",
			args: args{
				ctx:  context.Background(),
				task: &entity.Task{},
			},
			expected: expected{
				taskID: "",
				err:    test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectExec(regexp.QuoteMeta("UPDATE tasks SET content = $2, description = $3, is_completed = $4, due_date = $5, updated_at = $6 WHERE id = $1")).
					WithArgs("", "", "", false, nil, sqlmock.AnyArg()).
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error nil and task id when success update",
			args: args{
				ctx: context.Background(),
				task: &entity.Task{
					ID:          "task-xxxxx",
					UserID:      "user-xxxxx",
					Content:     "task_content",
					Description: "task_description",
					IsCompleted: true,
					DueDate:     &test.TimeAfterNow,
					CreatedAt:   test.TimeBeforeNow,
					UpdatedAt:   test.TimeBeforeNow,
				},
			},
			expected: expected{
				taskID: "task-xxxxx",
				err:    nil,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectExec(regexp.QuoteMeta("UPDATE tasks SET content = $2, description = $3, is_completed = $4, due_date = $5, updated_at = $6 WHERE id = $1")).
					WithArgs("task-xxxxx", "task_content", "task_description", true, test.TimeAfterNow, sqlmock.AnyArg()).
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

			deps := &dependency{
				mockDB: mockDB,
			}
			t.setup(deps)

			repository := New(db, nil)
			taskID, err := repository.Update(t.args.ctx, t.args.task)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.taskID, taskID)
		})
	}
}
