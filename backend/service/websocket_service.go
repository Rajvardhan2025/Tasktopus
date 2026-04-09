package service

import (
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/yourusername/project-management/models"
)

type Client struct {
	ID        string
	Conn      *websocket.Conn
	ProjectID string
	UserID    string
	Send      chan models.WSEvent
}

type WebSocketService struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan models.WSEvent
	mu         sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	svc := &WebSocketService{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan models.WSEvent, 256),
	}

	go svc.run()

	return svc
}

func (s *WebSocketService) run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client.ID] = client
			s.mu.Unlock()

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client.ID]; ok {
				delete(s.clients, client.ID)
				close(client.Send)
			}
			s.mu.Unlock()

		case event := <-s.broadcast:
			s.mu.RLock()
			for _, client := range s.clients {
				if client.ProjectID == event.ProjectID {
					select {
					case client.Send <- event:
					default:
						close(client.Send)
						delete(s.clients, client.ID)
					}
				}
			}
			s.mu.RUnlock()
		}
	}
}

func (s *WebSocketService) RegisterClient(conn *websocket.Conn, projectID, userID string) *Client {
	client := &Client{
		ID:        uuid.New().String(),
		Conn:      conn,
		ProjectID: projectID,
		UserID:    userID,
		Send:      make(chan models.WSEvent, 256),
	}

	s.register <- client

	// Send presence event
	s.BroadcastToProject(projectID, models.WSEvent{
		Type:      models.WSEventPresence,
		ProjectID: projectID,
		Data: models.PresenceData{
			UserID: userID,
			Action: "joined",
		},
		Timestamp: time.Now().Unix(),
	})

	return client
}

func (s *WebSocketService) UnregisterClient(client *Client) {
	s.unregister <- client

	// Send presence event
	s.BroadcastToProject(client.ProjectID, models.WSEvent{
		Type:      models.WSEventPresence,
		ProjectID: client.ProjectID,
		Data: models.PresenceData{
			UserID: client.UserID,
			Action: "left",
		},
		Timestamp: time.Now().Unix(),
	})
}

func (s *WebSocketService) BroadcastToProject(projectID string, event models.WSEvent) {
	event.Timestamp = time.Now().Unix()
	s.broadcast <- event
}

func (s *WebSocketService) GetConnectedUsers(projectID string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := []string{}
	for _, client := range s.clients {
		if client.ProjectID == projectID {
			users = append(users, client.UserID)
		}
	}
	return users
}
