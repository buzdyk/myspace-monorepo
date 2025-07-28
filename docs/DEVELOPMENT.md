# Development Guide

## Prerequisites

- **Go 1.21+** for backend development
- **Node.js 18+** for frontend development
- **Git** for version control

## Project Structure

This monorepo follows a clean architecture pattern:

```
myspace/
├── backend/           # Go API server
├── frontend/          # Next.js React app  
├── shared/            # Shared types/utilities
├── docs/              # Documentation
└── README.md
```

## Backend Development

### Architecture

The Go backend follows these principles:

1. **Interface-driven design**: All time trackers implement the `TimeTracker` interface
2. **Repository pattern**: `TrackersRepository` aggregates data from multiple providers
3. **Clean separation**: Handlers, repositories, and trackers are clearly separated
4. **Dependency injection**: Dependencies are injected at startup

### Key Components

- `internal/interfaces/time_tracker.go` - Core interface all providers implement
- `internal/trackers/` - Time tracking provider implementations
- `internal/repositories/trackers.go` - Aggregates multiple providers
- `internal/handlers/` - HTTP request handlers
- `internal/types/` - Core data structures

### Adding New Time Trackers

1. Create new file in `internal/trackers/`
2. Implement the `TimeTracker` interface:
   ```go
   type TimeTracker interface {
       GetUserID() string
       GetSeconds(from, to time.Time) (int, error)  
       GetRunningSeconds() (int, error)
       GetMonthlyTimeByProject(dayOfMonth time.Time) (*types.ProjectTimes, error)
       GetMonthIntervals(dayOfMonth time.Time) (*types.ProjectTimes, error)
   }
   ```
3. Add configuration to `internal/config/config.go`
4. Register in `repositories/trackers.go` hydrate method

### Running Tests

```bash
cd backend
go test ./...
```

### Building

```bash
cd backend
go build -o myspace cmd/main.go
```

## Frontend Development

### Architecture

The frontend uses Next.js 14 with:

- **App Router**: File-based routing with dynamic segments
- **TypeScript**: Type safety throughout
- **Tailwind CSS**: Utility-first styling
- **Axios**: HTTP client for API calls

### Key Features

- **Dynamic routing**: `[year]/[month]/[day]` for date-based URLs
- **API integration**: Proxies requests to Go backend via Next.js rewrites
- **Responsive design**: Mobile-first approach with Tailwind
- **Error handling**: Graceful loading states and error boundaries

### Development Server

```bash
cd frontend
npm run dev
```

### Building for Production

```bash
cd frontend
npm run build
npm start
```

## Environment Setup

### Backend Environment

Copy `backend/.env.example` to `backend/.env` and configure:

```bash
# Database
DB_PATH=./database.sqlite
PORT=8080

# Time tracking providers (configure as needed)
CLOCKIFY_TOKEN=
CLOCKIFY_WORKSPACE_ID=
CLOCKIFY_USER_ID=
EVERHOUR_TOKEN=
MAYVEN_AUTH=
```

### Frontend Environment

The frontend automatically proxies API requests to `localhost:8080` in development.

## Database

The application uses SQLite with GORM for ORM functionality.

### Models

- `User` - Authentication (placeholder for future auth)
- `Tracker` - Configured time tracking providers
- `Project` - Projects associated with trackers  
- `Track` - Individual time entries
- `Setting` - Application configuration

### Migrations

Database tables are auto-migrated on startup via GORM's `AutoMigrate`.

## Deployment

### Development

1. Start backend: `cd backend && go run cmd/main.go`
2. Start frontend: `cd frontend && npm run dev`  
3. Visit: `http://localhost:3000`

### Production

1. Build backend: `cd backend && go build -o myspace cmd/main.go`
2. Build frontend: `cd frontend && npm run build`
3. Deploy both services with appropriate environment variables

## Code Style

### Go

- Follow standard Go conventions
- Use `gofmt` for formatting
- Minimal commenting (self-documenting code preferred)
- Error handling on every operation

### TypeScript/React

- Use TypeScript interfaces for all data structures
- Functional components with hooks
- Tailwind for all styling
- Minimal commenting

## Git Workflow

1. Create feature branch from `main`
2. Make changes with descriptive commits
3. Test both backend and frontend  
4. Create pull request
5. Merge after review

## Troubleshooting

### Backend Issues

- Check environment variables are set correctly
- Verify time tracking API credentials are valid
- Check database file permissions
- Review logs for HTTP request errors

### Frontend Issues  

- Verify backend is running on port 8080
- Check browser network tab for API call failures
- Ensure Node.js version compatibility
- Clear browser cache if seeing stale data