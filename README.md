# Project Management Application

A full-stack project management application with secure authentication, built with Go, React, and MongoDB.

## Features

- 🔐 Secure JWT-based authentication with bcrypt password hashing
- 📱 Fully responsive design (375px mobile to 1280px+ desktop)
- 🌓 Dark mode with persistent theme preference
- 🐳 Single Docker image for frontend + backend
- 📊 Kanban board with drag-and-drop
- 🏃 Sprint management
- 💬 Comments and activity tracking
- 🔔 Real-time notifications via WebSocket
- 🔍 Advanced search functionality

## Quick Start with Docker

### Prerequisites

- Docker and Docker Compose installed
- No other services running on ports 8080 or 27017

### Setup

1. Clone the repository
2. Copy environment file:
   ```bash
   cp .env.example .env
   ```
3. Start the entire stack:
   ```bash
   docker compose up
   ```

That's it! The application will be available at:
- **Application**: http://localhost:8080 (frontend + backend in one)
- **MongoDB**: localhost:27017

### Architecture

The Docker setup uses a single container that runs:
- **Nginx** on port 8080 (serves frontend + proxies API requests)
- **Go Backend** on internal port 8081
- **MongoDB** in separate container

This centralized approach means:
- ✅ Single port to expose (8080)
- ✅ No CORS issues (same origin)
- ✅ Simpler deployment
- ✅ Smaller resource footprint

### Test Credentials

The database is automatically seeded with test data:
- **Email**: test@example.com
- **Password**: testpass123

The seed includes:
- 1 test user
- 1 demo project (DEMO)
- 3 sample tasks with different statuses

## Local Development

### Backend (Go)

```bash
cd backend
cp .env.example .env
# Update MONGO_URI if needed
go mod download
go run main.go
```

Backend runs on http://localhost:8080

### Frontend (React)

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Frontend runs on http://localhost:5173

## Environment Variables

### Root `.env` (for Docker Compose)
```env
# MongoDB Configuration
MONGO_ROOT_USER=admin
MONGO_ROOT_PASSWORD=admin123
MONGO_DATABASE=project_management

# JWT Configuration (CHANGE THIS!)
JWT_SECRET=your-secure-secret-key-here-min-32-chars

# Application Port
APP_PORT=8080
```

### Backend `.env` (for local dev)
```env
PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=project_management
JWT_SECRET=your-secure-secret-key-here
JWT_EXPIRY=24h
CORS_ORIGINS=http://localhost:5173,http://localhost:3000
```

### Frontend `.env` (for local dev)
```env
VITE_BACKEND_URL=http://localhost:8080
```

## Security Features

- Passwords hashed with bcrypt (cost factor 12)
- JWT tokens with configurable expiration
- CORS protection
- Input validation on both frontend and backend
- Password requirements (min 8 characters)
- Secure HTTP-only authentication flow
- Protected API routes with middleware

## Tech Stack

### Backend
- Go 1.21
- Fiber (web framework)
- MongoDB
- JWT authentication
- bcrypt password hashing

### Frontend
- React 18
- TypeScript
- TailwindCSS
- React Query
- React Router
- Shadcn/ui components

### Infrastructure
- Docker & Docker Compose
- MongoDB 7
- Nginx (reverse proxy)
- Multi-stage Docker build

## API Endpoints

All endpoints are available at `http://localhost:8080/api`

### Public
- `POST /api/auth/register` - Create account
- `POST /api/auth/login` - Login

### Protected (requires JWT token)
- `GET /api/auth/me` - Get current user
- `POST /api/auth/change-password` - Change password
- `GET /api/projects` - List projects
- `POST /api/projects` - Create project
- `GET /api/projects/:id/board` - Get project board
- And more...

## Production Deployment

1. Update `.env` with strong secrets:
   ```bash
   # Generate a strong JWT secret
   openssl rand -base64 32
   ```

2. Build and run:
   ```bash
   docker compose up -d
   ```

3. The app runs on a single port (8080) making it easy to:
   - Add SSL with Let's Encrypt
   - Put behind a reverse proxy
   - Deploy to any cloud provider

4. Consider adding:
   - SSL certificates (Let's Encrypt)
   - Rate limiting
   - Monitoring and logging
   - Backup strategy for MongoDB

## Development Commands

### Backend
```bash
cd backend
go test ./...        # Run tests
go build -o bin/server  # Build
go fmt ./...         # Format code
```

### Frontend
```bash
cd frontend
npm run dev          # Development server
npm run build        # Production build
npm run preview      # Preview production build
npm run lint         # Lint code
```

### Docker
```bash
docker compose up           # Start all services
docker compose up -d        # Start in background
docker compose down         # Stop all services
docker compose logs -f      # View logs
docker compose logs -f app  # View app logs only
docker compose restart app  # Restart app
```

## Responsive Design

The application is fully responsive and tested at:
- Mobile: 375px width
- Tablet: 768px width
- Desktop: 1280px+ width

All layouts adapt gracefully with no broken elements or console errors.

## Dark Mode

Dark mode is implemented with:
- System preference detection
- Manual toggle in navigation
- Persistent preference in localStorage
- Smooth transitions between themes

## Troubleshooting

### Port already in use
```bash
# Check what's using port 8080
lsof -i :8080

# Or change the port in .env
APP_PORT=3000
```

### MongoDB connection issues
```bash
# Check MongoDB is running
docker compose ps

# View MongoDB logs
docker compose logs mongodb
```

### Frontend not loading
```bash
# Rebuild the image
docker compose build --no-cache app
docker compose up
```

## License

MIT
