package dto

import (
	"time"

	"github.com/edwintantawi/taskit/internal/domain/entity"
)

// ProjectCreateIn represent create project input.
type ProjectCreateIn struct {
	UserID entity.UserID `json:"-"`
	Title  string        `json:"title"`
}

func (d *ProjectCreateIn) Validate() error {
	switch {
	case d.Title == "":
		return ErrTitleEmpty
	}
	return nil
}

// ProjectCreateIn represent create project input.
type ProjectCreateOut struct {
	ID entity.ProjectID `json:"id"`
}

// ProjectGetAllIn represent get all project input.
type ProjectGetAllIn struct {
	UserID entity.UserID `json:"-"`
}

// ProjectGetAllOut represent get all project output.
type ProjectGetAllOut struct {
	ID        entity.ProjectID `json:"id"`
	Title     string           `json:"title"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// ProjectGetByIDIn represent get by id project input.
type ProjectGetByIDIn struct {
	ID entity.ProjectID `json:"-"`
}

// ProjectGetByIDOut represent get by id project output.
type ProjectGetByIDOut struct {
	ID        entity.ProjectID `json:"id"`
	Title     string           `json:"title"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
