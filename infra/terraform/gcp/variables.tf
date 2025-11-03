variable "project_name" {
  description = "Project name"
  type        = string
  default     = "streaming-platform"
}

variable "environment" {
  description = "Environment (dev, stg, prd)"
  type        = string
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

variable "public_subnets" {
  description = "Public subnet CIDR blocks"
  type        = list(string)
  default     = ["10.1.101.0/24", "10.1.102.0/24", "10.1.103.0/24"]
}

variable "private_subnets" {
  description = "Private subnet CIDR blocks"
  type        = list(string)
  default     = ["10.1.1.0/24", "10.1.2.0/24", "10.1.3.0/24"]
}

variable "master_ipv4_cidr_block" {
  description = "CIDR for GKE master"
  type        = string
  default     = "172.16.0.0/28"
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

