package httpresponse

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(c echo.Context, statusCode int, message string, data any) error {
	return c.JSON(statusCode, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func OK(c echo.Context, message string, data any) error {
	return Success(c, http.StatusOK, message, data)
}

func Created(c echo.Context, message string, data any) error {
	return Success(c, http.StatusCreated, message, data)
}
