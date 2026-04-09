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

	issue, err := h.provider.IssueStore.FindByID(c.Context(), issueID)
	if err == nil {
		_ = h.provider.ActivityStore.Create(c.Context(), &models.Activity{
			ID:        uuid.New().String(),
			ProjectID: issue.ProjectID,
			IssueID:   issueID,
			UserID:    userID,
			Action:    models.ActivityCommentAdded,
			Changes: map[string]interface{}{
				"comment_id": comment.ID,
			},
		})

		h.provider.WebSocketService.BroadcastToProject(issue.ProjectID, models.WSEvent{
			Type:      models.WSEventCommentAdded,
			ProjectID: issue.ProjectID,
			Data:      comment,
		})
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
	userID := c.Locals("userID").(string)

	var req models.UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := h.provider.CommentStore.Update(c.Context(), commentID, req.Content); err != nil {
		return utils.InternalError(c, err.Error())
	}

	comment, err := h.provider.CommentStore.FindByID(c.Context(), commentID)
	if err == nil {
		issue, issueErr := h.provider.IssueStore.FindByID(c.Context(), comment.IssueID)
		if issueErr == nil {
			_ = h.provider.ActivityStore.Create(c.Context(), &models.Activity{
				ID:        uuid.New().String(),
				ProjectID: issue.ProjectID,
				IssueID:   comment.IssueID,
				UserID:    userID,
				Action:    models.ActivityCommentUpdated,
				Changes: map[string]interface{}{
					"comment_id": commentID,
				},
			})
		}
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Comment updated"})
}

func (h *CommentHandler) Delete(c *fiber.Ctx) error {
	commentID := c.Params("id")
	userID := c.Locals("userID").(string)

	comment, _ := h.provider.CommentStore.FindByID(c.Context(), commentID)

	if err := h.provider.CommentStore.Delete(c.Context(), commentID); err != nil {
		return utils.InternalError(c, err.Error())
	}

	if comment != nil {
		issue, issueErr := h.provider.IssueStore.FindByID(c.Context(), comment.IssueID)
		if issueErr == nil {
			_ = h.provider.ActivityStore.Create(c.Context(), &models.Activity{
				ID:        uuid.New().String(),
				ProjectID: issue.ProjectID,
				IssueID:   comment.IssueID,
				UserID:    userID,
				Action:    models.ActivityCommentDeleted,
				Changes: map[string]interface{}{
					"comment_id": commentID,
				},
			})
		}
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
