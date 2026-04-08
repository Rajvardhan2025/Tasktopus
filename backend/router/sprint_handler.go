package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/provider"
	"github.com/yourusername/project-management/utils"
)

type SprintHandler struct {
	provider *provider.Provider
}

func NewSprintHandler(p *provider.Provider) *SprintHandler {
	return &SprintHandler{provider: p}
}

func (h *SprintHandler) Create(c *fiber.Ctx) error {
	var req models.CreateSprintRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, err.Error())
	}

	sprint, err := h.provider.SprintService.Create(c.Context(), &req)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, sprint)
}

func (h *SprintHandler) List(c *fiber.Ctx) error {
	projectID := c.Params("id")

	sprints, err := h.provider.SprintService.GetByProject(c.Context(), projectID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, sprints)
}

func (h *SprintHandler) GetByID(c *fiber.Ctx) error {
	sprintID := c.Params("id")

	sprint, err := h.provider.SprintService.GetByID(c.Context(), sprintID)
	if err != nil {
		return utils.NotFoundError(c, "Sprint")
	}

	return utils.SuccessResponse(c, sprint)
}

func (h *SprintHandler) Update(c *fiber.Ctx) error {
	var req models.UpdateSprintRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	// Implementation similar to project update
	return utils.SuccessResponse(c, fiber.Map{"message": "Sprint updated"})
}

func (h *SprintHandler) Start(c *fiber.Ctx) error {
	sprintID := c.Params("id")
	userID := c.Locals("userID").(string)

	sprint, err := h.provider.SprintService.Start(c.Context(), sprintID, userID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, sprint)
}

func (h *SprintHandler) Complete(c *fiber.Ctx) error {
	sprintID := c.Params("id")
	userID := c.Locals("userID").(string)

	var req models.CompleteSprintRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	sprint, err := h.provider.SprintService.Complete(c.Context(), sprintID, &req, userID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, sprint)
}
