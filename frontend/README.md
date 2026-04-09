# Project Management Frontend

Modern, responsive project management UI built with React, Vite, TypeScript, Tailwind CSS, and shadcn/ui components.

## Features

- **Landing Page**: Beautiful hero section with clear CTAs and feature showcase
- **Kanban Board**: Drag-and-drop issue management across status columns
- **Real-Time Updates**: WebSocket integration for live collaboration
- **Project Management**: Create and manage multiple projects
- **Sprint Management**: Plan and track sprints with velocity metrics
- **Workflow Customization**: Design workflows that match your team's process
- **Issue Tracking**: Full CRUD operations with detailed issue views
- **Comments & Activity**: Threaded discussions and activity timeline
- **Search & Filter**: Quick search across issues
- **Responsive Design**: Optimized for desktop, tablet, and mobile
- **Modern UI**: Beautiful components from shadcn/ui with smooth animations
- **Accessibility**: Keyboard navigation and screen reader support

## Tech Stack

- **React 18**: Modern React with hooks
- **TypeScript**: Type-safe development
- **Vite**: Lightning-fast build tool
- **Tailwind CSS**: Utility-first styling
- **shadcn/ui**: High-quality React components
- **TanStack Query**: Powerful data fetching and caching
- **React Router**: Client-side routing
- **Axios**: HTTP client
- **date-fns**: Date formatting

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Backend server running on `http://localhost:8080`

### Installation

```bash
cd frontend
npm install
```

### Development

```bash
npm run dev
```

The app will be available at `http://localhost:5173`

### Build for Production

```bash
npm run build
npm run preview
```

## Project Structure

```
frontend/
├── src/
│   ├── components/
│   │   ├── ui/              # shadcn/ui components
│   │   ├── KanbanBoard.tsx  # Main board view
│   │   ├── IssueDialog.tsx  # Issue detail modal
│   │   ├── CreateIssueDialog.tsx
│   │   ├── SprintManagement.tsx
│   │   ├── WorkflowManagement.tsx
│   │   ├── SearchBar.tsx
│   │   └── Layout.tsx       # App layout with navigation
│   ├── pages/
│   │   ├── Landing.tsx      # Landing page with hero section
│   │   ├── ProjectList.tsx  # Projects overview
│   │   └── ProjectBoard.tsx # Project board view
│   ├── lib/
│   │   ├── api.ts           # API client and types
│   │   ├── websocket.ts     # WebSocket hook
│   │   └── utils.ts         # Utility functions
│   ├── App.tsx              # Main app component
│   ├── main.tsx             # Entry point
│   └── index.css            # Global styles & animations
├── UX_IMPROVEMENTS.md       # UX documentation
├── DESIGN_SYSTEM.md         # Design system reference
├── index.html
├── vite.config.ts
├── tailwind.config.js
└── package.json
```

## Features Overview

### Landing Page

- Hero section with compelling value proposition
- Feature showcase highlighting key capabilities
- Clear call-to-action buttons
- Smooth scroll navigation
- Responsive design with animations
- Direct path to project management

### Kanban Board

- Drag and drop issues between status columns
- Visual priority indicators
- Story points display
- Label tags
- Real-time updates via WebSocket

### Issue Management

- Create issues with type, priority, story points
- View detailed issue information
- Add comments with @mentions
- View activity history
- Update issue status via drag-and-drop

### Project Management

- Create and list projects
- Project-specific boards
- Real-time connection status indicator

### Real-Time Collaboration

- WebSocket connection for live updates
- Automatic reconnection
- Connection status indicator
- Instant issue updates across all connected clients

## API Integration

The frontend connects to the backend API at `http://localhost:8080/api`:

- **Projects**: `/api/projects`
- **Issues**: `/api/projects/:id/issues`, `/api/issues/:id`
- **Comments**: `/api/issues/:id/comments`
- **Activity**: `/api/issues/:id/activity`
- **WebSocket**: `ws://localhost:8080/api/ws/:projectId`

## Customization

### Theme

Edit `src/index.css` to customize the color scheme:

```css
:root {
  --primary: 221.2 83.2% 53.3%;
  --secondary: 210 40% 96.1%;
  /* ... */
}
```

### Components

All UI components are in `src/components/ui/` and can be customized using Tailwind classes.

## Development Tips

### Adding New Components

```bash
# Components follow shadcn/ui patterns
# See: https://ui.shadcn.com/docs/components
```

### Type Safety

All API responses are typed in `src/lib/api.ts`. Update types when backend changes.

### State Management

- **TanStack Query** for server state (issues, projects, etc.)
- **React hooks** for local UI state
- **WebSocket** for real-time updates

## Troubleshooting

### Backend Connection Issues

Make sure the backend is running on `http://localhost:8080`:

```bash
cd ../backend
go run main.go
```

### WebSocket Not Connecting

Check that:
1. Backend is running
2. WebSocket endpoint is accessible
3. No firewall blocking connections

### Build Errors

Clear node_modules and reinstall:

```bash
rm -rf node_modules package-lock.json
npm install
```

## Production Deployment

### Build

```bash
npm run build
```

Output will be in `dist/` directory.

### Deploy

Deploy the `dist/` folder to:
- **Vercel**: `vercel deploy`
- **Netlify**: `netlify deploy`
- **AWS S3 + CloudFront**
- **Any static hosting service**

### Environment Variables

For production, update the API base URL in `vite.config.ts` or use environment variables:

```typescript
// vite.config.ts
server: {
  proxy: {
    '/api': {
      target: process.env.VITE_API_URL || 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

## Contributing

1. Follow the existing code style
2. Use TypeScript for type safety
3. Keep components small and focused
4. Add comments for complex logic

## License

MIT
