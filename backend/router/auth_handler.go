package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/provider"
	"github.com/yourusername/project-management/service"
	"github.com/yourusername/project-management/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(p *provider.Provider) *AuthHandler {
	return &AuthHandler{
		authService: p.AuthService,
	}
}

// Register creates a new user account
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, err.Error())
	}

	authResp, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		if err == service.ErrUserExists {
			return utils.ErrorResponse(c, fiber.StatusConflict, "USER_EXISTS", "User with this email already exists")
		}
		return utils.InternalError(c, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    authResp,
	})
}

// Login authenticates a user
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, err.Error())
	}

	authResp, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
		}
		if err == service.ErrUserInactive {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "USER_INACTIVE", "User account is inactive")
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, authResp)
}

// GetMe returns the current authenticated user
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	user, err := h.authService.GetCurrentUser(c.Context(), userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "USER_NOT_FOUND", "User not found")
		}
		if err == service.ErrUserInactive {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "USER_INACTIVE", "User account is inactive")
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, user)
}

// ChangePassword updates the user's password
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, err.Error())
	}

	if err := h.authService.ChangePassword(c.Context(), userID, &req); err != nil {
		if err == service.ErrInvalidCredentials {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "INVALID_PASSWORD", "Current password is incorrect")
		}
		if err == service.ErrUserNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "USER_NOT_FOUND", "User not found")
		}
		return utils.InternalError(c, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password changed successfully",
	})
}
