# MySpace - Time Tracking Application

A Go/Next.js rewrite of the Laravel-based time tracking application that unifies time data from multiple providers (Clockify, Everhour, Mayven).

## Architecture

This is a monorepo containing:
- **backend/**: Go API server with time tracking integrations
- **frontend/**: Next.js React application 
- **shared/**: Shared types and utilities
- **docs/**: Architecture documentation

## Features

- **Multi-provider Support**: Aggregates data from Clockify, Everhour, and Mayven
- **Daily View**: Track current day hours and running timers
- **Monthly Views**: Project breakdown and calendar visualization
- **REST API**: Clean HTTP API for time tracking data
- **Real-time Updates**: Current running timer tracking

## Quick Start

### Backend (Go)

1. Navigate to backend directory:
```bash
cd backend
```

2. Copy environment configuration:
```bash
cp .env.example .env
```

3. Configure your time tracking providers in `.env`:
```bash
CLOCKIFY_TOKEN=your_token_here
CLOCKIFY_WORKSPACE_ID=your_workspace_id
CLOCKIFY_USER_ID=your_user_id
EVERHOUR_TOKEN=your_token_here
MAYVEN_AUTH=your_auth_here
```

4. Install dependencies and run:
```bash
go mod tidy
go run cmd/main.go
```

The API will be available at `http://localhost:8080`

### Frontend (Next.js)

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Run development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## API Endpoints

- `GET /` - Redirects to `/today`
- `GET /today` - Redirects to current date
- `GET /:year/:month/:day` - Daily time tracking data
- `GET /month` - Redirects to current month projects
- `GET /:year/:month/projects` - Monthly project breakdown
- `GET /:year/:month/calendar` - Monthly calendar view

## Development

### Backend Structure

```
backend/
├── cmd/main.go              # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database models and connection
│   ├── handlers/            # HTTP request handlers
│   ├── interfaces/          # TimeTracker interface
│   ├── repositories/        # Data aggregation layer
│   ├── trackers/            # Time tracking provider implementations
│   └── types/               # Core data types
└── go.mod
```

### Frontend Structure

```
frontend/
├── app/                     # Next.js 14 app router
│   ├── [year]/[month]/      # Dynamic routes for date-based views
│   ├── today/               # Today redirect page
│   └── month/               # Month redirect page
├── components/              # Reusable React components
└── lib/                     # Utilities and API client
```

## Time Tracking Providers

### Clockify
- Requires: API token, workspace ID, user ID
- Endpoint: `https://api.clockify.me`

### Everhour  
- Requires: API token
- Endpoint: `https://api.everhour.com`

### Mayven
- Requires: Authorization header
- Endpoint: `https://api.mayven.io` (configurable)

## Database Schema

The application uses SQLite with the following tables:
- `users` - User authentication 
- `trackers` - Configured time tracking providers
- `projects` - Projects associated with trackers
- `tracks` - Individual time entries
- `settings` - Application settings

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test both backend and frontend
5. Submit a pull request

## License

MIT License