package handler

import (
	"net/http"
	"strconv"

	"spotsync-api/dto"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ZoneHandler struct {
	zoneService service.ZoneService
	validate    *validator.Validate
}

func NewZoneHandler(zoneService service.ZoneService) *ZoneHandler {
	return &ZoneHandler{
		zoneService: zoneService,
		validate:    validator.New(),
	}
}

// CreateZone handles POST /api/v1/zones (admin only)
func (h *ZoneHandler) CreateZone(c echo.Context) error {
	var req dto.CreateZoneRequest

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

	zone, err := h.zoneService.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Success: false,
			Message: "Failed to create zone",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    zone,
	})
}

// GetAllZones handles GET /api/v1/zones (public)
func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	zones, err := h.zoneService.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Success: false,
			Message: "Failed to fetch zones",
		})
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    zones,
	})
}

// GetZoneByID handles GET /api/v1/zones/:id (public)
func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
	// URL param is always a string, so we must convert it to uint
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Success: false,
			Message: "Invalid zone id",
		})
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse{
			Success: false,
			Message: "Zone not found",
		})
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    zone,
	})
}