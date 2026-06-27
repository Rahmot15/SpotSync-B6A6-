package apperror

import "net/http"

type AppError struct {
	StatusCode int
	Message    string
	Details    any
}

func (e *AppError) Error() string {
	return e.Message
}

func New(statusCode int, message string, details any) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

func BadRequest(message string, details any) *AppError {
	return New(http.StatusBadRequest, message, details)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message, nil)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, message, nil)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message, nil)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, message, nil)
}

func Internal(message string) *AppError {
	return New(http.StatusInternalServerError, message, nil)
}
