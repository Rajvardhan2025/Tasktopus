package router

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/provider"
	"github.com/yourusername/project-management/utils"
)

type CommentHandler struct {
	provider *provider.Provider
}

func NewCommentHandler(p *provider.Provider) *CommentHandler {
	return &CommentHandler{provider: p}
}

func (h *CommentHandler) Create(c *fiber.Ctx) error {
	issueID := c.Params("id")
	userID := c.Locals("userID").(string)

	var req models.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, err.Error())
	}

	// Extract mentions
	mentions := extractMentions(req.Content)

	comment := &models.Comment{
		BaseModel: models.BaseModel{ID: uuid.New().String()},
		IssueID:   issueID,
		UserID:    userID,
		Content:   req.Content,
		Mentions:  mentions,
		ParentID:  req.ParentID,
	}

	if err := h.provider.CommentStore.Create(c.Context(), comment); err != nil {
		return utils.InternalError(c, err.Error())
	}

	// Notify mentioned users
	for _, mentionedUserID := range mentions {
		h.provider.NotificationService.NotifyMention(c.Context(), mentionedUserID, issueID, userID)
	}

	return utils.SuccessResponse(c, comment)
}

func (h *CommentHandler) List(c *fiber.Ctx) error {
	issueID := c.Params("id")

	comments, err := h.provider.CommentStore.FindByIssue(c.Context(), issueID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, comments)
}

func (h *CommentHandler) Update(c *fiber.Ctx) error {
	commentID := c.Params("id")

	var req models.UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := h.provider.CommentStore.Update(c.Context(), commentID, req.Content); err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Comment updated"})
}

func (h *CommentHandler) Delete(c *fiber.Ctx) error {
	commentID := c.Params("id")

	if err := h.provider.CommentStore.Delete(c.Context(), commentID); err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Comment deleted"})
}

func extractMentions(content string) []string {
	re := regexp.MustCompile(`@(\w+)`)
	matches := re.FindAllStringSubmatch(content, -1)

	mentions := []string{}
	for _, match := range matches {
		if len(match) > 1 {
			mentions = append(mentions, match[1])
		}
	}
	return mentions
}
