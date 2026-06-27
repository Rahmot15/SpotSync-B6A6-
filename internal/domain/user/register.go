package user

import "github.com/labstack/echo/v4"

func RegisterRoutes(group *echo.Group, handler *Handler) {
	authGroup := group.Group("/auth")

	authGroup.POST("/register", handler.Register)
	authGroup.POST("/login", handler.Login)
}
