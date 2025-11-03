# Terraform Infrastructure as Code

Multi-cloud infrastructure definitions for StreamVerse platform.

## Supported Clouds

- **AWS** - EKS, RDS, ElastiCache, MSK, S3, ECR
- **GCP** - GKE, Cloud SQL, Memorystore, Cloud Storage, Artifact Registry
- **Azure** - AKS, PostgreSQL Flexible Server, Redis Cache, Storage Account, ACR

## Structure

```
terraform/
├── aws/
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
├── gcp/
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
├── azure/
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
└── README.md
```

## Prerequisites

1. **Terraform** >= 1.5.0
2. Cloud provider CLI tools installed and configured:
   - AWS CLI with credentials
   - GCP `gcloud` CLI with project set
   - Azure CLI with `az login`

## Usage

### AWS Deployment

```bash
cd terraform/aws
terraform init
terraform plan
terraform apply
```

Required variables:
- `db_username` (database username)
- `db_password` (database password)

Optional overrides:
```bash
terraform apply -var="environment=production" -var="node_max_size=20"
```

### GCP Deployment

```bash
cd terraform/gcp
terraform init
terraform plan -var="gcp_project_id=your-project-id" \
               -var="db_username=admin" \
               -var="db_password=secret"
terraform apply
```

### Azure Deployment

```bash
cd terraform/azure
terraform init
terraform plan -var="db_username=admin" \
               -var="db_password=secret"
terraform apply
```

## Resources Created

### Common Resources
- **VPC/Network** - Virtual network with public/private subnets
- **Kubernetes Cluster** - Managed K8s (EKS/GKE/AKS)
- **Database** - PostgreSQL (RDS/Cloud SQL/Flexible Server)
- **Cache** - Redis (ElastiCache/Memorystore/Cache for Redis)
- **Message Queue** - Kafka (MSK/Pub/Sub/Service Bus)
- **Object Storage** - S3/Cloud Storage/Storage Account
- **Container Registry** - ECR/Artifact Registry/ACR

### AWS Specific
- **NAT Gateways** - For private subnet internet access
- **MSK Cluster** - Managed Kafka
- **ECR Repositories** - One per service
- **KMS Keys** - For encryption

### GCP Specific
- **Cloud NAT** - For private subnet internet access
- **Cloud Pub/Sub** - Event streaming
- **Workload Identity** - For service-to-service auth
- **Binary Authorization** - Image signing enforcement

### Azure Specific
- **Log Analytics Workspace** - Centralized logging
- **Private DNS Zones** - For private endpoints
- **Service Bus** - Event streaming
- **Geo-replication** - For ACR (production)

## Security Features

- **Network Isolation** - Private subnets for workloads
- **Encryption at Rest** - KMS keys for databases and storage
- **Encryption in Transit** - TLS for all services
- **Network Policies** - Calico/Cilium for pod-to-pod communication
- **Private Endpoints** - Databases accessible only from VPC
- **IAM Roles** - Least privilege access

## Outputs

Each cloud provider outputs:
- Cluster endpoint/connection info
- Database connection strings (stored securely)
- Storage bucket/container names
- Registry URLs

Access outputs:
```bash
terraform output
```

## State Management

State is stored remotely:
- **AWS** - S3 bucket `streamverse-terraform-state`
- **GCP** - GCS bucket `streamverse-terraform-state`
- **Azure** - Storage Account container `terraform-state`

## Cost Optimization

- **Spot Instances** - For non-production environments
- **Auto-scaling** - Nodes scale based on demand
- **Storage Optimization** - GP3/Ephemeral disks where possible
- **Regional Clusters** - High availability without extra cost

## Maintenance

### Updates
```bash
terraform plan  # Review changes
terraform apply  # Apply updates
```

### Destruction
```bash
terraform destroy  # ⚠️ Deletes all resources
```

## Troubleshooting

### Common Issues

1. **State Lock** - If Terraform is stuck:
   ```bash
   terraform force-unlock <lock-id>
   ```

2. **Provider Version** - Update provider versions in `required_providers`

3. **Credentials** - Ensure cloud credentials are configured:
   - AWS: `aws configure`
   - GCP: `gcloud auth application-default login`
   - Azure: `az login`

## Next Steps

After infrastructure is created:
1. Configure Kubernetes access (`aws eks update-kubeconfig`, etc.)
2. Deploy services using Kubernetes manifests
3. Set up CI/CD pipelines
4. Configure monitoring stack

