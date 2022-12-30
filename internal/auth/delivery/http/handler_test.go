package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/pkg/errorx"
	"github.com/edwintantawi/taskit/pkg/response"
	"github.com/edwintantawi/taskit/test"
)

type AuthHTTPHandlerTestSuite struct {
	suite.Suite
}

func TestAuthHTTPHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHTTPHandlerTestSuite))
}

type dependency struct {
	req         *http.Request
	authUsecase *mocks.AuthUsecase
}

func (s *AuthHTTPHandlerTestSuite) TestPost() {
	type args struct {
		requestBody []byte
	}
	type expected struct {
		contentType string
		statusCode  int
		message     string
		error       string
		payload     map[string]any
	}
	tests := []struct {
		name     string
		isError  bool
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name:    "it should response with error when request body is invalid or not provided",
			isError: true,
			args: args{
				requestBody: []byte(`{`),
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
				message:     http.StatusText(http.StatusBadRequest),
				error:       "Invalid request body",
			},
			setup: func(d *dependency) {},
		},
		{
			name:    "it should response with error when auth usecase Login return unexpected error",
			isError: true,
			args: args{
				requestBody: []byte(`{}`),
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusInternalServerError,
				message:     http.StatusText(http.StatusInternalServerError),
				error:       errorx.InternalServerErrorMessage,
			},
			setup: func(d *dependency) {
				d.authUsecase.On("Login", mock.Anything, &domain.LoginAuthIn{}).
					Return(domain.LoginAuthOut{}, test.ErrUnexpected)
			},
		},
		{
			name: "it should response with success when success",
			args: args{
				requestBody: []byte(`{}`),
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				message:     "Successfully logged in user",
				payload: map[string]any{
					"access_token":  "xxxxx.xxxxx.xxxxx",
					"refresh_token": "yyyyy.yyyyy.yyyyy",
				},
			},
			setup: func(d *dependency) {
				d.authUsecase.On("Login", mock.Anything, &domain.LoginAuthIn{}).
					Return(domain.LoginAuthOut{AccessToken: "xxxxx.xxxxx.xxxxx", RefreshToken: "yyyyy.yyyyy.yyyyy"}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			reqBody := bytes.NewReader(t.args.requestBody)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/", reqBody)

			deps := &dependency{
				authUsecase: &mocks.AuthUsecase{},
			}
			t.setup(deps)

			handler := New(deps.authUsecase)
			handler.Post(rr, req)

			s.Equal(t.expected.contentType, rr.Header().Get("Content-Type"))
			s.Equal(t.expected.statusCode, rr.Code)

			if t.isError {
				var resBody response.E
				json.NewDecoder(rr.Body).Decode(&resBody)

				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)
				s.Equal(t.expected.error, resBody.Error)
			} else {
				var resBody response.S
				json.NewDecoder(rr.Body).Decode(&resBody)
				payloadMap := resBody.Payload.(map[string]any)

				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)
				s.Equal(t.expected.payload, payloadMap)
			}
		})
	}
}

func (s *AuthHTTPHandlerTestSuite) TestDelete() {
	type args struct {
		requestBody []byte
	}
	type expected struct {
		contentType string
		statusCode  int
		message     string
		error       string
		payload     map[string]any
	}
	tests := []struct {
		name     string
		isError  bool
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name:    "it should response with error when request body is invalid or not provided",
			isError: true,
			args: args{
				requestBody: []byte(`{`),
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
				message:     http.StatusText(http.StatusBadRequest),
				error:       "Invalid request body",
			},
			setup: func(d *dependency) {},
		},
		{
			name:    "it should response with error when auth usecase Logout return unexpected error",
			isError: true,
			args: args{
				requestBody: []byte(`{}`),
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusInternalServerError,
				message:     http.StatusText(http.StatusInternalServerError),
				error:       errorx.InternalServerErrorMessage,
			},
			setup: func(d *dependency) {
				d.authUsecase.On("Logout", mock.Anything, &domain.LogoutAuthIn{}).
					Return(test.ErrUnexpected)
			},
		},
		{
			name:    "it should response with success when success",
			isError: false,
			args: args{
				requestBody: []byte(`{}`),
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				message:     "Successfully logout user",
				payload:     nil,
			},
			setup: func(d *dependency) {
				d.authUsecase.On("Logout", mock.Anything, &domain.LogoutAuthIn{}).
					Return(nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			reqBody := bytes.NewReader(t.args.requestBody)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/", reqBody)

			deps := &dependency{
				authUsecase: &mocks.AuthUsecase{},
			}
			t.setup(deps)

			handler := New(deps.authUsecase)
			handler.Delete(rr, req)

			s.Equal(t.expected.contentType, rr.Header().Get("Content-Type"))
			s.Equal(t.expected.statusCode, rr.Code)

			if t.isError {
				var resBody response.E
				json.NewDecoder(rr.Body).Decode(&resBody)

				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)
				s.Equal(t.expected.error, resBody.Error)
			} else {
				var resBody response.S
				json.NewDecoder(rr.Body).Decode(&resBody)

				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)
				s.Nil(resBody.Payload)
			}
		})
	}
}

func (s *AuthHTTPHandlerTestSuite) TestGet() {
	type expected struct {
		contentType string
		statusCode  int
		message     string
		error       string
		payload     map[string]any
	}
	tests := []struct {
		name     string
		isError  bool
		expected expected
		setup    func(d *dependency)
	}{
		{
			name:    "it should response with error when auth usecase GetProfile return unexpected error",
			isError: true,
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusInternalServerError,
				message:     http.StatusText(http.StatusInternalServerError),
				error:       errorx.InternalServerErrorMessage,
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.authUsecase.On("GetProfile", mock.Anything, &domain.GetProfileAuthIn{UserID: entity.UserID("user-xxxxx")}).
					Return(domain.GetProfileAuthOut{}, test.ErrUnexpected)
			},
		},
		{
			name:    "it should response with success when success",
			isError: false,
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				message:     http.StatusText(http.StatusOK),
				payload: map[string]any{
					"id":    "user-xxxxx",
					"name":  "Gopher",
					"email": "gopher@go.dev",
				},
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.authUsecase.On("GetProfile", mock.Anything, &domain.GetProfileAuthIn{UserID: entity.UserID("user-xxxxx")}).
					Return(domain.GetProfileAuthOut{ID: "user-xxxxx", Name: "Gopher", Email: "gopher@go.dev"}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			deps := &dependency{
				authUsecase: &mocks.AuthUsecase{},
				req:         req,
			}
			t.setup(deps)

			handler := New(deps.authUsecase)
			handler.Get(rr, deps.req)

			s.Equal(t.expected.contentType, rr.Header().Get("Content-Type"))
			s.Equal(t.expected.statusCode, rr.Code)

			if t.isError {
				var resBody response.E
				json.NewDecoder(rr.Body).Decode(&resBody)

				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)
				s.Equal(t.expected.error, resBody.Error)
			} else {
				var resBody response.S
				json.NewDecoder(rr.Body).Decode(&resBody)
				payloadMap := resBody.Payload.(map[string]any)

				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)
				s.Equal(t.expected.payload, payloadMap)
			}
		})
	}
}

func (s *AuthHTTPHandlerTestSuite) TestPut() {
	type args struct {
		requestBody []byte
	}
	type expected struct {
		contentType string
		statusCode  int
		message     string
		error       string
		payload     map[string]any
	}
	tests := []struct {
		name     string
		isError  bool
		args     args
		expected expected
		setup    func(d *dependency)
	}{
		{
			name:    "it should response with error when request body is invalid or not provided",
			isError: true,
			args: args{
				requestBody: []byte(`{`),
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
				message:     http.StatusText(http.StatusBadRequest),
				error:       "Invalid request body",
			},
			setup: func(d *dependency) {},
		},
		{
			name:    "it should response with error when auth usecase Refresh return unexpected error",
			isError: true,
			args: args{
				requestBody: []byte(`{}`),
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusInternalServerError,
				message:     http.StatusText(http.StatusInternalServerError),
				error:       errorx.InternalServerErrorMessage,
			},
			setup: func(d *dependency) {
				d.authUsecase.On("Refresh", mock.Anything, &domain.RefreshAuthIn{}).
					Return(domain.RefreshAuthOut{}, test.ErrUnexpected)
			},
		},
		{
			name:    "it should response with success when success",
			isError: false,
			args: args{
				requestBody: []byte(`{}`),
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				message:     "Successfully refreshed authentication token",
				payload: map[string]any{
					"access_token":  "xxxxx.xxxxx.xxxxx",
					"refresh_token": "yyyyy.yyyyy.yyyyy",
				},
			},
			setup: func(d *dependency) {
				d.authUsecase.On("Refresh", mock.Anything, &domain.RefreshAuthIn{}).
					Return(domain.RefreshAuthOut{AccessToken: "xxxxx.xxxxx.xxxxx", RefreshToken: "yyyyy.yyyyy.yyyyy"}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			reqBody := bytes.NewReader(t.args.requestBody)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/", reqBody)

			deps := &dependency{
				authUsecase: &mocks.AuthUsecase{},
			}
			t.setup(deps)

			handler := New(deps.authUsecase)
			handler.Put(rr, req)

			s.Equal(t.expected.contentType, rr.Header().Get("Content-Type"))
			s.Equal(t.expected.statusCode, rr.Code)

			if t.isError {
				var resBody response.E
				json.NewDecoder(rr.Body).Decode(&resBody)

				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)
				s.Equal(t.expected.error, resBody.Error)
			} else {
				var resBody response.S
				json.NewDecoder(rr.Body).Decode(&resBody)

				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)
				s.Equal(t.expected.payload, resBody.Payload)
			}
		})
	}
}
