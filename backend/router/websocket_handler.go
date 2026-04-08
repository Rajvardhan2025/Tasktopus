package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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

		return websocket.New(func(conn *websocket.Conn) {

			client := h.provider.WebSocketService.RegisterClient(conn, projectID, userID)
			defer h.provider.WebSocketService.UnregisterClient(client)

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
