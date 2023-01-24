package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/edwintantawi/taskit/internal/entity"
)

// Postgres represents a User postgres repository.
type Postgres struct {
	db *sql.DB
}

// NewPostgres creates a new User postgres repository.
func NewPostgres(db *sql.DB) Postgres {
	return Postgres{db}
}

// Save saves a new user to the database.
func (p *Postgres) Save(ctx context.Context, newUser entity.NewUser) (entity.AddedUser, error) {
	id := uuid.NewString()
	createdTime := time.Now()

	q := `INSERT INTO users (id, name, email, password, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $5)`

	_, err := p.db.ExecContext(ctx, q, id, newUser.Name, newUser.Email, newUser.Password, createdTime)
	if err != nil {
		return entity.AddedUser{}, err
	}

	return entity.AddedUser{ID: id, Email: newUser.Email}, nil
}
