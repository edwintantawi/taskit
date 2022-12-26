package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/pkg/response"
)

const authUserIDKey = entity.AuthUserIDKey("user_id")

type middleware func(next http.Handler) http.Handler

func New(jwtProvider domain.JWTProvider) middleware {
	return func(next http.Handler) http.Handler {
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
			userId, err := jwtProvider.VerifyAccessToken(rawToken)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				encoder.Encode(response.Error(http.StatusUnauthorized, "The access token provided is invalid"))
				return
			}

			ctx := context.WithValue(r.Context(), authUserIDKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}