package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain/dto"
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
	validator   *mocks.ValidatorProvider
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

				d.validator.On("Validate", mock.Anything).
					Return(test.ErrValidator)
			},
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

				d.validator.On("Validate", mock.Anything).
					Return(nil)

				d.taskUsecase.On("Create", mock.Anything, &dto.TaskCreateIn{UserID: "user-xxxxx"}).
					Return(dto.TaskCreateOut{}, test.ErrUnexpected)
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

				d.validator.On("Validate", mock.Anything).
					Return(nil)

				d.taskUsecase.On("Create", mock.Anything, &dto.TaskCreateIn{UserID: "user-xxxxx"}).
					Return(dto.TaskCreateOut{ID: "task-xxxxx"}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			reqBody := bytes.NewReader(t.args.requestBody)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/", reqBody)

			d := &dependency{
				req:         req,
				validator:   &mocks.ValidatorProvider{},
				taskUsecase: &mocks.TaskUsecase{},
			}
			t.setup(d)

			handler := New(d.validator, d.taskUsecase)
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

				d.taskUsecase.On("GetAll", mock.Anything, &dto.TaskGetAllIn{UserID: "user-xxxxx"}).
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

				d.taskUsecase.On("GetAll", mock.Anything, &dto.TaskGetAllIn{UserID: "user-xxxxx"}).
					Return([]dto.TaskGetAllOut{
						{ID: "task-xxxxx", Content: "task_xxxxx_content", Description: "task_xxxxx_description", IsCompleted: false, DueDate: entity.TaskDueDate{NullTime: sql.NullTime{Valid: false}}, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
						{ID: "task-yyyyy", Content: "task_yyyyy_content", Description: "task_yyyyy_description", IsCompleted: true, DueDate: entity.TaskDueDate{NullTime: sql.NullTime{Time: test.TimeAfterNow, Valid: true}}, CreatedAt: test.TimeBeforeNow, UpdatedAt: test.TimeBeforeNow},
					}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			d := &dependency{
				req:         req,
				taskUsecase: &mocks.TaskUsecase{},
			}
			t.setup(d)

			handler := New(nil, d.taskUsecase)
			handler.Get(rr, d.req)

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
	type args struct {
		params map[string]string
	}
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
		args     args
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

				d.taskUsecase.On("Remove", mock.Anything, &dto.TaskRemoveIn{TaskID: "", UserID: "user-xxxxx"}).
					Return(test.ErrUnexpected)
			},
		},
		{
			name:    "it should response with success when success",
			isError: false,
			args: args{
				params: map[string]string{
					"task_id": "task-xxxxx",
				},
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				message:     "Successfully deleted task",
				payload:     nil,
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.taskUsecase.On("Remove", mock.Anything, &dto.TaskRemoveIn{TaskID: "task-xxxxx", UserID: "user-xxxxx"}).
					Return(nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/{task_id}", nil)

			req = test.InjectChiRouterParams(req, t.args.params)

			d := &dependency{
				req:         req,
				taskUsecase: &mocks.TaskUsecase{},
			}
			t.setup(d)

			handler := New(nil, d.taskUsecase)
			handler.Delete(rr, d.req)

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
	type args struct {
		params map[string]string
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

				d.taskUsecase.On("GetByID", mock.Anything, &dto.TaskGetByIDIn{TaskID: "", UserID: "user-xxxxx"}).
					Return(dto.TaskGetByIDOut{}, test.ErrUnexpected)
			},
		},
		{
			name:    "it should response with success when success",
			isError: false,
			args: args{
				params: map[string]string{
					"task_id": "task-xxxxx",
				},
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				message:     http.StatusText(http.StatusOK),
				payload: map[string]any{
					"id": "task-xxxxx", "content": "task_xxxxx_content", "description": "task_xxxxx_description", "is_completed": true, "due_date": test.TimeAfterNow.Format(time.RFC3339Nano), "created_at": test.TimeBeforeNow.Format(time.RFC3339Nano), "updated_at": test.TimeBeforeNow.Format(time.RFC3339Nano),
				},
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.taskUsecase.On("GetByID", mock.Anything, &dto.TaskGetByIDIn{TaskID: "task-xxxxx", UserID: "user-xxxxx"}).
					Return(dto.TaskGetByIDOut{
						ID:          "task-xxxxx",
						Content:     "task_xxxxx_content",
						Description: "task_xxxxx_description",
						IsCompleted: true,
						DueDate:     entity.TaskDueDate{NullTime: sql.NullTime{Time: test.TimeAfterNow, Valid: true}},
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

			req = test.InjectChiRouterParams(req, t.args.params)

			d := &dependency{
				req:         req,
				taskUsecase: &mocks.TaskUsecase{},
			}
			t.setup(d)

			handler := New(nil, d.taskUsecase)
			handler.GetByID(rr, d.req)

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

func (s *TaskHTTPHandlerTestSuite) TestPut() {
	type args struct {
		requestBody []byte
		params      map[string]string
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

				d.validator.On("Validate", mock.Anything).
					Return(test.ErrValidator)
			},
		},
		{
			name:    "it should response with error when task usecase GetByID return unexpected error",
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

				d.validator.On("Validate", mock.Anything).
					Return(nil)

				d.taskUsecase.On("Update", mock.Anything, &dto.TaskUpdateIn{TaskID: "", UserID: "user-xxxxx"}).
					Return(dto.TaskUpdateOut{}, test.ErrUnexpected)
			},
		},
		{
			name:    "it should response with success when success",
			isError: false,
			args: args{
				requestBody: []byte(`{}`),
				params: map[string]string{
					"task_id": "task-xxxxx",
				},
			},
			expected: expected{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				message:     "Successfully updated task",
				payload: map[string]any{
					"id": "task-xxxxx",
				},
			},
			setup: func(d *dependency) {
				d.req = test.InjectAuthContext(d.req, entity.UserID("user-xxxxx"))

				d.validator.On("Validate", mock.Anything).
					Return(nil)

				d.taskUsecase.On("Update", mock.Anything, &dto.TaskUpdateIn{
					TaskID: "task-xxxxx",
					UserID: "user-xxxxx",
				}).Return(dto.TaskUpdateOut{
					ID: "task-xxxxx",
				}, nil)
			},
		},
	}

	for _, t := range tests {
		s.Run(t.name, func() {
			reqBody := bytes.NewReader(t.args.requestBody)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/{task_id}", reqBody)

			req = test.InjectChiRouterParams(req, t.args.params)

			d := &dependency{
				req:         req,
				validator:   &mocks.ValidatorProvider{},
				taskUsecase: &mocks.TaskUsecase{},
			}
			t.setup(d)

			handler := New(d.validator, d.taskUsecase)
			handler.Put(rr, d.req)

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
