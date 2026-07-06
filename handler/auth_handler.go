package handler

import (
	"net/http"

	"spotsync-api/dto"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// AuthHandler holds dependencies needed by auth-related HTTP handlers.
type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
}

// NewAuthHandler creates a new AuthHandler, injecting the service it needs.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest

	// 1. Bind incoming JSON to our struct
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// 2. Validate using the tags we defined in the DTO
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	// 3. Call the service to do the actual work
	user, err := h.authService.Register(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	// 4. Build a safe response (no password field)
	res := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    res,
	})
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest

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

	token, user, err := h.authService.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	res := dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse{
		Success: true,
		Message: "Login successful",
		Data:    res,
	})
}