variable "project_name" {
  description = "Project name"
  type        = string
  default     = "streaming-platform"
}

variable "environment" {
  description = "Environment (dev, stg, prd)"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnets" {
  description = "Public subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
}

variable "private_subnets" {
  description = "Private subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "kubernetes_version" {
  description = "Kubernetes version for EKS"
  type        = string
  default     = "1.28"
}

variable "cpu_ondemand_floor" {
  description = "On-demand CPU node floor configuration"
  type = object({
    count = number
  })
  default = {
    count = 1
  }
}

variable "cpu_spot_max" {
  description = "Maximum number of Spot CPU nodes"
  type        = number
  default     = 20
}

variable "ecr_repositories" {
  description = "List of ECR repositories to create"
  type        = list(string)
  default = [
    "media/ome",
    "media/gstreamer",
    "media/drm-proxy",
    "media/fast-scheduler",
    "cdn/ats",
    "cdn/varnish",
    "cdn/atc",
    "data/dragonfly",
    "data/scylla",
    "data/clickhouse",
    "telecom/kamailio",
    "telecom/freeswitch",
    "telecom/open5gs",
    "control/autoscaler",
    "control/config-pusher"
  ]
}

variable "enable_cross_cloud_replication" {
  description = "Enable S3 cross-cloud replication"
  type        = bool
  default     = false
}

variable "cross_cloud_bucket_name" {
  description = "Cross-cloud replication target bucket"
  type        = string
  default     = ""
}

