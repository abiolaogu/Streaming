# Production Runbook

Operational procedures for preemption handling, failover, database restore, and rollback.

## Table of Contents

1. [Preemption Handling](#preemption-handling)
2. [Failover Procedures](#failover-procedures)
3. [Database/Object Store Restore](#databaseobject-store-restore)
4. [Rollback to Baseline](#rollback-to-baseline)
5. [Incident Response](#incident-response)

## Preemption Handling

### Spot Instance Termination

**Detection**: AWS spot termination notice (2-minute warning)

```bash
# Monitor termination notices
kubectl logs -n kube-system -l app=aws-node-termination-handler --tail=100

# List nodes with termination notices
kubectl get nodes -o custom-columns=NAME:.metadata.name,TERMINATING:.metadata.labels."node\.kubernetes\.io/instance-id"
```

**Response**:

1. **Graceful Drain** (automatic via node-termination-handler)
   ```bash
   # Node termination handler automatically drains on 2-min notice
   # Pods redeploy on other nodes
   ```

2. **Manual Override** (if needed)
   ```bash
   # Force drain with 5-minute grace
   kubectl drain <node-name> --grace-period=300 --ignore-daemonsets --delete-emptydir-data
   
   # Cordon node to prevent rescheduling
   kubectl cordon <node-name>
   ```

3. **Verify Workload Redistribution**
   ```bash
   # Check pod distribution
   kubectl get pods -o wide --all-namespaces | grep <node-name>
   
   # Verify HPA scaling
   kubectl get hpa -n media
   ```

### GPU Preemption Recovery

**Scenario**: All Spot GPU nodes preempted during peak transcoding

**Automatic Response**:

1. HPA detects pod pending
2. Attempts fallback to on-demand GPU pool
3. If unavailable, queues transcoding jobs

**Manual Intervention**:

```bash
# Check GPU node status
kubectl get nodes -l workload-type=gpu

# Force scale-up on-demand GPU if needed
kubectl scale deployment ome-transcoder --replicas=1 -n media

# Monitor transcoding queue
kubectl logs -n media -l app=ome-transcoder --tail=1000 | grep "QUEUED"
```

## Failover Procedures

### Tiered CDN Failover

**Primary → Shield → Edge → Origin**

```bash
# Test CDN cascade
curl -I https://streaming.example.com/live/channel-1.m3u8

# Check Varnish health
curl http://varnish-shield.cdn.svc.cluster.local:6081 | grep -i varnish

# Check ATS health
curl http://ats-edge.cdn.svc.cluster.local/_hostdb

# Check origin MinIO
mc admin info minio-service.data.svc.cluster.local
```

**Manual Failover**:

```bash
# Update ATC topology to skip shield
kubectl edit configmap atc-topology -n cdn
# Set shield.max-hosts: 0

# Reload ATC
kubectl rollout restart deployment atc-controller -n cdn
```

### Kafka Replication Lag

**Detection**: MirrorMaker2 lag > 1 minute

```bash
# Check MM2 status
kubectl logs -n data deployment/mm2-replicator --tail=1000 | grep lag

# Check Kafka consumer lag
kafka-consumer-groups --bootstrap-server kafka-client.data.svc.cluster.local:9092 \
  --describe --all-groups | grep LAG
```

**Response**:

```bash
# Increase MM2 replicas
kubectl scale deployment mm2-replicator --replicas=4 -n data

# Restart stuck connectors
kubectl delete pod -n data -l app=mm2

# Verify catch-up
watch -n 5 'kafka-consumer-groups --bootstrap-server kafka-client.data.svc.cluster.local:9092 --describe --group mm2-consumer | tail -1'
```

### Cross-Region Failover

**Scenario**: Entire AWS region down

1. **DNS Failover** (Route53 active-active)
   ```bash
   # Update Route53 weighted routing
   aws route53 change-resource-record-sets --hosted-zone-id Z123456 --change-batch file://failover.json
   ```

2. **Kubernetes Cluster Failover**
   ```bash
   # Switch kubectl context
   kubectl config use-context gcp-streaming-platform-prd
   
   # Scale up GCP cluster
   kubectl scale deployment ats-edge --replicas=10 -n cdn
   ```

3. **Data Layer Failover**
   ```bash
   # Fail over ScyllaDB (multi-region replication)
   # Update client connection strings
   kubectl edit configmap scylla-config -n data
   # scylla-seeds: scylla-0.scylla-svc-gcp.data.svc.cluster.local
   
   # Restart deployments
   kubectl rollout restart -n data deployment/*
   ```

## Database/Object Store Restore

### ScyllaDB Restore

**Full Backup Restore**:

```bash
# List snapshots
scyllabackup list s3://s3-backup-bucket/scylla/snapshots/

# Restore latest snapshot
scyllabackup restore --snapshot 2024-01-15T00:00:00Z \
  --s3-region us-east-1 \
  --s3-bucket s3-backup-bucket \
  --s3-prefix scylla/snapshots/ \
  --in-place \
  --ks streaming_platform

# Verify restore
cqlsh -e "SELECT COUNT(*) FROM streaming_platform.catalog;"
```

**Point-in-Time Recovery**:

```bash
# Use CDC logs (requires CDC enabled)
scylla-cdc-tool restore \
  --snapshot 2024-01-15T00:00:00Z \
  --target-time 2024-01-15T12:30:00Z \
  --ks streaming_platform
```

### ClickHouse Restore

**Table Restore**:

```bash
# Clone from replica
clickhouse-client --query "
  CREATE TABLE streaming_platform.qoe_raw_new AS streaming_platform.qoe_raw;

  INSERT INTO streaming_platform.qoe_raw_new
  SELECT * FROM remote('clickhouse-replica:8123',
    streaming_platform.qoe_raw,
    'default',
    'password'
  ) WHERE timestamp >= '2024-01-15 00:00:00'
    AND timestamp <= '2024-01-15 12:30:00';

  RENAME TABLE streaming_platform.qoe_raw TO streaming_platform.qoe_raw_old,
    streaming_platform.qoe_raw_new TO streaming_platform.qoe_raw;

  DROP TABLE streaming_platform.qoe_raw_old;
"
```

**Full Restore**:

```bash
# Restore from S3 backup
clickhouse-client --query "
  BACKUP DATABASE streaming_platform TO S3('s3://backup-bucket/clickhouse/full/2024-01-15') ASYNC;
"
```

### MinIO Restore

**Bucket Restore**:

```bash
# List versions
mc ls --versions minio-svc/data/vod/

# Restore specific version
mc cp --version-id ABC123DEF456 \
  minio-svc/data/vod/video.mp4 \
  minio-svc/data/vod/video.mp4

# Bulk restore from backup
mc mirror --preserve --force \
  minio-backup/vod/ \
  minio-svc/data/vod/
```

## Rollback to Baseline

### Scale to Minimum

**Emergency Cost Reduction**:

```bash
# CPU scale-down
kubectl scale deployment ome-transcoder --replicas=1 -n media
kubectl scale deployment ats-edge --replicas=2 -n cdn
kubectl scale deployment shaka-packager --replicas=1 -n media

# GPU scale-to-zero
kubectl scale deployment ome-transcoder --replicas=0 -n media

# Wait for Spot node termination
sleep 300

# Verify baseline
kubectl get nodes --no-headers | wc -l  # Should be ~1 per cloud
```

### Terraform Rollback

**Infrastructure Rollback**:

```bash
# List terraform state versions
terraform state list

# Show current state
terraform show

# Revert to previous state
terraform state push tfstate-backup-2024-01-15.tfstate

# Re-apply with old configuration
terraform apply -var-file=dev-previous.tfvars
```

### Application Rollback

**Argo Rollouts**:

```bash
# Check rollout history
kubectl rollout history deployment/web-client -n platform

# Rollback to previous version
kubectl rollout undo deployment/web-client -n platform

# Rollback to specific revision
kubectl rollout undo deployment/web-client --to-revision=5 -n platform
```

**Helm Rollback**:

```bash
# List releases
helm list -n platform

# Rollback release
helm rollback streaming-platform-platform web-client-v1.2.0 -n platform
```

## Incident Response

### P1: Complete Service Outage

**Detection**: All players failing, CDN returning 500s

**Response**:

1. **Immediate** (0-5 minutes)
   - Check CloudWatch/Stackdriver dashboards
   - Verify DNS resolver
   - Test origin accessibility

2. **Assess** (5-15 minutes)
   - Review logs: `kubectl logs --since=15m -n media,cdn,platform`
   - Check node health: `kubectl get nodes`
   - Verify ingress controller: `kubectl get ingress -A`

3. **Mitigate** (15-30 minutes)
   - Scale up services: `kubectl scale deployment --all --replicas=3 -n media`
   - Restart stuck pods: `kubectl delete pod -l app=<service> -n <namespace>`
   - Enable maintenance mode if needed

4. **Resolve** (30+ minutes)
   - Root cause analysis
   - Permanent fix deployment
   - Post-mortem documentation

### P2: Partial Degradation

**Examples**: Slow startup, occasional stuttering, one region down

**Response**:

1. **Identify** affected service/region
2. **Scale** targeted deployment
3. **Failover** to backup region if needed
4. **Monitor** metrics for recovery

### P3: Quality Issues

**Examples**: Low bitrate streams, watermark visible, audio sync off

**Response**:

1. **Verify** OME transcoding configs
2. **Check** DRM key rotation
3. **Test** player compatibility
4. **Update** packaging profiles if needed

### Communication

**Stakeholder Updates**:

- **P1**: Every 15 minutes
- **P2**: Every 30 minutes
- **P3**: End of day

**Post-Mortem** (within 48 hours):

- Timeline
- Root cause
- Action items
- Prevention measures

## Monitoring Checklist

**Daily**:

- [ ] On-demand node count = 1 per cloud
- [ ] All HPA active (no stuck scaling)
- [ ] Kafka replication lag < 10 seconds
- [ ] ClickHouse ingest rate > 10k rows/sec
- [ ] CDN hit rate > 95%

**Weekly**:

- [ ] Review cost reports
- [ ] Analyze Spot preemption frequency
- [ ] Check DRM license expiry dates
- [ ] Verify backup integrity
- [ ] Review security scan results

**Monthly**:

- [ ] Capacity planning review
- [ ] Disaster recovery drill
- [ ] Terraform state backup verification
- [ ] Dependency updates assessment
- [ ] Billing reconciliation

