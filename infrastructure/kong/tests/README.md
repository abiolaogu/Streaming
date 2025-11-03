# Kong API Gateway Testing

This directory contains testing scripts and configurations for the Kong API Gateway.

## Test Files

- `integration_test.sh` - Integration tests for routing, JWT validation, rate limiting, CORS
- `load_test.js` - k6 load testing script (target: 1000 RPS with no errors)
- `health_check.sh` - Health check verification for all services

## Prerequisites

1. **Kong Gateway** running (via Docker Compose or Kubernetes)
2. **All services** running and accessible
3. **k6** installed (for load testing):
   ```bash
   # macOS
   brew install k6
   
   # Linux
   sudo gpg -k
   sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
   echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
   sudo apt-get update
   sudo apt-get install k6
   ```

## Running Tests

### Integration Tests

```bash
cd infrastructure/kong/tests
chmod +x integration_test.sh
./integration_test.sh
```

**Environment Variables**:
- `KONG_ADMIN_URL` - Kong Admin API URL (default: `http://localhost:8001`)
- `KONG_PROXY_URL` - Kong Proxy URL (default: `http://localhost:8000`)
- `BASE_URL` - Base URL for API requests (default: `http://localhost:8000`)

### Load Testing

```bash
cd infrastructure/kong/tests
k6 run load_test.js
```

**Environment Variables**:
- `BASE_URL` - Base URL for API requests (default: `http://localhost:8000`)

**Expected Results**:
- ✅ 1000 RPS target achieved
- ✅ P95 latency < 500ms
- ✅ P99 latency < 1000ms
- ✅ Error rate < 1%

### Health Check Verification

```bash
cd infrastructure/kong/tests
chmod +x health_check.sh
./health_check.sh
```

**Environment Variables**:
- `BASE_URL` - Base URL for API requests (default: `http://localhost:8000`)

## Test Coverage

### Integration Tests

1. ✅ Kong Admin API health
2. ✅ Kong Proxy accessibility
3. ✅ Service health checks (all 12 services)
4. ✅ JWT validation (protected endpoints)
5. ✅ Rate limiting
6. ✅ CORS headers
7. ✅ Route matching

### Load Tests

- Ramp up: 0 → 100 → 500 → 1000 concurrent users
- Duration: ~7 minutes total
- Tests all service endpoints randomly
- Measures: latency, error rate, throughput

### Health Checks

- Verifies all 12 services are accessible through Kong
- 5-second timeout per service
- Reports healthy/unhealthy status

## Acceptance Criteria

Based on Issue #22 requirements:

- ✅ All 12 services routable through Kong
- ✅ Rate limiting works
- ✅ JWT validation works
- ✅ Health checks pass
- ✅ Load test shows no errors under 1000 RPS

## Troubleshooting

### Services not accessible

1. Check if services are running:
   ```bash
   docker ps
   ```

2. Check Kong service configuration:
   ```bash
   curl http://localhost:8001/services
   ```

3. Check Kong routes:
   ```bash
   curl http://localhost:8001/routes
   ```

### JWT validation failing

1. Verify JWT secret is configured in Kong
2. Check Auth Service is running and accessible
3. Verify token format (Bearer token)

### Rate limiting not working

1. Check rate limiting plugin is enabled:
   ```bash
   curl http://localhost:8001/plugins
   ```

2. Verify rate limit configuration in `kong.yml`

## CI/CD Integration

These tests can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run Kong Integration Tests
  run: |
    cd infrastructure/kong/tests
    ./integration_test.sh

- name: Run Kong Load Tests
  run: |
    k6 run infrastructure/kong/tests/load_test.js
```

## Next Steps

- [ ] Add authentication token generation for protected endpoint tests
- [ ] Add performance benchmarking reports
- [ ] Add test results visualization
- [ ] Integrate with Prometheus for metrics collection

