terraform {
  required_version = ">= 1.5.0"
  
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }

  backend "gcs" {
    bucket = "streamverse-terraform-state"
    prefix = "gcp/terraform.tfstate"
  }
}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

# VPC Network
resource "google_compute_network" "main" {
  name                    = "${var.project_name}-vpc"
  auto_create_subnetworks = false

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

# Subnets
resource "google_compute_subnetwork" "private" {
  count = length(var.gke_node_zones)

  name          = "${var.project_name}-private-${count.index + 1}"
  ip_cidr_range = cidrsubnet(var.vpc_cidr, 8, count.index)
  region        = var.gcp_region
  network       = google_compute_network.main.id
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "pods"
    ip_cidr_range = cidrsubnet("10.1.0.0/16", 8, count.index)
  }

  secondary_ip_range {
    range_name    = "services"
    ip_cidr_range = cidrsubnet("10.2.0.0/16", 8, count.index)
  }
}

# Router for NAT
resource "google_compute_router" "router" {
  name    = "${var.project_name}-router"
  region  = var.gcp_region
  network = google_compute_network.main.id
}

# NAT Gateway
resource "google_compute_router_nat" "nat" {
  name   = "${var.project_name}-nat"
  router = google_compute_router.router.name
  region = var.gcp_region

  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}

# GKE Cluster
resource "google_container_cluster" "primary" {
  name     = "${var.project_name}-cluster"
  location = var.gcp_region

  # Regional cluster for HA
  node_locations = var.gke_node_zones

  remove_default_node_pool = true
  initial_node_count       = 1

  network    = google_compute_network.main.name
  subnetwork = google_compute_subnetwork.private[0].name

  # Workload Identity
  workload_identity_config {
    workload_pool = "${var.gcp_project_id}.svc.id.goog"
  }

  # Logging and Monitoring
  logging_config {
    enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
  }

  monitoring_config {
    enable_components = ["SYSTEM_COMPONENTS"]
    
    managed_prometheus {
      enabled = true
    }
  }

  # Binary Authorization
  binary_authorization {
    evaluation_mode = "PROJECT_SINGLETON_POLICY_ENFORCE"
  }

  # Network Policy
  network_policy {
    enabled  = true
    provider = "CALICO"
  }

  # Database encryption
  database_encryption {
    state    = "ENCRYPTED"
    key_name = google_kms_crypto_key.gke.id
  }

  # Maintenance window
  maintenance_policy {
    daily_maintenance_window {
      start_time = "03:00"
    }
  }

  # Release channel
  release_channel {
    channel = "REGULAR"
  }

  # IP allocation
  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }

  # Private cluster
  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block   = "172.16.0.0/28"
  }

  master_authorized_networks_config {
    cidr_blocks {
      cidr_block   = "0.0.0.0/0" # Restrict in production
      display_name = "All"
    }
  }
}

# Node Pool
resource "google_container_node_pool" "primary" {
  name       = "${var.project_name}-node-pool"
  location   = var.gcp_region
  cluster    = google_container_cluster.primary.name
  node_count = var.node_desired_size

  autoscaling {
    min_node_count = var.node_min_size
    max_node_count = var.node_max_size
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  node_config {
    preemptible  = var.environment != "production"
    machine_type = var.node_machine_type
    disk_size_gb = 50
    disk_type    = "pd-ssd"

    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]

    workload_metadata_config {
      mode = "GKE_METADATA"
    }

    labels = {
      workload = "general"
    }
  }
}

# Cloud SQL (PostgreSQL)
resource "google_sql_database_instance" "postgres" {
  name             = "${var.project_name}-db"
  database_version = "POSTGRES_15"
  region           = var.gcp_region

  settings {
    tier              = var.db_instance_tier
    availability_type = var.environment == "production" ? "REGIONAL" : "ZONAL"

    backup_configuration {
      enabled                        = true
      start_time                     = "03:00"
      point_in_time_recovery_enabled = true
      transaction_log_retention_days = 7
      backup_retention_settings {
        retained_backups = 7
      }
    }

    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.main.id
      require_ssl     = true
    }

    database_flags {
      name  = "max_connections"
      value = "200"
    }

    insights_config {
      query_insights_enabled  = true
      query_string_length     = 1024
      record_application_tags = true
      record_client_address   = true
    }

    deletion_protection = var.environment == "production"
  }

  deletion_protection = var.environment == "production"

  depends_on = [google_service_networking_connection.private_vpc]
}

resource "google_sql_database" "main" {
  name     = var.db_name
  instance = google_sql_database_instance.postgres.name
}

resource "google_sql_user" "main" {
  name     = var.db_username
  instance = google_sql_database_instance.postgres.name
  password = var.db_password
}

# Memorystore (Redis/DragonflyDB compatible)
resource "google_redis_instance" "redis" {
  name           = "${var.project_name}-redis"
  tier           = var.redis_tier
  memory_size_gb = var.redis_memory_size
  region         = var.gcp_region

  location_id             = var.gke_node_zones[0]
  alternative_location_id = var.gke_node_zones[1]

  authorized_network = google_compute_network.main.id

  redis_version    = "REDIS_7_0"
  display_name     = "${var.project_name} Redis"
  reserved_ip_range = "10.3.0.0/29"

  maintenance_policy {
    weekly_maintenance_window {
      day = "SUNDAY"
      start_time {
        hours   = 3
        minutes = 0
      }
    }
  }
}

# Cloud Storage Buckets
resource "google_storage_bucket" "media" {
  name          = "${var.project_name}-media-${var.environment}"
  location      = var.gcp_region
  force_destroy = var.environment != "production"

  uniform_bucket_level_access = true

  versioning {
    enabled = true
  }

  encryption {
    default_kms_key_name = google_kms_crypto_key.storage.id
  }

  lifecycle_rule {
    condition {
      age = 90
    }
    action {
      type = "Delete"
    }
  }
}

# Artifact Registry
resource "google_artifact_registry_repository" "services" {
  for_each = toset(var.services)

  location      = var.gcp_region
  repository_id = "${var.project_name}-${each.value}"
  description   = "Docker repository for ${each.value}"
  format        = "DOCKER"
}

# Cloud Pub/Sub (Kafka alternative)
resource "google_pubsub_topic" "events" {
  name = "${var.project_name}-events"
}

resource "google_pubsub_subscription" "events" {
  name  = "${var.project_name}-events-sub"
  topic = google_pubsub_topic.events.name

  ack_deadline_seconds = 20
  retain_acked_messages = false
}

# Cloud Logging Sink
resource "google_logging_project_sink" "audit_logs" {
  name        = "${var.project_name}-audit-logs"
  destination = "pubsub.googleapis.com/projects/${var.gcp_project_id}/topics/${google_pubsub_topic.events.name}"

  filter = "resource.type=\"k8s_cluster\" AND severity>=ERROR"

  unique_writer_identity = true
}

# KMS Keys
resource "google_kms_key_ring" "main" {
  name     = "${var.project_name}-keyring"
  location = var.gcp_region
}

resource "google_kms_crypto_key" "gke" {
  name            = "gke-key"
  key_ring        = google_kms_key_ring.main.id
  rotation_period = "7776000s" # 90 days

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_kms_crypto_key" "storage" {
  name            = "storage-key"
  key_ring        = google_kms_key_ring.main.id
  rotation_period = "7776000s"
}

# Service Networking Connection
resource "google_compute_global_address" "private_ip" {
  name          = "${var.project_name}-private-ip"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.main.id
}

resource "google_service_networking_connection" "private_vpc" {
  network                 = google_compute_network.main.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip.name]
}

