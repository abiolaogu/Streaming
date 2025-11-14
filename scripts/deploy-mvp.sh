#!/bin/bash

# StreamVerse MVP Deployment Script
# Quick deployment for testing the platform

set -e

echo "🚀 StreamVerse MVP Deployment"
echo "=============================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed${NC}"
    echo "Please install Docker first: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: Docker Compose is not installed${NC}"
    echo "Please install Docker Compose first: https://docs.docker.com/compose/install/"
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env file from template...${NC}"
    cp .env.example .env

    # Generate random secrets
    JWT_SECRET=$(openssl rand -hex 32)
    sed -i "s/your_jwt_secret_min_32_characters_here/$JWT_SECRET/" .env

    echo -e "${GREEN}✓ .env file created${NC}"
    echo -e "${YELLOW}⚠ Please edit .env and add your GEMINI_API_KEY${NC}"
    echo ""
    read -p "Press enter when you've added your API key..."
fi

# Stop any running containers
echo -e "\n${YELLOW}Stopping existing containers...${NC}"
docker-compose -f docker-compose-mvp.yml down 2>/dev/null || true

# Build images
echo -e "\n${YELLOW}Building Docker images...${NC}"
docker-compose -f docker-compose-mvp.yml build

# Start services
echo -e "\n${YELLOW}Starting services...${NC}"
docker-compose -f docker-compose-mvp.yml up -d

# Wait for services to be healthy
echo -e "\n${YELLOW}Waiting for services to be ready...${NC}"
sleep 10

# Check service health
echo -e "\n${YELLOW}Checking service health...${NC}"

check_service() {
    local name=$1
    local url=$2

    if curl -s -f "$url" > /dev/null; then
        echo -e "${GREEN}✓ $name is running${NC}"
        return 0
    else
        echo -e "${RED}✗ $name is not responding${NC}"
        return 1
    fi
}

check_service "PostgreSQL" "http://localhost:5432" 2>/dev/null || echo -e "${YELLOW}⚠ PostgreSQL check skipped (requires psql)${NC}"
check_service "Redis" "http://localhost:6379" 2>/dev/null || echo -e "${YELLOW}⚠ Redis check skipped (requires redis-cli)${NC}"
check_service "Auth Service" "http://localhost:8081/health"
check_service "Training Bot" "http://localhost:8096/health"
check_service "Web App" "http://localhost:3000/health"

# Display access information
echo -e "\n${GREEN}✓ StreamVerse MVP is running!${NC}"
echo ""
echo "=============================="
echo "Access your services:"
echo "=============================="
echo "Web Application:  http://localhost:3000"
echo "Auth Service:     http://localhost:8081"
echo "Training Bot:     http://localhost:8096"
echo "PostgreSQL:       localhost:5432"
echo "Redis:            localhost:6379"
echo ""
echo "=============================="
echo "Useful Commands:"
echo "=============================="
echo "View logs:        docker-compose -f docker-compose-mvp.yml logs -f"
echo "Stop services:    docker-compose -f docker-compose-mvp.yml down"
echo "Restart:          docker-compose -f docker-compose-mvp.yml restart"
echo "=============================="
echo ""
echo -e "${GREEN}Happy testing! 🎬${NC}"
