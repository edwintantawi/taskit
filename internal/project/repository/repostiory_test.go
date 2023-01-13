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

func (s *ProjectRepositoryTestSuite) TestFindAllByUserID() {
	type args struct {
		ctx    context.Context
		userID entity.UserID
	}
	type expected struct {
		projects      []entity.Project
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
				projects: nil,
				err:      test.ErrDatabase,
			},
			setup: func(d *dependency) {
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, title, created_at, updated_at FROM projects WHERE user_id = $1`)).
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
				projects:      nil,
				allowAnyError: true,
				err:           errors.New("anything"),
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "user_id", "title", "created_at", "updated_at"}).
					AddRow(nil, "user-xxxxx", "project_title", test.TimeBeforeNow, test.TimeBeforeNow)

				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, title, created_at, updated_at FROM projects WHERE user_id = $1`)).
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
				projects: nil,
				err:      test.ErrRows,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "user_id", "title", "created_at", "updated_at"}).
					AddRow("project-xxxxx", "user-xxxxx", "project_title_x", test.TimeBeforeNow, test.TimeBeforeNow).
					AddRow("project-yyyyy", "user-xxxxx", "project_title_y", test.TimeBeforeNow, test.TimeBeforeNow).
					RowError(1, test.ErrRows)

				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, title, created_at, updated_at FROM projects WHERE user_id = $1`)).
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
				projects: []entity.Project{},
				err:      nil,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "user_id", "title", "created_at", "updated_at"})
				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, title, created_at, updated_at FROM projects WHERE user_id = $1`)).
					WithArgs("user-xxxxx").
					WillReturnRows(mockRow)
			},
		},
		{
			name: "it should return error nil and all project when successfully query",
			args: args{
				ctx:    context.Background(),
				userID: "user-xxxxx",
			},
			expected: expected{
				projects: []entity.Project{
					{
						ID:        "project-xxxxx",
						UserID:    "user-xxxxx",
						Title:     "project_title_x",
						CreatedAt: test.TimeBeforeNow,
						UpdatedAt: test.TimeBeforeNow,
					},
					{
						ID:        "project-yyyyy",
						UserID:    "user-xxxxx",
						Title:     "project_title_y",
						CreatedAt: test.TimeBeforeNow,
						UpdatedAt: test.TimeBeforeNow,
					},
				},
				err: nil,
			},
			setup: func(d *dependency) {
				mockRow := sqlmock.NewRows([]string{"id", "user_id", "title", "created_at", "updated_at"}).
					AddRow("project-xxxxx", "user-xxxxx", "project_title_x", test.TimeBeforeNow, test.TimeBeforeNow).
					AddRow("project-yyyyy", "user-xxxxx", "project_title_y", test.TimeBeforeNow, test.TimeBeforeNow)

				d.mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, title, created_at, updated_at FROM projects WHERE user_id = $1`)).
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
			s.Equal(t.expected.projects, tasks)
		})
	}
}
