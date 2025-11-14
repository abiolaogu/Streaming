# StreamVerse Docker Deployment Guide

Complete guide for deploying StreamVerse using Docker and Docker Compose.

---

## Table of Contents

1. [Quick Start (MVP)](#quick-start-mvp)
2. [Production Deployment](#production-deployment)
3. [Docker Architecture](#docker-architecture)
4. [Configuration](#configuration)
5. [Scaling](#scaling)
6. [Monitoring](#monitoring)
7. [Troubleshooting](#troubleshooting)

---

## Quick Start (MVP)

Deploy a minimal viable product for testing in under 5 minutes.

### Prerequisites

- Docker 24.0+
- Docker Compose 2.0+
- 8GB RAM minimum
- 20GB disk space

### Installation

```bash
# 1. Clone repository
git clone https://github.com/yourusername/streamverse.git
cd streamverse

# 2. Create environment file
cp .env.example .env

# 3. Edit .env and add your Gemini API key
nano .env  # or vim/code

# 4. Run automated deployment
chmod +x scripts/deploy-mvp.sh
./scripts/deploy-mvp.sh
```

### What's Included in MVP

- ✅ PostgreSQL database
- ✅ Redis cache
- ✅ Auth Service (authentication)
- ✅ Training Bot Service (AI assistant)
- ✅ Web Application (React frontend)

### Access MVP

Once deployed:

- **Web App**: http://localhost:3000
- **Auth API**: http://localhost:8081
- **Training Bot**: http://localhost:8096

### Manual MVP Deployment

```bash
# Start services
docker-compose -f docker-compose-mvp.yml up -d

# View logs
docker-compose -f docker-compose-mvp.yml logs -f

# Stop services
docker-compose -f docker-compose-mvp.yml down
```

---

## Production Deployment

Full production deployment with all 15+ microservices.

### Prerequisites

- Docker 24.0+
- Docker Compose 2.20+
- 16GB RAM minimum (32GB recommended)
- 100GB disk space
- Domain name (for SSL)
- Cloud provider account (AWS/GCP/Azure)

### Step 1: Environment Setup

```bash
# Copy environment template
cp .env.example .env
```

**Edit `.env` and configure:**

```bash
# Required Variables
POSTGRES_PASSWORD=your_strong_password_here
MONGO_PASSWORD=your_strong_password_here
REDIS_PASSWORD=your_strong_password_here
JWT_SECRET=generate_32_character_secret_here
GEMINI_API_KEY=your_gemini_api_key_here

# Optional but Recommended
STRIPE_API_KEY=your_stripe_key
CDN_BASE_URL=https://cdn.yourdomain.com
API_BASE_URL=https://api.yourdomain.com
```

**Generate secure secrets:**

```bash
# Generate JWT secret
openssl rand -hex 32

# Generate passwords
openssl rand -base64 32
```

### Step 2: Build & Deploy

```bash
# Automated deployment
./scripts/deploy-production.sh
```

**Or manually:**

```bash
# Build all images
docker-compose build --parallel

# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

### Step 3: Initialize Database

```bash
# Wait for PostgreSQL to be ready
docker-compose exec postgres pg_isready

# Run migrations (if needed)
docker-compose exec postgres psql -U streamverse -d streamverse -f /docker-entrypoint-initdb.d/init.sql
```

### Step 4: Verify Deployment

```bash
# Check service health
curl http://localhost:8081/health  # Auth Service
curl http://localhost:8083/health  # Content Service
curl http://localhost:8096/health  # Training Bot
curl http://localhost/health       # Web App

# Check databases
docker-compose exec postgres psql -U streamverse -c "SELECT 1"
docker-compose exec mongodb mongosh --eval "db.runCommand({ping: 1})"
docker-compose exec redis redis-cli ping
```

### Services Included

**Core Services:**
- Auth Service (8081)
- User Service (8082)
- Content Service (8083)
- Streaming Service (8084)
- Payment Service (8085)
- Training Bot (8096)

**Data Services:**
- PostgreSQL (5432)
- MongoDB (27017)
- Redis (6379)
- Elasticsearch (9200)

**Messaging:**
- Kafka (9092)
- Zookeeper (2181)

**Monitoring:**
- Prometheus (9090)
- Grafana (3001)

**Frontend:**
- Web Application (80, 443)

---

## Docker Architecture

### Multi-Stage Builds

All services use multi-stage builds for optimized images:

```dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o service

# Stage 2: Production
FROM alpine:latest
COPY --from=builder /build/service .
CMD ["./service"]
```

**Benefits:**
- Smaller image sizes (20-50MB vs 1GB+)
- Faster deployments
- Improved security (no build tools in production)

### Service Templates

**Go Services:** `docker/Dockerfile.go-service`
- Used for: auth, user, content, streaming, payment, etc.
- Build time: 2-3 minutes
- Image size: ~30MB

**Python Services:** `docker/Dockerfile.python-service`
- Used for: recommendation, ML services
- Build time: 3-5 minutes
- Image size: ~150MB

**Node Services:** Root `Dockerfile`
- Used for: web frontend
- Build time: 5-8 minutes
- Image size: ~50MB (with nginx)

### Networks

```yaml
networks:
  default:
    name: streamverse
    driver: bridge
```

All services communicate on the `streamverse` network.

### Volumes

```yaml
volumes:
  postgres_data:     # PostgreSQL data
  mongodb_data:      # MongoDB data
  redis_data:        # Redis data
  kafka_data:        # Kafka logs
  prometheus_data:   # Metrics
  grafana_data:      # Dashboards
```

---

## Configuration

### Environment Variables

**Database:**
```bash
POSTGRES_DB=streamverse
POSTGRES_USER=streamverse
POSTGRES_PASSWORD=<secure-password>
MONGO_PASSWORD=<secure-password>
REDIS_PASSWORD=<secure-password>
```

**Security:**
```bash
JWT_SECRET=<32-character-secret>
JWT_EXPIRY=3600  # seconds
```

**External APIs:**
```bash
GEMINI_API_KEY=<your-key>
STRIPE_API_KEY=<your-key>
STRIPE_WEBHOOK_SECRET=<your-secret>
```

**Infrastructure:**
```bash
CDN_BASE_URL=https://cdn.streamverse.io
DRM_LICENSE_URL=https://drm.streamverse.io
API_BASE_URL=https://api.streamverse.io
```

### Service Configuration

Each service can be configured via environment variables:

```yaml
auth-service:
  environment:
    - PORT=8081
    - GRPC_PORT=9081
    - DB_HOST=postgres
    - REDIS_HOST=redis
    - JWT_SECRET=${JWT_SECRET}
```

### Resource Limits

Add resource limits to prevent service overconsumption:

```yaml
services:
  auth-service:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

---

## Scaling

### Horizontal Scaling

Scale specific services:

```bash
# Scale auth service to 3 replicas
docker-compose up -d --scale auth-service=3

# Scale multiple services
docker-compose up -d --scale auth-service=3 --scale content-service=3
```

### Load Balancing

Add nginx as load balancer:

```yaml
nginx-lb:
  image: nginx:alpine
  volumes:
    - ./nginx-lb.conf:/etc/nginx/nginx.conf
  ports:
    - "80:80"
  depends_on:
    - auth-service
```

**nginx-lb.conf:**
```nginx
upstream auth-backend {
    least_conn;
    server auth-service-1:8081;
    server auth-service-2:8081;
    server auth-service-3:8081;
}

server {
    location /api/auth {
        proxy_pass http://auth-backend;
    }
}
```

### Auto-Scaling

For cloud deployments, use orchestration platforms:
- **Kubernetes**: Full docs in `infrastructure/k8s/`
- **Docker Swarm**: Convert compose file
- **AWS ECS**: Task definitions
- **GCP Cloud Run**: Deploy from images

---

## Monitoring

### Prometheus

Access metrics at: http://localhost:9090

**Key Metrics:**
- Service health: `up{job="streamverse"}`
- Request rate: `rate(http_requests_total[5m])`
- Error rate: `rate(http_errors_total[5m])`
- Response time: `http_request_duration_seconds`

### Grafana

Access dashboards at: http://localhost:3001

**Default Credentials:**
- Username: `admin`
- Password: Set in `GRAFANA_PASSWORD` env var

**Pre-configured Dashboards:**
- System Overview
- Service Metrics
- Database Performance
- Error Tracking

### Logs

**View logs for specific service:**
```bash
docker-compose logs -f auth-service
```

**View logs for all services:**
```bash
docker-compose logs -f
```

**Save logs to file:**
```bash
docker-compose logs > logs.txt
```

### Health Checks

All services include health checks:

```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "http://localhost:8081/health"]
  interval: 30s
  timeout: 3s
  retries: 3
```

Check health status:
```bash
docker-compose ps
```

---

## Troubleshooting

### Common Issues

#### Services Won't Start

**Check logs:**
```bash
docker-compose logs <service-name>
```

**Check dependencies:**
```bash
docker-compose ps
```

Ensure databases are healthy before starting application services.

#### Database Connection Errors

**Test PostgreSQL:**
```bash
docker-compose exec postgres psql -U streamverse -c "SELECT 1"
```

**Test MongoDB:**
```bash
docker-compose exec mongodb mongosh --eval "db.runCommand({ping: 1})"
```

**Test Redis:**
```bash
docker-compose exec redis redis-cli -a $REDIS_PASSWORD ping
```

#### Port Conflicts

Check if ports are already in use:
```bash
lsof -i :8081  # Linux/Mac
netstat -ano | findstr :8081  # Windows
```

Stop conflicting services or change ports in `docker-compose.yml`.

#### Out of Memory

**Check Docker memory:**
```bash
docker stats
```

**Increase Docker memory limit:**
- Docker Desktop: Settings → Resources → Memory
- Recommended: 8GB minimum, 16GB+ for production

#### Image Build Failures

**Clear build cache:**
```bash
docker-compose build --no-cache
```

**Remove old images:**
```bash
docker image prune -a
```

#### Service Crashes

**Check service logs:**
```bash
docker-compose logs <service-name>
```

**Restart specific service:**
```bash
docker-compose restart <service-name>
```

**Restart all services:**
```bash
docker-compose restart
```

### Debugging

**Access service shell:**
```bash
docker-compose exec <service-name> sh
```

**Run commands in service:**
```bash
docker-compose exec postgres psql -U streamverse
docker-compose exec redis redis-cli
docker-compose exec mongodb mongosh
```

**Inspect service:**
```bash
docker inspect <container-name>
```

---

## Best Practices

### Production Deployment

1. **Use secrets management:**
   - Don't commit `.env` to git
   - Use Docker secrets or cloud secret managers
   - Rotate secrets regularly

2. **Enable SSL/TLS:**
   - Use Let's Encrypt for certificates
   - Configure nginx with HTTPS
   - Redirect HTTP to HTTPS

3. **Set resource limits:**
   - Prevent resource exhaustion
   - Enable OOM killer
   - Monitor resource usage

4. **Implement backups:**
   ```bash
   # PostgreSQL backup
   docker-compose exec postgres pg_dump -U streamverse streamverse > backup.sql

   # MongoDB backup
   docker-compose exec mongodb mongodump --out=/backup

   # Volume backups
   docker run --rm -v streamverse_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz /data
   ```

5. **Monitor and alert:**
   - Set up Prometheus alerts
   - Configure Grafana notifications
   - Monitor logs centrally

### Development

1. **Use MVP setup:**
   - Faster startup time
   - Minimal resource usage
   - Quick iteration

2. **Hot reload:**
   - Mount source code as volumes
   - Use development images
   - Enable debug mode

3. **Clean up regularly:**
   ```bash
   # Remove stopped containers
   docker-compose down

   # Remove unused images
   docker image prune

   # Remove unused volumes (CAUTION: deletes data)
   docker volume prune
   ```

---

## Cheat Sheet

```bash
# Start MVP
./scripts/deploy-mvp.sh

# Start Production
./scripts/deploy-production.sh

# Build services
docker-compose build

# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# Restart service
docker-compose restart <service>

# View logs
docker-compose logs -f <service>

# Check status
docker-compose ps

# Scale service
docker-compose up -d --scale <service>=3

# Access shell
docker-compose exec <service> sh

# Update service
docker-compose up -d --no-deps --build <service>

# Clean up
docker-compose down -v  # Removes volumes too

# Backup database
docker-compose exec postgres pg_dump -U streamverse > backup.sql

# Restore database
docker-compose exec -T postgres psql -U streamverse < backup.sql
```

---

## Additional Resources

- **Docker Documentation**: https://docs.docker.com
- **Docker Compose**: https://docs.docker.com/compose
- **StreamVerse Architecture**: [ARCHITECTURAL_BLUEPRINT.md](ARCHITECTURAL_BLUEPRINT.md)
- **Deployment Guide**: [DEPLOYMENT_AND_INTEGRATION_GUIDE.md](DEPLOYMENT_AND_INTEGRATION_GUIDE.md)
- **Kubernetes Deployment**: `infrastructure/k8s/README.md`

---

**Document Version**: 2.0
**Last Updated**: 2025
**Status**: Production Ready

**Happy Deploying! 🚀**
