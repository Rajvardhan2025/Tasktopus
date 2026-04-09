package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/provider"
	"github.com/yourusername/project-management/utils"
)

type WorkflowHandler struct {
	provider *provider.Provider
}

func NewWorkflowHandler(p *provider.Provider) *WorkflowHandler {
	return &WorkflowHandler{provider: p}
}

func (h *WorkflowHandler) GetByProject(c *fiber.Ctx) error {
	projectID := c.Params("id")

	workflow, err := h.provider.WorkflowStore.FindByProject(c.Context(), projectID)
	if err != nil {
		return utils.NotFoundError(c, "Workflow")
	}

	return utils.SuccessResponse(c, workflow)
}

func (h *WorkflowHandler) CreateDefault(c *fiber.Ctx) error {
	projectID := c.Params("id")

	log.Printf("[WorkflowHandler.CreateDefault] Creating default workflow for project: %s", projectID)

	// Check if project exists
	_, err := h.provider.ProjectStore.FindByID(c.Context(), projectID)
	if err != nil {
		return utils.NotFoundError(c, "Project")
	}

	// Check if workflow already exists
	existing, _ := h.provider.WorkflowStore.FindByProject(c.Context(), projectID)
	if existing != nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, "WORKFLOW_EXISTS", "Workflow already exists for this project")
	}

	// Create default workflow
	workflow := &models.Workflow{
		BaseModel: models.BaseModel{ID: uuid.New().String()},
		Name:      "Standard Software Development Workflow",
		ProjectID: projectID,
		Statuses:  []string{"to_do", "in_progress", "in_review", "done"},
		Transitions: []models.Transition{
			{
				From:       "to_do",
				To:         "in_progress",
				Conditions: []models.Condition{},
				Actions:    []models.Action{},
			},
			{
				From: "in_progress",
				To:   "in_review",
				Conditions: []models.Condition{
					{
						Field:    "assignee_id",
						Operator: "not_empty",
						Value:    "",
					},
				},
				Actions: []models.Action{
					{
						Type: "assign_reviewer",
					},
					{
						Type: "notify",
						Params: map[string]interface{}{
							"message": "Issue moved to review",
						},
					},
				},
			},
			{
				From:       "in_review",
				To:         "done",
				Conditions: []models.Condition{},
				Actions:    []models.Action{},
			},
			{
				From:       "in_review",
				To:         "in_progress",
				Conditions: []models.Condition{},
				Actions:    []models.Action{},
			},
		},
	}

	if err := h.provider.WorkflowStore.Create(c.Context(), workflow); err != nil {
		log.Printf("[WorkflowHandler.CreateDefault] Error: %v", err)
		return utils.InternalError(c, err.Error())
	}

	log.Printf("[WorkflowHandler.CreateDefault] Success - Workflow ID: %s", workflow.ID)
	return utils.SuccessResponse(c, workflow)
}
