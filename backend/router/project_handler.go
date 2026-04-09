package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/provider"
	"github.com/yourusername/project-management/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectMemberRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type ProjectHandler struct {
	provider *provider.Provider
}

func NewProjectHandler(p *provider.Provider) *ProjectHandler {
	return &ProjectHandler{provider: p}
}

func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	log.Printf("[ProjectHandler.Create] Starting project creation")

	var req models.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("[ProjectHandler.Create] Body parse error: %v", err)
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		log.Printf("[ProjectHandler.Create] Validation error: %v", err)
		return utils.ValidationError(c, err.Error())
	}

	project := &models.Project{
		BaseModel:    models.BaseModel{ID: uuid.New().String()},
		Name:         req.Name,
		Key:          req.Key,
		Description:  req.Description,
		CustomFields: req.CustomFields,
		Members:      []string{},
	}

	if err := h.provider.ProjectStore.Create(c.Context(), project); err != nil {
		log.Printf("[ProjectHandler.Create] Project store error: %v", err)
		return utils.InternalError(c, err.Error())
	}

	log.Printf("[ProjectHandler.Create] Project created: %s (ID: %s)", project.Name, project.ID)

	// Create default workflow for the project
	workflow := &models.Workflow{
		BaseModel: models.BaseModel{ID: uuid.New().String()},
		Name:      "Standard Software Development Workflow",
		ProjectID: project.ID,
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
		log.Printf("[ProjectHandler.Create] Workflow creation error: %v", err)
		// Don't fail project creation if workflow fails, but log it
	} else {
		log.Printf("[ProjectHandler.Create] Default workflow created for project: %s", project.ID)
	}

	return utils.SuccessResponse(c, project)
}

func (h *ProjectHandler) List(c *fiber.Ctx) error {
	projects, err := h.provider.ProjectStore.List(c.Context())
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, projects)
}

func (h *ProjectHandler) GetByID(c *fiber.Ctx) error {
	projectID := c.Params("id")

	project, err := h.provider.ProjectStore.FindByID(c.Context(), projectID)
	if err != nil {
		return utils.NotFoundError(c, "Project")
	}

	return utils.SuccessResponse(c, project)
}

func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	projectID := c.Params("id")

	var req models.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	update := bson.M{}
	if req.Name != nil {
		update["name"] = *req.Name
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}
	if req.CustomFields != nil {
		update["custom_fields"] = *req.CustomFields
	}

	if err := h.provider.ProjectStore.Update(c.Context(), projectID, update); err != nil {
		return utils.InternalError(c, err.Error())
	}

	project, _ := h.provider.ProjectStore.FindByID(c.Context(), projectID)
	return utils.SuccessResponse(c, project)
}

func (h *ProjectHandler) Delete(c *fiber.Ctx) error {
	projectID := c.Params("id")

	if err := h.provider.ProjectStore.Delete(c.Context(), projectID); err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Project deleted"})
}

func (h *ProjectHandler) GetMembers(c *fiber.Ctx) error {
	projectID := c.Params("id")

	project, err := h.provider.ProjectStore.FindByID(c.Context(), projectID)
	if err != nil {
		return utils.NotFoundError(c, "Project")
	}

	if len(project.Members) == 0 {
		return utils.SuccessResponse(c, []models.User{})
	}

	users, err := h.provider.UserStore.FindByIDs(c.Context(), project.Members)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, users)
}

func (h *ProjectHandler) AddMember(c *fiber.Ctx) error {
	projectID := c.Params("id")

	var req ProjectMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationError(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationError(c, err.Error())
	}

	if _, err := h.provider.ProjectStore.FindByID(c.Context(), projectID); err != nil {
		return utils.NotFoundError(c, "Project")
	}

	if _, err := h.provider.UserStore.FindByID(c.Context(), req.UserID); err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.NotFoundError(c, "User")
		}
		return utils.InternalError(c, err.Error())
	}

	if err := h.provider.ProjectStore.AddMember(c.Context(), projectID, req.UserID); err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Member added"})
}

func (h *ProjectHandler) RemoveMember(c *fiber.Ctx) error {
	projectID := c.Params("id")
	userID := c.Params("userId")

	if userID == "" {
		return utils.ValidationError(c, "User ID is required")
	}

	if err := h.provider.ProjectStore.RemoveMember(c.Context(), projectID, userID); err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Member removed"})
}
