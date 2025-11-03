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

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b", "us-east-1c"]
}

variable "allowed_cidr_blocks" {
  description = "Allowed CIDR blocks for EKS API"
  type        = list(string)
  default     = ["0.0.0.0/0"] # Restrict in production
}

variable "k8s_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.28"
}

variable "node_instance_types" {
  description = "EKS node instance types"
  type        = list(string)
  default     = ["t3.medium"]
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

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.medium"
}

variable "db_allocated_storage" {
  description = "RDS allocated storage (GB)"
  type        = number
  default     = 100
}

variable "db_max_storage" {
  description = "RDS max storage (GB)"
  type        = number
  default     = 500
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

variable "redis_node_type" {
  description = "ElastiCache node type"
  type        = string
  default     = "cache.t3.medium"
}

variable "kafka_instance_type" {
  description = "MSK broker instance type"
  type        = string
  default     = "kafka.m5.large"
}

variable "kafka_volume_size" {
  description = "MSK broker volume size (GB)"
  type        = number
  default     = 100
}

variable "services" {
  description = "List of services for ECR repositories"
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

