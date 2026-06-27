package httpresponse

import (
	"errors"
	"net/http"

	"spotsync/internal/apperror"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func Error(c echo.Context, statusCode int, message string, details any) error {
	return c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  details,
	})
}

func HandleError(c echo.Context, err error) error {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		return Error(c, appErr.StatusCode, appErr.Message, appErr.Details)
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		return Error(c, echoErr.Code, http.StatusText(echoErr.Code), echoErr.Message)
	}

	return Error(c, http.StatusInternalServerError, "Internal server error", nil)
}
