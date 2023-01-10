package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/pkg/errorx"
	"github.com/edwintantawi/taskit/pkg/response"
)

type Middleware struct {
	jwtProvider domain.JWTProvider
}

// New creates a new HTTP auth middleware.
func New(jwtProvider domain.JWTProvider) Middleware {
	return Middleware{jwtProvider: jwtProvider}
}

// Authenticate authenticates the request.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)

		bearerToken := r.Header.Get("Authorization")
		if !strings.Contains(bearerToken, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			encoder.Encode(response.Error(http.StatusUnauthorized, "Authentication bearer token are not provided"))
			return
		}

		rawToken := strings.TrimPrefix(bearerToken, "Bearer ")
		userId, err := m.jwtProvider.VerifyAccessToken(rawToken)
		if err != nil {
			code, msg := errorx.HTTPErrorTranslator(err)
			w.WriteHeader(code)
			encoder.Encode(response.Error(code, msg))
			return
		}

		ctx := context.WithValue(r.Context(), entity.AuthUserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
