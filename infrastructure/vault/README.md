# HashiCorp Vault Setup

Secrets management for StreamVerse platform.

## Deployment

### Using Helm

```bash
helm repo add hashicorp https://helm.releases.hashicorp.com
helm install vault hashicorp/vault \
  -n vault \
  --create-namespace \
  -f values.yaml
```

### Manual Deployment

```bash
kubectl apply -f k8s/vault/
```

## Initialization

1. **Initialize Vault**:
   ```bash
   kubectl exec -it vault-0 -n vault -- vault operator init
   ```

2. **Unseal Vault** (3 times with different keys):
   ```bash
   kubectl exec -it vault-0 -n vault -- vault operator unseal <key>
   ```

3. **Login**:
   ```bash
   kubectl exec -it vault-0 -n vault -- vault login <root-token>
   ```

## Configuration

### Enable KV Secrets Engine

```bash
vault secrets enable -version=2 kv
```

### Enable Database Secrets Engine

```bash
vault secrets enable database

vault write database/config/postgresql \
  plugin_name=postgresql-database-plugin \
  allowed_roles="*" \
  connection_url="postgresql://{{username}}:{{password}}@postgres:5432/streamverse?sslmode=disable" \
  username="vault" \
  password="vault-password"
```

### Create Policies

```bash
# Policy for services
vault policy write streamverse-service - <<EOF
path "kv/data/streamverse/*" {
  capabilities = ["read"]
}

path "database/creds/streamverse" {
  capabilities = ["read"]
}
EOF
```

### Create Service Account

```bash
vault auth enable kubernetes

vault write auth/kubernetes/config \
  kubernetes_host="https://kubernetes.default.svc" \
  kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt

vault write auth/kubernetes/role/streamverse-service \
  bound_service_account_names=default \
  bound_service_account_namespaces=streamverse \
  policies=streamverse-service \
  ttl=1h
```

## Usage in Services

### Kubernetes Integration

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: auth-service
  namespace: streamverse
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
spec:
  template:
    spec:
      serviceAccountName: auth-service
      containers:
      - name: auth-service
        image: streamverse/auth-service:latest
        env:
        - name: MONGODB_URI
          value: vault:kv/data/streamverse/mongodb
        volumeMounts:
        - name: vault-token
          mountPath: /var/run/secrets/vault
          readOnly: true
```

### Vault Agent Sidecar

```yaml
initContainers:
- name: vault-agent
  image: hashicorp/vault:latest
  command:
    - vault
    - agent
    - -config=/vault/config
  volumeMounts:
    - name: vault-config
      mountPath: /vault/config
    - name: vault-secrets
      mountPath: /vault/secrets
```

## Secrets Rotation

### Database Credentials

```bash
# Rotate root password
vault write -force database/rotate-root/postgresql

# Generate new credentials
vault read database/creds/streamverse
```

### Application Secrets

```bash
# Rotate JWT secret
vault kv put kv/streamverse/jwt secret-key="$(openssl rand -base64 32)"
```

## Backup and Restore

### Backup

```bash
vault operator raft snapshot save /backup/vault-$(date +%Y%m%d).snapshot
```

### Restore

```bash
vault operator raft snapshot restore /backup/vault-20240101.snapshot
```

## Production Considerations

1. **High Availability**: Run 3+ Vault instances
2. **Auto-unseal**: Use AWS KMS for automatic unsealing
3. **Audit Logging**: Enable audit logs
4. **Network Policies**: Restrict access to Vault
5. **TLS**: Use certificates for all connections
6. **Monitoring**: Monitor Vault health metrics

