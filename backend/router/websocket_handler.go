package router

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/provider"
)

type WebSocketHandler struct {
	provider *provider.Provider
}

func NewWebSocketHandler(p *provider.Provider) *WebSocketHandler {
	return &WebSocketHandler{provider: p}
}

func (h *WebSocketHandler) HandleWebSocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		// Extract params BEFORE WebSocket upgrade
		projectID := c.Params("projectId")
		userID := c.Query("userId", "anonymous")
		sinceRaw := c.Query("since", "")

		return websocket.New(func(conn *websocket.Conn) {

			client := h.provider.WebSocketService.RegisterClient(conn, projectID, userID)
			defer h.provider.WebSocketService.UnregisterClient(client)

			if sinceRaw != "" {
				sinceUnix, err := strconv.ParseInt(sinceRaw, 10, 64)
				if err == nil {
					since := time.Unix(sinceUnix, 0)
					activities, replayErr := h.provider.ActivityStore.FindAfterTimestamp(context.Background(), projectID, since)
					if replayErr == nil {
						for _, activity := range activities {
							eventType := mapActivityToEventType(activity.Action)
							if eventType == "" {
								continue
							}
							if err := conn.WriteJSON(map[string]interface{}{
								"type":       eventType,
								"project_id": projectID,
								"data":       activity,
								"timestamp":  activity.Timestamp.Unix(),
							}); err != nil {
								break
							}
						}
					}
				}
			}

			// Read messages (for presence updates, etc.)
			go func() {
				for {
					_, _, err := conn.ReadMessage()
					if err != nil {
						break
					}
				}
			}()

			// Write messages
			for event := range client.Send {
				if err := conn.WriteJSON(event); err != nil {
					log.Printf("WebSocket write error: %v", err)
					break
				}
			}
		})(c)
	}

	return fiber.ErrUpgradeRequired
}

func mapActivityToEventType(action models.ActivityAction) models.WSEventType {
	switch action {
	case models.ActivityIssueCreated:
		return models.WSEventIssueCreated
	case models.ActivityIssueUpdated, models.ActivityIssueDeleted:
		return models.WSEventIssueUpdated
	case models.ActivityStatusChanged:
		return models.WSEventIssueMoved
	case models.ActivityCommentAdded, models.ActivityCommentUpdated, models.ActivityCommentDeleted:
		return models.WSEventCommentAdded
	case models.ActivitySprintStarted, models.ActivitySprintCompleted:
		return models.WSEventSprintUpdated
	default:
		return ""
	}
}
