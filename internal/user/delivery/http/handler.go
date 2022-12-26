package http

import (
	"encoding/json"
	"net/http"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/pkg/errorx"
	"github.com/edwintantawi/taskit/pkg/response"
)

type HTTPHandler struct {
	userUsecase domain.UserUsecase
}

// New creates a new user handler.
func New(userUsecase domain.UserUsecase) *HTTPHandler {
	return &HTTPHandler{userUsecase: userUsecase}
}

// POST /users to create new user.
func (h *HTTPHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload domain.CreateUserIn
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(response.Error(http.StatusBadRequest, "Invalid request body"))
		return
	}

	result, err := h.userUsecase.Create(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(response.Error(code, msg))
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(response.Success(http.StatusCreated, "Successfully registered user", result))
}
