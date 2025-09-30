#!/bin/bash

echo "🚀 Setting up Jobfair Local Development Environment"

# Check if we're in the right directory
if [ ! -d "../jobfair-auth-service" ] || [ ! -d "../jobfair-company-service" ]; then
    echo "❌ Please run this script from jobfair-infrastructure directory"
    echo "📁 Expected structure:"
    echo "   jobfair-ecosystem/"
    echo "   ├── jobfair-auth-service/"
    echo "   ├── jobfair-"
    echo "   ├── jobfair-infrastructure/  <- You are here"
    echo "   └── ..."
    exit 1
fi

echo "✅ Repository structure looks good"

# Create monitoring directories if they don't exist
echo "📁 Setting up monitoring directories..."
mkdir -p monitoring/{prometheus,grafana/{provisioning/datasources,dashboards},logstash/pipeline,filebeat}

# Check Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed"
    exit 1
fi

echo "✅ Docker is available"

# Stop existing services
echo "🛑 Stopping existing services..."
docker-compose -f docker-compose/development.yml down

# Start development environment
echo "🔨 Starting development environment..."
docker-compose -f docker-compose/development.yml up --build -d

echo "⏳ Waiting for services to start..."
sleep 30

# Health check
echo "🏥 Checking service health..."
check_service() {
    local service_name=$1
    local url=$2
    
    if curl -s "$url" > /dev/null; then
        echo "✅ $service_name is healthy"
    else
        echo "❌ $service_name is not responding"
    fi
}

check_service "Auth Service" "http://localhost:8080/health"
check_service "Company Service" "http://localhost:8081/health"
check_service "API Gateway" "http://localhost:8000/health"
check_service "Kibana" "http://localhost:5601/api/status"
check_service "Grafana" "http://localhost:3000/api/health"

echo ""
echo "🎯 Development Environment URLs:"
echo "┌─────────────────────────────────────────────────────────────┐"
echo "│ 🔐 Auth Service:      http://localhost:8080                 │"
echo "│ 🏢 Company Service:   http://localhost:8081                 │"
echo "│ 🌐 API Gateway:       http://localhost:8000                 │"
echo "│                                                             │"
echo "│ 📊 Kibana:           http://localhost:5601                 │"
echo "│ 📈 Grafana:          http://localhost:3000 (admin/jobfair123)│"
echo "│ 🔍 Prometheus:       http://localhost:9090                 │"
echo "└─────────────────────────────────────────────────────────────┘"
echo ""
echo "💡 Tip: Use 'docker-compose -f docker-compose/development.yml logs -f' to view logs"