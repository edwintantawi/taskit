package http

import (
	"encoding/json"
	"net/http"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/pkg/errorx"
)

type HTTPHandler struct {
	userUsecase domain.UserUsecase
	validator   domain.ValidatorProvider
}

// New creates a new user handler.
func New(validator domain.ValidatorProvider, userUsecase domain.UserUsecase) HTTPHandler {
	return HTTPHandler{validator: validator, userUsecase: userUsecase}
}

// POST /users to create new user.
func (h *HTTPHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.UserCreateIn
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(domain.NewErrorResponse(http.StatusBadRequest, "Invalid request body"))
		return
	}
	if err := h.validator.Validate(&payload); err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(domain.NewErrorResponse(code, msg))
		return
	}

	output, err := h.userUsecase.Create(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(domain.NewErrorResponse(code, msg))
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(domain.NewSuccessResponse(http.StatusCreated, "Successfully registered user", output))
}
