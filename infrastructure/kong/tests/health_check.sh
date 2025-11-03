#!/bin/bash

# Health check verification script for all services through Kong

set -e

BASE_URL="${BASE_URL:-http://localhost:8000}"

echo "🏥 Kong API Gateway Health Check Verification"
echo "=============================================="

# Services to check
services=(
    "auth:/auth/health"
    "user:/users/health"
    "content:/content/health"
    "streaming:/streaming/health"
    "transcoding:/transcode/health"
    "payment:/payments/health"
    "search:/search/health"
    "analytics:/analytics/health"
    "recommendation:/recommendations/health"
    "notification:/notifications/health"
    "admin:/admin/health"
    "scheduler:/scheduler/health"
)

PASSED=0
FAILED=0

for service in "${services[@]}"; do
    name="${service%%:*}"
    path="${service#*:}"
    
    echo -n "Checking $name service... "
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL$path" --max-time 5)
    
    if [ "$response" -eq 200 ]; then
        echo "✅ Healthy (HTTP $response)"
        ((PASSED++))
    else
        echo "❌ Unhealthy (HTTP $response)"
        ((FAILED++))
    fi
done

echo ""
echo "=============================================="
echo "Health Check Summary"
echo "=============================================="
echo "✅ Healthy: $PASSED"
echo "❌ Unhealthy: $FAILED"
echo "Total: $((PASSED + FAILED))"

if [ $FAILED -eq 0 ]; then
    echo ""
    echo "✅ All services are healthy!"
    exit 0
else
    echo ""
    echo "❌ Some services are unhealthy. Please check the services."
    exit 1
fi

