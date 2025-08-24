.PHONY: dev prod up down restart logs clean

# Development commands
dev:
	docker-compose -f docker-compose.dev.yml up

dev-build:
	docker-compose -f docker-compose.dev.yml up --build

dev-down:
	docker-compose -f docker-compose.dev.yml down

# Production commands
prod:
	docker-compose up -d

prod-build:
	docker-compose up --build -d

down:
	docker-compose down

# Utility commands
restart: down prod

logs:
	docker-compose logs -f

logs-backend:
	docker-compose logs -f backend

logs-frontend:
	docker-compose logs -f frontend

clean:
	docker-compose down -v
	rm -rf backend/tmp
	rm -rf frontend/.next
	rm -rf frontend/node_modules

# Local development (without Docker)
local-backend:
	cd backend && air

local-frontend:
	cd frontend && npm run dev

local-install:
	cd backend && go mod download
	cd frontend && npm install

# Install air globally for local development
install-air:
	go install github.com/air-verse/air@latest