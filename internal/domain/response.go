package domain

import "net/http"

// SuccessResponse represents the success response with payload.
type SuccessResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Payload    any    `json:"payload"`
}

// ErrorResponse represents the error response.
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

// NewSuccessResponse creates a new ResponseSuccess.
func NewSuccessResponse(statusCode int, message string, payload any) SuccessResponse {
	return SuccessResponse{
		StatusCode: statusCode,
		Message:    message,
		Payload:    payload,
	}
}

// NewErrorResponse creates a new ResponseError and translate an error.
func NewErrorResponse(statusCode int, errorMessage string) ErrorResponse {
	return ErrorResponse{
		StatusCode: statusCode,
		Message:    http.StatusText(statusCode),
		Error:      errorMessage,
	}
}
