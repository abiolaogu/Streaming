# OvenMediaEngine (OME) Live Ingest Setup

## Overview

OvenMediaEngine configuration for live streaming ingest (RTMP, SRT, WebRTC) and real-time transcoding to LL-HLS output.

## Features

- **RTMP Ingest** (Port 1935): Primary ingest protocol from encoders (OBS, Wirecast)
- **SRT Ingest** (Port 9999): Backup ingest protocol, resilient over lossy networks
- **WebRTC Ingest** (Port 3333): Browser-based ingest for user-generated content
- **LL-HLS Output** (Port 8080): Low-latency HLS with 2s segments, 6 parts
- **DASH Output** (Port 8081): DASH streaming (optional)
- **WebRTC Output** (Port 3478): Near-real-time subscribers
- **Failover Support**: Backup OME instance for seamless failover

## Transcoding Profiles

- **240p**: 512 kbps, H.264
- **360p**: 1.5 Mbps, H.264
- **480p**: 2.5 Mbps, H.264
- **720p**: 5 Mbps, H.264, 30fps
- **1080p**: 8 Mbps, H.265, 30fps
- **4K**: 15 Mbps, H.265, 60fps

## Deployment

### Docker Compose

```bash
cd infrastructure/ome
docker-compose up -d
```

### Kubernetes

```bash
kubectl apply -f infrastructure/ome/kubernetes/ome-deployment.yaml
```

## Ingest URLs

### RTMP
```
rtmp://ome-host:1935/live/stream-key
```

### SRT
```
srt://ome-host:9999?streamid=live/stream-key
```

### WebRTC
```
ws://ome-host:3333/live/stream-key
```

## Output URLs

### LL-HLS
```
http://ome-host:8080/live/stream-key_llhls/index.m3u8
```

### DASH
```
http://ome-host:8081/live/stream-key/index.mpd
```

## Monitoring

OME exposes metrics at:
```
http://ome-host:8080/v1/stats/current
```

## Failover

Primary ingest fails → Secondary OME instance takes over (<3s interruption)

## Acceptance Criteria

- ✅ RTMP ingest works
- ✅ Real-time ABR generation (6 profiles)
- ✅ LL-HLS segments generated (2s, 6 parts)
- ✅ <2s glass-to-glass latency
- ✅ Failover tested and working

