# StreamVerse SaaS - On-Premise Deployment Guide

## Overview

StreamVerse SaaS is deployed on **dedicated physical servers** with a **hybrid GPU architecture** that combines:
- **Local NVIDIA GPUs** for baseline transcoding workloads
- **Runpod.io cloud GPUs** for elastic scaling during peak demand

This guide covers the complete deployment process from bare metal to production-ready streaming platform.

---

## Hardware Requirements

### Minimum Production Setup

**Servers** (Bare Metal):
- **Quantity**: 3+ servers for high availability
- **CPU**: Intel Xeon or AMD EPYC (32+ cores per server)
- **RAM**: 128GB+ per server
- **Storage**:
  - Hot Storage: 2TB+ NVMe SSD array (RAID 10)
  - Warm Storage: 20TB+ HDD array (RAID 6)
  - Cold Storage: 100TB+ (Ceph cluster)
- **Network**: 10Gbps NIC minimum (25Gbps recommended)

**GPU Configuration** (Per Server):
- **Option 1 (Budget)**: 2-4x NVIDIA RTX 4090 (24GB VRAM each)
- **Option 2 (Professional)**: 2-4x NVIDIA RTX A6000 (48GB VRAM each)
- **Option 3 (Enterprise)**: 2-4x NVIDIA H100 PCIe (80GB VRAM each)

**Network Infrastructure**:
- 10Gbps switches
- Load balancer (HAProxy or NGINX)
- Firewall with DDoS protection

### Recommended Production Setup

- **6 servers** (3 for compute, 3 for storage/databases)
- **12-16 local GPUs** total across compute nodes
- **100TB+ hot storage** on NVMe arrays
- **500TB+ warm/cold storage** on HDD arrays
- **Dedicated 25Gbps network** backbone

---

## Software Prerequisites

### Operating System
- **Ubuntu Server 22.04 LTS** (recommended)
- **Rocky Linux 9** (alternative)
- **Debian 12** (alternative)

### Core Components
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Install NVIDIA Docker Runtime
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | \
  sudo tee /etc/apt/sources.list.d/nvidia-docker.list
sudo apt update
sudo apt install -y nvidia-docker2
sudo systemctl restart docker

# Install Kubernetes (Rancher RKE2)
curl -sfL https://get.rke2.io | sh -
sudo systemctl enable rke2-server.service
sudo systemctl start rke2-server.service

# Install kubectl
sudo snap install kubectl --classic

# Install Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

### NVIDIA Drivers and CUDA
```bash
# Install NVIDIA drivers (version 535+)
sudo apt install -y nvidia-driver-535

# Verify GPU detection
nvidia-smi

# Expected output:
# +-----------------------------------------------------------------------------+
# | NVIDIA-SMI 535.xx.xx    Driver Version: 535.xx.xx    CUDA Version: 12.2     |
# |-------------------------------+----------------------+----------------------+
# | GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
# |   0  NVIDIA RTX 4090     Off  | 00000000:01:00.0 Off |                  N/A |
# +-------------------------------+----------------------+----------------------+
```

---

## Runpod.io Setup

### 1. Create Runpod.io Account
1. Go to https://www.runpod.io
2. Sign up for an account
3. Navigate to **Console** → **User Settings** → **API Keys**
4. Create a new API key and save it securely

### 2. Configure Runpod Endpoint (Optional)
For serverless GPU pods, create a custom endpoint:

1. Go to **Serverless** → **Endpoints**
2. Click **New Endpoint**
3. Configure:
   - **Name**: `streamverse-transcoding`
   - **GPU Type**: RTX 4090 / A100 / H100
   - **Container Image**: `registry.streamverse.io/streamverse/transcoding-service:latest`
   - **Container Port**: `8101`
   - **Environment Variables**:
     ```
     GPU_ACCELERATION=true
     NVIDIA_VISIBLE_DEVICES=all
     ```

4. Save the **Endpoint ID** for configuration

### 3. Test Runpod.io Connection
```bash
# Install Runpod CLI
pip install runpod

# Test authentication
export RUNPOD_API_KEY="your_api_key_here"
runpod list gpu-types

# Expected output: List of available GPU types
```

---

## Deployment Steps

### Step 1: Clone Repository
```bash
git clone https://github.com/streamverse/streaming-saas.git
cd streaming-saas
```

### Step 2: Configure Environment
```bash
# Copy environment template
cp .env.example .env

# Edit configuration
nano .env

# Required variables:
# - RUNPOD_API_KEY (from Runpod.io)
# - YOUTUBE_API_KEY, TWITCH_CLIENT_ID, etc. (platform integrations)
# - CLOUDFLARE_API_KEY (for CDN)
# - LOCAL_GPU_COUNT (number of GPUs on this server)
```

### Step 3: Option A - Docker Compose (Single Server)

**For development/testing or single-server deployments:**

```bash
# Start all services
docker-compose -f streaming-saas/docker-compose.streaming-saas.yml up -d

# Verify services
docker ps

# Check logs
docker-compose logs -f transcoding-service

# Access services:
# - API Gateway: http://localhost:8080
# - Ingestion: http://localhost:8100
# - Transcoding: http://localhost:8101
# - Grafana: http://localhost:3000
# - Prometheus: http://localhost:9090
```

### Step 4: Option B - Kubernetes (Production)

**For multi-server production deployments:**

#### 4.1 Setup Jenkins
```bash
# Install Jenkins
docker run -d -p 8090:8080 -p 50000:50000 \
  -v jenkins_home:/var/jenkins_home \
  --name jenkins \
  jenkins/jenkins:lts

# Get initial admin password
docker exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword

# Access Jenkins: http://your-server:8090
# Install required plugins:
# - Kubernetes Plugin
# - Docker Pipeline
# - Ansible Plugin
```

#### 4.2 Import Jenkinsfile
1. Create new Pipeline job: `StreamVerse-Deploy`
2. Point to: `ci-cd/jenkins/Jenkinsfile`
3. Configure parameters:
   - `environment`: staging / production
   - `version`: latest / specific tag
   - `run_tests`: true / false

#### 4.3 Setup AWX/Ansible
```bash
# Install AWX (Ansible Tower open-source)
git clone https://github.com/ansible/awx.git
cd awx
make docker-compose-build
make docker-compose

# Access AWX: http://your-server:80
# Default credentials: admin / password

# Import playbook: ci-cd/ansible/streamverse-deploy.yml
```

#### 4.4 Deploy via Jenkins
1. Go to Jenkins → `StreamVerse-Deploy`
2. Click **Build with Parameters**
3. Set:
   - Environment: `production`
   - Version: `latest`
   - Run Tests: `Yes`
4. Click **Build**

**Pipeline will automatically:**
- Build Docker images for all services
- Run security scans (Trivy, Snyk)
- Execute end-to-end tests
- Deploy to Kubernetes via Tekton
- Verify deployment health

**Average deployment time:** 8-12 minutes

---

## Infrastructure Architecture

### Multi-Server Layout

```
┌─────────────────────────────────────────────────────────────┐
│                     LOAD BALANCER                           │
│                  (HAProxy / NGINX)                          │
│                  External IP: X.X.X.X                       │
└────────────────────┬────────────────────────────────────────┘
                     │
        ┌────────────┴────────────┬─────────────┐
        │                         │             │
┌───────▼────────┐    ┌──────────▼────────┐   ┌▼──────────────┐
│  COMPUTE-1     │    │  COMPUTE-2        │   │  COMPUTE-3    │
│  ───────────── │    │  ──────────────── │   │  ────────────│
│  • 4x RTX 4090 │    │  • 4x RTX 4090    │   │  • 4x RTX 4090│
│  • Ingestion   │    │  • Ingestion      │   │  • Ingestion  │
│  • Transcoding │    │  • Transcoding    │   │  • Transcoding│
│  • AI Enhance  │    │  • AI Enhance     │   │  • AI Enhance │
│  • DRM         │    │  • DRM            │   │  • DRM        │
└────────────────┘    └───────────────────┘   └───────────────┘
        │                      │                     │
        └──────────────┬───────┴─────────────────────┘
                       │
        ┌──────────────▼──────────────────┐
        │      STORAGE CLUSTER            │
        │  ─────────────────────────      │
        │  • PostgreSQL (metadata)        │
        │  • ScyllaDB (metrics)           │
        │  • Redis (cache)                │
        │  • MinIO (object storage)       │
        │  • Ceph (archive)               │
        └─────────────────────────────────┘
                       │
        ┌──────────────▼──────────────────┐
        │      MONITORING STACK           │
        │  • Prometheus                   │
        │  • Grafana                      │
        │  • Jaeger                       │
        └─────────────────────────────────┘
```

### GPU Scaling Strategy

**Baseline Load** (< 40 jobs):
- Use local GPUs only
- Zero cloud costs
- Predictable performance

**Peak Load** (> 40 jobs):
- Local GPUs: Handle baseline (40 jobs)
- Runpod.io: Handle overflow automatically
- Auto-scaling based on queue depth

**Example:**
- **Local capacity**: 12 GPUs × 10 jobs each = 120 jobs/hour
- **Queue depth**: 150 jobs
- **Runpod.io spawns**: 3 GPU pods to handle 30 overflow jobs
- **Cost**: Only pay for 3 GPU-hours on Runpod.io

---

## Monitoring and Operations

### Grafana Dashboards
Access: http://your-server:3000 (admin / admin)

**Key Dashboards**:
1. **Ingestion Metrics**
   - Active streams
   - Bitrate statistics
   - Error rates
   - Protocol breakdown (RTMP/SRT/WebRTC)

2. **GPU Utilization**
   - Local GPU usage (%)
   - Runpod.io pod count
   - Transcoding queue depth
   - Jobs per hour

3. **Storage Metrics**
   - Disk usage (hot/warm/cold)
   - I/O operations
   - Cache hit rates

4. **Platform Integration Health**
   - YouTube upload success rate
   - Twitch stream uptime
   - API call latency by platform

### Prometheus Alerts

**Critical Alerts**:
- GPU temperature > 85°C
- Disk usage > 90%
- Transcoding queue > 200 jobs
- API error rate > 5%

**Notifications**:
- Slack: `#streaming-alerts`
- Email: `devops@streamverse.io`

### Log Aggregation
```bash
# View real-time logs
kubectl logs -f deployment/transcoding-service -n streamverse-production

# Search logs (ELK Stack)
# Access Kibana: http://your-server:5601
```

---

## Scaling Guidelines

### Horizontal Scaling (Add More Servers)

**When to scale:**
- Consistent > 80% GPU utilization
- Ingestion streams > 8K per server
- API latency > 50ms

**How to scale:**
1. Provision new bare metal server
2. Install NVIDIA drivers + Kubernetes
3. Join server to cluster:
   ```bash
   rke2 agent --server https://master-node:9345 --token <token>
   ```
4. Label node for GPU workloads:
   ```bash
   kubectl label node new-server dedicated=streaming
   ```

### Vertical Scaling (Add More GPUs)

**Adding GPUs to existing server:**
1. Power down server
2. Install additional GPU(s)
3. Power on, verify with `nvidia-smi`
4. Update environment: `LOCAL_GPU_COUNT=8`
5. Restart transcoding service

### Runpod.io Scaling

**Adjust cloud GPU limits:**
```bash
# Edit .env
RUNPOD_MAX_PODS=20        # Increase from 10
RUNPOD_GPU_TYPE=A100      # Upgrade to more powerful GPUs
```

**Cost optimization:**
- Use RTX 4090 for standard workloads ($0.39/hr)
- Use A100 for AI enhancement ($1.89/hr)
- Use H100 for 8K transcoding ($4.50/hr)

---

## Cost Analysis

### Monthly Infrastructure Costs (Example Setup)

**On-Premise (CapEx):**
- 3 servers @ $8,000 each = $24,000
- 12x RTX 4090 @ $1,600 each = $19,200
- Storage (200TB) = $15,000
- Network equipment = $5,000
- **Total CapEx**: $63,200 (amortized over 3 years = $1,756/month)

**Operational Costs (OpEx):**
- Power (12 GPUs + servers) = $800/month
- Bandwidth (10TB/month @ Cloudflare) = $400/month
- Runpod.io (average 100 GPU-hours/month) = $390/month
- Total OpEx = $1,590/month

**Total Monthly Cost**: $3,346

**Revenue Example** (10,000 hours streamed/month):
- 10,000 hours × 60 mins = 600,000 minutes
- 600,000 ÷ 1,000 × $0.40 = $240 revenue

*Note: This is just transcoding. Add delivery, storage, and platform fees to customer pricing.*

---

## Troubleshooting

### GPU Not Detected
```bash
# Check NVIDIA drivers
nvidia-smi

# If missing, reinstall:
sudo apt purge nvidia-*
sudo apt install nvidia-driver-535
sudo reboot
```

### Runpod.io Connection Failed
```bash
# Test API key
curl -H "Authorization: Bearer $RUNPOD_API_KEY" \
  https://api.runpod.io/v2/user

# Expected: {"id": "...", "email": "..."}
```

### High Transcoding Queue
```bash
# Check local GPU utilization
nvidia-smi

# Check Runpod.io pod count
curl -H "Authorization: Bearer $RUNPOD_API_KEY" \
  https://api.runpod.io/v2/pods | jq '.pods | length'

# Manually trigger scaling
# Increase RUNPOD_MAX_PODS in .env
```

### Service Unhealthy
```bash
# Check service status
kubectl get pods -n streamverse-production

# View logs
kubectl logs deployment/transcoding-service -n streamverse-production

# Restart service
kubectl rollout restart deployment/transcoding-service -n streamverse-production
```

---

## Security Best Practices

1. **Network Segmentation**
   - Isolate GPU servers on dedicated VLAN
   - Firewall rules: Only allow ingestion ports (1935, 9999)

2. **Secret Management**
   - Store API keys in Kubernetes secrets
   - Rotate Runpod.io API keys quarterly
   - Use Vault for DRM keys

3. **GPU Access Control**
   - Limit container GPU access with `NVIDIA_VISIBLE_DEVICES`
   - Monitor GPU usage for anomalies

4. **Regular Updates**
   - Weekly security patches
   - Monthly NVIDIA driver updates
   - Quarterly Kubernetes upgrades

---

## Support

**Documentation**: https://docs.streamverse.io
**Issues**: https://github.com/streamverse/streaming-saas/issues
**Community**: https://discord.gg/streamverse
**Enterprise Support**: enterprise@streamverse.io

**SLA**: 99.99% uptime (52 minutes downtime/year)
