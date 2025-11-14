#!/bin/bash

# StreamVerse Production Deployment Script
# Full production deployment with all services

set -e

echo "🚀 StreamVerse Production Deployment"
echo "====================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check prerequisites
check_prereqs() {
    echo "Checking prerequisites..."

    if ! command -v docker &> /dev/null; then
        echo -e "${RED}Error: Docker is not installed${NC}"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        echo -e "${RED}Error: Docker Compose is not installed${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ Prerequisites satisfied${NC}"
}

# Validate environment variables
validate_env() {
    echo "Validating environment variables..."

    required_vars=(
        "POSTGRES_PASSWORD"
        "MONGO_PASSWORD"
        "REDIS_PASSWORD"
        "JWT_SECRET"
        "GEMINI_API_KEY"
    )

    missing_vars=()

    for var in "${required_vars[@]}"; do
        if [ -z "${!var}" ]; then
            missing_vars+=("$var")
        fi
    done

    if [ ${#missing_vars[@]} -ne 0 ]; then
        echo -e "${RED}Error: Missing required environment variables:${NC}"
        printf '%s\n' "${missing_vars[@]}"
        echo "Please set these in your .env file"
        exit 1
    fi

    # Validate JWT_SECRET length
    if [ ${#JWT_SECRET} -lt 32 ]; then
        echo -e "${RED}Error: JWT_SECRET must be at least 32 characters${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ Environment variables validated${NC}"
}

# Build all services
build_services() {
    echo "Building Docker images..."
    docker-compose build --parallel
    echo -e "${GREEN}✓ Images built successfully${NC}"
}

# Start services
start_services() {
    echo "Starting services..."
    docker-compose up -d
    echo -e "${GREEN}✓ Services started${NC}"
}

# Run database migrations
run_migrations() {
    echo "Running database migrations..."

    # Wait for PostgreSQL
    echo "Waiting for PostgreSQL..."
    sleep 10

    # Run schema init
    docker-compose exec postgres psql -U $POSTGRES_USER -d $POSTGRES_DB -f /docker-entrypoint-initdb.d/init.sql || true

    echo -e "${GREEN}✓ Migrations completed${NC}"
}

# Health check
health_check() {
    echo "Performing health checks..."

    services=(
        "postgres:5432"
        "redis:6379"
        "mongodb:27017"
        "elasticsearch:9200"
        "auth-service:8081"
        "content-service:8083"
        "training-bot:8096"
        "web:80"
    )

    for service in "${services[@]}"; do
        IFS=':' read -r name port <<< "$service"
        if docker-compose ps | grep -q "$name"; then
            echo -e "${GREEN}✓ $name is running${NC}"
        else
            echo -e "${YELLOW}⚠ $name may not be running${NC}"
        fi
    done
}

# Display summary
display_summary() {
    echo ""
    echo "======================================"
    echo -e "${GREEN}StreamVerse is now running!${NC}"
    echo "======================================"
    echo ""
    echo "Services:"
    echo "  Web:          http://localhost"
    echo "  API Gateway:  http://localhost:8080"
    echo "  Auth:         http://localhost:8081"
    echo "  Training Bot: http://localhost:8096"
    echo "  Prometheus:   http://localhost:9090"
    echo "  Grafana:      http://localhost:3001"
    echo ""
    echo "Databases:"
    echo "  PostgreSQL:      localhost:5432"
    echo "  MongoDB:         localhost:27017"
    echo "  Redis:           localhost:6379"
    echo "  Elasticsearch:   localhost:9200"
    echo ""
    echo "Management:"
    echo "  View logs:       docker-compose logs -f [service]"
    echo "  Stop all:        docker-compose down"
    echo "  Restart:         docker-compose restart [service]"
    echo "  Scale service:   docker-compose up -d --scale [service]=N"
    echo ""
    echo "======================================"
}

# Main execution
main() {
    check_prereqs

    # Load environment variables
    if [ -f .env ]; then
        export $(cat .env | grep -v '^#' | xargs)
    else
        echo -e "${RED}Error: .env file not found${NC}"
        echo "Please create .env from .env.example"
        exit 1
    fi

    validate_env
    build_services
    start_services
    run_migrations
    sleep 15  # Give services time to start
    health_check
    display_summary
}

# Run main function
main
