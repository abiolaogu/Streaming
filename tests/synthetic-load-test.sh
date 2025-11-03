#!/bin/bash

# Synthetic Load Test
# Trigger autoscaling to validate Spot node provisioning

set -e

NAMESPACE="media"
TEST_DURATION="${1:-300}"  # 5 minutes default

echo "=========================================="
echo "Synthetic Load Test"
echo "Duration: ${TEST_DURATION}s"
echo "=========================================="

# Create test job
cat <<EOF | kubectl apply -f -
apiVersion: batch/v1
kind: Job
metadata:
  name: load-test-$(date +%s)
  namespace: ${NAMESPACE}
spec:
  parallelism: 50
  completions: 500
  template:
    spec:
      tolerations:
      - key: spot-enabled
        operator: Equal
        value: "true"
        effect: NoSchedule
      restartPolicy: Never
      containers:
      - name: cpu-stress
        image: progrium/stress:latest
        command: ["/usr/local/bin/stress"]
        args: ["--cpu", "4", "--timeout", "${TEST_DURATION}s"]
        resources:
          requests:
            cpu: "4000m"
            memory: "1Gi"
          limits:
            cpu: "4000m"
            memory: "1Gi"
EOF

echo "✓ Test job created"

# Monitor scaling
echo "Monitoring autoscaling..."
END_TIME=$((SECONDS + TEST_DURATION))

while [ $SECONDS -lt $END_TIME ]; do
    NODE_COUNT=$(kubectl get nodes --no-headers | wc -l)
    SPOT_COUNT=$(kubectl get nodes -l spot-enabled=true --no-headers 2>/dev/null | wc -l)
    PENDING_PODS=$(kubectl get pods -n ${NAMESPACE} --field-selector=status.phase=Pending --no-headers 2>/dev/null | wc -l)
    
    echo "[$(date +%H:%M:%S)] Nodes: $NODE_COUNT (Spot: $SPOT_COUNT) | Pending Pods: $PENDING_PODS"
    
    sleep 10
done

# Cleanup
echo "Cleaning up test job..."
kubectl delete job -n ${NAMESPACE} --selector job-name -l app=load-test

echo "Waiting for scale-down..."
sleep 60

FINAL_NODE_COUNT=$(kubectl get nodes --no-headers | wc -l)
echo "=========================================="
echo "Test complete"
echo "Final node count: $FINAL_NODE_COUNT"
echo "=========================================="

