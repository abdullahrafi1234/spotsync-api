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
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	// 2. Connect database
	config.ConnectDatabase()

	// 3. Dependency Injection: Repository -> Service -> Handler
	userRepo := repository.NewUserRepository(config.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// 4. Create Echo instance
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "SpotSync API is running 🚗⚡",
		})
	})

	// 5. Register routes
	routes.SetupRoutes(e, authHandler)

	// 6. Start server
	e.Logger.Fatal(e.Start(":8080"))
}