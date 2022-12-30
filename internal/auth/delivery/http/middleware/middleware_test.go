package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/pkg/response"
	"github.com/edwintantawi/taskit/test"
)

type HTTPAuthMiddlewareTestSuite struct {
	suite.Suite
}

func TestHTTPAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(HTTPAuthMiddlewareTestSuite))
}

type dependency struct {
	jwtProvider *mocks.JWTProvider
	req         *http.Request
}

func (s *HTTPAuthMiddlewareTestSuite) TestAuthentication() {
	type args struct {
		handler http.Handler
	}
	type expected struct {
		body        string
		message     string
		error       string
		statusCode  int
		contentType string
	}
	tests := []struct {
		name     string
		isError  bool
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name:    "it should response with error when authorization header is not provided",
			isError: true,
			args: args{
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			expected: expected{
				statusCode:  http.StatusUnauthorized,
				contentType: "application/json",
				message:     http.StatusText(http.StatusUnauthorized),
				error:       "Authentication bearer token are not provided",
			},
			setup: func(d *dependency) {
				d.req.Header.Del("Authorization")
			},
		},
		{
			name:    "it should return error when authorization header token is not valid",
			isError: true,
			args: args{
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			expected: expected{
				statusCode:  http.StatusUnauthorized,
				contentType: "application/json",
				message:     http.StatusText(http.StatusUnauthorized),
				error:       "The access token provided is invalid",
			},
			setup: func(d *dependency) {
				d.req.Header.Set("Authorization", "Bearer xxxxx.xxxxx.xxxxx")

				d.jwtProvider.On("VerifyAccessToken", "xxxxx.xxxxx.xxxxx").
					Return(entity.UserID(""), test.ErrUnexpected)
			},
		},
		{
			name:    "it should forward to next handler when authorization header is valid",
			isError: false,
			args: args{
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					userID := entity.GetAuthContext(r.Context())
					w.Write([]byte(userID))
				}),
			},
			expected: expected{
				statusCode: http.StatusOK,
				body:       "user-xxxxx",
			},
			setup: func(d *dependency) {
				d.req.Header.Set("Authorization", "Bearer xxxxx.xxxxx.xxxxx")

				d.jwtProvider.On("VerifyAccessToken", "xxxxx.xxxxx.xxxxx").
					Return(entity.UserID("user-xxxxx"), nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			req := httptest.NewRequest("GET", "/", nil)
			dep := &dependency{
				jwtProvider: &mocks.JWTProvider{},
				req:         req,
			}
			t.setup(dep)

			rr := httptest.NewRecorder()
			middleware := New(dep.jwtProvider)
			handler := middleware(t.args.handler)

			handler.ServeHTTP(rr, dep.req)

			if t.isError {
				var resBody response.E
				json.NewDecoder(rr.Body).Decode(&resBody)

				s.Equal(t.expected.contentType, rr.Header().Get("Content-Type"))
				s.Equal(t.expected.statusCode, rr.Code)
				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)
				s.Equal(t.expected.error, resBody.Error)
			} else {
				s.Equal(t.expected.statusCode, rr.Code)
				s.Equal(t.expected.body, rr.Body.String())
			}
		})
	}
}
