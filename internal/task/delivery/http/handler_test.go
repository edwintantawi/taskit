package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
				req:         req,
				taskUsecase: &mocks.TaskUsecase{},
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

func (s *TaskHTTPHandlerTestSuite) TestGet() {
	type expected struct {
		contentType string
		statusCode  int
		message     string
		error       string
		payload     []map[string]any
	}
	tests := []struct {
		name     string
		isError  bool
		expected expected
		setup    func(d *dependency)
	}{
		{
			name:    "it should response with error when task usecase return unexpected error",
			isError: true,
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusInternalServerError,
				message:     http.StatusText(http.StatusInternalServerError),
				error:       errorx.InternalServerErrorMessage,
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.taskUsecase.On("GetAll", mock.Anything, &domain.GetAllTaskIn{UserID: "user-xxxxx"}).
					Return(nil, test.ErrUnexpected)
			},
		},
		{
			name:    "it should response with success when success",
			isError: false,
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				message:     http.StatusText(http.StatusOK),
				payload: []map[string]any{
					{"id": "task-xxxxx", "content": "task_xxxxx_content", "description": "task_xxxxx_description", "is_completed": false, "due_date": nil, "created_at": test.TimeBeforeNow.Format(time.RFC3339Nano), "updated_at": test.TimeBeforeNow.Format(time.RFC3339Nano)},
					{"id": "task-yyyyy", "content": "task_yyyyy_content", "description": "task_yyyyy_description", "is_completed": true, "due_date": test.TimeAfterNow.Format(time.RFC3339Nano), "created_at": test.TimeBeforeNow.Format(time.RFC3339Nano), "updated_at": test.TimeBeforeNow.Format(time.RFC3339Nano)},
				},
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.taskUsecase.On("GetAll", mock.Anything, &domain.GetAllTaskIn{UserID: "user-xxxxx"}).
					Return([]domain.GetAllTaskOut{
						{ID: "task-xxxxx", Content: "task_xxxxx_content", Description: "task_xxxxx_description", IsCompleted: false, DueDate: nil, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
						{ID: "task-yyyyy", Content: "task_yyyyy_content", Description: "task_yyyyy_description", IsCompleted: true, DueDate: &test.TimeAfterNow, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
					}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			deps := &dependency{
				req:         req,
				taskUsecase: &mocks.TaskUsecase{},
			}
			t.setup(deps)

			handler := New(deps.taskUsecase)
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
				payloadList := resBody.Payload.([]any)

				s.Equal(t.expected.statusCode, resBody.StatusCode)
				s.Equal(t.expected.message, resBody.Message)

				for i, payload := range t.expected.payload {
					s.Equal(payload, payloadList[i].(map[string]any))
				}
			}
		})
	}
}

func (s *TaskHTTPHandlerTestSuite) TestDelete() {
	type expected struct {
		contentType string
		statusCode  int
		message     string
		error       string
		payload     any
	}
	tests := []struct {
		name     string
		isError  bool
		expected expected
		setup    func(d *dependency)
	}{
		{
			name:    "it should response with error when task usecase Remove return unexpected error",
			isError: true,
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusInternalServerError,
				message:     http.StatusText(http.StatusInternalServerError),
				error:       errorx.InternalServerErrorMessage,
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.taskUsecase.On("Remove", mock.Anything, &domain.RemoveTaskIn{TaskID: "", UserID: "user-xxxxx"}).
					Return(test.ErrUnexpected)
			},
		},
		{
			name:    "it should response with success when success",
			isError: false,
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				message:     "Successfully deleted task",
				payload:     nil,
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.taskUsecase.On("Remove", mock.Anything, &domain.RemoveTaskIn{TaskID: "", UserID: "user-xxxxx"}).
					Return(nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/", nil)

			deps := &dependency{
				req:         req,
				taskUsecase: &mocks.TaskUsecase{},
			}
			t.setup(deps)

			handler := New(deps.taskUsecase)
			handler.Delete(rr, deps.req)

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

				s.Equal(t.expected.payload, nil)
			}
		})
	}
}

func (s *TaskHTTPHandlerTestSuite) TestGetByID() {
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
			name:    "it should response with error when task usecase GetByID return unexpected error",
			isError: true,
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusInternalServerError,
				message:     http.StatusText(http.StatusInternalServerError),
				error:       errorx.InternalServerErrorMessage,
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.taskUsecase.On("GetByID", mock.Anything, &domain.GetTaskByIDIn{TaskID: "", UserID: "user-xxxxx"}).
					Return(domain.GetTaskByIDOut{}, test.ErrUnexpected)
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
					"id": "task-yyyyy", "content": "task_yyyyy_content", "description": "task_yyyyy_description", "is_completed": true, "due_date": test.TimeAfterNow.Format(time.RFC3339Nano), "created_at": test.TimeBeforeNow.Format(time.RFC3339Nano), "updated_at": test.TimeBeforeNow.Format(time.RFC3339Nano),
				},
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.taskUsecase.On("GetByID", mock.Anything, &domain.GetTaskByIDIn{TaskID: "", UserID: "user-xxxxx"}).
					Return(domain.GetTaskByIDOut{
						ID:          "task-yyyyy",
						Content:     "task_yyyyy_content",
						Description: "task_yyyyy_description",
						IsCompleted: true,
						DueDate:     &test.TimeAfterNow,
						CreatedAt:   test.TimeBeforeNow,
						UpdatedAt:   test.TimeBeforeNow,
					}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/{task_id}", nil)

			deps := &dependency{
				req:         req,
				taskUsecase: &mocks.TaskUsecase{},
			}
			t.setup(deps)

			handler := New(deps.taskUsecase)
			handler.GetByID(rr, deps.req)

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
