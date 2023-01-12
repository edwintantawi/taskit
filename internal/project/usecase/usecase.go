package usecase

import (
	"context"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type Usecase struct {
	projectRepository domain.ProjectRepository
}

// New create a new project usecase.
func New(projectRepository domain.ProjectRepository) Usecase {
	return Usecase{projectRepository: projectRepository}
}

// Create create a new project
func (u *Usecase) Create(ctx context.Context, payload *dto.ProjectCreateIn) (dto.ProjectCreateOut, error) {
	project := &entity.Project{UserID: payload.UserID, Title: payload.Title}
	id, err := u.projectRepository.Store(ctx, project)
	if err != nil {
		return dto.ProjectCreateOut{}, err
	}
	return dto.ProjectCreateOut{ID: id}, nil
}
