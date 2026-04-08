package service

import (
	"context"
	"fmt"

	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/store"
)

type WorkflowService struct {
	workflowStore *store.WorkflowStore
}

func NewWorkflowService(workflowStore *store.WorkflowStore) *WorkflowService {
	return &WorkflowService{
		workflowStore: workflowStore,
	}
}

func (s *WorkflowService) GetByProject(ctx context.Context, projectID string) (*models.Workflow, error) {
	return s.workflowStore.FindByProject(ctx, projectID)
}

func (s *WorkflowService) ValidateTransition(ctx context.Context, projectID string, issue *models.Issue, toStatus string) error {
	workflow, err := s.workflowStore.FindByProject(ctx, projectID)
	if err != nil {
		return fmt.Errorf("workflow not found: %w", err)
	}

	// Find transition
	var transition *models.Transition
	for _, t := range workflow.Transitions {
		if t.From == issue.Status && t.To == toStatus {
			transition = &t
			break
		}
	}

	if transition == nil {
		return fmt.Errorf("transition from '%s' to '%s' is not allowed", issue.Status, toStatus)
	}

	// Validate conditions
	for _, condition := range transition.Conditions {
		if err := s.validateCondition(issue, condition); err != nil {
			return err
		}
	}

	return nil
}

func (s *WorkflowService) validateCondition(issue *models.Issue, condition models.Condition) error {
	switch condition.Field {
	case "assignee_id":
		if condition.Operator == "not_empty" && issue.AssigneeID == "" {
			return fmt.Errorf("assignee is required for this transition")
		}
	case "story_points":
		if condition.Operator == "not_empty" && issue.StoryPoints == 0 {
			return fmt.Errorf("story points are required for this transition")
		}
	}
	return nil
}

func (s *WorkflowService) ExecuteActions(ctx context.Context, projectID, fromStatus, toStatus string, issue *models.Issue) error {
	workflow, err := s.workflowStore.FindByProject(ctx, projectID)
	if err != nil {
		return err
	}

	// Find transition
	for _, t := range workflow.Transitions {
		if t.From == fromStatus && t.To == toStatus {
			for _, action := range t.Actions {
				if err := s.executeAction(ctx, action, issue); err != nil {
					return err
				}
			}
			break
		}
	}

	return nil
}

func (s *WorkflowService) executeAction(ctx context.Context, action models.Action, issue *models.Issue) error {
	switch action.Type {
	case "assign_reviewer":
		// Implementation would assign a reviewer based on params
		// For now, this is a placeholder
	case "notify":
		// Send notification
	case "set_field":
		// Set a field value
	}
	return nil
}
