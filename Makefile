.PHONY: help deps build run test clean docker-build docker-run db-init dev-setup dev-up dev-down dev-run dev-logs prod-setup prod-build prod-deploy

# Configuration
APP_NAME = file-pub
DOCKER_IMAGE = $(APP_NAME):latest
DOCKER_CONTAINER = $(APP_NAME)-container

help:
	@echo "File Pub Makefile — Commands:"
	@echo ""
	@echo "Development (Local):"
	@echo "  dev-setup    - Start Docker MySQL in detached mode"
	@echo "  dev-up       - Start Docker MySQL"
	@echo "  dev-down     - Stop Docker MySQL"
	@echo "  dev-run      - Run application with dev environment"
	@echo "  dev-logs     - View MySQL logs"
	@echo "  dev-compose  - Run everything with docker-compose"
	@echo ""
	@echo "Production (AWS):"
	@echo "  prod-setup   - Validate production configuration"
	@echo "  prod-build   - Build production binary"
	@echo "  prod-deploy  - Deploy to EC2 (requires SSH_HOST)"
	@echo ""
	@echo "General:"
	@echo "  deps         - Download Go dependencies"
	@echo "  build        - Build the application binary"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run application in Docker"
	@echo "  docker-stop  - Stop Docker container"
	@echo ""
	@echo "Database:"
	@echo "  db-init      - Initialize database schema"

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

build: deps
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) main.go

run: deps
	@echo "Running $(APP_NAME)..."
	go run main.go

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run || echo "golangci-lint not installed. Run: brew install golangci-lint"

docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	@echo "Running Docker container..."
	docker run -d \
		--name $(DOCKER_CONTAINER) \
		-p 8080:8080 \
		--env-file .env \
		$(DOCKER_IMAGE)
	@echo "Container running at http://localhost:8080"

docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(DOCKER_CONTAINER) || true
	docker rm $(DOCKER_CONTAINER) || true

docker-logs:
	docker logs -f $(DOCKER_CONTAINER)

db-init:
	@echo "Initializing database..."
	@if [ -f db/init.sql ]; then \
		echo "Run this SQL against your RDS instance:"; \
		cat db/init.sql; \
	else \
		echo "db/init.sql not found"; \
	fi

# Development Commands
dev-setup:
	@echo "Setting up development environment..."
	@chmod +x scripts/setup-dev.sh
	@./scripts/setup-dev.sh

dev-up:
	@echo "Starting Docker MySQL..."
	docker-compose up -d mysql
	@echo "Waiting for MySQL to be ready..."
	@sleep 5
	@echo "MySQL is ready at localhost:3306"

dev-down:
	@echo "Stopping Docker services..."
	docker-compose down

dev-run:
	@echo "Running application in development mode..."
	@if [ ! -f .env.dev ]; then \
		echo "Error: .env.dev not found"; \
		exit 1; \
	fi
	@echo "Checking if MySQL is running..."
	@if ! docker ps | grep -q filepub-mysql-dev; then \
		echo "MySQL not running. Starting it now..."; \
		docker-compose up -d mysql; \
		echo "Waiting for MySQL to be ready..."; \
		sleep 5; \
	else \
		echo "MySQL is already running."; \
	fi
	@echo "Starting application..."
	@export $$(cat .env.dev | grep -v '^#' | xargs) && go run main.go

dev-logs:
	@echo "Viewing MySQL logs..."
	docker-compose logs -f mysql

dev-compose:
	@echo "Starting full development environment with docker-compose..."
	docker-compose up

dev-compose-build:
	@echo "Building and starting development environment..."
	docker-compose up --build

# Production Commands
prod-setup:
	@echo "Setting up production environment..."
	@chmod +x scripts/setup-prod.sh
	@./scripts/setup-prod.sh

prod-build:
	@echo "Building production binary..."
	@if [ ! -f .env.prod ]; then \
		echo "Warning: .env.prod not found"; \
	fi
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/$(APP_NAME) main.go
	@echo "Production binary built: bin/$(APP_NAME)"

prod-deploy:
	@echo "Deploying to production..."
	@if [ -z "$(SSH_HOST)" ]; then \
		echo "Error: SSH_HOST not set. Usage: make prod-deploy SSH_HOST=ec2-user@1.2.3.4"; \
		exit 1; \
	fi
	@echo "Building production binary..."
	@make prod-build
	@echo "Copying files to $(SSH_HOST)..."
	scp bin/$(APP_NAME) $(SSH_HOST):~/
	scp .env.prod $(SSH_HOST):~/.env
	scp -r templates $(SSH_HOST):~/
	@echo "Deployment complete!"
	@echo "SSH to server and run:"
	@echo "  export \$$(cat .env | xargs)"
	@echo "  ./$(APP_NAME)"
