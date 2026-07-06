package handler

import (
	"net/http"

	"spotsync-api/dto"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return utils.NewAppError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		return utils.NewAppError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	user, err := h.authService.Register(req)
	if err != nil {
		// Business errors from the service (e.g. "email already registered")
		return utils.NewAppError(http.StatusBadRequest, err.Error())
	}

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
		return utils.NewAppError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		return utils.NewAppError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	token, user, err := h.authService.Login(req)
	if err != nil {
		return utils.NewAppError(http.StatusUnauthorized, err.Error())
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