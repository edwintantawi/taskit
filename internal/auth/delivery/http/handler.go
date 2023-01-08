package http

import (
	"encoding/json"
	"net/http"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/pkg/errorx"
	"github.com/edwintantawi/taskit/pkg/response"
)

type HTTPHandler struct {
	validator   domain.ValidatorProvider
	authUsecase domain.AuthUsecase
}

// New creates a new auth handler
func New(validator domain.ValidatorProvider, authUsecase domain.AuthUsecase) *HTTPHandler {
	return &HTTPHandler{validator: validator, authUsecase: authUsecase}
}

// POST /authentications to login user
func (h *HTTPHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.AuthLoginIn
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(response.Error(http.StatusBadRequest, "Invalid request body"))
		return
	}
	if err := h.validator.Validate(&payload); err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(response.Error(code, msg))
		return
	}

	output, err := h.authUsecase.Login(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(response.Error(code, msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(response.Success(http.StatusOK, "Successfully logged in user", output))
}

// DELETE /authentications to logout from current authentication
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.AuthLogoutIn
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(response.Error(http.StatusBadRequest, "Invalid request body"))
		return
	}
	if err := h.validator.Validate(&payload); err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(response.Error(code, msg))
		return
	}

	err := h.authUsecase.Logout(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(response.Error(code, msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(response.Success(http.StatusOK, "Successfully logout user", nil))
}

// GET /authentications to get user authenticated profile
func (h *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.AuthProfileIn
	payload.UserID = entity.GetAuthContext(r.Context())

	output, err := h.authUsecase.GetProfile(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(response.Error(code, msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(response.Success(http.StatusOK, http.StatusText(http.StatusOK), output))
}

// PUT /authentications to refresh authentication token
func (h *HTTPHandler) Put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.AuthRefreshIn
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(response.Error(http.StatusBadRequest, "Invalid request body"))
		return
	}
	if err := h.validator.Validate(&payload); err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(response.Error(code, msg))
		return
	}

	output, err := h.authUsecase.Refresh(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(response.Error(code, msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(response.Success(http.StatusOK, "Successfully refreshed authentication token", output))
}
