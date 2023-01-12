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

type ProjectRepositoryTestSuite struct {
	suite.Suite
}

func TestProjectRepositorySuite(t *testing.T) {
	suite.Run(t, new(ProjectRepositoryTestSuite))
}

type dependency struct {
	mockDB     sqlmock.Sqlmock
	idProvider *mocks.IDProvider
}

func (s *ProjectRepositoryTestSuite) TestCreate() {
	type args struct {
		ctx     context.Context
		project *entity.Project
	}
	type expected struct {
		projectID entity.ProjectID
		err       error
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
				ctx:     context.Background(),
				project: &entity.Project{UserID: "user-xxxxx", Title: "project_title"},
			},
			expected: expected{
				projectID: "",
				err:       test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.idProvider.On("Generate").Return("project-xxxxx")
				d.mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO projects (id, user_id, title) VALUES ($1, $2, $3)`)).
					WithArgs("project-xxxxx", "user-xxxxx", "project_title").
					WillReturnError(test.ErrDatabase)
			},
		},
		{
			name: "it should return error nil and user id when successfully store",
			args: args{
				ctx:     context.Background(),
				project: &entity.Project{UserID: "user-xxxxx", Title: "project_title"},
			},
			expected: expected{
				projectID: "project-xxxxx",
				err:       nil,
			},
			setup: func(d *dependency) {
				d.idProvider.On("Generate").Return("project-xxxxx")
				d.mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO projects (id, user_id, title) VALUES ($1, $2, $3)`)).
					WithArgs("project-xxxxx", "user-xxxxx", "project_title").
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
			userID, err := repository.Store(t.args.ctx, t.args.project)

			s.Equal(t.expected.projectID, userID)
			s.Equal(t.expected.err, err)
		})
	}
}
