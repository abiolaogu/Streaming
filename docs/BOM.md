# Bill of Materials (BOM)

Hardware and rental equivalents for Tier-1, Tier-2, Tier-3, Satellite headend, STB/home edge, and network infrastructure.

## Tier-1 Origin (Primary Location)

| Component | Spec | Count | Power (W) | Rack (U) | Monthly Cost |
|-----------|------|-------|-----------|----------|--------------|
| **Compute Nodes** |
| Dell R7525 (dual EPYC) | 128 cores, 512GB RAM, 4x NVMe | 10 | 750 | 42U | $8,000 |
| GPU Nodes | 8x A100, 256GB RAM | 5 | 3000 | 42U | $25,000 |
| **Storage** |
| MinIO Cluster | 50TB NVMe, 3x replica | 3 | 400 | 6U | $2,000 |
| ScyllaDB | 100TB SSD, 3-node | 3 | 500 | 9U | $5,000 |
| **Network** |
| Arista 7280CR3 | 32x100G + 2x400G | 2 | 900 | 8U | $3,500 |
| **Support** | Rack, UPS, cooling, PDU | - | 10kW | 42U | $1,500 |
| **Total** | | **26 servers** | **~35kW** | **~100U** | **$45,000/mo** |

### Rental Equivalent (OVH/Ubiquity)

- **Bare Metal**: 10x Rise-3 (128vCPU, 512GB), 5x Rise-LE (8x A100), 3x Rise Storage
- **Cost**: ~$18,000/month
- **Savings vs owned**: ~60% (ignoring CapEx)

## Tier-2 Edge (Per Region)

| Component | Spec | Count | Power (W) | Rack (U) | Monthly Cost |
|-----------|------|-------|-----------|----------|--------------|
| **Compute** |
| Dell R6525 (single EPYC) | 64 cores, 256GB RAM | 5 | 450 | 21U | $2,000 |
| GPU Nodes | 4x A100, 128GB RAM | 2 | 2000 | 17U | $8,000 |
| **CDN Cache** |
| ATS Cluster | 20TB NVMe | 10 | 300 | 21U | $3,000 |
| Varnish Shield | 10TB NVMe | 3 | 250 | 6U | $1,000 |
| **Network** |
| Arista 7280SR3 | 32x100G | 2 | 700 | 6U | $2,000 |
| **Support** | | - | 8kW | 35U | $1,200 |
| **Total** | | **22 servers** | **~20kW** | **~70U** | **$17,200/mo** |

### Multi-Region Scale

- **Regions**: 10 (US-East, US-West, EU-West, EU-Central, Asia-Pacific, etc.)
- **Total Edge Cost**: $172,000/month
- **Rental Equivalent (Zenlayer/Voxility)**: ~$90,000/month

## Tier-3 PoP (Per City)

| Component | Spec | Count | Power (W) | Rack (U) | Monthly Cost |
|-----------|------|-------|-----------|----------|--------------|
| **Compute** |
| 1U Dell Server | 32 cores, 128GB, 4TB NVMe | 3 | 300 | 3U | $400 |
| **Cache** |
| ATS Edge | 5TB NVMe | 1 | 150 | 1U | $200 |
| **Network** |
| 10G Switch | 24x1G + 4x10G | 1 | 50 | 1U | $150 |
| **Support** | Colo, power, transit | - | 500W | 5U | $300 |
| **Total** | | **5 servers** | **~1kW** | **~10U** | **$1,050/mo** |

### PoP Coverage

- **Cities**: 50
- **Total PoP Cost**: $52,500/month
- **Rental Equivalent (Vultr/BareMetal)**: ~$30,000/month

## Satellite Headend

| Component | Spec | Count | Power (W) | Rack (U) | Monthly Cost |
|-----------|------|-------|-----------|----------|--------------|
| **Encoding** |
| ATEME TITAN Live | 2xHEVC, 4xAVC | 2 | 1500 | 8U | $15,000 |
| Harmonic Electra XVM | 4xAV1 | 2 | 1200 | 6U | $12,000 |
| **Mux & Carousel** |
| ATEME CAROUSEL | DVB-NIP/I, DSM-CC | 1 | 800 | 4U | $5,000 |
| Rohde & Schwarz M8000 | Modulator | 1 | 1000 | 2U | $8,000 |
| **Uplink** |
| iDirect Quantum | Satellite Modem | 1 | 400 | 2U | $3,000 |
| **Support** | Rack, power, VSAT | - | 5kW | 22U | $2,000 |
| **Total** | | **7 appliances** | **~10kW** | **~25U** | **$45,000/mo** |

### Satellite OPEX

- **Transponder Lease**: 36 MHz Ku-band: $180,000/month
- **Hub/Gateway**: iDirect-style hub: $50,000/month
- **Ground Segment**: VSAT + terrestrial repair: $20,000/month
- **Total Satellite**: **$295,000/month**

## STB/Home Edge

| Component | Spec | Count | Unit Cost | Monthly Cost |
|-----------|------|-------|-----------|--------------|
| **Hardware** |
| STB (ODM) | ARM A78, 4GB RAM, 32GB eMMC, DVB-S2X tuner | 1M | $50 | $5M CapEx |
| **Software** |
| Linux + cache daemon | Per-device license | 1M | $2 | $2M/year |
| **Total** | | **1M devices** | **$52** | **~$500k/mo** |

### TCO (3-year)

- **CapEx**: $52M (amortized: $1.4M/month)
- **OpEx**: $500k/month (SW + support)
- **Per-device**: $174 total, $58/year

## Network Gear Summary

| Tier | Switch Model | Ports | Power | Monthly |
|------|--------------|-------|-------|---------|
| Tier-1 | Arista 7280CR3 | 32x100G + 2x400G | 900W | $3,500 |
| Tier-2 | Arista 7280SR3 | 32x100G | 700W | $2,000 |
| Tier-3 | 10G Edge Switch | 24x1G + 4x10G | 50W | $150 |
| PoP | 1G Access Switch | 48x1G | 30W | $100 |

### Rental vs Owned

| Deployment | Owned (monthly) | Rental (monthly) | Savings |
|------------|----------------|------------------|---------|
| Tier-1 Origin | $45,000 | $18,000 | 60% |
| Tier-2 Edge (10 regions) | $172,000 | $90,000 | 48% |
| Tier-3 PoP (50 cities) | $52,500 | $30,000 | 43% |
| Satellite Headend | $45,000 | N/A | N/A |
| **Subtotal (Infra)** | **$314,500** | **$138,000** | **56%** |
| **Satellite OPEX** | **$295,000** | **$295,000** | 0% |
| **Total** | **$609,500** | **$433,000** | **29%** |

Note: Rental pricing from OVH, Zenlayer, Voxility, Vultr (as of 2024). Owned includes 3-year depreciation + maintenance.

## Power Estimates (Total)

| Location | Servers | Network | Cooling | UPS | Total |
|----------|---------|---------|---------|-----|-------|
| Tier-1 | 12kW | 2kW | 6kW | 2kW | **22kW** |
| Tier-2 (per region) | 8kW | 1kW | 3kW | 1kW | **13kW** |
| Tier-3 (per PoP) | 1kW | 0.1kW | 0.5kW | 0.1kW | **1.7kW** |
| Satellite Headend | 5kW | 0.5kW | 2kW | 1kW | **8.5kW** |

**Global Total**: ~500kW (all sites)

## Rack Summary

| Location | Racks | Servers | Switches | PDU/UPS | Total U |
|----------|-------|---------|----------|---------|---------|
| Tier-1 Origin | 2x 42U | 26 | 2 | 2 | **84U** |
| Tier-2 Edge (per region) | 2x 42U | 22 | 2 | 2 | **84U** |
| Tier-3 PoP (per site) | 1x 42U | 5 | 1 | 1 | **42U** |
| Satellite Headend | 1x 42U | 7 | 0 | 1 | **42U** |

**Total Racks**: ~65 (Tier-1: 2, Tier-2: 20, Tier-3: 50, Satellite: 1)

