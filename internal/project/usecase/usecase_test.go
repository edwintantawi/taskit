package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/test"
)

type ProjectUsecaseTestSuite struct {
	suite.Suite
}

func TestProjectUsecaseSuite(t *testing.T) {
	suite.Run(t, new(ProjectUsecaseTestSuite))
}

type dependency struct {
	validator         *mocks.ValidatorProvider
	ProjectRepository *mocks.ProjectRepository
}

func (s *ProjectUsecaseTestSuite) TestCreate() {
	type args struct {
		ctx     context.Context
		payload *dto.ProjectCreateIn
	}
	type expected struct {
		output dto.ProjectCreateOut
		err    error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name: "it should return error when project repository Store return unexpected error",
			args: args{
				ctx:     context.Background(),
				payload: &dto.ProjectCreateIn{UserID: "user-xxxxx", Title: "project_title"},
			},
			expected: expected{
				output: dto.ProjectCreateOut{},
				err:    test.ErrUnexpected,
			},
			setup: func(d *dependency) {
				d.ProjectRepository.On("Store", context.Background(), &entity.Project{UserID: "user-xxxxx", Title: "project_title"}).
					Return(entity.ProjectID(""), test.ErrUnexpected)
			},
		},
		{
			name: "it should return error nil and project id when success",
			args: args{
				ctx:     context.Background(),
				payload: &dto.ProjectCreateIn{UserID: "user-xxxxx", Title: "project_title"},
			},
			expected: expected{
				output: dto.ProjectCreateOut{ID: "project-xxxxx"},
				err:    nil,
			},
			setup: func(d *dependency) {
				d.ProjectRepository.On("Store", context.Background(), &entity.Project{UserID: "user-xxxxx", Title: "project_title"}).
					Return(entity.ProjectID("project-xxxxx"), nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			d := &dependency{
				validator:         &mocks.ValidatorProvider{},
				ProjectRepository: &mocks.ProjectRepository{},
			}
			t.setup(d)

			usecase := New(d.ProjectRepository)
			output, err := usecase.Create(t.args.ctx, t.args.payload)

			s.Equal(t.expected.err, err)
			s.Equal(t.expected.output, output)
		})
	}
}
