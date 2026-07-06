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
		return utils.NewAppError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		return utils.NewAppError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	zone, err := h.zoneService.CreateZone(req)
	if err != nil {
		return err // will be logged + shown as generic 500 by central handler
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
		return err
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    zones,
	})
}

// GetZoneByID handles GET /api/v1/zones/:id (public)
func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return utils.NewAppError(http.StatusBadRequest, "Invalid zone id")
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		// If it's gorm.ErrRecordNotFound, CentralErrorHandler turns it into 404 automatically
		return err
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    zone,
	})
}