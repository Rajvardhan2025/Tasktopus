package service

import (
	"context"
	"strings"

	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/store"
	"go.mongodb.org/mongo-driver/bson"
)

type SearchService struct {
	issueStore *store.IssueStore
}

func NewSearchService(issueStore *store.IssueStore) *SearchService {
	return &SearchService{
		issueStore: issueStore,
	}
}

type SearchQuery struct {
	Text       string
	ProjectID  string
	Status     string
	AssigneeID string
	Priority   string
	Labels     []string
	Limit      int
	Skip       int
}

func (s *SearchService) Search(ctx context.Context, query SearchQuery) ([]*models.Issue, error) {
	filters := bson.M{}

	if query.ProjectID != "" {
		filters["project_id"] = query.ProjectID
	}

	if query.Status != "" {
		filters["status"] = query.Status
	}

	if query.AssigneeID != "" {
		filters["assignee_id"] = query.AssigneeID
	}

	if query.Priority != "" {
		filters["priority"] = query.Priority
	}

	if len(query.Labels) > 0 {
		filters["labels"] = bson.M{"$in": query.Labels}
	}

	if query.Limit == 0 {
		query.Limit = 50
	}

	return s.issueStore.Search(ctx, query.Text, filters, query.Limit, query.Skip)
}

func (s *SearchService) ParseQueryString(queryStr string) SearchQuery {
	query := SearchQuery{Limit: 50}
	parts := strings.Fields(queryStr)

	var textParts []string
	for _, part := range parts {
		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			key := strings.ToLower(kv[0])
			value := kv[1]

			switch key {
			case "status":
				query.Status = value
			case "assignee":
				query.AssigneeID = value
			case "priority":
				query.Priority = value
			case "project":
				query.ProjectID = value
			}
		} else {
			textParts = append(textParts, part)
		}
	}

	if len(textParts) > 0 {
		query.Text = strings.Join(textParts, " ")
	}

	return query
}
