# Disaster Recovery & Backup Strategy

Comprehensive DR and backup procedures for StreamVerse platform.

## Backup Strategy

### Database Backups

#### MongoDB
- **Frequency**: Hourly incremental, daily full
- **Retention**: 30 days
- **Location**: S3/GCS (multi-region)
- **Method**: `mongodump` with compression

#### PostgreSQL
- **Frequency**: Hourly WAL, daily full
- **Retention**: 30 days
- **Location**: S3/GCS (multi-region)
- **Method**: `pg_dump` with WAL archiving

#### Redis/DragonflyDB
- **Frequency**: Daily snapshots
- **Retention**: 7 days
- **Location**: S3/GCS
- **Method**: RDB snapshots

### Application Backups

#### Kubernetes Resources
- **Frequency**: Daily
- **Retention**: 30 days
- **Method**: `kubectl get all -o yaml`

#### Container Images
- **Location**: Container Registry (ECR/GCR/ACR)
- **Retention**: All versions tagged
- **Method**: Automated via CI/CD

### Configuration Backups

- **Terraform State**: Stored in S3/GCS with versioning
- **Secrets**: Vault snapshots daily
- **ConfigMaps/Secrets**: Exported to Git

## Recovery Procedures

### RTO/RPO Targets

| Component | RTO | RPO |
|-----------|-----|-----|
| Critical Services | 15 min | 5 min |
| Database | 30 min | 1 hour |
| Storage | 1 hour | 1 hour |
| Entire Platform | 4 hours | 1 hour |

### Recovery Steps

#### 1. Database Recovery

```bash
# MongoDB
mongorestore --uri="mongodb://..." /backup/mongodb-20240101

# PostgreSQL
pg_restore -h postgres -U user -d streamverse /backup/postgres-20240101.dump
```

#### 2. Service Recovery

```bash
# Restore from Git
git checkout <commit-before-disaster>
kubectl apply -f k8s/base/

# Or restore from backup
kubectl apply -f /backup/k8s-20240101.yaml
```

#### 3. Data Recovery

```bash
# Restore from object storage
aws s3 sync s3://streamverse-backups/20240101/ /restore/
```

### Disaster Recovery Runbooks

#### Complete Platform Failure

1. **Assess Damage**: Identify affected components
2. **Provision Infrastructure**: Use Terraform to recreate
3. **Restore Databases**: From latest backups
4. **Deploy Services**: From container registry
5. **Restore Data**: From object storage backups
6. **Validate**: Run smoke tests
7. **Traffic Cutover**: Update DNS/load balancer

#### Database Corruption

1. **Stop Services**: Prevent further writes
2. **Identify Corruption**: Check logs and metrics
3. **Restore from Backup**: Use point-in-time recovery
4. **Validate Data**: Run integrity checks
5. **Resume Services**: Gradual rollout

#### Regional Failure

1. **Failover to Secondary Region**: Update DNS
2. **Scale Up**: Increase capacity in secondary region
3. **Data Sync**: Sync from primary when available
4. **Failback**: Once primary is restored

## Testing

### Backup Verification

```bash
# Test backup restoration in staging
./scripts/test-backup-restore.sh
```

### DR Drill

Quarterly disaster recovery drills:
1. Simulate failure scenario
2. Execute recovery procedures
3. Measure RTO/RPO
4. Document lessons learned

## Automation

### Automated Backups

```yaml
# Kubernetes CronJob
apiVersion: batch/v1
kind: CronJob
metadata:
  name: backup-databases
spec:
  schedule: "0 */6 * * *"  # Every 6 hours
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: postgres:15
            command:
              - /bin/sh
              - -c
              - |
                pg_dump -h postgres -U user streamverse | \
                gzip > /backup/postgres-$(date +%Y%m%d-%H%M%S).dump.gz
                aws s3 cp /backup/ s3://streamverse-backups/ --recursive
```

### Backup Monitoring

- **Success/Failure Alerts**: Via Prometheus/Alertmanager
- **Backup Size Monitoring**: Track growth trends
- **Retention Policy Enforcement**: Automatically delete old backups

## Multi-Region Replication

### Database Replication

- **MongoDB**: Replica sets across regions
- **PostgreSQL**: Streaming replication
- **Redis**: Redis Sentinel or Cluster mode

### Storage Replication

- **S3**: Cross-region replication enabled
- **GCS**: Multi-region buckets
- **Azure**: Geo-redundant storage

## Compliance

- **GDPR**: Data retention policies
- **SOC 2**: Backup and recovery procedures documented
- **ISO 27001**: Regular backup testing

