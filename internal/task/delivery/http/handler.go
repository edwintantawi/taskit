package http

import (
	"encoding/json"
	"net/http"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/pkg/errorx"
	"github.com/edwintantawi/taskit/pkg/response"
)

type HTTPHandler struct {
	taskUsecase domain.TaskUsecase
}

// New creates a new HTTPHandler.
func New(taskUsecase domain.TaskUsecase) *HTTPHandler {
	return &HTTPHandler{taskUsecase: taskUsecase}
}

// POST /tasks to create new task
func (h *HTTPHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload domain.CreateTaskIn
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(response.Error(http.StatusBadRequest, "Invalid request body"))
		return
	}
	payload.UserID = entity.GetAuthContext(r.Context())

	output, err := h.taskUsecase.Create(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(response.Error(code, msg))
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(response.Success(http.StatusCreated, "Successfully created new task", output))
}
