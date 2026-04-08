package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/provider"
	"github.com/yourusername/project-management/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type IssueHandler struct {
	provider *provider.Provider
}

func NewIssueHandler(p *provider.Provider) *IssueHandler {
	return &IssueHandler{provider: p}
}

func (h *IssueHandler) Create(c *fiber.Ctx) error {
	projectID := c.Params("id")
	userID := c.Locals("userID").(string) // From auth middleware

	log.Printf("[IssueHandler.Create] Starting - ProjectID: %s, UserID: %s", projectID, userID)

	var req models.CreateIssueRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[IssueHandler.Create] Body parse error: %v", err)
		return utils.ValidationError(c, "Invalid request body")
	}

	log.Printf("[IssueHandler.Create] Request parsed: %+v", req)

	req.ProjectID = projectID

	if err := utils.ValidateStruct(&req); err != nil {
		log.Printf("[IssueHandler.Create] Validation error: %v", err)
		return utils.ValidationError(c, err.Error())
	}

	issue, err := h.provider.IssueService.Create(c.Context(), &req, userID)
	if err != nil {
		log.Printf("[IssueHandler.Create] Service error: %v", err)
		return utils.InternalError(c, err.Error())
	}

	log.Printf("[IssueHandler.Create] Success - IssueID: %s, IssueKey: %s", issue.ID, issue.IssueKey)
	return utils.SuccessResponse(c, issue)
}

func (h *IssueHandler) GetByID(c *fiber.Ctx) error {
	issueID := c.Params("id")

	issue, err := h.provider.IssueService.GetByID(c.Context(), issueID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.NotFoundError(c, "Issue")
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, issue)
}

func (h *IssueHandler) GetBoard(c *fiber.Ctx) error {
	projectID := c.Params("id")

	issues, err := h.provider.IssueService.GetByProject(c.Context(), projectID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"issues": issues,
	})
}

func (h *IssueHandler) Update(c *fiber.Ctx) error {
	issueID := c.Params("id")
	userID := c.Locals("userID").(string)

	var req models.UpdateIssueRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, err.Error())
	}

	issue, err := h.provider.IssueService.Update(c.Context(), issueID, &req, userID)
	if err != nil {
		if err.Error() == "version conflict: issue has been modified" {
			return utils.ConflictError(c, "Issue has been modified by another user")
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, issue)
}

func (h *IssueHandler) Transition(c *fiber.Ctx) error {
	issueID := c.Params("id")
	userID := c.Locals("userID").(string)

	var req models.TransitionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, err.Error())
	}

	issue, err := h.provider.IssueService.Transition(c.Context(), issueID, &req, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, "TRANSITION_ERROR", err.Error())
	}

	return utils.SuccessResponse(c, issue)
}

func (h *IssueHandler) Delete(c *fiber.Ctx) error {
	issueID := c.Params("id")

	if err := h.provider.IssueStore.Delete(c.Context(), issueID); err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Issue deleted"})
}

func (h *IssueHandler) AddWatcher(c *fiber.Ctx) error {
	issueID := c.Params("id")
	userID := c.Locals("userID").(string)

	if err := h.provider.IssueStore.AddWatcher(c.Context(), issueID, userID); err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Watcher added"})
}

func (h *IssueHandler) RemoveWatcher(c *fiber.Ctx) error {
	issueID := c.Params("id")
	userID := c.Locals("userID").(string)

	if err := h.provider.IssueStore.RemoveWatcher(c.Context(), issueID, userID); err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Watcher removed"})
}
