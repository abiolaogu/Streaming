# Operational Runbooks

Step-by-step procedures for common operations and incidents.

## Table of Contents

1. [Deployment Procedures](#deployment-procedures)
2. [Troubleshooting](#troubleshooting)
3. [Incident Response](#incident-response)
4. [Maintenance](#maintenance)
5. [Scaling](#scaling)

## Deployment Procedures

### Standard Deployment

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/new-feature
   ```

2. **Make Changes**
   ```bash
   # Edit code
   git add .
   git commit -m "feat: Add new feature"
   ```

3. **Create PR**
   - Push branch to GitHub
   - Create Pull Request
   - Wait for CI/CD to pass

4. **Merge to Main**
   - Merge PR after approval
   - CI/CD automatically deploys

5. **Verify Deployment**
   ```bash
   kubectl get pods -n streamverse
   ./tests/smoke-tests.sh
   ```

### Rollback Procedure

1. **Identify Bad Deployment**
   ```bash
   kubectl rollout history deployment/auth-service -n streamverse
   ```

2. **Rollback**
   ```bash
   kubectl rollout undo deployment/auth-service -n streamverse
   ```

3. **Verify Rollback**
   ```bash
   kubectl rollout status deployment/auth-service -n streamverse
   ```

## Troubleshooting

### Service Not Responding

1. **Check Pod Status**
   ```bash
   kubectl get pods -n streamverse | grep auth-service
   ```

2. **Check Logs**
   ```bash
   kubectl logs -f deployment/auth-service -n streamverse
   ```

3. **Check Events**
   ```bash
   kubectl describe pod auth-service-xxx -n streamverse
   ```

4. **Check Health Endpoint**
   ```bash
   curl http://auth-service.streamverse.svc.cluster.local:8080/health
   ```

### High Error Rate

1. **Identify Affected Service**
   ```bash
   # Prometheus query
   topk(10, sum(rate(http_requests_total{status=~"5.."}[5m])) by (service))
   ```

2. **Check Logs for Errors**
   ```bash
   kubectl logs -f deployment/auth-service -n streamverse | grep ERROR
   ```

3. **Check Resource Usage**
   ```bash
   kubectl top pods -n streamverse
   ```

4. **Scale Up if Needed**
   ```bash
   kubectl scale deployment/auth-service --replicas=5 -n streamverse
   ```

### Database Connection Issues

1. **Check Database Status**
   ```bash
   kubectl get pods -n streamverse | grep postgres
   ```

2. **Check Connection String**
   ```bash
   kubectl get secret mongodb-secret -n streamverse -o yaml
   ```

3. **Test Connection**
   ```bash
   kubectl run -it --rm debug --image=mongo --restart=Never -- \
     mongo "mongodb://..."
   ```

## Incident Response

### Service Outage

1. **Declare Incident**
   - Create incident in PagerDuty/Opsgenie
   - Notify on-call engineer

2. **Assess Impact**
   - Identify affected users
   - Check service status pages
   - Review monitoring dashboards

3. **Contain Issue**
   - Scale up services
   - Enable circuit breakers
   - Block problematic traffic

4. **Investigate Root Cause**
   - Check logs
   - Review recent deployments
   - Analyze metrics

5. **Resolve Issue**
   - Apply fix
   - Deploy hotfix if needed
   - Verify resolution

6. **Post-Mortem**
   - Document incident
   - Identify action items
   - Update runbooks

### Data Breach

1. **Contain Breach**
   - Isolate affected systems
   - Revoke compromised credentials
   - Enable enhanced logging

2. **Assess Impact**
   - Identify affected data
   - Determine scope of breach
   - Check for lateral movement

3. **Notify Stakeholders**
   - Legal team
   - Security team
   - Management
   - Customers (if required)

4. **Investigation**
   - Forensic analysis
   - Log review
   - Timeline reconstruction

5. **Remediation**
   - Patch vulnerabilities
   - Strengthen security controls
   - Update procedures

## Maintenance

### Database Maintenance

1. **Backup Before Maintenance**
   ```bash
   ./scripts/backup-databases.sh
   ```

2. **Perform Maintenance**
   ```bash
   # Vacuum PostgreSQL
   kubectl exec -it postgres-0 -n streamverse -- \
     psql -U postgres -d streamverse -c "VACUUM ANALYZE;"
   ```

3. **Verify After Maintenance**
   ```bash
   ./tests/smoke-tests.sh
   ```

### Certificate Rotation

1. **Generate New Certificates**
   ```bash
   cert-manager certs generate
   ```

2. **Update Secrets**
   ```bash
   kubectl create secret tls api-tls \
     --cert=cert.pem --key=key.pem \
     --dry-run=client -o yaml | kubectl apply -f -
   ```

3. **Restart Services**
   ```bash
   kubectl rollout restart deployment/kong -n streamverse
   ```

### Version Updates

1. **Update Base Images**
   ```bash
   # Update Dockerfile
   FROM golang:1.21-alpine
   ```

2. **Rebuild Images**
   ```bash
   docker build -t streamverse/auth-service:latest .
   ```

3. **Deploy**
   ```bash
   kubectl set image deployment/auth-service \
     auth-service=streamverse/auth-service:latest \
     -n streamverse
   ```

## Scaling

### Horizontal Scaling

```bash
# Scale deployment
kubectl scale deployment/auth-service --replicas=10 -n streamverse

# Or use HPA
kubectl autoscale deployment/auth-service \
  --min=3 --max=20 --cpu-percent=70 -n streamverse
```

### Vertical Scaling

```yaml
# Update resources in deployment
resources:
  requests:
    memory: "1Gi"
    cpu: "1000m"
  limits:
    memory: "2Gi"
    cpu: "2000m"
```

## Emergency Contacts

- **On-Call Engineer**: +1-XXX-XXX-XXXX
- **Security Team**: security@streamverse.io
- **DevOps Lead**: devops@streamverse.io

