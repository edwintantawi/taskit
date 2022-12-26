package repository

import (
	"context"
	"database/sql"

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
