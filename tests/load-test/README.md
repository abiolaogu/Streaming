# Load Testing & Performance Optimization

Comprehensive load testing strategy for StreamVerse platform.

## Tools

- **k6**: HTTP load testing
- **JMeter**: GUI-based testing (optional)
- **Gatling**: Scala-based testing (optional)
- **Locust**: Python-based testing (optional)

## k6 Load Testing

### Installation

```bash
# macOS
brew install k6

# Linux
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D6
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

### Running Tests

```bash
# Basic load test
k6 run load-test.js

# With environment variables
BASE_URL=https://api.streamverse.io AUTH_TOKEN=xxx k6 run load-test.js

# Distributed testing
k6 run --out cloud load-test.js
```

## Test Scenarios

### 1. Smoke Test
- **Users**: 1
- **Duration**: 1 minute
- **Purpose**: Verify system works

### 2. Load Test
- **Users**: 100-500
- **Duration**: 10-30 minutes
- **Purpose**: Verify normal load capacity

### 3. Stress Test
- **Users**: 500-2000
- **Duration**: 30-60 minutes
- **Purpose**: Find breaking point

### 4. Spike Test
- **Users**: 0 -> 1000 -> 0
- **Duration**: 5 minutes
- **Purpose**: Test sudden traffic spikes

### 5. Endurance Test
- **Users**: 200
- **Duration**: 24 hours
- **Purpose**: Test for memory leaks

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Response Time (P95) | < 500ms | HTTP request duration |
| Response Time (P99) | < 1000ms | HTTP request duration |
| Error Rate | < 1% | Failed requests / Total requests |
| Throughput | > 1000 req/s | Requests per second |
| CPU Usage | < 70% | Average CPU utilization |
| Memory Usage | < 80% | Average memory utilization |

## Monitoring During Tests

```bash
# Prometheus queries during load test
# Request rate
sum(rate(http_requests_total[1m])) by (service)

# Error rate
sum(rate(http_requests_total{status=~"5.."}[1m])) by (service)

# Latency
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[1m])) by (le, service))
```

## Optimization Strategies

### 1. Database Optimization
- Add indexes on frequently queried fields
- Use connection pooling
- Enable query caching
- Optimize slow queries

### 2. Caching
- Redis/DragonflyDB for frequently accessed data
- CDN for static assets
- Application-level caching

### 3. Load Balancing
- Round-robin or least-connections
- Health checks
- Circuit breakers

### 4. Auto-scaling
- Horizontal Pod Autoscaling (HPA)
- Cluster Autoscaling
- Predictive scaling

### 5. Database Sharding
- Partition by tenant_id
- Partition by content_category
- Read replicas for read-heavy workloads

## Benchmark Results

Store results in `tests/load-test/results/`:
- Response times
- Throughput
- Error rates
- Resource utilization

## Continuous Load Testing

Integrate into CI/CD:

```yaml
# .github/workflows/load-test.yml
name: Load Test
on:
  schedule:
    - cron: '0 2 * * *'  # Daily at 2 AM
jobs:
  load-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run k6
        run: k6 run tests/load-test/k6-load-test.js
```

## Performance Budgets

Define performance budgets in CI:

```javascript
// performance-budget.js
export const budgets = {
  'p95_latency': 500,    // ms
  'p99_latency': 1000,   // ms
  'error_rate': 0.01,     // 1%
  'throughput': 1000,     // req/s
};
```

