package dto

import "github.com/edwintantawi/taskit/internal/domain/entity"

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
