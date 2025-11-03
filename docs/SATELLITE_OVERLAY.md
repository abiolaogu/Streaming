# Satellite Overlay - DVB-NIP/I/MABR Implementation

Production-grade satellite delivery overlay for VoD + FAST content distribution.

## Overview

The satellite overlay leverages DVB-IPTV (DVB-I) for service discovery, DVB-NIP (Network Independent Protocol) for file distribution, and DVB-MABR (Multicast Adaptive Bitrate) for efficient bandwidth utilization.

**Key Features**:

- DVB-NIP carousel for VoD content
- DVB-I catalog for service discovery
- DVB-MABR for live channels
- Terrestrial repair via CDN
- STB cache daemon for local HTTP serving

## Architecture

### Headend Components

```
Internet/CDN
    ↓
[Transcoding Pipeline] → [Packaging] → [DVB Carousel Generator]
                                               ↓
                                    [Uplink Encoder]
                                               ↓
                                    [Satellite → STB]
```

### Edge Components

```
[STB DVB-S2X Tuner] → [Cache Daemon] → [HTTP Local Service]
                                           ↓
                                    [Player/App]
```

### Carousel Flow

1. **Content Ingestion**: VoD files from MinIO/CDN
2. **DSM-CC Packaging**: Wrap in MPEG-TS, create carousel
3. **NIP Table Generation**: Create content catalog
4. **Carousel Injection**: Periodic broadcast (every 30-60 min)
5. **STB Reception**: Tune to NIP table, download files
6. **Cache Management**: LRU eviction, 10GB local storage
7. **HTTP Serving**: STB app requests via localhost

## DVB-NIP Configuration

### Headend Setup

**Carousel Generator** (`satellite/headend/dvb-nip-config.xml`):

- Service ID: 5001
- Transponder: 11727.5 MHz (Ku-band)
- Symbol rate: 27500 ksps
- Modulation: 8PSK, FEC 3/4
- Cycle time: 3600 seconds

**Content Selection**:

- **Priority High**: Popular VoD (200 most watched)
- **Priority Medium**: EPG data, channel metadata
- **Priority Low**: Promotional content

**NIP Table Updates**:

- Refresh interval: 15 minutes
- Content expiry: 24 hours
- Compression: gzip for efficient bandwidth

### Edge Setup

**STB Cache Daemon** (`satellite/edge/stb-cache-daemon.c`):

- Multithreaded: 1 receiver + 1 HTTP server thread
- Cache: 10GB LRU, disk-backed
- HTTP: localhost:8080
- Logging: journald

**Installation**:

```bash
cd satellite/edge
make install
sudo systemctl enable stb-cache-daemon
sudo systemctl start stb-cache-daemon
```

## DVB-I Service Catalog

### Catalog Structure

**Service Metadata** (`satellite/headend/dvb-i-catalog.json`):

- Service ID, name, genre
- Broadcast schedule (EPG link)
- Locators (HTTP + DVB identifiers)
- Presentation (aspect, color, format)

### Discovery Flow

1. STB boots, tunes to NIP channel
2. Receives DVB-I catalog via NIP carousel
3. Parses JSON, displays channel list
4. User selects service
5. STB looks up locator (HTTP or DVB)
6. If HTTP, request via CDN
7. If DVB, check local cache first
8. If miss, tune to DVB carousel, download, serve

## DVB-MABR Implementation

### ABR Adaptation

**Bandwidth Prediction**:

- LTE-aware: Expect 10-50 Mbps (variable)
- Satellite link: Stable bandwidth per STB
- Hybrid: Fallback to terrestrial CDN if sat link degraded

**ABR Algorithm**:

- Min bitrate: 500 kbps (mobile-friendly)
- Max bitrate: 15 Mbps (4K HDR)
- Ladder: 500k, 1M, 2M, 4M, 8M, 15M
- Adaptation: CBR-based (satellite-friendly)

**Segment Carousel**:

- Each ABR variant repeated in carousel
- Segment duration: 2-4 seconds (short for satellite)
- Playlist update: Via NIP for reduced overhead

## Terrestrial Repair

### Hybrid Architecture

**Primary**: Satellite carousel (one-way broadcast)

**Secondary**: CDN (bidirectional, repair)

**Flow**:

1. STB plays from satellite cache
2. If cache miss, check NIP table for HTTP locator
3. Request missing segment from CDN
4. Concurrent play: sat segments + CDN segments
5. Prefetch: Download next 3 segments from CDN

**Config** (`satellite/edge/repair-config.json`):

```json
{
  "cdn_endpoint": "https://cdn.streaming.example.com",
  "timeout_ms": 5000,
  "prefetch_count": 3,
  "fallback_enabled": true
}
```

## Partner Integration

### LEO Gateway

**SpaceX Starlink / OneWeb Integration**:

- Gateway: OAM uplink to LEO constellation
- Backhaul: LEO → terrestrial fiber
- Delivery: Consumer terminal (flat-panel antenna)

**Placeholders**:

- `satellite/headend/leo-gateway-config.xml` (TODO: partner config)
- `satellite/headend/s2x-uplink.yaml` (TODO: K8s deployment)

### iDirect Hub

**Traditional Hub Model**:

- Hub: HTS iDirect Evolution hub
- Outroutes: DVB-S2X to small VSAT terminals
- Management: Web UI + API

**Placeholders**:

- `satellite/headend/idirect-hub-config.yaml` (TODO: hub API config)

## T+2y Rollout Plan

### Phase 1: Lab (Months 1-3)

- [x] Headend carousel generator PoC
- [x] STB cache daemon development
- [ ] DVB-S2X modulator integration
- [ ] EPG integration testing
- [ ] Multi-client playback validation

### Phase 2: Field Trial (Months 4-6)

- [ ] 100 STBs deployed in rural markets
- [ ] 10-channel FAST lineup
- [ ] VoD catalog: 500 titles
- [ ] QoS monitoring (startup, rebuffer, availability)

### Phase 3: Regional Expansion (Months 7-12)

- [ ] 1000 STBs in 3 regions
- [ ] 30-channel lineup
- [ ] VoD catalog: 2000 titles
- [ ] Subscription tiers (basic/premium)

### Phase 4: Scale (Months 13-24)

- [ ] 100K STBs deployed
- [ ] 60-channel lineup
- [ ] Full VoD library (10000+ titles)
- [ ] International expansion

## Testing

### Lab Validation

**Test Coverage**:

```bash
# Carousel integrity
make -C satellite/edge test-carousel

# Cache daemon stress test
make -C satellite/edge stress-test-cache

# Hybrid playback
make -C satellite/edge test-hybrid

# Recovery scenarios
make -C satellite/edge test-recovery
```

**SLOs**:

- Cache hit rate: > 85%
- Terrestrial repair latency: < 500ms
- Startup time (satellite): < 10s
- Startup time (terrestrial fallback): < 3s

### Field Testing

**Rural Deployment** (planned Q2 2025):

- Location: Montana, Alaska, Appalachia
- STB count: 100-500
- Duration: 90 days
- Metrics: QoS, user satisfaction, cost efficiency

## Cost Analysis

### Capital

- **Headend Equipment**: $45k (see BOM.md)
- **STB (ODM)**: $50/unit (1M units = $50M CapEx)
- **VSAT Terminals** (if used): $200/unit
- **Total CapEx**: ~$60M at scale

### Operating

- **Transponder Lease**: $180k/month (36 MHz Ku-band)
- **Hub/Gateway**: $50k/month
- **Ground Segment**: $20k/month
- **STB Software**: $2/unit/year (~$200k/month at 1M units)
- **Total Monthly OpEx**: ~$450k (excluding one-time CapEx)

### ROI Timeline

- **Break-even**: 18-24 months (assuming subscription revenue $10-15/month/STB)
- **Payback period**: 3-4 years
- **Addressable market**: Rural broadband underserved (10M+ households US)

## References

- DVB-NIP: ETSI TS 102 323
- DVB-I: ETSI TS 103 770
- DVB-MABR: ETSI TS 103 772
- DSM-CC: ISO/IEC 13818-6
- iDirect Hub API: https://idirect.net

