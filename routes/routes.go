package routes

import (
	"spotsync-api/handler"
	"spotsync-api/middleware"

	"github.com/labstack/echo/v4"
)

// SetupRoutes registers all API routes onto the Echo instance.
func SetupRoutes(e *echo.Echo, authHandler *handler.AuthHandler, zoneHandler *handler.ZoneHandler) {
	api := e.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Zone routes
	zones := api.Group("/zones")
	zones.GET("", zoneHandler.GetAllZones)         // public
	zones.GET("/:id", zoneHandler.GetZoneByID)     // public

	// Admin-only: create zone (needs JWT + admin role)
	zones.POST("", zoneHandler.CreateZone, middleware.JWTMiddleware, middleware.RequireAdmin)
}