package test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/edwintantawi/taskit/internal/domain/entity"
)

var (
	ErrUnexpected  = errors.New("test.unexpected")
	ErrRowAffected = errors.New("test.row_affected")
	ErrDatabase    = errors.New("test.database")
	ErrRowScan     = errors.New("test.rowscan")
	ErrRows        = errors.New("test.rows")
	ErrValidator   = errors.New("test.validator")

	TimeAfterNow  = time.Now().Add(time.Hour)
	TimeBeforeNow = time.Now().Add(-time.Hour)
)

// InjectAuthContext injects the user ID into the request context.
func InjectAuthContext(r *http.Request, userID entity.UserID) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), entity.AuthUserIDKey, userID))
}

// InjectChiRouterParams injects the chi router params into the request context.
func InjectChiRouterParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for key, value := range params {
		rctx.URLParams.Add(key, value)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// NewNullString create new null time.
func NewNullTime(id any) entity.NullTime {
	var isValid bool
	var timeValue time.Time
	if id != nil {
		timeValue = id.(time.Time)
		isValid = true
	}
	return entity.NullTime{
		NullTime: sql.NullTime{Time: timeValue, Valid: isValid},
	}
}

// NewNullString create new null string.
func NewNullString(id any) entity.NullString {
	var isValid bool
	var stringValue string
	if id != nil {
		stringValue = id.(string)
		isValid = true
	}
	return entity.NullString{
		NullString: sql.NullString{String: stringValue, Valid: isValid},
	}
}
