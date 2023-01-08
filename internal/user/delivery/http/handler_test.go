package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/pkg/errorx"
	"github.com/edwintantawi/taskit/pkg/response"
	"github.com/edwintantawi/taskit/test"
)

type UserHTTPHandlerTestSuite struct {
	suite.Suite
}

func TestUserHTTPHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHTTPHandlerTestSuite))
}

type dependency struct {
	req         *http.Request
	validator   *mocks.ValidatorProvider
	userUsecase *mocks.UserUsecase
}

func (s *UserHTTPHandlerTestSuite) TestPost() {
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
			name:    "it should response with error when payload is not valid",
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
				d.validator.On("Validate", &dto.UserCreateIn{}).
					Return(test.ErrValidator)
			},
		},
		{
			name:    "it should response with error when user usecase Create return unexpected error",
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
				d.validator.On("Validate", &dto.UserCreateIn{}).
					Return(nil)

				d.userUsecase.On("Create", mock.Anything, &dto.UserCreateIn{}).
					Return(dto.UserCreateOut{}, test.ErrUnexpected)
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
				statusCode:  http.StatusCreated,
				message:     "Successfully registered user",
				payload: map[string]any{
					"id":    "user-xxxxx",
					"email": "gopher@go.dev",
				},
			},
			setup: func(d *dependency) {
				d.validator.On("Validate", &dto.UserCreateIn{}).
					Return(nil)

				d.userUsecase.On("Create", mock.Anything, &dto.UserCreateIn{}).
					Return(dto.UserCreateOut{ID: "user-xxxxx", Email: "gopher@go.dev"}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			reqBody := bytes.NewReader(t.args.requestBody)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/", reqBody)

			d := &dependency{
				validator:   &mocks.ValidatorProvider{},
				userUsecase: &mocks.UserUsecase{},
				req:         req,
			}
			t.setup(d)

			handler := New(d.validator, d.userUsecase)
			handler.Post(rr, d.req)

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
