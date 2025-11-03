# GStreamer Worker Pool Setup

## Overview

GStreamer worker pool for transcoding VOD jobs. Workers consume jobs from Kafka and process them using GStreamer with optional GPU acceleration.

## Features

- **Kafka Consumer**: Consumes transcoding jobs from `transcoding` topic
- **GPU Acceleration**: NVIDIA CUDA support for H.264/H.265 encoding (nvh264enc, nvh265enc)
- **CPU Fallback**: x264enc, x265enc if GPU not available
- **Job Checkpointing**: Redis-based checkpointing for resume on failure
- **Progress Updates**: Real-time progress updates to Transcoding Service
- **Auto-scaling**: Kubernetes HPA based on queue depth and CPU utilization

## Pipeline Examples

### H.264 (CPU)
```bash
gst-launch-1.0 filesrc location=input.mp4 ! \
  qtdemux ! avdec_h264 ! \
  videoscale ! video/x-raw,width=1280,height=720 ! \
  x264enc speed-preset=fast bitrate=5000000 ! \
  h264parse ! mpegtsmux ! \
  filesink location=output.ts
```

### H.265 (GPU - NVIDIA)
```bash
gst-launch-1.0 filesrc location=input.mp4 ! \
  qtdemux ! avdec_h264 ! \
  videoscale ! video/x-raw,width=1920,height=1080 ! \
  nvh265enc preset=fast bitrate=8000000 ! \
  h265parse ! mpegtsmux ! \
  filesink location=output.ts
```

## Deployment

### Docker

```bash
cd infrastructure/gstreamer
docker build -t gstreamer-worker .
docker run -e KAFKA_BROKERS=localhost:9092 \
  -e KAFKA_TOPIC=transcoding \
  -e REDIS_HOST=localhost \
  -e GPU_ENABLED=true \
  gstreamer-worker
```

### Kubernetes

```bash
kubectl apply -f infrastructure/gstreamer/kubernetes/gstreamer-deployment.yaml
```

## Configuration

### Environment Variables

- `KAFKA_BROKERS`: Kafka broker addresses (comma-separated)
- `KAFKA_TOPIC`: Kafka topic name (default: `transcoding`)
- `REDIS_HOST`: Redis host for checkpointing
- `REDIS_PORT`: Redis port (default: 6379)
- `TRANSCODING_SERVICE_URL`: Transcoding Service API URL
- `GPU_ENABLED`: Enable GPU acceleration (true/false)

## Auto-scaling

HPA configuration:
- Min replicas: 3
- Max replicas: 20
- Scale on CPU > 70%
- Scale on Kafka consumer lag > 100

## Job Checkpointing

Checkpoints stored in Redis:
- Key: `transcode:checkpoint:{job_id}`
- Contains: status, progress, timestamps

Resume on failure:
1. Worker restarts
2. Reads checkpoint from Redis
3. Resumes from last checkpoint

## Metrics

Worker exposes metrics:
- Encoding speed (frames/sec vs. real-time)
- CPU/GPU utilization
- Job completion time
- Error rate

## Acceptance Criteria

- ✅ Worker consumes jobs from Kafka
- ✅ Transcoding produces correct output
- ✅ GPU acceleration verified (if available)
- ✅ Job metrics logged
- ✅ Auto-scaling on queue depth works
- ✅ Checkpointing tested with failure scenario

