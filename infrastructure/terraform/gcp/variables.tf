variable "project_name" {
  description = "Project name"
  type        = string
  default     = "streamverse"
}

variable "environment" {
  description = "Environment (dev/staging/production)"
  type        = string
  default     = "dev"
}

variable "gcp_project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "gcp_region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

variable "gke_node_zones" {
  description = "GKE node zones"
  type        = list(string)
  default     = ["us-central1-a", "us-central1-b", "us-central1-c"]
}

variable "node_machine_type" {
  description = "GKE node machine type"
  type        = string
  default     = "e2-medium"
}

variable "node_desired_size" {
  description = "Desired number of nodes"
  type        = number
  default     = 3
}

variable "node_min_size" {
  description = "Minimum number of nodes"
  type        = number
  default     = 1
}

variable "node_max_size" {
  description = "Maximum number of nodes"
  type        = number
  default     = 10
}

variable "db_instance_tier" {
  description = "Cloud SQL instance tier"
  type        = string
  default     = "db-f1-micro"
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "streamverse"
}

variable "db_username" {
  description = "Database username"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "redis_tier" {
  description = "Memorystore tier"
  type        = string
  default     = "STANDARD_HA"
}

variable "redis_memory_size" {
  description = "Memorystore memory size (GB)"
  type        = number
  default     = 1
}

variable "services" {
  description = "List of services for Artifact Registry"
  type        = list(string)
  default = [
    "auth-service",
    "user-service",
    "content-service",
    "streaming-service",
    "transcoding-service",
    "payment-service",
    "search-service",
    "analytics-service",
    "recommendation-service",
    "notification-service",
    "admin-service",
    "scheduler-service",
  ]
}

