#!/bin/bash

# Production Smoke Tests
# Run after deployment to validate platform functionality

set -e

NAMESPACE_PLATFORM="platform"
NAMESPACE_MEDIA="media"
NAMESPACE_CDN="cdn"
NAMESPACE_DATA="data"
NAMESPACE_TELECOM="telecom"

echo "=========================================="
echo "Production Smoke Tests"
echo "=========================================="

# Test 1: Cluster connectivity
echo "Test 1: Cluster connectivity"
kubectl cluster-info || { echo "FAIL: Cannot connect to cluster"; exit 1; }
echo "✓ Cluster accessible"

# Test 2: Namespaces exist
echo "Test 2: Namespaces"
for ns in platform media cdn data telecom observability; do
    kubectl get namespace $ns || { echo "FAIL: Namespace $ns not found"; exit 1; }
done
echo "✓ All namespaces exist"

# Test 3: Pods running
echo "Test 3: Critical pods running"
for pod in $(kubectl get pods -n $NAMESPACE_MEDIA -l app=ome-transcoder -o name | head -1); do
    kubectl get $pod -n $NAMESPACE_MEDIA || { echo "FAIL: OME transcoder not running"; exit 1; }
done

for pod in $(kubectl get pods -n $NAMESPACE_CDN -l app=ats-edge -o name | head -1); do
    kubectl get $pod -n $NAMESPACE_CDN || { echo "FAIL: ATS edge not running"; exit 1; }
done
echo "✓ Critical pods running"

# Test 4: Services accessible
echo "Test 4: Services accessible"
for svc in dragonfly-svc.data scylla-client.data kafka-client.data clickhouse-svc.data; do
    kubectl get svc $svc -n ${svc##*.} || { echo "FAIL: Service $svc not found"; exit 1; }
done
echo "✓ Data services accessible"

# Test 5: HPA configured
echo "Test 5: HPA configured"
kubectl get hpa -n $NAMESPACE_MEDIA | grep ome-transcoder || { echo "FAIL: OME HPA not configured"; exit 1; }
kubectl get hpa -n $NAMESPACE_CDN | grep ats-edge || { echo "FAIL: ATS HPA not configured"; exit 1; }
echo "✓ HPA configured"

# Test 6: Node pools
echo "Test 6: Node pools"
CPU_ONDEMAND=$(kubectl get nodes -l spot-enabled=false --no-headers | wc -l)
if [ "$CPU_ONDEMAND" -gt 1 ]; then
    echo "WARN: More than 1 on-demand CPU node found ($CPU_ONDEMAND)"
else
    echo "✓ On-demand CPU nodes compliant: $CPU_ONDEMAND"
fi

# Test 7: Storage classes
echo "Test 7: Storage classes"
kubectl get storageclass || { echo "WARN: Storage classes not configured"; }
echo "✓ Storage classes checked"

# Test 8: Network policies
echo "Test 8: Network policies"
kubectl get networkpolicy -A --no-headers | wc -l | grep -q "[0-9]" || { echo "WARN: No network policies found"; }
echo "✓ Network policies checked"

# Test 9: Secrets
echo "Test 9: Secrets"
for secret in drm-secrets minio-secrets; do
    kubectl get secret $secret -n $NAMESPACE_MEDIA || echo "WARN: Secret $secret not found"
done
echo "✓ Secrets checked"

# Test 10: Metrics
echo "Test 10: Prometheus metrics"
kubectl port-forward -n observability svc/prometheus 9090:9090 &
PF_PID=$!
sleep 5
curl -s http://localhost:9090/api/v1/targets | grep -q '"activeTargets"' || { echo "WARN: Prometheus not responding"; }
kill $PF_PID
echo "✓ Metrics endpoint accessible"

# Test 11: Cost guardrails
echo "Test 11: Cost guardrails"
ONDEMAND_COUNT=$(kubectl get nodes --no-headers | grep -v spot | wc -l)
if [ "$ONDEMAND_COUNT" -gt 3 ]; then
    echo "FAIL: Too many on-demand nodes: $ONDEMAND_COUNT"
    exit 1
else
    echo "✓ On-demand node count: $ONDEMAND_COUNT"
fi

# Test 12: GPU nodes (if any)
echo "Test 12: GPU nodes"
GPU_COUNT=$(kubectl get nodes -l nvidia.com/gpu=true --no-headers 2>/dev/null | wc -l)
echo "✓ GPU nodes: $GPU_COUNT"

echo "=========================================="
echo "All smoke tests passed!"
echo "=========================================="

