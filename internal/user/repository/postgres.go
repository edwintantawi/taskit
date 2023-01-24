package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/edwintantawi/taskit/internal/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// Postgres represents a User postgres repository.
type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) Postgres {
	return Postgres{db}
}

// Save saves a new user to the database.
func (p Postgres) Save(ctx context.Context, newUser entity.NewUser) (entity.AddedUser, error) {
	id := uuid.NewString()
	createdTime := time.Now().UTC()

	q := `INSERT INTO users (id, name, email, password, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $5)`

	_, err := p.db.ExecContext(ctx, q, id, newUser.Name, newUser.Email, newUser.Password, createdTime)
	if err != nil {
		return entity.AddedUser{}, err
	}

	return entity.AddedUser{ID: id, Email: newUser.Email}, nil
}

// FindByEmail find user by email.
func (p Postgres) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	q := `SELECT id, name, email, password, created_at, updated_at
			FROM users WHERE email = $1`

	row := p.db.QueryRowContext(ctx, q, email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.User{}, ErrUserNotFound
	} else if err != nil {
		return entity.User{}, err
	}

	user.CreatedAt = user.CreatedAt.UTC()
	user.UpdatedAt = user.UpdatedAt.UTC()

	return user, nil
}
