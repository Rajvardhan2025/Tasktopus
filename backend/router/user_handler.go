package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/provider"
	"github.com/yourusername/project-management/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	provider *provider.Provider
}

func NewUserHandler(p *provider.Provider) *UserHandler {
	return &UserHandler{provider: p}
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	users, err := h.provider.UserStore.List(c.Context())
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, users)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, err.Error())
	}

	existingUser, err := h.provider.UserStore.FindByEmail(c.Context(), req.Email)
	if err == nil && existingUser != nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, "USER_EXISTS", "User with this email already exists")
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return utils.InternalError(c, err.Error())
	}

	user := &models.User{
		BaseModel:   models.BaseModel{ID: uuid.New().String()},
		Email:       req.Email,
		DisplayName: req.DisplayName,
		IsActive:    true,
	}

	if err := h.provider.UserStore.Create(c.Context(), user); err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, user)
}
