# Project Management Platform

A complete, production-ready project management system with a Go backend and React frontend.

## Overview

This is a full-stack project management platform that enables engineering teams to plan, track, and deliver software collaboratively. It features real-time updates, flexible workflows, sprint management, and a modern Kanban board interface.

## Architecture

```
project-management/
├── backend/          # Go + Fiber + MongoDB
│   ├── main.go
│   ├── config/
│   ├── provider/
│   ├── service/
│   ├── store/
│   ├── router/
│   └── models/
└── frontend/         # React + Vite + TypeScript + Tailwind
    ├── src/
    │   ├── components/
    │   ├── pages/
    │   └── lib/
    └── package.json
```

## Features

### Backend
- RESTful API with Fiber framework
- MongoDB for flexible data storage
- WebSocket for real-time updates
- Optimistic concurrency control
- Workflow engine with validation
- Sprint management
- Activity logging and audit trail
- Full-text search

### Frontend
- Modern React with TypeScript
- Beautiful landing page with hero section
- Enhanced navigation with mobile support
- Drag-and-drop Kanban board
- Real-time collaboration via WebSocket
- shadcn/ui components with smooth animations
- Fully responsive design (mobile, tablet, desktop)
- TanStack Query for data management
- Accessibility-first approach

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- MongoDB 5.0+

### 1. Start MongoDB

```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

### 2. Start Backend

```bash
cd backend
cp .env.example .env
go mod download
go run main.go
```

Backend will run on `http://localhost:8080`

### 3. Start Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend will run on `http://localhost:5173`

### 4. Open Browser

Navigate to `http://localhost:5173` and start creating projects!

## Detailed Setup

See individual README files:
- [Backend Setup](./backend/README.md)
- [Frontend Setup](./frontend/README.md)
- [UX Improvements](./UX_IMPLEMENTATION_SUMMARY.md) - New landing page and UX enhancements

## API Documentation

### Projects
- `POST /api/projects` - Create project
- `GET /api/projects` - List projects
- `GET /api/projects/:id` - Get project

### Issues
- `POST /api/projects/:id/issues` - Create issue
- `GET /api/projects/:id/board` - Get board
- `PATCH /api/issues/:id` - Update issue
- `POST /api/issues/:id/transitions` - Change status

### Real-Time
- `WS /api/ws/:projectId` - WebSocket connection

Full API documentation in [backend/README.md](./backend/README.md)

## Development

### Backend Development

```bash
cd backend
make run          # Run with hot reload
make test         # Run tests
make build        # Build binary
```

### Frontend Development

```bash
cd frontend
npm run dev       # Development server
npm run build     # Production build
npm run lint      # Lint code
```

## Production Deployment

### Backend

```bash
cd backend
make build-prod
# Deploy binary to your server
```

Or use Docker:

```bash
cd backend
docker build -t pm-backend .
docker run -p 8080:8080 --env-file .env pm-backend
```

### Frontend

```bash
cd frontend
npm run build
# Deploy dist/ folder to Vercel, Netlify, or S3
```

## Key Technologies

**Backend:**
- Go 1.21
- Fiber v2 (HTTP framework)
- MongoDB (database)
- WebSocket (real-time)

**Frontend:**
- React 18
- TypeScript
- Vite (build tool)
- Tailwind CSS
- shadcn/ui components
- TanStack Query

## Project Structure

### Backend Layers

```
Router → Service → Store → MongoDB
```

- **Router**: HTTP handlers
- **Service**: Business logic
- **Store**: Data access
- **Provider**: Dependency injection

### Frontend Architecture

```
Pages → Components → API → Backend
```

- **Pages**: Route components
- **Components**: Reusable UI
- **API**: Backend integration
- **WebSocket**: Real-time updates

## Features Showcase

### Kanban Board
- Drag-and-drop issues between columns
- Real-time updates across all users
- Visual priority and story points
- Label management

### Issue Management
- Create issues with type, priority, story points
- Detailed issue view with comments
- Activity timeline
- @mentions in comments

### Sprint Management
- Create and manage sprints
- Start/complete sprints
- Velocity tracking
- Carry-over incomplete items

### Real-Time Collaboration
- WebSocket connection
- Live issue updates
- Presence tracking
- Automatic reconnection

## Testing

### Backend Tests

```bash
cd backend
go test ./...
```

### Frontend Tests

```bash
cd frontend
npm run test
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT

## Support

For issues and questions:
- Backend: See [backend/README.md](./backend/README.md)
- Frontend: See [frontend/README.md](./frontend/README.md)
- Architecture: See [backend/ARCHITECTURE.md](./backend/ARCHITECTURE.md)
