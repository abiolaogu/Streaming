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
- ✅ **Client Apps**: Complete TV apps (8 platforms) and mobile apps (iOS + Android)

## 📱 Client Applications

### TV Apps (8 Platforms)
- ✅ Android TV / Google TV
- ✅ Samsung Tizen TV
- ✅ LG webOS TV
- ✅ Roku OS
- ✅ Amazon Fire TV
- ✅ Apple tvOS
- ✅ VIDAA TV
- ✅ KaiOS

### Mobile Apps
- ✅ iOS Mobile App (Swift + SwiftUI)
- ✅ Android Mobile App (Kotlin + Jetpack Compose)

See `apps/clients/` for all client applications.

## 📋 Prerequisites

- Terraform >= 1.5.0
- kubectl >= 1.28
- AWS/GCP/Azure accounts with admin access
- RunPod API key for GPU burst
- GitHub Actions configured with secrets

## 🚀 Quick Start

### 1. Clone and Setup

```bash
git clone https://github.com/abiolaogu/https-github.com-abiolaogu-Video-Streaming_AI-Studio.git
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
├── apps/
│   ├── media/           # OME, Shaka, DRM, FAST, SSAI
│   ├── cdn/             # ATS, Varnish, ATC, Rust purge
│   ├── data/            # DragonflyDB, ScyllaDB, Kafka, ClickHouse
│   ├── telecom/         # Kamailio, FreeSWITCH, Open5GS
│   ├── clients/         # All client apps
│   │   ├── web/         # Next.js web app
│   │   ├── ios-app/     # iOS mobile app
│   │   ├── android-app/ # Android mobile app
│   │   └── tv-apps/     # All TV platform apps
│   └── control/         # Autoscaler, config pusher
├── satellite/           # DVB-NIP/I/MABR configs
├── tests/               # Smoke tests, synthetic load
└── docs/                # Documentation
```

## 🎛️ Configuration

See individual app README files in `apps/clients/` for platform-specific configuration.

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

## 📚 Documentation

- **[README](docs/README.md)**: Quick start, configuration, troubleshooting
- **[BOM](docs/BOM.md)**: Bill of Materials
- **[RUNBOOK](docs/RUNBOOK.md)**: Operations runbook
- **[SATELLITE_OVERLAY](docs/SATELLITE_OVERLAY.md)**: DVB-NIP/I/MABR implementation
- **[SLOs](docs/SLOs.md)**: QoE and DORA metrics, alerts
- **[TV Apps Complete](apps/clients/tv-apps/TV_APPS_COMPLETE.md)**: All TV apps overview
- **[Mobile Apps Complete](apps/clients/MOBILE_APPS_COMPLETE.md)**: Mobile apps overview
- **[Build Guide](apps/clients/BUILD_APK_IPA_GUIDE.md)**: APK/IPA build instructions

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

- **Issues**: https://github.com/abiolaogu/https-github.com-abiolaogu-Video-Streaming_AI-Studio/issues
- **Email**: support@streamverse.com

---

**Status**: ✅ Production Ready | **Latest**: v1.0.0
