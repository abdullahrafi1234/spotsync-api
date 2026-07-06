package utils

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CentralErrorHandler replaces Echo's default error handler.
// Every error returned from any handler ends up here, in ONE place.
func CentralErrorHandler(err error, c echo.Context) {
	// Avoid double-writing a response if headers are already sent
	if c.Response().Committed {
		return
	}

	var appErr *AppError
	var echoErr *echo.HTTPError

	switch {
	// 1. Our own custom AppError (e.g. from services)
	case errors.As(err, &appErr):
		c.JSON(appErr.Code, ErrorResponse{
			Success: false,
			Message: appErr.Message,
		})

	// 2. GORM "record not found" -> always map to 404, never leak raw GORM text
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Message: "Resource not found",
		})

	// 3. Echo's own HTTP errors (e.g. wrong method, malformed request)
	case errors.As(err, &echoErr):
		c.JSON(echoErr.Code, ErrorResponse{
			Success: false,
			Message: echoErr.Message.(string),
		})

	// 4. Anything else: log the real error internally, but show a generic message
	default:
		log.Println("Unhandled error:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Something went wrong, please try again later",
		})
	}
}