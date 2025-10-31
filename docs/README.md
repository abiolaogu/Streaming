# Streaming Platform - Production Grade Infrastructure

A production-grade streaming platform implementing VoD + FAST + PayTV with global CDN, multi-cloud infrastructure, live transcoding, telecom backend, and satellite overlay.

## Architecture Overview

This platform implements:

- **Video Services**: OME transcoding, GStreamer/FFmpeg GPU acceleration, Shaka Packager, Widevine/PlayReady/FairPlay DRM, FAST scheduler, SSAI/CSAI
- **CDN**: Apache Traffic Server + Varnish shield, ATC topology, Rust-based purge bus with Kafka fan-out, H3/QUIC, TCP BBR
- **Data Layer**: MinIO multi-site origin, DragonflyDB (ephemeral cache), ScyllaDB (durable catalog/entitlement), Kafka with MirrorMaker2, ClickHouse analytics
- **Telecom**: Kamailio/FreeSWITCH/Open5GS, RTPengine, WebRTC GW, lawful intercept
- **Satellite**: DVB-NIP/I/MABR carousel, STB cache daemon, terrestrial repair
- **Multi-Cloud**: AWS/GCP/Azure with Spot/Preemptible autoscaling

## Quick Start

### Prerequisites

- Terraform >= 1.5.0
- kubectl >= 1.28
- AWS/GCP/Azure accounts with credentials configured
- GitHub Actions secrets configured

### Bootstrap

```bash
# Clone repository
git clone https://github.com/abiolaogu/Streaming2.git
cd Streaming2

# Checkout branch
git checkout infra/video-sat-overlay

# Initialize Terraform
cd infra/terraform/aws
terraform init

# Plan deployment
terraform plan -var-file=dev.tfvars

# Apply infrastructure (requires approval)
terraform apply -var-file=dev.tfvars

# Configure kubectl
aws eks update-kubeconfig --name streaming-platform-dev --region us-east-1

# Apply Kubernetes manifests
kubectl apply -f ../../k8s/

# Deploy applications
kubectl apply -f ../../apps/
```

### Free-Tier Configuration

The dev environment is configured for free-tier usage:

- **1 on-demand CPU node** per cloud (AWS: t3.medium, GCP: e2-medium, Azure: Standard_B2s)
- **0-10 Spot CPU nodes** (auto-scale based on load)
- **0-5 Spot GPU nodes** (auto-scale based on transcoding demand)
- **0-1 On-demand GPU fallback** (only if Spot unavailable)
- **Nightly GPU scale-to-zero** when idle

### Secrets Configuration

Required secrets for GitHub Actions:

```bash
# AWS
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY

# GCP
GCP_SA_KEY (service account JSON)

# Azure (if using)
AZURE_CLIENT_ID
AZURE_CLIENT_SECRET
AZURE_TENANT_ID

# DRM
WIDEVINE_URL
PLAYREADY_URL
FAIRPLAY_CERT

# MinIO
MINIO_ACCESS_KEY
MINIO_SECRET_KEY

# ScyllaDB
SCYLLA_USERNAME
SCYLLA_PASSWORD

# Kafka
KAFKA_BOOTSTRAP_SERVERS

# Telecom
TELECOM_AUTH_TOKEN
```

Add these to GitHub repository secrets: Settings → Secrets and variables → Actions

## Scale Policies

### CPU Autoscaling

- **HPA thresholds**: CPU 70%, Memory 80%
- **Scale-up**: Aggressive (100% increase every 30s, max 5 pods at once)
- **Scale-down**: Conservative (50% decrease every 60s, 5min stabilization)
- **Preemption**: Spot nodes drain gracefully with 5-minute notice

### GPU Autoscaling

- **Triggers**: OME transcoding workload detected
- **Scale-out**: Spot GPU nodes (g4dn.xlarge, g5.xlarge, NC6s_v3)
- **Fallback**: On-demand GPU if Spot capacity unavailable
- **Nightly reset**: All GPU nodes scale to 0 at 2 AM if idle

### CDN Autoscaling

- **ATS Edge**: DaemonSet on all nodes, scales with cluster
- **Varnish Shield**: 3 dedicated replicas, auto-restart
- **Purge bus**: Rate-limited (100 req/s, 1000 burst)

## Testing

### Synthetic Load Test

```bash
# Trigger Spot CPU scale-out
./tests/synthetic-load-test.sh --type cpu --duration 300

# Trigger GPU scale-out
./tests/synthetic-load-test.sh --type gpu --duration 300

# Monitor autoscaling
watch kubectl get nodes

# Verify return to baseline after test
kubectl get nodes -l spot-enabled=true
```

### End-to-End Verification

```bash
# Test HLS playback
curl https://streaming.example.com/live/channel-1.m3u8

# Test DRM entitlement
curl -X POST https://api.example.com/drm/license \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"video_id": "test-001"}'

# Test FAST channel
curl https://streaming.example.com/fast/channel-1.m3u8

# Verify ClickHouse QoE metrics
clickhouse-client --query "SELECT * FROM streaming_platform.qoe_1min ORDER BY minute DESC LIMIT 10"
```

## Monitoring

### Key Dashboards

- **QoE SLOs**: Startup time < 2s, rebuffer rate < 1%
- **DORA Metrics**: Lead time, deploy frequency, change fail rate, MTTR
- **Cost**: On-demand node count (must be ≤ 1 per cloud)
- **Capacity**: Spot node utilization, GPU quota usage

### Alerts

- On-demand CPU count > 1
- GPU nodes idle > 4 hours
- QoE SLO breach (startup > 5s, rebuffer > 3%)
- Kafka replication lag > 1 minute
- Terraform drift detected

Access dashboards:
```bash
kubectl port-forward -n observability svc/grafana 3000:80
```

## Troubleshooting

### Pods Not Scheduling on Spot

```bash
# Check taints
kubectl get nodes -o custom-columns=NAME:.metadata.name,TAINTS:.spec.taints

# Add toleration to workload
kubectl patch deployment <name> -n <namespace> -p '{"spec":{"template":{"spec":{"tolerations":[{"key":"spot-enabled","operator":"Equal","value":"true","effect":"NoSchedule"}]}}}}'
```

### GPU Preemption

```bash
# Monitor node termination notices
kubectl get nodes -w

# Check termination handler logs
kubectl logs -n kube-system -l app=aws-node-termination-handler

# Manual drain (5-minute grace period)
kubectl drain <node> --grace-period=300 --ignore-daemonsets
```

### Kafka Replication Issues

```bash
# Check MirrorMaker2 status
kubectl logs -n data deployment/mm2-replicator

# Verify topic replication
kafka-console-consumer --bootstrap-server kafka-client.data.svc.cluster.local:9092 \
  --topic video-upload --from-beginning --max-messages 10
```

## Cost Optimization

### Baseline Costs (per month)

- **Dev**: ~$50-100 (1 on-demand CPU, minimal storage)
- **Staging**: ~$500-1000 (spot burst to 10 CPU, 5 GPU)
- **Production**: ~$2000-5000 (peak 50 CPU, 20 GPU, multi-region)

### Spot Savings

- **CPU**: ~70% discount vs on-demand
- **GPU**: ~50-60% discount vs on-demand
- **Nightly GPU off**: ~$500/month savings per GPU node

### Storage Costs

- **ScyllaDB**: $100/TB/month (3x replica)
- **ClickHouse**: $50/TB/month
- **MinIO**: $25/TB/month (standard tier)

## Security

### Compliance

- **OpenSCAP**: CIS Level 1 Server baseline scans
- **Trivy**: Container vulnerability scanning
- **Cosign**: SLSA attestations for all images
- **OPA**: PodSecurity, NetworkPolicy enforcement

### Best Practices

- All images signed with Cosign
- NetworkPolicies deny-by-default
- mTLS between services (Linkerd/Istio)
- Secrets in KeyVault/Cloud KMS
- RBAC least-privilege

## Support

For issues or questions:

- **Issues**: https://github.com/abiolaogu/Streaming2/issues
- **Discussions**: https://github.com/abiolaogu/Streaming2/discussions
- **Email**: support@streaming-platform.com

## License

See [LICENSE](../LICENSE) file for details.

