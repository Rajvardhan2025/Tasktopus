package router

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/project-management/provider"
	"github.com/yourusername/project-management/utils"
)

type SearchHandler struct {
	provider *provider.Provider
}

func NewSearchHandler(p *provider.Provider) *SearchHandler {
	return &SearchHandler{provider: p}
}

func (h *SearchHandler) Search(c *fiber.Ctx) error {
	queryStr := c.Query("q", "")

	query := h.provider.SearchService.ParseQueryString(queryStr)

	// Override with explicit query params
	if projectID := c.Query("project"); projectID != "" {
		query.ProjectID = projectID
	}
	if status := c.Query("status"); status != "" {
		query.Status = status
	}
	if assignee := c.Query("assignee"); assignee != "" {
		query.AssigneeID = assignee
	}
	if priority := c.Query("priority"); priority != "" {
		query.Priority = priority
	}
	if priorityGTE := c.Query("priority_gte"); priorityGTE != "" {
		query.PriorityGTE = priorityGTE
	}
	if cursor := c.Query("cursor"); cursor != "" {
		query.Cursor = cursor
	}

	if limit := c.Query("limit"); limit != "" {
		query.Limit, _ = strconv.Atoi(limit)
	}
	issues, nextCursor, err := h.provider.SearchService.Search(c.Context(), query)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"items":       issues,
		"next_cursor": nextCursor,
		"query":       query,
	})
}
