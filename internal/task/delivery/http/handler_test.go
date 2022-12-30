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

type TaskHTTPHandlerTestSuite struct {
	suite.Suite
}

func TestTaskHTTPHandlerSuite(t *testing.T) {
	suite.Run(t, new(TaskHTTPHandlerTestSuite))
}

type dependency struct {
	req         *http.Request
	taskUsecase *mocks.TaskUsecase
}

func (s *TaskHTTPHandlerTestSuite) TestPost() {
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
			name:    "it should response with error when taskUsecase Create returns unexpected error",
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
				d.taskUsecase.On("Create", mock.Anything, &domain.CreateTaskIn{UserID: "user-xxxxx"}).
					Return(domain.CreateTaskOut{}, test.ErrUnexpected)
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
				message:     "Successfully created new task",
				payload: map[string]any{
					"id": "task-xxxxx",
				},
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))
				d.taskUsecase.On("Create", mock.Anything, &domain.CreateTaskIn{UserID: "user-xxxxx"}).
					Return(domain.CreateTaskOut{ID: "task-xxxxx"}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			reqBody := bytes.NewReader(t.args.requestBody)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/", reqBody)

			deps := &dependency{
				taskUsecase: &mocks.TaskUsecase{},
				req:         req,
			}
			t.setup(deps)

			handler := New(deps.taskUsecase)
			handler.Post(rr, deps.req)

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
