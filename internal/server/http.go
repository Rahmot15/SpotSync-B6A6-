package server

import (
	"net/http"

	"spotsync/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func NewHTTPServer(cfg config.Config, db *gorm.DB) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"success": true,
			"message": "SpotSync API is running",
		})
	})

	return e
}
