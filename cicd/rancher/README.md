# Rancher Deployment Configuration

GitOps-based deployment using Rancher Fleet.

## Setup

### 1. Rancher Management Cluster

```bash
# Install Rancher
helm repo add rancher-stable https://releases.rancher.com/server-charts/stable
helm install rancher rancher-stable/rancher \
  --namespace cattle-system \
  --create-namespace \
  --set hostname=rancher.streamverse.io \
  --set bootstrapPassword=admin
```

### 2. Fleet Configuration

Fleet automatically syncs Git repositories to clusters.

```yaml
# fleet.yaml
apiVersion: fleet.cattle.io/v1alpha1
kind: GitRepo
metadata:
  name: streamverse
  namespace: fleet-default
spec:
  repo: https://github.com/streamverse/streamverse-platform
  branch: main
  paths:
    - k8s/base
    - k8s/overlays/production
```

### 3. Cluster Registration

1. Navigate to Rancher UI
2. Go to "Cluster Management" > "Import Existing"
3. Follow instructions to register cluster
4. Cluster will appear in Fleet

## Deployment Flow

1. **Code Push** → GitHub
2. **CI/CD** → Jenkins/Tekton builds and pushes images
3. **GitOps** → Fleet detects changes in Git
4. **Sync** → Fleet applies manifests to clusters
5. **Monitoring** → Prometheus/Grafana track deployment

## Multi-Cluster Deployment

Deploy to multiple clusters using Fleet:

```yaml
# k8s/overlays/production/fleet.yaml
fleet:
  clusterLabels:
    environment: production
    region: us-east-1
```

## Rollback

```bash
# Manual rollback via kubectl
kubectl rollout undo deployment/auth-service -n streamverse

# Git-based rollback
git revert <commit-hash>
git push origin main
# Fleet automatically syncs
```

## Monitoring Deployments

- **Rancher UI**: View deployments in cluster
- **Fleet Status**: Check Git sync status
- **Argo CD**: Optional GitOps UI

