package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
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

func (s *AuthHTTPHandlerTestSuite) TestGet() {
	s.Run("it should return error response when fail find the profile", func() {
		userProfile := domain.GetProfileAuthOut{
			ID:    "xxxxx",
			Name:  "Gopher",
			Email: "gopher@go.dev",
		}

		usecase := &mocks.AuthUsecase{}
		usecase.On("GetProfile", mock.Anything, mock.Anything).Return(domain.GetProfileAuthOut{}, errors.New("unexpected error"))

		handler := New(usecase)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/authentications", nil)
		req = req.WithContext(context.WithValue(req.Context(), entity.AuthUserIDKey("user_id"), userProfile.ID))

		handler.Get(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(500, rr.Code)
		s.Equal(500, resBody.StatusCode)
		s.Equal(http.StatusText(500), resBody.Message)
		s.NotEmpty(resBody.Error)
	})

	s.Run("it should return user profile response when success find the profile", func() {
		userProfile := domain.GetProfileAuthOut{
			ID:    "xxxxx",
			Name:  "Gopher",
			Email: "gopher@go.dev",
		}

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/authentications", nil)
		req = req.WithContext(context.WithValue(req.Context(), entity.AuthUserIDKey("user_id"), userProfile.ID))

		usecase := &mocks.AuthUsecase{}
		usecase.On("GetProfile", req.Context(), &domain.GetProfileAuthIn{UserID: userProfile.ID}).Return(userProfile, nil)

		handler := New(usecase)

		handler.Get(rr, req)

		var resBody response.S
		json.NewDecoder(rr.Body).Decode(&resBody)
		resPayload := resBody.Payload.(map[string]any)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(200, rr.Code)
		s.Equal(200, resBody.StatusCode)
		s.Equal(http.StatusText(200), resBody.Message)
		s.Equal(string(userProfile.ID), resPayload["id"])
		s.Equal(userProfile.Name, resPayload["name"])
		s.Equal(userProfile.Email, resPayload["email"])
	})
}

func (s *AuthHTTPHandlerTestSuite) TestPut() {
	s.Run("it should return error response when request body is invalid", func() {
		handler := New(nil)

		reqBody := bytes.NewReader(nil)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/authentications", reqBody)

		handler.Put(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(400, rr.Code)
		s.Equal(400, resBody.StatusCode)
		s.Equal(http.StatusText(400), resBody.Message)
		s.Equal("Invalid request body", resBody.Error)
	})

	s.Run("it should return error response when fail to refresh authentication token", func() {
		usecase := &mocks.AuthUsecase{}
		usecase.On("Refresh", mock.Anything, mock.Anything).Return(domain.RefreshAuthOut{}, errors.New("unexpected error"))

		handler := New(usecase)

		reqBody := bytes.NewReader([]byte(`{}`))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/authentications", reqBody)

		handler.Put(rr, req)

		var resBody response.E
		json.NewDecoder(rr.Body).Decode(&resBody)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(500, rr.Code)
		s.Equal(500, resBody.StatusCode)
		s.Equal(http.StatusText(500), resBody.Message)
		s.NotEmpty(resBody.Error)
	})

	s.Run("it should return authentication token response when success to refresh authentication token", func() {
		result := domain.RefreshAuthOut{
			AccessToken:  "xxxxx.xxxxx.xxxxx",
			RefreshToken: "yyyyy.yyyyy.yyyyy",
		}
		usecase := &mocks.AuthUsecase{}
		usecase.On("Refresh", mock.Anything, mock.Anything).Return(result, nil)

		handler := New(usecase)

		reqBody := bytes.NewReader([]byte(`{}`))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/authentications", reqBody)

		handler.Put(rr, req)

		var resBody response.S
		json.NewDecoder(rr.Body).Decode(&resBody)
		resPayload := resBody.Payload.(map[string]any)

		s.Equal("application/json", rr.Header().Get("Content-Type"))
		s.Equal(200, rr.Code)
		s.Equal(200, resBody.StatusCode)
		s.Equal("Successfully refreshed authentication token", resBody.Message)
		s.Equal(result.AccessToken, resPayload["access_token"])
		s.Equal(result.RefreshToken, resPayload["refresh_token"])
	})
}
