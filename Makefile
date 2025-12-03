.PHONY: help docker-up docker-down run-backend run-frontend migrate-up migrate-down test clean validate-config

# Default database connection variables (can be overridden via environment)
DB_HOST ?= localhost
DB_PORT ?= 5434
DB_USER ?= cyclingstream
DB_PASSWORD ?= cyclingstream_dev
DB_NAME ?= cyclingstream
DB_SSLMODE ?= disable

# Build database connection string
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ''
	@echo 'Environment variables (with defaults):'
	@echo '  DB_HOST=$(DB_HOST)'
	@echo '  DB_PORT=$(DB_PORT)'
	@echo '  DB_USER=$(DB_USER)'
	@echo '  DB_NAME=$(DB_NAME)'

docker-up: ## Start Docker containers (Postgres, pgAdmin, and Owncast)
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Waiting for Postgres to be ready..."
	@timeout=30; \
	while [ $$timeout -gt 0 ]; do \
		if docker exec cyclingstream_postgres pg_isready -U $(DB_USER) > /dev/null 2>&1; then \
			echo "Postgres is ready!"; \
			break; \
		fi; \
		sleep 1; \
		timeout=$$((timeout - 1)); \
	done; \
	if [ $$timeout -eq 0 ]; then \
		echo "ERROR: Postgres failed to start within 30 seconds"; \
		exit 1; \
	fi
	@echo "pgAdmin available at http://localhost:$${PGADMIN_PORT:-5050}"
	@echo "Owncast available at http://localhost:$${OWNCAST_HTTP_PORT:-8081}"
	@echo "Owncast RTMP endpoint: rtmp://localhost:$${OWNCAST_RTMP_PORT:-1935}/live"

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker container logs
	docker-compose logs -f

run-backend: ## Run the backend server
	@echo "Starting backend server..."
	cd backend && go run cmd/api/main.go

run-frontend: ## Run the frontend dev server
	@echo "Starting frontend dev server..."
	cd frontend && npm run dev

validate-config: ## Validate backend configuration (checks env vars)
	@echo "Validating backend configuration..."
	@cd backend && go run -tags validation cmd/api/main.go 2>&1 | head -1 || \
		(echo "Configuration validation: Starting server to check config..." && \
		 timeout 2 go run cmd/api/main.go 2>&1 | grep -q "configuration validation failed" && \
		 echo "ERROR: Configuration validation failed. Check your environment variables." && exit 1 || \
		 echo "Configuration appears valid.")

migrate-up: ## Run database migrations (uses DB_* env vars)
	@echo "Running database migrations..."
	@if ! command -v migrate > /dev/null 2>&1; then \
		echo "ERROR: migrate CLI not found. Install from https://github.com/golang-migrate/migrate"; \
		exit 1; \
	fi
	cd backend && migrate -path migrations -database "$(DB_URL)" up
	@echo "Migrations completed successfully."

migrate-down: ## Rollback last database migration (uses DB_* env vars)
	@echo "Rolling back last migration..."
	@if ! command -v migrate > /dev/null 2>&1; then \
		echo "ERROR: migrate CLI not found. Install from https://github.com/golang-migrate/migrate"; \
		exit 1; \
	fi
	cd backend && migrate -path migrations -database "$(DB_URL)" down 1
	@echo "Migration rolled back successfully."

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "ERROR: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@if ! command -v migrate > /dev/null 2>&1; then \
		echo "ERROR: migrate CLI not found. Install from https://github.com/golang-migrate/migrate"; \
		exit 1; \
	fi
	cd backend && migrate create -ext sql -dir migrations -seq $(NAME)
	@echo "Migration file created: backend/migrations/$(shell ls -t backend/migrations/*.sql | head -1 | xargs basename)"

test: ## Run all tests
	@echo "Running backend tests..."
	@cd backend && go test ./... || (echo "Backend tests failed!" && exit 1)
	@echo "Running frontend tests..."
	@cd frontend && npm test || (echo "Frontend tests failed!" && exit 1)
	@echo "All tests passed!"

test-backend: ## Run backend tests only
	@cd backend && go test ./...

test-frontend: ## Run frontend tests only
	@cd frontend && npm test

test-api: ## Run comprehensive API test suite (requires server running)
	@echo "Running API test suite..."
	@if [ ! -f backend/scripts/test_all_api.sh ]; then \
		echo "ERROR: Test script not found: backend/scripts/test_all_api.sh"; \
		exit 1; \
	fi
	@chmod +x backend/scripts/test_all_api.sh
	cd backend && ./scripts/test_all_api.sh

test-api-quick: ## Run quick API smoke tests
	@echo "Running quick API tests..."
	@if [ ! -f backend/scripts/test_security.sh ]; then \
		echo "ERROR: Test script not found: backend/scripts/test_security.sh"; \
		exit 1; \
	fi
	@chmod +x backend/scripts/test_security.sh
	cd backend && ./scripts/test_security.sh

lint-backend: ## Lint backend code
	@echo "Running go vet..."
	@cd backend && go vet ./... || (echo "go vet found issues!" && exit 1)
	@if command -v golangci-lint > /dev/null 2>&1; then \
		echo "Running golangci-lint..."; \
		cd backend && golangci-lint run ./... || (echo "golangci-lint found issues!" && exit 1); \
	else \
		echo "WARNING: golangci-lint not installed, skipping..."; \
		echo "Install from: https://golangci-lint.run/usage/install/"; \
	fi

lint-frontend: ## Lint frontend code
	@cd frontend && npm run lint

build-backend: ## Build backend binary
	@echo "Building backend..."
	@cd backend && go build -o bin/cyclingstream-api ./cmd/api/main.go
	@echo "Backend built: backend/bin/cyclingstream-api"

build-frontend: ## Build frontend for production
	@echo "Building frontend..."
	@cd frontend && npm run build
	@echo "Frontend built successfully."

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@cd backend && go clean
	@cd backend && rm -rf bin/
	@cd frontend && rm -rf .next out
	@echo "Clean complete."

clean-all: clean ## Clean build artifacts and Docker volumes
	@echo "Cleaning Docker volumes..."
	@docker-compose down -v
	@echo "Clean complete (including Docker volumes)."

setup: ## Initial setup (install dependencies)
	@echo "Setting up backend..."
	@cd backend && go mod download
	@echo "Setting up frontend..."
	@cd frontend && npm install
	@echo "Setup complete!"

backup-db: ## Backup the database
	@echo "Backing up database..."
	@chmod +x scripts/backup-db.sh
	@./scripts/backup-db.sh

validate-schema: ## Validate database schema against code expectations
	@echo "Validating database schema..."
	@chmod +x scripts/validate-schema.sh
	@./scripts/validate-schema.sh

inspect-db: ## Inspect database tables and schema
	@echo "Inspecting database..."
	@chmod +x scripts/inspect-db.sh
	@./scripts/inspect-db.sh

fix-schema: ## Attempt to fix database schema issues
	@echo "Fixing database schema..."
	@chmod +x scripts/fix-schema.sh
	@./scripts/fix-schema.sh

owncast-generate-key: ## Generate a secure stream key for Owncast
	@echo "Generating secure stream key..."
	@if command -v openssl > /dev/null 2>&1; then \
		KEY=$$(openssl rand -hex 32); \
		echo "Generated stream key: $$KEY"; \
		echo ""; \
		echo "Add this to your .env file or export it:"; \
		echo "  export OWNCAST_STREAM_KEY=$$KEY"; \
		echo ""; \
		echo "Or add to docker-compose.yml environment section:"; \
		echo "  - OWNCAST_STREAM_KEY=$$KEY"; \
	else \
		echo "ERROR: openssl not found. Install openssl or generate key manually."; \
		exit 1; \
	fi

owncast-logs: ## View Owncast container logs
	docker-compose logs -f owncast

owncast-status: ## Check Owncast service status
	@echo "Checking Owncast status..."
	@docker-compose ps owncast
	@echo ""
	@echo "Testing Owncast API..."
	@curl -s http://localhost:$${OWNCAST_HTTP_PORT:-8081}/api/status | python3 -m json.tool 2>/dev/null || echo "Owncast API not responding"
