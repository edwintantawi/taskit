package test

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/edwintantawi/taskit/internal/domain/entity"
)

var (
	ErrUnexpected  = errors.New("test.unexpected")
	ErrRowAffected = errors.New("test.row_affected")
	ErrDatabase    = errors.New("test.database")
	ErrRowScan     = errors.New("test.rowscan")
	ErrRows        = errors.New("test.rows")

	TimeAfterNow  = time.Now().Add(time.Hour)
	TimeBeforeNow = time.Now().Add(-time.Hour)
)

// InjectAuthContext injects the user ID into the request context.
func InjectAuthContext(r *http.Request, userID entity.UserID) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), entity.AuthUserIDKey, userID))
}
