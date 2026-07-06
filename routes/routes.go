package routes

import (
	"spotsync-api/handler"
	"spotsync-api/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, authHandler *handler.AuthHandler, zoneHandler *handler.ZoneHandler, reservationHandler *handler.ReservationHandler) {
	api := e.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Zone routes
	zones := api.Group("/zones")
	zones.GET("", zoneHandler.GetAllZones)
	zones.GET("/:id", zoneHandler.GetZoneByID)
	zones.POST("", zoneHandler.CreateZone, middleware.JWTMiddleware, middleware.RequireAdmin)

	// Reservation routes (all require authentication)
	reservations := api.Group("/reservations", middleware.JWTMiddleware)
	reservations.POST("", reservationHandler.Reserve)
	reservations.GET("/my-reservations", reservationHandler.GetMyReservations)
	reservations.DELETE("/:id", reservationHandler.Cancel)
	reservations.GET("", reservationHandler.GetAllReservations, middleware.RequireAdmin)
}