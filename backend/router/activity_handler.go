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
	skip, _ := strconv.Atoi(c.Query("skip", "0"))

	activities, err := h.provider.ActivityStore.FindByProject(c.Context(), projectID, limit, skip)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, activities)
}

func (h *ActivityHandler) GetIssueActivity(c *fiber.Ctx) error {
	issueID := c.Params("id")

	activities, err := h.provider.ActivityStore.FindByIssue(c.Context(), issueID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, activities)
}
