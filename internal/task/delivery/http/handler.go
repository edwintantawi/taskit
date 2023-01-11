package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/edwintantawi/taskit/internal/domain"
	"github.com/edwintantawi/taskit/internal/domain/dto"
	"github.com/edwintantawi/taskit/internal/domain/entity"
	"github.com/edwintantawi/taskit/pkg/errorx"
)

type HTTPHandler struct {
	validator   domain.ValidatorProvider
	taskUsecase domain.TaskUsecase
}

// New creates a new HTTPHandler.
func New(validator domain.ValidatorProvider, taskUsecase domain.TaskUsecase) HTTPHandler {
	return HTTPHandler{validator: validator, taskUsecase: taskUsecase}
}

// POST /tasks to create new task.
func (h *HTTPHandler) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.TaskCreateIn
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

	output, err := h.taskUsecase.Create(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(domain.NewErrorResponse(code, msg))
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(domain.NewSuccessResponse(http.StatusCreated, "Successfully created new task", output))
}

// GET /tasks to get all tasks.
func (h *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.TaskGetAllIn
	payload.UserID = entity.GetAuthContext(r.Context())

	output, err := h.taskUsecase.GetAll(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(domain.NewErrorResponse(code, msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(domain.NewSuccessResponse(http.StatusOK, http.StatusText(http.StatusOK), output))
}

// DELETE /tasks/{task_id} to remove task.
func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.TaskRemoveIn
	payload.UserID = entity.GetAuthContext(r.Context())
	payload.TaskID = entity.TaskID(chi.URLParam(r, "task_id"))

	if err := h.taskUsecase.Remove(r.Context(), &payload); err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(domain.NewErrorResponse(code, msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(domain.NewSuccessResponse(http.StatusOK, "Successfully deleted task", nil))
}

// GET /tasks/{task_id} to get task by task id.
func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.TaskGetByIDIn
	payload.UserID = entity.GetAuthContext(r.Context())
	payload.TaskID = entity.TaskID(chi.URLParam(r, "task_id"))

	output, err := h.taskUsecase.GetByID(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(domain.NewErrorResponse(code, msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(domain.NewSuccessResponse(http.StatusOK, http.StatusText(http.StatusOK), output))
}

// PUT /tasks/{task_id} to update task by task id.
func (h *HTTPHandler) Put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	var payload dto.TaskUpdateIn
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(domain.NewErrorResponse(http.StatusBadRequest, "Invalid request body"))
		return
	}
	payload.UserID = entity.GetAuthContext(r.Context())
	payload.TaskID = entity.TaskID(chi.URLParam(r, "task_id"))

	if err := h.validator.Validate(&payload); err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(domain.NewErrorResponse(code, msg))
		return
	}

	output, err := h.taskUsecase.Update(r.Context(), &payload)
	if err != nil {
		code, msg := errorx.HTTPErrorTranslator(err)
		w.WriteHeader(code)
		encoder.Encode(domain.NewErrorResponse(code, msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(domain.NewSuccessResponse(http.StatusOK, "Successfully updated task", output))
}
