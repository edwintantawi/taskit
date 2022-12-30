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

// New create a new auth repository.
func New(db *sql.DB, idProvider domain.IDProvider) domain.AuthRepository {
	return &repository{db: db, idProvider: idProvider}
}

// Store save a new auth to database.
func (r *repository) Store(ctx context.Context, a *entity.Auth) error {
	id := r.idProvider.Generate()
	q := `INSERT INTO authentications (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, q, id, a.UserID, a.Token, a.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

// VerifyAvailableByToken check if a authentication is available by token.
func (r *repository) VerifyAvailableByToken(ctx context.Context, token string) error {
	var id entity.AuthID
	q := `SELECT id FROM authentications WHERE token = $1`
	row := r.db.QueryRowContext(ctx, q, token)
	err := row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrAuthNotFound
	} else if err != nil {
		return err
	}
	return nil
}

// Delete remove an auth from database.
func (r *repository) Delete(ctx context.Context, a *entity.Auth) error {
	q := `DELETE FROM authentications WHERE token = $1`
	result, err := r.db.ExecContext(ctx, q, a.Token)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return domain.ErrAuthNotFound
	}
	return nil
}

// FindByToken find an auth by token.
func (r *repository) FindByToken(ctx context.Context, token string) (entity.Auth, error) {
	var a entity.Auth
	q := `SELECT id, user_id, token, expires_at FROM authentications WHERE token = $1`
	row := r.db.QueryRowContext(ctx, q, token)
	err := row.Scan(&a.ID, &a.UserID, &a.Token, &a.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return a, domain.ErrAuthNotFound
	} else if err != nil {
		return a, err
	}
	return a, nil
}
