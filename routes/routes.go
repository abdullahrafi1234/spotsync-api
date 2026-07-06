package routes

import (
	"spotsync-api/handler"

	"github.com/labstack/echo/v4"
)

// SetupRoutes registers all API routes onto the Echo instance.
func SetupRoutes(e *echo.Echo, authHandler *handler.AuthHandler) {
	api := e.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}