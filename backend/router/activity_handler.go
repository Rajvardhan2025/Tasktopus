package router

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/project-management/provider"
	"github.com/yourusername/project-management/utils"
)

type ActivityHandler struct {
	provider *provider.Provider
}

func NewActivityHandler(p *provider.Provider) *ActivityHandler {
	return &ActivityHandler{provider: p}
}

func (h *ActivityHandler) GetProjectActivity(c *fiber.Ctx) error {
	projectID := c.Params("id")

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	cursor := c.Query("cursor", "")
	action := c.Query("action", "")
	userID := c.Query("user_id", "")

	activities, nextCursor, err := h.provider.ActivityStore.FindByProjectCursor(c.Context(), projectID, action, userID, limit, cursor)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"items":       activities,
		"next_cursor": nextCursor,
	})
}

func (h *ActivityHandler) GetIssueActivity(c *fiber.Ctx) error {
	issueID := c.Params("id")

	activities, err := h.provider.ActivityStore.FindByIssue(c.Context(), issueID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, activities)
}
