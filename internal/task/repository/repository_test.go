package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
)

type TaskRepositoryTestSuite struct {
	suite.Suite
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}

func (s *TaskRepositoryTestSuite) TestCreate() {
	s.Run("it should return error when database fail", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}

		task := entity.Task{
			ID:          "xxxxx",
			UserID:      "yyyyyy",
			Content:     "task content",
			Description: "task description",
			DueDate:     &time.Time{},
		}

		mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO tasks (id, user_id, content, description, due_date)`)).
			WithArgs(task.ID, task.UserID, task.Content, task.Description, task.DueDate).
			WillReturnError(errors.New("database error"))

		mockIDProvider := &mocks.IDProvider{}
		mockIDProvider.On("Generate").Return(string(task.ID))

		repo := New(db, mockIDProvider)
		id, err := repo.Store(context.Background(), &task)

		s.Equal(errors.New("database error"), err)
		s.Empty(id)
	})

	s.Run("it should return task id when store task successfully", func() {
		db, mockDB, err := sqlmock.New()
		if err != nil {
			s.FailNow("an error '%s' was not expected when opening a stub database connection", err)
		}

		task := entity.Task{
			ID:          "xxxxx",
			UserID:      "yyyyyy",
			Content:     "task content",
			Description: "task description",
			DueDate:     &time.Time{},
		}

		mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO tasks (id, user_id, content, description, due_date)`)).
			WithArgs(task.ID, task.UserID, task.Content, task.Description, task.DueDate).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockIDProvider := &mocks.IDProvider{}
		mockIDProvider.On("Generate").Return(string(task.ID))

		repo := New(db, mockIDProvider)
		id, err := repo.Store(context.Background(), &task)

		s.NoError(err)
		s.Equal(task.ID, id)
	})
}
