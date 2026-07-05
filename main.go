package main

import (
	"log"
	"spotsync-api/config"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	// Connect to database
	config.ConnectDatabase()

	// Create Echo instance
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "SpotSync API is running 🚗⚡",
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}