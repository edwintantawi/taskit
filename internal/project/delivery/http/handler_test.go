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
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/pkg/errorx"
	"github.com/edwintantawi/taskit/test"
)

type ProjectHTTPHandlerTestSuite struct {
	suite.Suite
}

func TestProjectHTTPHandlerSuite(t *testing.T) {
	suite.Run(t, new(ProjectHTTPHandlerTestSuite))
}

type dependency struct {
	req            *http.Request
	validator      *mocks.ValidatorProvider
	projectUsecase *mocks.ProjectUsecase
}

func (s *ProjectHTTPHandlerTestSuite) TestPost() {
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
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.validator.On("Validate", &dto.ProjectCreateIn{UserID: "user-xxxxx"}).
					Return(test.ErrValidator)
			},
		},
		{
			name:    "it should response with error when project usecase Create return unexpected error",
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
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.validator.On("Validate", &dto.ProjectCreateIn{UserID: "user-xxxxx"}).
					Return(nil)

				d.projectUsecase.On("Create", mock.Anything, &dto.ProjectCreateIn{UserID: "user-xxxxx"}).
					Return(dto.ProjectCreateOut{}, test.ErrUnexpected)
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
				message:     "Successfully created project",
				payload: map[string]any{
					"id": "project-xxxxx",
				},
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.validator.On("Validate", &dto.ProjectCreateIn{UserID: "user-xxxxx"}).
					Return(nil)

				d.projectUsecase.On("Create", mock.Anything, &dto.ProjectCreateIn{UserID: "user-xxxxx"}).
					Return(dto.ProjectCreateOut{ID: "project-xxxxx"}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			reqBody := bytes.NewReader(t.args.requestBody)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/projects", reqBody)

			d := &dependency{
				validator:      &mocks.ValidatorProvider{},
				projectUsecase: &mocks.ProjectUsecase{},
				req:            req,
			}
			t.setup(d)

			handler := New(d.validator, d.projectUsecase)
			handler.Post(rr, d.req)

			s.Equal(t.expected.contentType, rr.Header().Get("Content-Type"))
			s.Equal(t.expected.statusCode, rr.Code)

			if t.isError {
				var resBody domain.ErrorResponse
				json.NewDecoder(rr.Body).Decode(&resBody)

				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)
				s.Equal(t.expected.error, resBody.Error)
			} else {
				var resBody domain.SuccessResponse
				json.NewDecoder(rr.Body).Decode(&resBody)
				payloadMap := resBody.Payload.(map[string]any)

				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)
				s.Equal(t.expected.payload, payloadMap)
			}
		})
	}
}
