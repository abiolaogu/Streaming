# Monitoring Stack - Prometheus, Grafana, Loki, Jaeger

Complete observability stack for StreamVerse platform.

## Components

### Prometheus
- **Metrics Collection**: Time-series database for metrics
- **Alerting**: Rule-based alerting via Alertmanager
- **Service Discovery**: Automatic discovery of Kubernetes pods/services

### Grafana
- **Visualization**: Dashboards for metrics and logs
- **Alerting**: Alert notifications via multiple channels
- **Dashboards**: Pre-built dashboards for services

### Loki
- **Log Aggregation**: Centralized log collection
- **Label-based Indexing**: Efficient log queries
- **Integration**: Native Grafana integration

### Jaeger
- **Distributed Tracing**: Request tracing across services
- **Performance Analysis**: Identify bottlenecks
- **Service Map**: Visual service dependencies

## Deployment

### Using Helm

```bash
# Prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack \
  -n monitoring \
  --create-namespace

# Loki
helm repo add grafana https://grafana.github.io/helm-charts
helm install loki grafana/loki-stack \
  -n monitoring

# Jaeger
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm install jaeger jaegertracing/jaeger \
  -n monitoring
```

### Using Kubernetes Manifests

```bash
kubectl apply -f k8s/monitoring/
```

## Configuration

### Prometheus

- **Scrape Interval**: 15s
- **Evaluation Interval**: 15s
- **Retention**: 30 days (configurable)

### Grafana

- **Default Admin Password**: Set via secret
- **Dashboards**: Auto-imported from ConfigMap
- **Data Sources**: Auto-configured for Prometheus and Loki

### Loki

- **Retention**: 30 days
- **Storage**: Filesystem (use S3/GCS for production)
- **Stream Labels**: namespace, service, pod

### Jaeger

- **Storage**: Elasticsearch (configured separately)
- **Sampling**: 1% (configurable)
- **UI**: Available at http://jaeger-query:16686

## Usage

### Accessing Dashboards

- **Grafana**: `http://grafana.monitoring.svc.cluster.local:3000`
- **Prometheus**: `http://prometheus.monitoring.svc.cluster.local:9090`
- **Jaeger UI**: `http://jaeger-query.monitoring.svc.cluster.local:16686`

### Querying Metrics

**Prometheus Query Language (PromQL)**:

```promql
# Request rate
rate(http_requests_total[5m])

# Error rate
rate(http_requests_total{status=~"5.."}[5m])

# P95 latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

### Querying Logs

**LogQL** (Loki Query Language):

```logql
# All logs from auth-service
{app="auth-service"}

# Error logs
{app="auth-service"} |= "error"

# Logs with specific label
{namespace="streamverse", service="auth-service"}
```

### Viewing Traces

1. Navigate to Jaeger UI
2. Select service from dropdown
3. Click "Find Traces"
4. View trace timeline and spans

## Alerts

Pre-configured alerts:
- **High Error Rate**: > 5% errors
- **High Latency**: P99 > 1s
- **Pod Crash Looping**: Restarts > 0
- **Pod Not Ready**: Not ready for > 10m
- **High CPU Usage**: > 80%
- **High Memory Usage**: > 85%
- **Disk Space Low**: < 15% free
- **Database Connection Failure**: Cannot connect

## Customization

### Adding New Dashboards

1. Create dashboard JSON in `grafana/dashboards/`
2. Apply ConfigMap:
   ```bash
   kubectl apply -f k8s/monitoring/grafana-dashboards-configmap.yaml
   ```

### Adding New Alerts

1. Edit `prometheus/alerts/alerts.yml`
2. Reload Prometheus:
   ```bash
   kubectl exec -it prometheus-0 -n monitoring -- kill -HUP 1
   ```

### Service Instrumentation

Add Prometheus client library to services:

**Go**:
```go
import "github.com/prometheus/client_golang/prometheus"

var httpRequests = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total HTTP requests",
    },
    []string{"method", "endpoint", "status"},
)
```

**Python**:
```python
from prometheus_client import Counter

http_requests = Counter('http_requests_total', 'Total HTTP requests', ['method', 'endpoint', 'status'])
```

## Production Considerations

1. **High Availability**: Run 3+ replicas of Prometheus
2. **Storage**: Use object storage (S3/GCS) for long-term retention
3. **Scaling**: Use Thanos for horizontal scaling
4. **Security**: Enable authentication for Grafana
5. **Retention**: Configure appropriate retention policies
6. **Backup**: Regular backups of Prometheus data

