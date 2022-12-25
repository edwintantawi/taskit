package response

import "net/http"

// S represents the success response with payload.
type S struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Payload    any    `json:"payload"`
}

// E represents the error response.
type E struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

// Success creates a new ResponseSuccess.
func Success(statusCode int, message string, payload any) *S {
	return &S{
		StatusCode: statusCode,
		Message:    message,
		Payload:    payload,
	}
}

// Error creates a new ResponseError and translate an error.
func Error(statusCode int, errorMessage string) *E {
	return &E{
		StatusCode: statusCode,
		Message:    http.StatusText(statusCode),
		Error:      errorMessage,
	}
}
