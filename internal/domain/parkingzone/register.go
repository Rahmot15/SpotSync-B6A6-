package parkingzone

import (
	"spotsync/internal/auth"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(group *echo.Group, handler *Handler, jwtService *auth.JWTService) {
	zonesGroup := group.Group("/zones")

	zonesGroup.GET("", handler.GetAll)
	zonesGroup.GET("/:id", handler.GetByID)

	adminGroup := zonesGroup.Group("")
	adminGroup.Use(middlewares.JWTAuth(jwtService), middlewares.RequireAdmin)
	adminGroup.POST("", handler.Create)
	adminGroup.PATCH("/:id", handler.Update)
	adminGroup.DELETE("/:id", handler.Delete)
}
