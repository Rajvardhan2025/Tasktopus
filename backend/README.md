# Project Management Platform - Backend

A scalable, production-ready project management backend built with Go, Fiber, and MongoDB. Supports real-time collaboration, flexible workflows, and high concurrency.

## Features

- **Project & Issue Management**: Full CRUD operations with parent-child relationships
- **Flexible Workflow Engine**: Configurable status transitions with validation rules
- **Sprint Management**: Plan, start, and complete sprints with velocity tracking
- **Real-Time Updates**: WebSocket support for live board synchronization
- **Activity Logging**: Complete audit trail for all changes
- **Search & Filtering**: Full-text search with structured queries
- **Notifications**: In-app notifications for assignments, mentions, and watchers
- **Optimistic Locking**: Version-based concurrency control
- **Scalable Architecture**: Stateless design with MongoDB clustering support

## Architecture

```
backend/
├── main.go              # Application entry point
├── config/              # Configuration management
├── provider/            # Dependency injection
├── store/               # Data access layer (MongoDB)
├── service/             # Business logic layer
├── router/              # HTTP handlers
├── models/              # Domain models
├── middleware/          # HTTP middleware
└── utils/               # Utilities
```

See [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed design decisions.

## Tech Stack

- **Go 1.21+**: High-performance runtime
- **Fiber v2**: Fast HTTP framework
- **MongoDB**: Flexible document database
- **WebSocket**: Real-time communication
- **JWT**: Authentication (placeholder implementation)

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MongoDB 5.0 or higher
- Make (optional)

### Installation

1. Clone the repository
```bash
git clone <repository-url>
cd backend
```

2. Install dependencies
```bash
go mod download
```

3. Configure environment
```bash
cp .env.example .env
# Edit .env with your MongoDB connection string
```

4. Run the server
```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Using Make

```bash
# Run the server
make run

# Build binary
make build

# Run tests
make test

# Clean build artifacts
make clean
```

## API Endpoints

### Projects

- `POST /api/projects` - Create project
- `GET /api/projects` - List projects
- `GET /api/projects/:id` - Get project
- `PATCH /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Delete project

### Issues

- `POST /api/projects/:id/issues` - Create issue
- `GET /api/projects/:id/board` - Get board state
- `GET /api/issues/:id` - Get issue
- `PATCH /api/issues/:id` - Update issue
- `POST /api/issues/:id/transitions` - Transition status
- `POST /api/issues/:id/watch` - Watch issue
- `DELETE /api/issues/:id/watch` - Unwatch issue

### Sprints

- `POST /api/sprints` - Create sprint
- `GET /api/projects/:id/sprints` - List sprints
- `GET /api/sprints/:id` - Get sprint
- `POST /api/sprints/:id/start` - Start sprint
- `POST /api/sprints/:id/complete` - Complete sprint

### Comments

- `GET /api/issues/:id/comments` - List comments
- `POST /api/issues/:id/comments` - Add comment
- `PATCH /api/comments/:id` - Update comment
- `DELETE /api/comments/:id` - Delete comment

### Activity

- `GET /api/projects/:id/activity` - Project activity feed
- `GET /api/issues/:id/activity` - Issue activity history

### Search

- `GET /api/search?q=...` - Search issues

### WebSocket

- `WS /api/ws/:projectId?userId=...` - Real-time updates

## Example Requests

### Create Project

```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Backend Services",
    "key": "BACKEND",
    "description": "Core backend infrastructure"
  }'
```

### Create Issue

```bash
curl -X POST http://localhost:8080/api/projects/{projectId}/issues \
  -H "Content-Type: application/json" \
  -d '{
    "type": "story",
    "title": "Add OAuth authentication",
    "description": "Implement OAuth 2.0 flow",
    "priority": "high",
    "story_points": 5
  }'
```

### Transition Issue

```bash
curl -X POST http://localhost:8080/api/issues/{issueId}/transitions \
  -H "Content-Type: application/json" \
  -d '{
    "to_status": "in_progress",
    "version": 1
  }'
```

### Search Issues

```bash
curl "http://localhost:8080/api/search?q=authentication&status=in_progress&priority=high"
```

## WebSocket Events

Connect to `/api/ws/:projectId?userId=yourUserId` to receive real-time events:

```javascript
const ws = new WebSocket('ws://localhost:8080/api/ws/proj_123?userId=user_456');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Event:', data.type, data.data);
};
```

Event types:
- `issue_created` - New issue created
- `issue_updated` - Issue modified
- `issue_moved` - Issue moved to different status
- `comment_added` - New comment added
- `sprint_updated` - Sprint modified
- `presence` - User joined/left

## Configuration

Environment variables (`.env`):

```env
PORT=8080
ENV=development
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=project_management
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
CORS_ORIGINS=http://localhost:5173,http://localhost:3000
```

## Key Design Decisions

### Optimistic Concurrency Control

Issues use version-based locking to prevent lost updates:

```json
{
  "title": "Updated title",
  "version": 3
}
```

Returns `409 Conflict` if version mismatch.

### Workflow Engine

Workflows define allowed transitions and validation rules:

```json
{
  "from": "in_progress",
  "to": "in_review",
  "conditions": [
    {"field": "assignee_id", "operator": "not_empty"}
  ],
  "actions": [
    {"type": "assign_reviewer"}
  ]
}
```

### Activity Logging

All changes are logged immutably for audit trail and event replay.

### Real-Time Architecture

WebSocket connections are managed per project with automatic presence tracking.

## Scalability

- **Horizontal Scaling**: Stateless API design
- **Database**: MongoDB replica sets and sharding
- **Caching**: Redis for frequently accessed data (future)
- **Load Balancing**: Multiple instances behind load balancer
- **WebSocket Distribution**: Redis pub/sub for multi-instance (future)

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./service/...
```

## Production Deployment

1. Build optimized binary:
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
```

2. Use Docker:
```bash
docker build -t project-management-backend .
docker run -p 8080:8080 --env-file .env project-management-backend
```

3. Deploy to cloud (AWS, GCP, Azure) with:
   - MongoDB Atlas for database
   - Load balancer for multiple instances
   - Redis for caching and WebSocket distribution

## Security Considerations

- Implement proper JWT authentication
- Add rate limiting per user
- Validate all inputs
- Use HTTPS in production
- Sanitize search queries
- Implement RBAC for project access


## License

MIT

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request
