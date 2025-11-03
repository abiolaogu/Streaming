# StreamVerse Platform - Deployment Guide

Complete deployment guide for StreamVerse streaming platform.

## Prerequisites

- Kubernetes cluster (1.28+)
- kubectl configured
- Terraform installed
- Docker installed
- Access to container registry

## Architecture Overview

StreamVerse is a microservices-based streaming platform with:
- 12 Go microservices
- 2 Python services
- 2 Node.js services
- Kong API Gateway
- MongoDB, PostgreSQL, Redis
- Kafka message queue
- Monitoring stack (Prometheus, Grafana, Loki, Jaeger)

## Deployment Steps

### 1. Infrastructure Setup

#### Using Terraform (AWS)

```bash
cd infrastructure/terraform/aws
terraform init
terraform plan
terraform apply
```

#### Using Terraform (GCP)

```bash
cd infrastructure/terraform/gcp
terraform init
terraform plan -var="gcp_project_id=your-project-id"
terraform apply
```

#### Using Terraform (Azure)

```bash
cd infrastructure/terraform/azure
terraform init
terraform plan
terraform apply
```

### 2. Configure Kubernetes Access

```bash
# AWS EKS
aws eks update-kubeconfig --name streamverse-cluster --region us-east-1

# GCP GKE
gcloud container clusters get-credentials streamverse-cluster --region us-central1

# Azure AKS
az aks get-credentials --resource-group streamverse-rg --name streamverse-aks
```

### 3. Create Namespaces

```bash
kubectl apply -f k8s/base/namespaces.yaml
```

### 4. Set Up Secrets

```bash
# MongoDB
kubectl create secret generic mongodb-secret \
  --from-literal=uri="mongodb://..." \
  -n streamverse

# PostgreSQL
kubectl create secret generic postgres-secret \
  --from-literal=uri="postgresql://..." \
  -n streamverse

# JWT
kubectl create secret generic jwt-secret \
  --from-literal=secret-key="$(openssl rand -base64 32)" \
  -n streamverse

# Redis
kubectl create secret generic redis-secret \
  --from-literal=uri="redis://..." \
  -n streamverse
```

### 5. Deploy Services

```bash
# Deploy all services
kubectl apply -f k8s/base/auth-service/
kubectl apply -f k8s/base/content-service/
# ... deploy other services

# Or use Kustomize
kubectl apply -k k8s/overlays/production/
```

### 6. Deploy API Gateway

```bash
kubectl apply -f infrastructure/kong/kubernetes/kong-deployment.yaml
```

### 7. Deploy Monitoring Stack

```bash
# Using Helm
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack \
  -n monitoring --create-namespace

# Or using manifests
kubectl apply -f infrastructure/monitoring/
```

### 8. Verify Deployment

```bash
# Check all pods are running
kubectl get pods -n streamverse

# Run smoke tests
./tests/smoke-tests.sh

# Check service health
curl http://api.streamverse.io/health
```

## Configuration

### Environment Variables

Each service requires:
- Database connection strings
- JWT secret key
- Redis connection
- Service URLs

### Resource Limits

Default resource requests/limits:
- **Requests**: 256Mi memory, 250m CPU
- **Limits**: 512Mi memory, 500m CPU

Adjust based on load testing results.

## Scaling

### Manual Scaling

```bash
kubectl scale deployment/auth-service --replicas=5 -n streamverse
```

### Auto-scaling

HPA is configured for all services:
```bash
kubectl get hpa -n streamverse
```

## Monitoring

- **Grafana**: http://grafana.monitoring.svc.cluster.local:3000
- **Prometheus**: http://prometheus.monitoring.svc.cluster.local:9090
- **Jaeger**: http://jaeger-query.monitoring.svc.cluster.local:16686

## Troubleshooting

See [Runbooks](runbooks/README.md) for detailed troubleshooting procedures.

## Rollback

```bash
kubectl rollout undo deployment/auth-service -n streamverse
```

## Production Checklist

- [ ] All secrets configured
- [ ] Resource limits set appropriately
- [ ] Network policies applied
- [ ] Monitoring stack deployed
- [ ] Backups configured
- [ ] Disaster recovery tested
- [ ] Security scanning completed
- [ ] Load testing performed
- [ ] Documentation reviewed

