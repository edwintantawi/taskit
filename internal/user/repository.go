package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
)

var (
	ErrEmailNotAvailable = errors.New("email_not_available")
)

type repository struct {
	db         *sql.DB
	idProvider domain.IDProvider
}

func NewRepository(db *sql.DB, idProvider domain.IDProvider) domain.UserRepository {
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
		return ErrEmailNotAvailable
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}
