.PHONY: help setup dev prod monitor status logs stop clean

help: ## Show help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Initial setup for local development
	@echo "🏗️ Setting up local development environment..."
	@chmod +x scripts/setup-local-dev.sh
	@./scripts/setup-local-dev.sh

dev: ## Start development environment
	@echo "🚀 Starting development environment..."
	docker-compose -f docker-compose/development.yml up --build

dev-bg: ## Start development environment in background
	@echo "🚀 Starting development environment in background..."
	docker-compose -f docker-compose/development.yml up --build -d

prod: ## Start production environment
	@echo "🚀 Starting production environment..."
	docker-compose -f docker-compose/production.yml up --build -d

monitor: ## Start monitoring stack only
	@echo "📊 Starting monitoring stack..."
	docker-compose -f docker-compose/monitoring.yml up --build -d

status: ## Show services status
	@echo "📊 Services Status:"
	docker-compose -f docker-compose/development.yml ps

logs: ## Show logs from all services
	docker-compose -f docker-compose/development.yml logs -f

logs-app: ## Show application logs only
	docker-compose -f docker-compose/development.yml logs -f auth-service company-service api-gateway

logs-monitor: ## Show monitoring logs only
	docker-compose -f docker-compose/development.yml logs -f kibana grafana prometheus

stop: ## Stop all services
	@echo "🛑 Stopping all services..."
	docker-compose -f docker-compose/development.yml down

clean: ## Clean up everything
	@echo "🧹 Cleaning up..."
	docker-compose -f docker-compose/development.yml down -v
	docker-compose -f docker-compose/production.yml down -v
	docker system prune -f

health: ## Check health of all services
	@echo "🏥 Checking service health..."
	@curl -s http://localhost:8080/health && echo "✅ Auth Service" || echo "❌ Auth Service"
	@curl -s http://localhost:8081/health && echo "✅ Company Service" || echo "❌ Company Service"  
	@curl -s http://localhost:8000/health && echo "✅ API Gateway" || echo "❌ API Gateway"
	@curl -s http://localhost:5601/api/status && echo "✅ Kibana" || echo "❌ Kibana"
	@curl -s http://localhost:3000/api/health && echo "✅ Grafana" || echo "❌ Grafana"

update: ## Update all services
	@echo "🔄 Updating all services..."
	@cd ../jobfair-auth-service && git pull
	@cd ../jobfair-company-service && git pull
	@cd ../jobfair-api-gateway && git pull
	@cd ../jobfair-shared-libs && git pull
	@echo "✅ All repositories updated"
	$(MAKE) dev-bg