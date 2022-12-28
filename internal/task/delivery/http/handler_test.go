package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/mocks"
	"github.com/edwintantawi/taskit/pkg/response"
	"github.com/edwintantawi/taskit/test"
)

type TaskHTTPHandlerTestSuite struct {
	suite.Suite
}

func TestTaskHTTPHandlerSuite(t *testing.T) {
	suite.Run(t, new(TaskHTTPHandlerTestSuite))
}

func (s *TaskHTTPHandlerTestSuite) TestPost() {
	s.Run("it should return error response when request body is invalid", func() {
		handler := New(nil)

		reqBody := bytes.NewReader(nil)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/tasks", reqBody)

		handler.Post(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(400, rr.Code)
		s.Equal(400, resBody.StatusCode)
		s.Equal(http.StatusText(400), resBody.Message)
		s.Equal("Invalid request body", resBody.Error)
	})

	s.Run("it should return error response when fail to create new task", func() {
		mockTaskUsecase := &mocks.TaskUsecase{}
		mockTaskUsecase.On("Create", mock.Anything, mock.Anything).Return(domain.CreateTaskOut{}, errors.New("some error"))

		handler := New(mockTaskUsecase)

		reqBody := bytes.NewReader([]byte(`{}`))
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/tasks", reqBody)
		req = test.InjectAuthContext(req, "xxxxx")

		handler.Post(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(500, rr.Code)
		s.Equal(500, resBody.StatusCode)
		s.Equal(http.StatusText(500), resBody.Message)
		s.NotEmpty(resBody.Error)
	})

	s.Run("it should return success response when successfully create new task", func() {
		result := domain.CreateTaskOut{ID: "xxxxx"}
		mockTaskUsecase := &mocks.TaskUsecase{}
		mockTaskUsecase.On("Create", mock.Anything, mock.Anything).Return(result, nil)

		handler := New(mockTaskUsecase)

		reqBody := bytes.NewReader([]byte(`{}`))
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/tasks", reqBody)
		req = test.InjectAuthContext(req, "xxxxx")

		handler.Post(rr, req)

		var resBody response.S
		json.NewDecoder(rr.Body).Decode(&resBody)
		resPayload := resBody.Payload.(map[string]any)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(201, rr.Code)
		s.Equal(201, resBody.StatusCode)
		s.Equal("Successfully created new task", resBody.Message)
		s.NotEmpty(string(result.ID), resPayload["id"])
	})
}
