package server

import (
	"log"
	"net/http"

	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/domain/parkingzone"
	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/user"
	"spotsync/internal/httpresponse"
	customvalidator "spotsync/internal/validator"

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
	e.Validator = customvalidator.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if handleErr := httpresponse.HandleError(c, err); handleErr != nil {
			e.DefaultHTTPErrorHandler(handleErr, c)
		}
	}

	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return httpresponse.Success(c, http.StatusOK, "SpotSync API is running", nil)
	})

	jwtService := auth.NewJWTService(cfg.JWTSecret)
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository, jwtService)
	userHandler := user.NewHandler(userService)

	zoneRepository := parkingzone.NewRepository(db)
	zoneService := parkingzone.NewService(zoneRepository)
	zoneHandler := parkingzone.NewHandler(zoneService)

	api := e.Group("/api/v1")
	user.RegisterRoutes(api, userHandler)
	parkingzone.RegisterRoutes(api, zoneHandler, jwtService)

	return e
}
