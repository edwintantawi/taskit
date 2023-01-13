package http

import (
	"encoding/json"
	"net/http"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/pkg/errorx"
)

type HTTPHandler struct {
	validator      domain.ValidatorProvider
	projectUsecase domain.ProjectUsecase
}

// New creates a new project handler.
func New(validator domain.ValidatorProvider, projectUsecase domain.ProjectUsecase) HTTPHandler {
	return HTTPHandler{validator: validator, projectUsecase: projectUsecase}
}

// POST /projects to create new project.
func (h *HTTPHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.ProjectCreateIn
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(domain.NewErrorResponse(http.StatusBadRequest, "Invalid request body"))
		return
	}
	payload.UserID = entity.GetAuthContext(r.Context())

	if err := h.validator.Validate(&payload); err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(domain.NewErrorResponse(code, msg))
		return
	}

	output, err := h.projectUsecase.Create(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(domain.NewErrorResponse(code, msg))
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(domain.NewSuccessResponse(http.StatusCreated, "Successfully created project", output))
}

// GET /projects to get all project.
func (h *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.ProjectGetAllIn
	payload.UserID = entity.GetAuthContext(r.Context())

	output, err := h.projectUsecase.GetAll(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(domain.NewErrorResponse(code, msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(domain.NewSuccessResponse(http.StatusOK, http.StatusText(http.StatusOK), output))
}
