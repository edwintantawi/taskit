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
)

type UserHTTPHandlerTestSuite struct {
	suite.Suite
}

func TestUserHTTPHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHTTPHandlerTestSuite))
}

func (s *UserHTTPHandlerTestSuite) TestPost() {
	s.Run("it should return error when request body is invalid", func() {
		handler := New(nil)

		reqBody := bytes.NewReader(nil)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/users", reqBody)

		handler.Post(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(400, rr.Code)
		s.Equal(400, resBody.StatusCode)
		s.Equal(http.StatusText(400), resBody.Message)
		s.Equal("Invalid request body", resBody.Error)
	})

	s.Run("it should return error when fail to create user", func() {
		usecase := &mocks.UserUsecase{}
		usecase.On("Create", mock.Anything, mock.Anything).Return(domain.CreateUserOut{}, errors.New("some error"))

		handler := New(usecase)

		reqBody := bytes.NewReader([]byte(`{}`))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/users", reqBody)

		handler.Post(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(500, rr.Code)
		s.Equal(500, resBody.StatusCode)
		s.Equal(http.StatusText(500), resBody.Message)
		s.Equal("Something went wrong", resBody.Error)
	})

	s.Run("it should return success when successfully create user", func() {
		usecaseResult := domain.CreateUserOut{
			ID:    "xxxxx",
			Email: "gopher@go.dev",
		}

		usecase := &mocks.UserUsecase{}
		usecase.On("Create", mock.Anything, mock.Anything).Return(usecaseResult, nil)

		handler := New(usecase)

		reqBody := bytes.NewReader([]byte(`{}`))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/users", reqBody)

		handler.Post(rr, req)

		var resBody response.S
		json.NewDecoder(rr.Body).Decode(&resBody)
		resPayload := resBody.Payload.(map[string]any)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(201, rr.Code)
		s.Equal(201, resBody.StatusCode)
		s.Equal("Successfully registered user", resBody.Message)
		s.Equal(string(usecaseResult.ID), resPayload["id"])
		s.Equal(usecaseResult.Email, resPayload["email"])
	})
}
