package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/pkg/response"
)

type HTTPAuthMiddlewareTestSuite struct {
	suite.Suite
}

func TestHTTPAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(HTTPAuthMiddlewareTestSuite))
}

func (s *HTTPAuthMiddlewareTestSuite) TestAuthentication() {
	s.Run("it should return error resposne when invalid authorization header", func() {
		mockJWTProvider := &mocks.JWTProvider{}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		authMiddleware := New(mockJWTProvider)
		handler := authMiddleware(mainHandler)

		handler.ServeHTTP(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(401, rr.Code)
		s.Equal(401, resBody.StatusCode)
		s.Equal(http.StatusText(http.StatusUnauthorized), resBody.Message)
		s.Equal("Authentication bearer token are not provided", resBody.Error)
	})

	s.Run("it should return error resposne fail to verify access token", func() {
		bearerToken := "Bearer xxxxx.xxxxx.xxxxx"
		rawToken := "xxxxx.xxxxx.xxxxx"
		mockJWTProvider := &mocks.JWTProvider{}
		mockJWTProvider.On("VerifyAccessToken", rawToken).Return(entity.UserID(""), errors.New("invalid access token"))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", bearerToken)

		mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		authMiddleware := New(mockJWTProvider)
		handler := authMiddleware(mainHandler)

		handler.ServeHTTP(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(401, rr.Code)
		s.Equal(401, resBody.StatusCode)
		s.Equal(http.StatusText(http.StatusUnauthorized), resBody.Message)
		s.Equal("The access token provided is invalid", resBody.Error)
	})

	s.Run("it should return error resposne fail to verify access token", func() {
		bearerToken := "Bearer xxxxx.xxxxx.xxxxx"
		rawToken := "xxxxx.xxxxx.xxxxx"
		mockJWTProvider := &mocks.JWTProvider{}
		mockJWTProvider.On("VerifyAccessToken", rawToken).Return(entity.UserID("xxxxx"), nil)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", bearerToken)

		mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			userID := r.Context().Value(entity.AuthUserIDKey("user_id")).(entity.UserID)
			w.Write([]byte(userID))
		})

		authMiddleware := New(mockJWTProvider)
		handler := authMiddleware(mainHandler)

		handler.ServeHTTP(rr, req)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(200, rr.Code)

		userIDInCtx := rr.Body.String()
		s.Equal("xxxxx", userIDInCtx)
	})

}
