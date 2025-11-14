# StreamVerse - Complete Deployment and Integration Guide

This guide provides end-to-end instructions for deploying the StreamVerse platform across all infrastructure, backend services, and client applications.

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Infrastructure Deployment](#infrastructure-deployment)
4. [Backend Services Deployment](#backend-services-deployment)
5. [Client Applications Deployment](#client-applications-deployment)
6. [Integration & Testing](#integration--testing)
7. [Monitoring & Operations](#monitoring--operations)
8. [Scaling & Performance](#scaling--performance)
9. [Troubleshooting](#troubleshooting)

---

## Overview

StreamVerse is a comprehensive streaming platform with:
- **Backend**: 15+ microservices (Go, Python, TypeScript)
- **Frontend**: Web app (React/Next.js)
- **Mobile**: Flutter app (iOS & Android)
- **TV**: 10 TV platform apps

### Architecture Summary

```
┌─────────────┐
│   Clients   │ → Web, Mobile (iOS/Android), TV (10 platforms)
└──────┬──────┘
       │
┌──────▼──────────┐
│   API Gateway   │ → Kong (Authentication, Rate Limiting, Routing)
└──────┬──────────┘
       │
┌──────▼───────────┐
│  Microservices   │ → 15+ services (Auth, Content, Streaming, etc.)
└──────┬───────────┘
       │
┌──────▼───────────┐
│   Data Layer     │ → PostgreSQL, MongoDB, Redis, Elasticsearch
└──────────────────┘
```

---

## Prerequisites

### Required Tools

- **Cloud Accounts**:
  - AWS account (primary) or GCP/Azure
  - Domain name and DNS management

- **Development Tools**:
  - Docker 24.0+
  - Kubernetes 1.28+
  - Terraform 1.6+
  - Helm 3.13+
  - kubectl CLI

- **Programming Environments**:
  - Go 1.21+
  - Python 3.11+
  - Node.js 18+
  - Flutter 3.16+ (for mobile)

### Access Requirements

- Cloud provider credentials (AWS/GCP/Azure)
- GitHub repository access
- Docker Hub or private container registry
- CDN provider account (CloudFlare, CloudFront)
- DRM provider credentials (for content protection)

---

## Infrastructure Deployment

### Phase 1: Core Infrastructure

#### Step 1: Clone Repository

```bash
git clone https://github.com/yourusername/streamverse.git
cd streamverse
```

#### Step 2: Configure Terraform

```bash
cd infrastructure/terraform/aws  # or gcp/azure

# Copy and edit variables
cp terraform.tfvars.example terraform.tfvars
vim terraform.tfvars
```

**terraform.tfvars**:
```hcl
# AWS Configuration
aws_region = "us-east-1"
environment = "production"
project_name = "streamverse"

# Network Configuration
vpc_cidr = "10.0.0.0/16"
availability_zones = ["us-east-1a", "us-east-1b", "us-east-1c"]

# EKS Configuration
eks_cluster_version = "1.28"
eks_node_instance_types = ["t3.xlarge", "t3.2xlarge"]
eks_min_nodes = 3
eks_max_nodes = 20

# RDS Configuration
rds_instance_class = "db.r6g.2xlarge"
rds_allocated_storage = 500
rds_engine_version = "14.10"

# ElastiCache Configuration
elasticache_node_type = "cache.r6g.xlarge"
elasticache_num_nodes = 3

# Domain Configuration
domain_name = "streamverse.io"
```

#### Step 3: Deploy Infrastructure

```bash
# Initialize Terraform
terraform init

# Plan deployment
terraform plan -out=tfplan

# Apply infrastructure
terraform apply tfplan
```

**Expected Resources Created**:
- VPC with public/private subnets
- EKS cluster with node groups
- RDS PostgreSQL (multi-AZ)
- ElastiCache Redis cluster
- MSK Kafka cluster
- S3 buckets for media storage
- ECR repositories
- CloudFront CDN
- Route53 DNS zones
- Security groups & IAM roles

#### Step 4: Configure kubectl

```bash
# Get EKS cluster credentials
aws eks update-kubeconfig \
  --region us-east-1 \
  --name streamverse-production

# Verify connection
kubectl cluster-info
kubectl get nodes
```

### Phase 2: Kubernetes Setup

#### Step 1: Install Helm

```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

#### Step 2: Install Core Components

```bash
cd infrastructure/k8s

# Install NGINX Ingress Controller
kubectl apply -f ingress-nginx/

# Install cert-manager for SSL
kubectl apply -f cert-manager/

# Install Kong API Gateway
kubectl apply -f kong/

# Install Istio Service Mesh (optional but recommended)
istioctl install --set profile=production -y
kubectl label namespace default istio-injection=enabled
```

#### Step 3: Set up Monitoring

```bash
# Install Prometheus
kubectl apply -f monitoring/prometheus/

# Install Grafana
kubectl apply -f monitoring/grafana/

# Install Loki for logging
kubectl apply -f monitoring/loki/

# Install Jaeger for tracing
kubectl apply -f monitoring/jaeger/
```

#### Step 4: Configure Secrets

```bash
# Create namespace
kubectl create namespace streamverse

# Install Vault for secrets management
kubectl apply -f vault/

# Or use Kubernetes secrets
kubectl create secret generic api-secrets \
  --from-literal=jwt-secret=YOUR_JWT_SECRET \
  --from-literal=db-password=YOUR_DB_PASSWORD \
  --from-literal=redis-password=YOUR_REDIS_PASSWORD \
  -n streamverse
```

---

## Backend Services Deployment

### Phase 1: Build Docker Images

#### Step 1: Set Up Container Registry

```bash
# AWS ECR example
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin \
  ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com
```

#### Step 2: Build All Services

```bash
# Build script
./scripts/build-all-services.sh
```

Or manually for each service:

```bash
# Example: Auth Service
cd services/auth-service
docker build -t streamverse/auth-service:latest .
docker tag streamverse/auth-service:latest \
  ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/auth-service:latest
docker push ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/auth-service:latest
```

Repeat for all services:
- auth-service
- user-service
- content-service
- streaming-service
- payment-service
- transcoding-service
- search-service
- recommendation-service
- analytics-service
- notification-service
- ad-service
- admin-service
- websocket-service
- scheduler-service

### Phase 2: Deploy Services to Kubernetes

#### Step 1: Deploy Databases

```bash
# Deploy PostgreSQL (if not using managed RDS)
kubectl apply -f infrastructure/k8s/databases/postgresql.yaml

# Deploy MongoDB
kubectl apply -f infrastructure/k8s/databases/mongodb.yaml

# Deploy Redis
kubectl apply -f infrastructure/k8s/databases/redis.yaml

# Deploy Elasticsearch
kubectl apply -f infrastructure/k8s/databases/elasticsearch.yaml
```

#### Step 2: Deploy Message Brokers

```bash
# Deploy Kafka (if not using MSK)
kubectl apply -f infrastructure/k8s/kafka/

# Deploy NATS
kubectl apply -f infrastructure/k8s/nats/
```

#### Step 3: Deploy Microservices

```bash
# Deploy all services
kubectl apply -f infrastructure/k8s/services/

# Or deploy individually
kubectl apply -f infrastructure/k8s/services/auth-service.yaml
kubectl apply -f infrastructure/k8s/services/content-service.yaml
# ... repeat for all services
```

**Example Service Deployment (auth-service.yaml)**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: streamverse
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/auth-service:latest
        ports:
        - containerPort: 8081
        - containerPort: 9081
        env:
        - name: DB_HOST
          value: postgresql.streamverse.svc.cluster.local
        - name: REDIS_HOST
          value: redis.streamverse.svc.cluster.local
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: jwt-secret
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: auth-service
  namespace: streamverse
spec:
  selector:
    app: auth-service
  ports:
  - name: http
    port: 80
    targetPort: 8081
  - name: grpc
    port: 9081
    targetPort: 9081
  type: ClusterIP
```

#### Step 4: Verify Deployments

```bash
# Check all pods are running
kubectl get pods -n streamverse

# Check services
kubectl get svc -n streamverse

# Check logs
kubectl logs -f deployment/auth-service -n streamverse
```

### Phase 3: Configure API Gateway

```bash
# Apply Kong configuration
kubectl apply -f infrastructure/kong/kong-config.yaml

# Configure routes
kubectl apply -f infrastructure/kong/routes/

# Test API Gateway
curl https://api.streamverse.io/v1/health
```

---

## Client Applications Deployment

### 1. Web Application (React/Next.js)

#### Build & Deploy

```bash
cd apps/clients/web

# Install dependencies
npm install

# Build for production
npm run build

# Deploy to Vercel (recommended)
vercel deploy --prod

# Or deploy to S3 + CloudFront
aws s3 sync out/ s3://streamverse-web/
aws cloudfront create-invalidation \
  --distribution-id DISTRIBUTION_ID \
  --paths "/*"
```

### 2. Mobile Application (Flutter)

#### iOS Deployment

```bash
cd apps/clients/mobile-flutter

# Build iOS
flutter build ios --release

# Archive and upload to App Store Connect
# Use Xcode or Fastlane for automation
fastlane ios release
```

#### Android Deployment

```bash
# Build Android App Bundle
flutter build appbundle --release

# Upload to Google Play Console
# Use Fastlane for automation
fastlane android deploy
```

### 3. TV Applications

#### Android TV
```bash
cd apps/clients/tv-apps/android-tv

# Build APK
./gradlew assembleRelease

# Upload to Google Play Console (TV section)
```

#### Samsung Tizen
```bash
cd apps/clients/tv-apps/samsung-tizen

# Package app
tizen package -t wgt -s YOUR_CERTIFICATE

# Submit to Samsung TV App Store
```

#### LG webOS
```bash
cd apps/clients/tv-apps/lg-webos

# Package app
ares-package .

# Submit to LG Content Store
```

#### Roku
```bash
cd apps/clients/tv-apps/roku

# Package via Roku Developer Dashboard
# Upload to Roku Channel Store
```

#### Apple tvOS
```bash
cd apps/clients/tv-apps/apple-tvos

# Build and archive
xcodebuild -workspace StreamVerse.xcworkspace \
           -scheme StreamVerse \
           -archivePath StreamVerse.xcarchive \
           archive

# Submit to App Store Connect
```

#### Fire TV
```bash
cd apps/clients/tv-apps/amazon-fire-tv

# Build APK
./gradlew assembleRelease

# Submit to Amazon Appstore
```

#### Other HTML5-Based Platforms
Follow platform-specific packaging instructions in respective READMEs.

---

## Integration & Testing

### 1. End-to-End Testing

```bash
# Run integration tests
cd tests
npm run test:integration

# Run load tests
cd load-test
k6 run load-test.js
```

### 2. Verify All Integrations

```bash
# Test authentication
curl -X POST https://api.streamverse.io/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'

# Test content API
curl https://api.streamverse.io/v1/content/home \
  -H "Authorization: Bearer YOUR_TOKEN"

# Test streaming
curl https://api.streamverse.io/v1/streaming/CONTENT_ID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. Cross-Platform Testing

Test on all platforms:
- [ ] Web (Chrome, Firefox, Safari, Edge)
- [ ] iOS app
- [ ] Android app
- [ ] Android TV
- [ ] Samsung Tizen TV
- [ ] LG webOS TV
- [ ] Roku TV
- [ ] Apple TV
- [ ] Fire TV
- [ ] Other TV platforms

---

## Monitoring & Operations

### Accessing Dashboards

```bash
# Grafana
kubectl port-forward svc/grafana 3000:80 -n monitoring
# Access: http://localhost:3000

# Prometheus
kubectl port-forward svc/prometheus 9090:9090 -n monitoring
# Access: http://localhost:9090

# Jaeger
kubectl port-forward svc/jaeger-query 16686:16686 -n monitoring
# Access: http://localhost:16686
```

### Setting Up Alerts

Configure alerts in `infrastructure/monitoring/alerts/`:
- High error rates
- Slow API response times
- Service downtime
- Database connection issues
- High memory/CPU usage

### Log Aggregation

```bash
# Query logs with Loki
logcli query '{app="auth-service"}' --since=1h
```

---

## Scaling & Performance

### Horizontal Pod Autoscaling

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-service-hpa
  namespace: streamverse
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Database Scaling

- **Read Replicas**: Add read replicas for heavy read workloads
- **Sharding**: Implement sharding for horizontal scaling
- **Connection Pooling**: Use PgBouncer for PostgreSQL

### CDN Optimization

- Configure cache rules in CloudFront/CloudFlare
- Set appropriate TTLs for different content types
- Enable HTTP/2 and HTTP/3

---

## Troubleshooting

### Common Issues

#### 1. Pods not starting
```bash
# Check pod status
kubectl describe pod POD_NAME -n streamverse

# Check logs
kubectl logs POD_NAME -n streamverse

# Common causes:
# - Image pull errors
# - Resource limits exceeded
# - Configuration issues
```

#### 2. Service connectivity issues
```bash
# Test service connectivity
kubectl run -it --rm debug --image=curlimages/curl --restart=Never -- \
  curl http://auth-service.streamverse.svc.cluster.local/health

# Check service endpoints
kubectl get endpoints -n streamverse
```

#### 3. Database connection errors
```bash
# Check database pods
kubectl get pods -n streamverse | grep postgres

# Test connection
kubectl exec -it auth-service-POD -n streamverse -- \
  psql -h postgresql.streamverse.svc.cluster.local -U streamverse
```

#### 4. High latency
- Check API Gateway logs
- Review service logs for slow queries
- Monitor database performance
- Check network policies

### Health Check Endpoints

All services expose:
- `/health` - Liveness probe
- `/ready` - Readiness probe
- `/metrics` - Prometheus metrics

---

## CI/CD Pipeline

### GitHub Actions Workflow

```yaml
name: Deploy to Production

on:
  push:
    branches: [ main ]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Build Docker Images
        run: ./scripts/build-all-services.sh

      - name: Push to ECR
        run: ./scripts/push-to-ecr.sh

      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/auth-service \
            auth-service=ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/auth-service:$GITHUB_SHA \
            -n streamverse
          kubectl rollout status deployment/auth-service -n streamverse
```

---

## Backup & Disaster Recovery

### Database Backups

```bash
# Automated daily backups
kubectl apply -f infrastructure/disaster-recovery/backup-cronjob.yaml

# Manual backup
kubectl exec -it postgresql-0 -n streamverse -- \
  pg_dump streamverse > backup.sql
```

### Disaster Recovery Plan

1. **RTO**: 1 hour
2. **RPO**: 15 minutes
3. **Multi-region**: Active-active in us-east-1 and us-west-2
4. **Failover**: Automatic DNS failover with Route53

---

## Security Checklist

- [ ] All services use TLS/SSL
- [ ] Secrets stored in Vault/Kubernetes Secrets
- [ ] Network policies configured
- [ ] Pod security policies enabled
- [ ] DRM configured for content protection
- [ ] Rate limiting enabled
- [ ] Input validation on all APIs
- [ ] Regular security audits
- [ ] Compliance certifications (SOC 2, GDPR, ISO 27001)

---

## Cost Optimization

- Use reserved instances for predictable workloads
- Implement autoscaling to scale down during low traffic
- Use spot instances for batch jobs
- Optimize S3 storage with lifecycle policies
- Monitor costs with AWS Cost Explorer

---

## Support & Documentation

- **Documentation**: [docs.streamverse.io](https://docs.streamverse.io)
- **API Reference**: [api.streamverse.io/docs](https://api.streamverse.io/docs)
- **Status Page**: [status.streamverse.io](https://status.streamverse.io)
- **Support Email**: support@streamverse.io

---

**Deployment Guide Version**: 2.0
**Last Updated**: 2025
**Status**: Production Ready

---

## Quick Reference Commands

```bash
# Check all services
kubectl get all -n streamverse

# Restart a service
kubectl rollout restart deployment/auth-service -n streamverse

# Scale a service
kubectl scale deployment/auth-service --replicas=10 -n streamverse

# View logs
kubectl logs -f deployment/auth-service -n streamverse --tail=100

# Execute command in pod
kubectl exec -it POD_NAME -n streamverse -- /bin/sh

# Port forward for debugging
kubectl port-forward svc/auth-service 8081:80 -n streamverse
```

---

This guide provides a complete deployment workflow for the StreamVerse platform. For platform-specific details, refer to individual service and client application READMEs.
