package test

import (
	"context"
	"net/http"

	"github.com/edwintantawi/taskit/internal/domain/entity"
)

// InjectAuthContext injects the user ID into the request context.
func InjectAuthContext(r *http.Request, userID entity.UserID) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), entity.AuthUserIDKey, userID))
}
