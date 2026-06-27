package reservation

import (
	"spotsync/internal/auth"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(group *echo.Group, handler *Handler, jwtService *auth.JWTService) {
	reservationsGroup := group.Group("/reservations")

	authGroup := reservationsGroup.Group("")
	authGroup.Use(middlewares.JWTAuth(jwtService))
	authGroup.POST("", handler.Reserve)
	authGroup.GET("/my-reservations", handler.GetMyReservations)
	authGroup.DELETE("/:id", handler.Cancel)

	adminGroup := reservationsGroup.Group("")
	adminGroup.Use(middlewares.JWTAuth(jwtService), middlewares.RequireAdmin)
	adminGroup.GET("", handler.GetAll)
}
