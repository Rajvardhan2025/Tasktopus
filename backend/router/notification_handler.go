package router

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/project-management/provider"
	"github.com/yourusername/project-management/utils"
)

type NotificationHandler struct {
	provider *provider.Provider
}

func NewNotificationHandler(p *provider.Provider) *NotificationHandler {
	return &NotificationHandler{provider: p}
}

func (h *NotificationHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	notifications, err := h.provider.NotificationService.GetUserNotifications(c.Context(), userID, limit)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, notifications)
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	if notificationID == "" {
		return utils.ValidationError(c, "Notification ID is required")
	}

	if err := h.provider.NotificationService.MarkAsRead(c.Context(), notificationID); err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	if err := h.provider.NotificationService.MarkAllAsRead(c.Context(), userID); err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "All notifications marked as read"})
}
