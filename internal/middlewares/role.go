package middlewares

import (
	"spotsync/internal/apperror"

	"github.com/labstack/echo/v4"
)

const RoleAdmin = "admin"

func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := GetRole(c)
		if !ok {
			return apperror.Unauthorized("Authenticated user role not found")
		}

		if role != RoleAdmin {
			return apperror.Forbidden("Admin access required")
		}

		return next(c)
	}
}
