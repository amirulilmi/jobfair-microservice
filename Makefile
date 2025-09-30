# File: Makefile (di root project)
.PHONY: help dev prod status logs stop clean test

help: ## Show help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Start all services in development
	@echo "🚀 Starting all services..."
	docker-compose up --build

dev-bg: ## Start all services in background
	@echo "🚀 Starting all services in background..."
	docker-compose up --build -d

status: ## Show services status
	@echo "📊 Services Status:"
	docker-compose ps

logs: ## Show all services logs
	docker-compose logs -f

logs-auth: ## Show auth service logs
	docker-compose logs -f auth-service

logs-company: ## Show company service logs
	docker-compose logs -f company-service

stop: ## Stop all services
	@echo "🛑 Stopping all services..."
	docker-compose down

clean: ## Clean up everything
	@echo "🧹 Cleaning up..."
	docker-compose down -v
	docker system prune -f

test-all: ## Test all services endpoints
	@echo "🧪 Testing all services..."
	@echo "Auth Service:"
	@curl -s http://localhost:8080/health || echo "Auth service not responding"
	@echo "\nCompany Service:"
	@curl -s http://localhost:8081/health || echo "Company service not responding"