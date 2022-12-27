package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

type repository struct {
	db         *sql.DB
	idProvider domain.IDProvider
}

// New create a new user repository.
func New(db *sql.DB, idProvider domain.IDProvider) domain.UserRepository {
	return &repository{db: db, idProvider: idProvider}
}

// Store save a new user to database.
func (r *repository) Store(ctx context.Context, u *entity.User) (entity.UserID, error) {
	id := entity.UserID(r.idProvider.Generate())
	q := `INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, q, id, u.Name, u.Email, u.Password)
	if err != nil {
		return "", err
	}
	return id, nil
}

// VerifyAvailableEmail check if the email is available.
func (r *repository) VerifyAvailableEmail(ctx context.Context, email string) error {
	var id entity.UserID
	q := `SELECT id FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, q, email).Scan(&id)
	if err == nil {
		return domain.ErrEmailNotAvailable
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

// FindByEmail find a user by email.
func (r *repository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	var u entity.User
	q := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, q, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return u, domain.ErrUserEmailNotExist
	}
	if err != nil {
		return u, err
	}
	return u, nil
}

// FindByID find a user by id.
func (r *repository) FindByID(ctx context.Context, id entity.UserID) (entity.User, error) {
	var u entity.User
	q := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return u, domain.ErrUserIDNotExist
	}
	if err != nil {
		return u, err
	}
	return u, nil
}
