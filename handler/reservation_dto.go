package handler

import (
	"net/http"
	"strconv"

	"spotsync-api/dto"
	"spotsync-api/repository"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ReservationHandler struct {
	reservationService service.ReservationService
	validate           *validator.Validate
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
		validate:           validator.New(),
	}
}

// Reserve handles POST /api/v1/reservations
func (h *ReservationHandler) Reserve(c echo.Context) error {
	var req dto.CreateReservationRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	// user_id was injected into context by JWTMiddleware.
	// It comes from jwt.MapClaims as float64 (JSON numbers decode to float64 in Go),
	// so we must convert it carefully.
	userIDFloat := c.Get("user_id").(float64)
	userID := uint(userIDFloat)

	reservation, err := h.reservationService.Reserve(userID, req)
	if err != nil {
		if err == repository.ErrZoneFull {
			return c.JSON(http.StatusConflict, utils.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    reservation,
	})
}

// GetMyReservations handles GET /api/v1/reservations/my-reservations
func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userIDFloat := c.Get("user_id").(float64)
	userID := uint(userIDFloat)

	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Success: false,
			Message: "Failed to fetch reservations",
		})
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    reservations,
	})
}

// GetAllReservations handles GET /api/v1/reservations (admin only)
func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	reservations, err := h.reservationService.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Success: false,
			Message: "Failed to fetch reservations",
		})
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "All reservations retrieved successfully",
		Data:    reservations,
	})
}

// Cancel handles DELETE /api/v1/reservations/:id
func (h *ReservationHandler) Cancel(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Success: false,
			Message: "Invalid reservation id",
		})
	}

	userIDFloat := c.Get("user_id").(float64)
	requesterID := uint(userIDFloat)
	requesterRole := c.Get("role").(string)

	err = h.reservationService.Cancel(uint(id), requesterID, requesterRole)
	if err != nil {
		if err == service.ErrForbidden {
			return c.JSON(http.StatusForbidden, utils.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
		}
		if err == service.ErrNotFound {
			return c.JSON(http.StatusNotFound, utils.ErrorResponse{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Success: false,
			Message: "Failed to cancel reservation",
		})
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}