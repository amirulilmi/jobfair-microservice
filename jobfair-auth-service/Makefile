.PHONY: help build test docker-build docker-push deploy

# Add migration commands
.PHONY: migrate-up migrate-down migrate-create migrate-status

DB_URL := postgres://jobfair_user:jobfair_pass@localhost:5432/jobfair?sslmode=disable

migrate-up: ## Run all up migrations
	@echo "Running migrations up..."
	migrate -path migrations -database "$(DB_URL)" up

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-down-all: ## Rollback all migrations
	@echo "Rolling back all migrations..."
	migrate -path migrations -database "$(DB_URL)" down

migrate-create: ## Create new migration file (usage: make migrate-create name=create_users)
	@echo "Creating new migration..."
	migrate create -ext sql -dir migrations -seq $(name)

migrate-status: ## Show migration status
	@echo "Migration status:"
	migrate -path migrations -database "$(DB_URL)" version

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@echo "Forcing migration version..."
	migrate -path migrations -database "$(DB_URL)" force $(version)



SERVICE_NAME := jobfair-auth-service
VERSION := $(shell git describe --tags --always --dirty)
DOCKER_REGISTRY := your-registry.com
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(VERSION)

help: ## Show help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the service
	go build -o bin/$(SERVICE_NAME) cmd/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) -f deployments/docker/Dockerfile .

docker-push: docker-build ## Push Docker image
	docker push $(DOCKER_IMAGE)

dev: ## Run development server
	go run cmd/main.go

deploy-dev: ## Deploy to development
	kubectl apply -f deployments/k8s/ --namespace=dev

deploy-staging: ## Deploy to staging
	kubectl apply -f deployments/k8s/ --namespace=staging

deploy-prod: ## Deploy to production
	kubectl apply -f deployments/k8s/ --namespace=prod

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

deps: ## Download dependencies
	go mod download
	go mod tidy

