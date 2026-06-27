package middlewares

import (
	"strings"

	"spotsync/internal/apperror"
	"spotsync/internal/auth"

	"github.com/labstack/echo/v4"
)

const (
	ContextUserID = "user_id"
	ContextRole   = "role"
)

func JWTAuth(jwtService *auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString, err := extractBearerToken(c.Request().Header.Get("Authorization"))
			if err != nil {
				return err
			}

			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return apperror.Unauthorized("Invalid or expired token")
			}

			c.Set(ContextUserID, claims.UserID)
			c.Set(ContextRole, claims.Role)

			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (uint, bool) {
	userID, ok := c.Get(ContextUserID).(uint)
	return userID, ok
}

func GetRole(c echo.Context) (string, bool) {
	role, ok := c.Get(ContextRole).(string)
	return role, ok
}

func extractBearerToken(authorizationHeader string) (string, error) {
	if authorizationHeader == "" {
		return "", apperror.Unauthorized("Authorization token is required")
	}

	parts := strings.SplitN(authorizationHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", apperror.Unauthorized("Authorization header must be Bearer token")
	}

	return parts[1], nil
}
