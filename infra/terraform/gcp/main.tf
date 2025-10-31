terraform {
  required_version = ">= 1.5.0"
  
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.22.0"
    }
  }
  
  backend "gcs" {
    bucket = "${var.project_name}-terraform-state-${var.environment}"
    prefix = "gcp/${var.environment}"
  }
}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

# VPC Network
resource "google_compute_network" "main" {
  name                    = "${var.project_name}-${var.environment}-vpc"
  auto_create_subnetworks = false
  
  depends_on = [google_project_service.compute]
}

# Subnets
resource "google_compute_subnetwork" "public" {
  count = length(var.public_subnets)
  
  name          = "${var.project_name}-${var.environment}-public-${count.index + 1}"
  ip_cidr_range = var.public_subnets[count.index]
  region        = var.gcp_region
  network       = google_compute_network.main.id
  
  private_ip_google_access = false
}

resource "google_compute_subnetwork" "private" {
  count = length(var.private_subnets)
  
  name          = "${var.project_name}-${var.environment}-private-${count.index + 1}"
  ip_cidr_range = var.private_subnets[count.index]
  region        = var.gcp_region
  network       = google_compute_network.main.id
  
  private_ip_google_access = true
}

# GKE Cluster
resource "google_container_cluster" "main" {
  name     = "${var.project_name}-${var.environment}"
  location = var.gcp_region
  
  # Network
  network    = google_compute_network.main.name
  subnetwork = google_compute_subnetwork.private[0].name
  
  # Default node pool for system workloads
  remove_default_node_pool = true
  initial_node_count       = var.cpu_ondemand_floor.count
  
  # Enable features
  enable_shielded_nodes          = true
  enable_binary_authorization    = true
  enable_autopilot               = false
  
  # Logging
  logging_service    = "logging.googleapis.com/kubernetes"
  monitoring_service = "monitoring.googleapis.com/kubernetes"
  
  # Private cluster
  private_cluster_config {
    enable_private_endpoint = false
    enable_private_nodes    = true
    master_ipv4_cidr_block  = var.master_ipv4_cidr_block
  }
  
  # Master authorized networks
  master_authorized_networks_config {
    cidr_blocks {
      cidr_block = var.private_subnets[0]
    }
  }
  
  # Release channel
  release_channel {
    channel = "STABLE"
  }
  
  # Workload Identity
  workload_identity_config {
    workload_pool = "${var.gcp_project_id}.svc.id.goog"
  }
  
  # Network policy
  network_policy {
    enabled = true
  }
  
  # Vertical Pod Autoscaler
  vertical_pod_autoscaling {
    enabled = true
  }
  
  # Addons
  addons_config {
    horizontal_pod_autoscaling {
      disabled = false
    }
    
    http_load_balancing {
      disabled = false
    }
    
    network_policy_config {
      disabled = false
    }
    
    gce_persistent_disk_csi_driver_config {
      enabled = true
    }
  }
  
  depends_on = [
    google_project_service.container,
    google_project_service.compute
  ]
}

# On-demand node pool
resource "google_container_node_pool" "ondemand_cpu" {
  name       = "ondemand-cpu-pool"
  cluster    = google_container_cluster.main.name
  location   = var.gcp_region
  node_count = var.cpu_ondemand_floor.count
  
  management {
    auto_repair  = true
    auto_upgrade = true
  }
  
  node_config {
    machine_type = "e2-medium"
    disk_size_gb = 50
    disk_type    = "pd-standard"
    
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
    
    labels = {
      workload-type = "general"
      preemptible   = "false"
    }
    
    workload_metadata_config {
      mode = "GKE_METADATA"
    }
  }
}

# Spot CPU node pool
resource "google_container_node_pool" "spot_cpu" {
  name     = "spot-cpu-pool"
  cluster  = google_container_cluster.main.name
  location = var.gcp_region
  
  autoscaling {
    min_node_count = 0
    max_node_count = var.cpu_spot_max
  }
  
  management {
    auto_repair  = true
    auto_upgrade = true
  }
  
  node_config {
    machine_type = "e2-medium"
    disk_size_gb = 50
    disk_type    = "pd-standard"
    preemptible  = true
    
    spot = true
    
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
    
    labels = {
      workload-type = "general"
      preemptible   = "true"
    }
    
    taint {
      key    = "preemptible"
      value  = "true"
      effect = "NO_SCHEDULE"
    }
    
    workload_metadata_config {
      mode = "GKE_METADATA"
    }
  }
}

# Cloud KMS
resource "google_kms_key_ring" "main" {
  name     = "${var.project_name}-${var.environment}"
  location = var.gcp_region
}

resource "google_kms_crypto_key" "main" {
  name            = "${var.project_name}-encryption"
  key_ring        = google_kms_key_ring.main.id
  rotation_period = "2592000s" # 30 days
  
  version_template {
    algorithm = "GOOGLE_SYMMETRIC_ENCRYPTION"
  }
}

# Artifact Registry
resource "google_artifact_registry_repository" "main" {
  location      = var.gcp_region
  repository_id = "${var.project_name}-${var.environment}"
  description   = "Docker repository for ${var.project_name}"
  format        = "DOCKER"
}

# Cloud Storage for MinIO
resource "google_storage_bucket" "minio" {
  name          = "${var.project_name}-minio-${var.environment}-${random_id.bucket_suffix.hex}"
  location      = var.gcp_region
  force_destroy = false
  
  versioning {
    enabled = true
  }
  
  encryption {
    default_kms_key_name = google_kms_crypto_key.main.id
  }
  
  uniform_bucket_level_access = true
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# API enabling
resource "google_project_service" "compute" {
  service = "compute.googleapis.com"
}

resource "google_project_service" "container" {
  service = "container.googleapis.com"
}

resource "google_project_service" "cloudkms" {
  service = "cloudkms.googleapis.com"
}

resource "google_project_service" "artifactregistry" {
  service = "artifactregistry.googleapis.com"
}

resource "google_project_service" "storage" {
  service = "storage.googleapis.com"
}

# Outputs
output "cluster_name" {
  value = google_container_cluster.main.name
}

output "cluster_endpoint" {
  value     = google_container_cluster.main.endpoint
  sensitive = true
}

output "cluster_ca_certificate" {
  value     = google_container_cluster.main.master_auth[0].cluster_ca_certificate
  sensitive = true
}

output "network" {
  value = google_compute_network.main.name
}

output "subnetworks" {
  value = google_compute_subnetwork.private[*].name
}

output "artifact_registry" {
  value = google_artifact_registry_repository.main.name
}

output "minio_bucket" {
  value = google_storage_bucket.minio.name
}

output "kms_key_ring" {
  value = google_kms_key_ring.main.name
}

