# CDN Infrastructure - Apache Traffic Control/Server

CDN infrastructure configuration for StreamVerse platform.

## Components

- Traffic Ops (Control Plane API and UI)
- Traffic Router (DNS/HTTP routing)
- Traffic Server (L1/L2 cache nodes)
- Traffic Monitor (Health checking)

## Deployment

### Using Docker Compose

```bash
cd infrastructure/cdn
docker-compose up -d
```

### Using Terraform

```bash
cd infrastructure/cdn/terraform
terraform init
terraform plan
terraform apply
```

### Using Ansible

```bash
cd infrastructure/cdn/ansible
ansible-playbook site.yml
```

## Configuration

Edit configuration files in `config/` directory.

