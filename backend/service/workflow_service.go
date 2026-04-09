package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/store"
	"go.mongodb.org/mongo-driver/bson"
)

type WorkflowService struct {
	workflowStore   *store.WorkflowStore
	issueStore      *store.IssueStore
	notificationSvc *NotificationService
}

type TransitionValidationError struct {
	From    string
	To      string
	Allowed []string
}

func (e *TransitionValidationError) Error() string {
	if len(e.Allowed) == 0 {
		return fmt.Sprintf("transition from '%s' to '%s' is not allowed", e.From, e.To)
	}
	return fmt.Sprintf("transition from '%s' to '%s' is not allowed; allowed transitions: %s", e.From, e.To, strings.Join(e.Allowed, ", "))
}

func NewWorkflowService(workflowStore *store.WorkflowStore, issueStore *store.IssueStore, notificationSvc *NotificationService) *WorkflowService {
	return &WorkflowService{
		workflowStore:   workflowStore,
		issueStore:      issueStore,
		notificationSvc: notificationSvc,
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

	allowed := make([]string, 0)
	var transition *models.Transition
	for _, t := range workflow.Transitions {
		if t.From == issue.Status {
			allowed = append(allowed, t.To)
		}
		if t.From == issue.Status && t.To == toStatus {
			transition = &t
			break
		}
	}

	if transition == nil {
		return &TransitionValidationError{From: issue.Status, To: toStatus, Allowed: allowed}
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

func (s *WorkflowService) ExecuteActions(ctx context.Context, projectID, fromStatus, toStatus, actorID string, issue *models.Issue) (bson.M, error) {
	workflow, err := s.workflowStore.FindByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	updates := bson.M{}

	// Find transition
	for _, t := range workflow.Transitions {
		if t.From == fromStatus && t.To == toStatus {
			for _, action := range t.Actions {
				actionUpdates, err := s.executeAction(ctx, action, issue, actorID)
				if err != nil {
					return nil, err
				}
				for key, value := range actionUpdates {
					updates[key] = value
				}
			}
			break
		}
	}

	return updates, nil
}

func (s *WorkflowService) executeAction(ctx context.Context, action models.Action, issue *models.Issue, actorID string) (bson.M, error) {
	updates := bson.M{}

	switch action.Type {
	case "assign_reviewer":
		reviewerID, _ := action.Params["reviewer_id"].(string)
		if reviewerID == "" {
			reviewerID = issue.ReporterID
		}
		if reviewerID != "" {
			updates["assignee_id"] = reviewerID
			issue.AssigneeID = reviewerID
			if reviewerID != actorID && s.notificationSvc != nil {
				_ = s.notificationSvc.NotifyAssignment(ctx, reviewerID, issue.ID, actorID)
			}
		}
	case "notify":
		if s.notificationSvc == nil {
			break
		}
		for _, watcherID := range issue.Watchers {
			if watcherID != actorID {
				_ = s.notificationSvc.NotifyWatcher(ctx, watcherID, issue.ID, actorID, "workflow transition")
			}
		}
	case "set_field":
		fieldName, _ := action.Params["field"].(string)
		if fieldName == "" {
			return nil, fmt.Errorf("set_field action requires 'field' param")
		}
		value, ok := action.Params["value"]
		if !ok {
			return nil, fmt.Errorf("set_field action requires 'value' param")
		}
		updates[fieldName] = value
	default:
		return nil, fmt.Errorf("unsupported workflow action: %s", action.Type)
	}

	return updates, nil
}
