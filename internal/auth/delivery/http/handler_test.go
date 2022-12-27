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

type AuthHTTPHandlerTestSuite struct {
	suite.Suite
}

func TestAuthHTTPHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHTTPHandlerTestSuite))
}

func (s *AuthHTTPHandlerTestSuite) TestPost() {
	s.Run("it should return error response when request body is invalid", func() {
		handler := New(nil)

		reqBody := bytes.NewReader(nil)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/authentications", reqBody)

		handler.Post(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(400, rr.Code)
		s.Equal(400, resBody.StatusCode)
		s.Equal(http.StatusText(400), resBody.Message)
		s.Equal("Invalid request body", resBody.Error)
	})

	s.Run("it should return error response when fail to login with existing user", func() {
		usecase := &mocks.AuthUsecase{}
		usecase.On("Login", mock.Anything, mock.Anything).Return(domain.LoginAuthOut{}, errors.New("some error"))

		handler := New(usecase)

		reqBody := bytes.NewReader([]byte(`{}`))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/authentications", reqBody)

		handler.Post(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(500, rr.Code)
		s.Equal(500, resBody.StatusCode)
		s.Equal(http.StatusText(500), resBody.Message)
		s.Equal("Something went wrong", resBody.Error)
	})

	s.Run("it should return success response when successfully login with existing user", func() {
		usecaseResult := domain.LoginAuthOut{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
		}

		usecase := &mocks.AuthUsecase{}
		usecase.On("Login", mock.Anything, mock.Anything).Return(usecaseResult, nil)

		handler := New(usecase)

		reqBody := bytes.NewReader([]byte(`{}`))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/authentications", reqBody)

		handler.Post(rr, req)

		var resBody response.S
		json.NewDecoder(rr.Body).Decode(&resBody)
		resPayload := resBody.Payload.(map[string]any)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(200, rr.Code)
		s.Equal(200, resBody.StatusCode)
		s.Equal("Successfully logged in user", resBody.Message)
		s.Equal(string(usecaseResult.AccessToken), resPayload["access_token"])
		s.Equal(usecaseResult.RefreshToken, resPayload["refresh_token"])
	})
}

func (s *AuthHTTPHandlerTestSuite) TestDelete() {
	s.Run("it should return error response when request body is invalid", func() {
		handler := New(nil)

		reqBody := bytes.NewReader(nil)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/authentications", reqBody)

		handler.Delete(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(400, rr.Code)
		s.Equal(400, resBody.StatusCode)
		s.Equal(http.StatusText(400), resBody.Message)
		s.Equal("Invalid request body", resBody.Error)
	})

	s.Run("it should return error response when fail to logout", func() {
		usecase := &mocks.AuthUsecase{}
		usecase.On("Logout", mock.Anything, mock.Anything).Return(errors.New("unexpected error"))

		handler := New(usecase)

		reqBody := bytes.NewReader([]byte(`{}`))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/authenticaitons", reqBody)

		handler.Delete(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(500, rr.Code)
		s.Equal(500, resBody.StatusCode)
		s.Equal(http.StatusText(500), resBody.Message)
		s.NotEmpty(resBody.Error)
	})

	s.Run("it should return error response when success to logout", func() {
		usecase := &mocks.AuthUsecase{}
		usecase.On("Logout", mock.Anything, mock.Anything).Return(nil)

		handler := New(usecase)

		reqBody := bytes.NewReader([]byte(`{}`))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/authenticaitons", reqBody)

		handler.Delete(rr, req)

		var resBody response.S
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(200, rr.Code)
		s.Equal(200, resBody.StatusCode)
		s.Equal("Successfully logout user", resBody.Message)
		s.Nil(resBody.Payload)
	})
}
