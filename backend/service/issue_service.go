package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/store"
	"go.mongodb.org/mongo-driver/bson"
)

type IssueService struct {
	issueStore      *store.IssueStore
	userStore       *store.UserStore
	projectStore    *store.ProjectStore
	activityStore   *store.ActivityStore
	workflowSvc     *WorkflowService
	notificationSvc *NotificationService
	wsSvc           *WebSocketService
}

func NewIssueService(
	issueStore *store.IssueStore,
	userStore *store.UserStore,
	projectStore *store.ProjectStore,
	activityStore *store.ActivityStore,
	workflowSvc *WorkflowService,
	notificationSvc *NotificationService,
	wsSvc *WebSocketService,
) *IssueService {
	return &IssueService{
		issueStore:      issueStore,
		userStore:       userStore,
		projectStore:    projectStore,
		activityStore:   activityStore,
		workflowSvc:     workflowSvc,
		notificationSvc: notificationSvc,
		wsSvc:           wsSvc,
	}
}

func (s *IssueService) Create(ctx context.Context, req *models.CreateIssueRequest, reporterID string) (*models.Issue, error) {
	log.Printf("[IssueService.Create] Starting - ProjectID: %s, ReporterID: %s", req.ProjectID, reporterID)

	// Validate project exists
	project, err := s.projectStore.FindByID(ctx, req.ProjectID)
	if err != nil {
		log.Printf("[IssueService.Create] Project not found error: %v", err)
		return nil, fmt.Errorf("project not found: %w", err)
	}
	log.Printf("[IssueService.Create] Project found: %s (Key: %s)", project.Name, project.Key)

	// Get workflow and default status
	workflow, err := s.workflowSvc.GetByProject(ctx, req.ProjectID)
	if err != nil {
		log.Printf("[IssueService.Create] Workflow fetch error: %v", err)
		return nil, fmt.Errorf("workflow not found: %w", err)
	}
	if len(workflow.Statuses) == 0 {
		log.Printf("[IssueService.Create] Workflow has no statuses")
		return nil, fmt.Errorf("invalid workflow configuration")
	}
	log.Printf("[IssueService.Create] Workflow found with %d statuses, default: %s", len(workflow.Statuses), workflow.Statuses[0])

	if err := s.validateCustomFields(req.CustomFields, project.CustomFields); err != nil {
		return nil, err
	}

	if err := s.validateAssignee(ctx, project, req.AssigneeID); err != nil {
		return nil, err
	}

	if err := s.validateParentRelationship(ctx, req.ProjectID, req.Type, req.ParentID); err != nil {
		return nil, err
	}

	// Generate issue key
	issueNum, err := s.issueStore.GetNextIssueNumber(ctx, project.Key)
	if err != nil {
		log.Printf("[IssueService.Create] Issue number generation error: %v", err)
		return nil, err
	}
	log.Printf("[IssueService.Create] Generated issue number: %d", issueNum)

	issue := &models.Issue{
		BaseModel:    models.BaseModel{ID: uuid.New().String()},
		IssueKey:     fmt.Sprintf("%s-%d", project.Key, issueNum),
		ProjectID:    req.ProjectID,
		Type:         req.Type,
		Title:        req.Title,
		Description:  req.Description,
		Status:       workflow.Statuses[0], // Default to first status
		Priority:     req.Priority,
		AssigneeID:   req.AssigneeID,
		ReporterID:   reporterID,
		SprintID:     req.SprintID,
		ParentID:     req.ParentID,
		Labels:       req.Labels,
		StoryPoints:  req.StoryPoints,
		CustomFields: req.CustomFields,
		Watchers:     []string{reporterID},
	}

	log.Printf("[IssueService.Create] Creating issue: %s", issue.IssueKey)
	if err := s.issueStore.Create(ctx, issue); err != nil {
		log.Printf("[IssueService.Create] Store create error: %v", err)
		return nil, err
	}

	// Log activity
	log.Printf("[IssueService.Create] Logging activity for issue: %s", issue.ID)
	s.logActivity(ctx, issue.ProjectID, issue.ID, reporterID, models.ActivityIssueCreated, nil)

	// Send notifications
	if issue.AssigneeID != "" && issue.AssigneeID != reporterID {
		log.Printf("[IssueService.Create] Sending assignment notification to: %s", issue.AssigneeID)
		s.notificationSvc.NotifyAssignment(ctx, issue.AssigneeID, issue.ID, reporterID)
	}

	// Broadcast WebSocket event
	log.Printf("[IssueService.Create] Broadcasting WebSocket event for project: %s", issue.ProjectID)
	s.wsSvc.BroadcastToProject(issue.ProjectID, models.WSEvent{
		Type:      models.WSEventIssueCreated,
		ProjectID: issue.ProjectID,
		Data:      issue,
	})

	log.Printf("[IssueService.Create] Success - Created issue: %s", issue.IssueKey)
	return issue, nil
}

func (s *IssueService) Update(ctx context.Context, issueID string, req *models.UpdateIssueRequest, userID string) (*models.Issue, error) {
	log.Printf("[IssueService.Update] Starting - IssueID: %s, UserID: %s", issueID, userID)

	issue, err := s.issueStore.FindByID(ctx, issueID)
	if err != nil {
		log.Printf("[IssueService.Update] Issue not found: %v", err)
		return nil, err
	}

	update := bson.M{}
	changes := make(map[string]interface{})

	project, err := s.projectStore.FindByID(ctx, issue.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if req.Title != nil {
		update["title"] = *req.Title
		changes["title"] = map[string]interface{}{"old": issue.Title, "new": *req.Title}
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}
	if req.Priority != nil {
		update["priority"] = *req.Priority
		changes["priority"] = map[string]interface{}{"old": issue.Priority, "new": *req.Priority}
	}
	if req.AssigneeID != nil {
		if err := s.validateAssignee(ctx, project, *req.AssigneeID); err != nil {
			return nil, err
		}

		update["assignee_id"] = *req.AssigneeID
		changes["assignee_id"] = map[string]interface{}{"old": issue.AssigneeID, "new": *req.AssigneeID}

		// Notify new assignee
		if *req.AssigneeID != "" && *req.AssigneeID != userID {
			log.Printf("[IssueService.Update] Notifying new assignee: %s", *req.AssigneeID)
			s.notificationSvc.NotifyAssignment(ctx, *req.AssigneeID, issueID, userID)
		}
	}
	if req.SprintID != nil {
		update["sprint_id"] = *req.SprintID
		changes["sprint_id"] = map[string]interface{}{"old": issue.SprintID, "new": *req.SprintID}
	}
	if req.Labels != nil {
		update["labels"] = *req.Labels
	}
	if req.StoryPoints != nil {
		update["story_points"] = *req.StoryPoints
	}
	if req.CustomFields != nil {
		if err := s.validateCustomFields(*req.CustomFields, project.CustomFields); err != nil {
			return nil, err
		}
		update["custom_fields"] = *req.CustomFields
	}

	log.Printf("[IssueService.Update] Applying %d updates", len(update))

	// Optimistic locking
	if err := s.issueStore.UpdateWithVersion(ctx, issueID, req.Version, update); err != nil {
		if !isVersionConflict(err) {
			log.Printf("[IssueService.Update] Update error: %v", err)
			return nil, err
		}

		latest, latestErr := s.issueStore.FindByID(ctx, issueID)
		if latestErr != nil {
			return nil, err
		}

		if retryErr := s.issueStore.UpdateWithVersion(ctx, issueID, latest.Version, update); retryErr != nil {
			log.Printf("[IssueService.Update] Retry update error: %v", retryErr)
			return nil, retryErr
		}
	}

	// Log activity
	if len(changes) > 0 {
		s.logActivity(ctx, issue.ProjectID, issueID, userID, models.ActivityIssueUpdated, changes)
	}

	// Get updated issue
	updatedIssue, _ := s.issueStore.FindByID(ctx, issueID)

	// Broadcast WebSocket event
	s.wsSvc.BroadcastToProject(issue.ProjectID, models.WSEvent{
		Type:      models.WSEventIssueUpdated,
		ProjectID: issue.ProjectID,
		Data:      updatedIssue,
	})

	log.Printf("[IssueService.Update] Success - Updated issue: %s", issueID)
	return updatedIssue, nil
}

func (s *IssueService) Transition(ctx context.Context, issueID string, req *models.TransitionRequest, userID string) (*models.Issue, error) {
	log.Printf("[IssueService.Transition] Starting - IssueID: %s, ToStatus: %s, UserID: %s", issueID, req.ToStatus, userID)

	issue, err := s.issueStore.FindByID(ctx, issueID)
	if err != nil {
		log.Printf("[IssueService.Transition] Issue not found: %v", err)
		return nil, err
	}

	log.Printf("[IssueService.Transition] Current status: %s", issue.Status)

	// Validate transition
	if err := s.workflowSvc.ValidateTransition(ctx, issue.ProjectID, issue, req.ToStatus); err != nil {
		log.Printf("[IssueService.Transition] Validation failed: %v", err)
		return nil, err
	}

	// Execute transition actions
	actionUpdates, err := s.workflowSvc.ExecuteActions(ctx, issue.ProjectID, issue.Status, req.ToStatus, userID, issue)
	if err != nil {
		log.Printf("[IssueService.Transition] Action execution failed: %v", err)
		return nil, err
	}

	// Update status
	update := bson.M{"status": req.ToStatus}
	for key, value := range actionUpdates {
		update[key] = value
	}
	if err := s.issueStore.UpdateWithVersion(ctx, issueID, req.Version, update); err != nil {
		log.Printf("[IssueService.Transition] Update failed: %v", err)
		return nil, err
	}

	// Log activity
	changes := map[string]interface{}{
		"status": map[string]interface{}{"old": issue.Status, "new": req.ToStatus},
	}
	s.logActivity(ctx, issue.ProjectID, issueID, userID, models.ActivityStatusChanged, changes)

	// Notify watchers
	log.Printf("[IssueService.Transition] Notifying %d watchers", len(issue.Watchers))
	for _, watcherID := range issue.Watchers {
		if watcherID != userID {
			s.notificationSvc.NotifyWatcher(ctx, watcherID, issueID, userID, "status changed")
		}
	}

	updatedIssue, _ := s.issueStore.FindByID(ctx, issueID)

	// Broadcast WebSocket event
	s.wsSvc.BroadcastToProject(issue.ProjectID, models.WSEvent{
		Type:      models.WSEventIssueMoved,
		ProjectID: issue.ProjectID,
		Data:      updatedIssue,
	})

	log.Printf("[IssueService.Transition] Success - Transitioned to: %s", req.ToStatus)
	return updatedIssue, nil
}

func (s *IssueService) GetByID(ctx context.Context, issueID string) (*models.Issue, error) {
	return s.issueStore.FindByID(ctx, issueID)
}

func (s *IssueService) GetByProject(ctx context.Context, projectID string) ([]*models.Issue, error) {
	return s.issueStore.FindByProject(ctx, projectID)
}

func (s *IssueService) Delete(ctx context.Context, issueID, userID string) error {
	issue, err := s.issueStore.FindByID(ctx, issueID)
	if err != nil {
		return err
	}

	if err := s.issueStore.Delete(ctx, issueID); err != nil {
		return err
	}

	s.logActivity(ctx, issue.ProjectID, issueID, userID, models.ActivityIssueDeleted, map[string]interface{}{
		"issue_key": issue.IssueKey,
	})

	s.wsSvc.BroadcastToProject(issue.ProjectID, models.WSEvent{
		Type:      models.WSEventIssueUpdated,
		ProjectID: issue.ProjectID,
		Data:      fiberLikeMap("deleted_issue_id", issueID),
	})

	return nil
}

func (s *IssueService) AddWatcher(ctx context.Context, issueID, userID string) error {
	issue, err := s.issueStore.FindByID(ctx, issueID)
	if err != nil {
		return err
	}

	if err := s.issueStore.AddWatcher(ctx, issueID, userID); err != nil {
		return err
	}

	s.logActivity(ctx, issue.ProjectID, issueID, userID, models.ActivityWatcherAdded, map[string]interface{}{
		"watcher_id": userID,
	})
	return nil
}

func (s *IssueService) RemoveWatcher(ctx context.Context, issueID, userID string) error {
	issue, err := s.issueStore.FindByID(ctx, issueID)
	if err != nil {
		return err
	}

	if err := s.issueStore.RemoveWatcher(ctx, issueID, userID); err != nil {
		return err
	}

	s.logActivity(ctx, issue.ProjectID, issueID, userID, models.ActivityWatcherRemoved, map[string]interface{}{
		"watcher_id": userID,
	})
	return nil
}

func (s *IssueService) logActivity(ctx context.Context, projectID, issueID, userID string, action models.ActivityAction, changes map[string]interface{}) {
	activity := &models.Activity{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		IssueID:   issueID,
		UserID:    userID,
		Action:    action,
		Changes:   changes,
	}
	s.activityStore.Create(ctx, activity)
}

func (s *IssueService) validateParentRelationship(ctx context.Context, projectID string, issueType models.IssueType, parentID string) error {
	if parentID == "" {
		if issueType == models.IssueTypeSubtask {
			return fmt.Errorf("subtask issues require a parent")
		}
		return nil
	}

	parent, err := s.issueStore.FindByID(ctx, parentID)
	if err != nil {
		return fmt.Errorf("parent issue not found: %w", err)
	}
	if parent.ProjectID != projectID {
		return fmt.Errorf("parent issue must belong to the same project")
	}

	allowedChildren := map[models.IssueType][]models.IssueType{
		models.IssueTypeEpic:    {models.IssueTypeStory, models.IssueTypeTask, models.IssueTypeBug},
		models.IssueTypeStory:   {models.IssueTypeSubtask},
		models.IssueTypeTask:    {models.IssueTypeSubtask},
		models.IssueTypeBug:     {models.IssueTypeSubtask},
		models.IssueTypeSubtask: {},
	}

	children := allowedChildren[parent.Type]
	for _, child := range children {
		if child == issueType {
			return nil
		}
	}

	return fmt.Errorf("invalid parent-child relationship: parent '%s' cannot contain child '%s'", parent.Type, issueType)
}

func (s *IssueService) validateCustomFields(values map[string]interface{}, defs []models.CustomField) error {
	if len(values) == 0 {
		return nil
	}

	defByName := map[string]models.CustomField{}
	for _, def := range defs {
		defByName[def.Name] = def
	}

	for key, value := range values {
		def, ok := defByName[key]
		if !ok {
			return fmt.Errorf("unknown custom field: %s", key)
		}

		switch def.Type {
		case "text":
			if _, ok := value.(string); !ok {
				return fmt.Errorf("custom field '%s' must be text", key)
			}
		case "number":
			switch value.(type) {
			case int, int32, int64, float32, float64:
			default:
				return fmt.Errorf("custom field '%s' must be a number", key)
			}
		case "dropdown":
			selected, ok := value.(string)
			if !ok {
				return fmt.Errorf("custom field '%s' must be a dropdown option", key)
			}
			allowed := false
			for _, option := range def.Options {
				if option == selected {
					allowed = true
					break
				}
			}
			if !allowed {
				return fmt.Errorf("custom field '%s' has invalid option '%s'", key, selected)
			}
		case "date":
			dateString, ok := value.(string)
			if !ok {
				return fmt.Errorf("custom field '%s' must be a date string", key)
			}
			if _, err := time.Parse("2006-01-02", dateString); err != nil {
				if _, err := time.Parse(time.RFC3339, dateString); err != nil {
					return fmt.Errorf("custom field '%s' must be ISO date (YYYY-MM-DD or RFC3339)", key)
				}
			}
		default:
			return fmt.Errorf("unsupported custom field type '%s'", def.Type)
		}
	}

	return nil
}

func fiberLikeMap(key, value string) map[string]string {
	return map[string]string{key: value}
}

func isVersionConflict(err error) bool {
	return err != nil && strings.Contains(err.Error(), "version conflict")
}

func (s *IssueService) validateAssignee(ctx context.Context, project *models.Project, assigneeID string) error {
	if assigneeID == "" {
		return nil
	}

	if _, err := s.userStore.FindByID(ctx, assigneeID); err != nil {
		return fmt.Errorf("assignee user not found")
	}

	isMember := false
	for _, memberID := range project.Members {
		if memberID == assigneeID {
			isMember = true
			break
		}
	}

	if !isMember {
		return fmt.Errorf("assignee must be a member of this project")
	}

	return nil
}
