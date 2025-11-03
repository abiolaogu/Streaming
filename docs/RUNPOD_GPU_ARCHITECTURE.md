# RunPod GPU Architecture

## Overview

The platform uses a **hybrid GPU architecture** combining RunPod API-driven cloud GPU burst with local GPU baseline at Tier-1 PoPs for maximum cost efficiency and ultra-low-latency performance.

## Architecture Components

```
┌─────────────────────────────────────────────────────────────┐
│                    Control Plane                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │     KEDA Triggers (Kafka queue, latency SLO)        │   │
│  └─────────────────────────────────────────────────────┘   │
│                        ↓                                     │
│  ┌─────────────────────────────────────────────────────┐   │
│  │      RunPod Autoscaler (Go microservice)            │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                           ↓
        ┌──────────────────┴──────────────────┐
        ↓                                     ↓
┌────────────────────┐              ┌─────────────────────┐
│   RunPod Cloud     │              │  Local GPU Baseline │
│   GPU Burst        │              │  (Tier-1 PoPs)      │
│                    │              │                     │
│  ┌──────────────┐  │              │  ┌──────────────┐   │
│  │ RTX 6000 Ada │  │              │  │ A100 2x      │   │
│  │ A100         │  │              │  │ OpenStack    │   │
│  │ 0-50 nodes   │  │              │  │ 1-2 per PoP  │   │
│  └──────────────┘  │              │  └──────────────┘   │
│                    │              │                     │
│  • SSH to K8s      │              │  • Always-on        │
│  • Auto-join       │              │  • Ultra-low-lat    │
│  • Checkpoint MinIO│              │  • Checkpoint MinIO │
└────────────────────┘              └─────────────────────┘
```

## Trigger Mechanisms

### 1. Kafka Queue Depth

```yaml
# KEDA ScaledObject
triggers:
- type: kafka
  metadata:
    bootstrapServers: kafka-client.data.svc.cluster.local:9092
    consumerGroup: ome-transcoder-consumer
    topic: video-transcode
    lagThreshold: '5'  # Trigger at 5 jobs queued
```

### 2. Latency SLO Breach

```yaml
- type: prometheus
  metadata:
    serverAddress: http://prometheus.observability.svc.cluster.local:9090
    metricName: transcoding_latency_p95
    threshold: '30000'  # 30 seconds
```

### 3. Custom RunPod Trigger

```yaml
- type: external
  metadata:
    scalerAddress: runpod-autoscaler.platform.svc.cluster.local:8080
    queueName: gpu-queue
    queueLength: '5'
```

## Autoscaling Flow

### Scale-Up

1. **Detection**: KEDA detects queue depth > 5 OR latency > 30s
2. **Request**: KEDA triggers RunPod autoscaler via HTTP
3. **Spin-up**: RunPod API creates GPU instance (RTX 6000 Ada, A100, etc.)
4. **SSH Setup**: Autoscaler configures SSH tunnel to instance
5. **K8s Join**: Instance auto-joins Kubernetes as worker node
6. **Labeling**: Instance labeled `nvidia.com/gpu=true`, `gpu-type=runpod`
7. **Pod Scheduling**: OME transcoder pods schedule on RunPod instance
8. **Checkpoint**: Work begins, checkpointed to MinIO
9. **Egress**: Processed chunks pushed via RunPod network

### Scale-Down

1. **Detection**: Queue depth = 0 AND idle > 60s
2. **Drain**: Pods gracefully drain from RunPod instances
3. **Checkpoint**: Final state saved to MinIO
4. **Terminate**: RunPod API deletes instances
5. **Mesh Cleanup**: SSH tunnels torn down
6. **Baseline**: Return to local GPU baseline only

## Local GPU Baseline

### Purpose

- **Ultra-low-latency**: < 5ms regional availability
- **Media AI**: Real-time upscaling, noise reduction
- **Backup encoding**: Graceful degradation if RunPod unavailable
- **Persistent workloads**: Always-on transcoding pipelines

### Deployment

```bash
# Deploy via Terraform (OpenStack)
cd infra/terraform/openstack
terraform apply -var-file=local-gpu-baseline.tfvars
```

### Configuration

```hcl
resource "openstack_compute_instance_v2" "local_gpu_baseline" {
  count     = 1  # 1-2 per Tier-1 PoP
  flavor_name = "gpu.a100.2x"  # A100 2x GPUs
  
  # Labels
  workload-type = "gpu"
  gpu-type      = "local-baseline"
  preemptible   = "false"
  
  # Taints
  taints = ["gpu-type=local-baseline:NoSchedule"]
}
```

## Checkpoint & Resume

### Checkpoint Location

- **Backend**: MinIO S3-compatible storage
- **Path**: `s3://checkpoint/<job-id>/`
- **Format**: HLS segments, DASH manifests, DRM keys

### Resume Logic

```go
// On RunPod termination
func checkpointJob(jobID string) error {
    // Save current progress
    checkpoint := Checkpoint{
        SegmentsComplete: currentSegment,
        BitrateLevel:     currentBitrate,
        Timestamp:        time.Now(),
    }
    
    // Upload to MinIO
    minioClient.PutObject("checkpoint", jobID, checkpoint)
}

// On resume (any GPU node)
func resumeJob(jobID string) error {
    // Download from MinIO
    checkpoint := downloadCheckpoint(jobID)
    
    // Resume from checkpoint
    return transcoder.ResumeFrom(checkpoint)
}
```

## Egress Optimization

### Problem

RunPod egress costs can be high (~$0.05-0.10/GB). To minimize:

### Solution

1. **Push processed chunks** via RunPod network (free internal)
2. **CDN pull** from RunPod edge locations
3. **Batch uploading** to origin (MinIO) during off-peak

```bash
# Egress-optimized workflow
OME transcoder → Process segment → Local cache → CDN edge pull → Origin
                                   ↓
                           RunPod internal (free)
```

## Cost Analysis

### RunPod Pricing

- **RTX 6000 Ada**: ~$0.69/hour (spot: ~$0.35/hour)
- **A100 80GB**: ~$1.89/hour (spot: ~$0.95/hour)
- **A6000**: ~$0.89/hour (spot: ~$0.44/hour)

### Local GPU Baseline

- **A100 2x OpenStack**: ~$500/month (amortized CapEx)
- **Power**: ~2kW @ $0.15/kWh = $216/month
- **Total**: ~$716/month per baseline server

### Break-Even

- **RunPod 100 hrs/month** vs local baseline: **~$140/month savings**
- **RunPod 500 hrs/month** vs local baseline: **~$340/month savings**

### Recommendation

- **< 200 hrs/month**: Use RunPod only
- **200-500 hrs/month**: Hybrid (RunPod burst + 1 local baseline)
- **> 500 hrs/month**: 2+ local baseline + RunPod for peak

## Monitoring

### Key Metrics

```promql
# RunPod instance count
sum(kube_node_labels{label_nvidia_com_gpu="true",label_gpu_type="runpod"})

# Queue depth
keda_scaler_metrics_value{scalerName="runpod-gpu-scaler"}

# Latency SLO
histogram_quantile(0.95, sum(rate(ome_transcoding_duration_seconds_bucket[5m])) by (le))

# Cost tracking
sum(runpod_instance_cost_per_hour * runpod_instance_uptime_hours)
```

### Alerts

- **Queue depth > 20**: Scale-up trigger failing
- **Latency > 60s**: SLO breach, add more RunPod instances
- **GPU idle > 4 hours**: Scale-down not working
- **Checkpoint failure**: MinIO issues, manual intervention

## Troubleshooting

### RunPod Not Joining K8s

```bash
# Check SSH connectivity
kubectl exec -it runpod-autoscaler -n platform -- ssh root@<runpod-ip>

# Check bootstrap token
kubectl get secret bootstrap-token -n kube-system

# Manual join
kubeadm join <master> --token <token> --discovery-token-unsafe-skip-ca-verification
```

### High Egress Costs

```bash
# Enable CDN pull mode
kubectl set env deployment/ome-transcoder -n media \
  EGRESS_MODE=cdn-pull \
  CDN_ENDPOINT=https://cdn.example.com

# Monitor egress
curl http://prometheus:9090/api/v1/query?query=runpod_egress_bytes_total
```

### Checkpoint Failures

```bash
# Test MinIO connectivity
mc ping s3://checkpoint

# Check disk space
kubectl exec -it minio-0 -n data -- df -h

# Retry failed checkpoints
kubectl exec -it ome-transcoder-0 -n media -- retry-checkpoints
```

## Future Enhancements

- [ ] Multi-region RunPod clusters
- [ ] Advanced checkpoint compression (zstd)
- [ ] Preemptive resume (estimate GPU availability)
- [ ] GPU affinity scheduling (same-model preference)
- [ ] RunPod spot price optimization
- [ ] Federated learning integration
- [ ] Edge GPU (NVIDIA Jetson) support

