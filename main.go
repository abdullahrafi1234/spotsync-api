package main

import (
	"log"

	"spotsync-api/config"
	"spotsync-api/handler"
	"spotsync-api/repository"
	"spotsync-api/routes"
	"spotsync-api/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	config.ConnectDatabase()

	// Repositories
	userRepo := repository.NewUserRepository(config.DB)
	zoneRepo := repository.NewZoneRepository(config.DB)
	reservationRepo := repository.NewReservationRepository(config.DB)

	// Services
	authService := service.NewAuthService(userRepo)
	zoneService := service.NewZoneService(zoneRepo)
	reservationService := service.NewReservationService(reservationRepo, zoneRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	zoneHandler := handler.NewZoneHandler(zoneService)
	reservationHandler := handler.NewReservationHandler(reservationService)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "SpotSync API is running 🚗⚡",
		})
	})

	routes.SetupRoutes(e, authHandler, zoneHandler, reservationHandler)

	e.Logger.Fatal(e.Start(":8080"))
}