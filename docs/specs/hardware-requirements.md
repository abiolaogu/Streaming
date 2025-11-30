# StreamVerse Hardware Requirements

## Minimum System Requirements (Development)
-   **CPU**: 4 Cores (Intel Core i5 / AMD Ryzen 5 or equivalent)
-   **RAM**: 16 GB
-   **Storage**: 50 GB SSD available space
-   **OS**: macOS Ventura+, Linux (Ubuntu 22.04+), or Windows 11 (WSL2)

## Recommended Production Specifications (Per Node)
### Kubernetes Worker Node
-   **CPU**: 8 vCPUs
-   **RAM**: 32 GB
-   **Storage**: 100 GB NVMe SSD
-   **Network**: 10 Gbps

### Database Node (MongoDB/PostgreSQL)
-   **CPU**: 16 vCPUs
-   **RAM**: 64 GB
-   **Storage**: 1 TB NVMe SSD (RAID 10 recommended)
-   **Network**: 10 Gbps

### Cache Node (Redis/DragonflyDB)
-   **CPU**: 4 vCPUs
-   **RAM**: 32 GB (High Frequency)
-   **Storage**: 100 GB SSD
-   **Network**: 10 Gbps

## Network Requirements
-   **Public IP**: Static IP for Load Balancer/Ingress.
-   **Ports**:
    -   80/443 (HTTP/HTTPS)
    -   6443 (Kubernetes API)
    -   27017 (MongoDB - Internal)
    -   6379 (Redis - Internal)
