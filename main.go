package main

import (
	"log"

	"spotsync-api/config"
	"spotsync-api/handler"
	"spotsync-api/repository"
	"spotsync-api/routes"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	// 2. Connect database
	config.ConnectDatabase()

	// 3. Repositories
	userRepo := repository.NewUserRepository(config.DB)
	zoneRepo := repository.NewZoneRepository(config.DB)
	reservationRepo := repository.NewReservationRepository(config.DB)

	// 4. Services
	authService := service.NewAuthService(userRepo)
	zoneService := service.NewZoneService(zoneRepo)
	reservationService := service.NewReservationService(reservationRepo, zoneRepo)

	// 5. Handlers
	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	reservationHandler := handler.NewReservationHandler(reservationService)

	// 6. Create Echo instance
	e := echo.New()

	// Register our centralized error handler — every returned error
	// from any handler funnels through here.
	e.HTTPErrorHandler = utils.CentralErrorHandler

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "SpotSync API is running 🚗⚡",
		})
	})

	// 7. Register routes
	routes.SetupRoutes(e, authHandler, zoneHandler, reservationHandler)

	// 8. Start server
	e.Logger.Fatal(e.Start(":8080"))
}