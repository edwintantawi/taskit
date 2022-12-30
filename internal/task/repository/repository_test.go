package repository

import (
	"context"
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
