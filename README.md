# Streaming2 - Production-Grade Multi-Cloud Streaming Platform

🚀 **A complete production-ready infrastructure for VoD + FAST + PayTV delivery with global CDN, RunPod GPU plane, multi-cloud Spot/Preemptible autoscaling, live GPU transcoding, telecom backend, and satellite overlay.**

## 🎯 Overview

This repository implements a comprehensive streaming platform designed for maximum cost efficiency through hybrid GPU architecture: RunPod API-driven cloud GPU burst + local GPU baseline at Tier-1 PoPs, with multi-cloud Spot/Preemptible CPU autoscaling.

### Key Features

- ✅ **Video Services**: OME transcoding, GStreamer/FFmpeg GPU, Shaka Packager, DRM (Widevine/PlayReady/FairPlay), FAST scheduler, SSAI/CSAI
- ✅ **CDN**: Apache Traffic Server, Varnish shield, ATC topology, Rust purge bus, H3/QUIC, TCP BBR
- ✅ **Data Layer**: MinIO origin, DragonflyDB, ScyllaDB, Kafka/MM2, ClickHouse analytics
- ✅ **Telecom**: Kamailio, FreeSWITCH, Open5GS (5G Core), WebRTC GW, RTPengine
- ✅ **Satellite Overlay**: DVB-NIP/I/MABR carousel, STB cache daemon, terrestrial repair
- ✅ **GPU Architecture**: RunPod cloud burst + local GPU baseline (Tier-1 PoPs) via KEDA triggers
- ✅ **Multi-Cloud**: AWS (EKS), GCP (GKE), Azure (AKS) with Spot CPU, OpenStack for on-prem
- ✅ **Cost Optimized**: ≤1 on-demand CPU per cloud, Spot/Preemptible burst, GPU=0 when idle

## 📋 Prerequisites

- Terraform >= 1.5.0
- kubectl >= 1.28
- AWS/GCP/Azure accounts with admin access
- RunPod API key for GPU burst
- GitHub Actions configured with secrets

## 🚀 Quick Start

### 1. Clone and Setup

```bash
git clone https://github.com/abiolaogu/Streaming2.git
cd Streaming2
git checkout infra/video-sat-overlay
```

### 2. Deploy Infrastructure

```bash
# AWS
cd infra/terraform/aws
terraform init
terraform plan -var-file=dev.tfvars
terraform apply -var-file=dev.tfvars

# GCP
cd ../gcp
terraform init
terraform plan -var-file=dev.tfvars
terraform apply -var-file=dev.tfvars

# Azure
cd ../azure
terraform init
terraform plan -var-file=dev.tfvars
terraform apply -var-file=dev.tfvars
```

### 3. Deploy Applications

```bash
# Configure kubectl
aws eks update-kubeconfig --name streaming-platform-dev --region us-east-1

# Apply manifests
kubectl apply -f k8s/
kubectl apply -f apps/media/
kubectl apply -f apps/cdn/
kubectl apply -f apps/data/
kubectl apply -f apps/telecom/
kubectl apply -f apps/clients/
```

### 4. Run Smoke Tests

```bash
./tests/smoke-tests.sh
```

## 📁 Repository Structure

```
Streaming2/
├── infra/
│   └── terraform/
│       ├── aws/          # AWS EKS infrastructure
│       ├── gcp/          # GCP GKE infrastructure
│       ├── azure/        # Azure AKS infrastructure
│       ├── openstack/    # OpenStack infrastructure
│       └── global/       # Cross-cloud DNS/CDN
├── k8s/                  # Kubernetes manifests
│   ├── namespaces.yaml
│   ├── resource-quotas.yaml
│   ├── priority-classes.yaml
│   ├── network-policies.yaml
│   ├── hpa-config.yaml
│   └── pdbs.yaml
├── apps/
│   ├── media/           # OME, Shaka, DRM, FAST, SSAI
│   ├── cdn/             # ATS, Varnish, ATC, Rust purge
│   ├── data/            # DragonflyDB, ScyllaDB, Kafka, ClickHouse
│   ├── telecom/         # Kamailio, FreeSWITCH, Open5GS
│   ├── clients/         # Web (Next.js), Mobile (Flutter), STB
│   └── control/         # Autoscaler, config pusher
├── satellite/
│   ├── headend/         # DVB-NIP/I/MABR configs
│   └── edge/            # STB cache daemon
├── tests/               # Smoke tests, synthetic load
├── docs/                # Documentation
│   ├── README.md        # Main guide
│   ├── BOM.md           # Bill of Materials
│   ├── RUNBOOK.md       # Operations runbook
│   ├── SATELLITE_OVERLAY.md
│   └── SLOs.md          # SLOs and dashboards
└── .github/workflows/   # CI/CD pipelines
    ├── ci.yml
    ├── deploy.yml
    └── drift-detect.yml
```

## 🎛️ Configuration

### Environment Variables

Required secrets (add to GitHub Actions):

```bash
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
GCP_SA_KEY
AZURE_CLIENT_ID
AZURE_CLIENT_SECRET
AZURE_TENANT_ID
RUNPOD_API_KEY
RUNPOD_NETWORK_VOLUME_ID
RUNPOD_TEMPLATE_ID
WIDEVINE_URL
PLAYREADY_URL
FAIRPLAY_CERT
MINIO_ACCESS_KEY
MINIO_SECRET_KEY
SCYLLA_USERNAME
SCYLLA_PASSWORD
KAFKA_BOOTSTRAP_SERVERS
```

### Node Pool Configuration

**Dev environment**:

- **On-Demand CPU**: 1 node per cloud (t3.medium, e2-medium, B2s) - control plane only
- **Spot CPU**: 0-10 nodes (auto-scale)
- **RunPod GPU**: 0-20 instances (API-driven burst)
- **Local GPU Baseline**: 1 GPU server per Tier-1 PoP (OpenStack)

**Production environment**:

- **On-Demand CPU**: 1 node per cloud
- **Spot CPU**: 0-50 nodes
- **RunPod GPU**: 0-50 instances
- **Local GPU Baseline**: 1-2 GPU servers per Tier-1 PoP

## 🔄 Autoscaling

### CPU Autoscaling

- **Trigger**: CPU > 70% or Memory > 80% for 60s
- **Scale-up**: +100% replicas per 30s, max +5 pods
- **Scale-down**: -50% replicas per 60s, 5min stabilization
- **Preemption**: 5-minute grace period

### GPU Autoscaling

- **Trigger**: Kafka queue depth > 5 OR latency SLO breach OR KEDA metrics
- **RunPod Burst**: API-driven spin-up (RTX 6000 Ada, A100, etc.)
- **SSH Attachment**: Auto-join to Kubernetes cluster, appear as worker nodes
- **Local Baseline**: Tier-1 PoPs for ultra-low-latency media AI
- **Checkpoint**: MinIO-backed, jobs resume anywhere (RunPod or local)
- **Egress Optimization**: Push processed chunks via RunPod network
- **Scale-down**: GPU = 0 when idle (60s threshold)

### CDN Autoscaling

- **ATS Edge**: DaemonSet, scales with cluster
- **Varnish Shield**: 3 replicas, auto-restart
- **Purge bus**: Rate-limited (100 req/s, 1000 burst)

## 📊 Monitoring & Observability

### Key Metrics

- **QoE**: Startup time, rebuffer rate, bitrate adaptation
- **DORA**: Lead time, deploy frequency, change fail rate, MTTR
- **Cost**: On-demand node count (must be ≤1 per cloud)
- **Capacity**: Spot utilization, GPU quota, storage

### Dashboards

Access Grafana:

```bash
kubectl port-forward -n observability svc/grafana 3000:80
```

Login: `admin` / `admin` (change on first login)

## 🛡️ Security

- **OpenSCAP**: CIS Level 1 Server baseline scans
- **Trivy**: Container vulnerability scanning
- **Cosign**: SLSA attestations for all images
- **OPA**: PodSecurity, NetworkPolicy enforcement
- **mTLS**: Linkerd/Istio between services

## 📈 Cost Optimization

### Dev Environment

- **Monthly Cost**: ~$50-100
- **On-demand**: 1 CPU node per cloud
- **Spot**: Minimal usage during non-business hours

### Production Environment

- **Monthly Cost**: ~$2k-5k
- **On-demand**: 1 CPU node per cloud
- **Spot**: 50 CPU + 20 GPU peak
- **Savings**: ~70% vs all on-demand

### Nightly GPU Off

- **Savings**: ~$500/month per GPU node
- **Auto-scale**: GPU → 0 at 2 AM if idle

## 🧪 Testing

### Smoke Tests

```bash
./tests/smoke-tests.sh
```

Validates:
- Cluster connectivity
- Namespace existence
- Pod health
- Service accessibility
- HPA configuration
- Node pool compliance
- Cost guardrails

### Synthetic Load Test

```bash
./tests/synthetic-load-test.sh --duration 300
```

Triggers:
- Spot CPU scale-out
- GPU auto-scaling
- Verifies return to baseline

## 📚 Documentation

- **[README](docs/README.md)**: Quick start, configuration, troubleshooting
- **[BOM](docs/BOM.md)**: Hardware and rental equivalents
- **[RUNBOOK](docs/RUNBOOK.md)**: Preemption handling, failover, restore, rollback
- **[SATELLITE_OVERLAY](docs/SATELLITE_OVERLAY.md)**: DVB-NIP/I/MABR implementation
- **[SLOs](docs/SLOs.md)**: QoE and DORA metrics, alerts

## 🐛 Troubleshooting

### Pods Not Scheduling on Spot

```bash
# Check taints
kubectl get nodes -o custom-columns=NAME:.metadata.name,TAINTS:.spec.taints

# Add toleration
kubectl patch deployment <name> -n <namespace> -p '{"spec":{"template":{"spec":{"tolerations":[{"key":"spot-enabled","operator":"Equal","value":"true","effect":"NoSchedule"}]}}}}'
```

### GPU Preemption

```bash
# Monitor termination notices
kubectl get nodes -w

# Check logs
kubectl logs -n kube-system -l app=aws-node-termination-handler
```

### Kafka Replication Issues

```bash
# Check MirrorMaker2
kubectl logs -n data deployment/mm2-replicator

# Verify topics
kafka-console-consumer --bootstrap-server kafka-client.data.svc.cluster.local:9092 --topic video-upload --from-beginning
```

See [RUNBOOK](docs/RUNBOOK.md) for detailed troubleshooting.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Apache OvenMedia Engine
- Apache Traffic Server
- Varnish Cache
- Apache Kafka
- ScyllaDB
- DragonflyDB
- ClickHouse
- Kamailio
- FreeSWITCH
- Open5GS

## 📞 Support

- **Issues**: https://github.com/abiolaogu/Streaming2/issues
- **Discussions**: https://github.com/abiolaogu/Streaming2/discussions
- **Email**: support@streaming-platform.com

---

**Status**: ✅ Production Ready | **Latest**: v1.0.0

