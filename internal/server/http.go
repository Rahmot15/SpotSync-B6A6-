package server

import (
	"log"
	"net/http"

	"spotsync/internal/config"
	"spotsync/internal/domain/parkingzone"
	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func NewHTTPServer(cfg config.Config, db *gorm.DB) *echo.Echo {
	if err := config.AutoMigrate(
		db,
		&user.User{},
		&parkingzone.ParkingZone{},
		&reservation.Reservation{},
	); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

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
