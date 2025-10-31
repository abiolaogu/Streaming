variable "project_name" {
  description = "Project name"
  type        = string
  default     = "streaming-platform"
}

variable "environment" {
  description = "Environment (dev, stg, prd)"
  type        = string
}

variable "azure_region" {
  description = "Azure region"
  type        = string
  default     = "eastus"
}

variable "vnet_address_space" {
  description = "Address space for VNet"
  type        = string
  default     = "10.2.0.0/16"
}

variable "public_subnets" {
  description = "Public subnet CIDR blocks"
  type        = list(string)
  default     = ["10.2.101.0/24", "10.2.102.0/24", "10.2.103.0/24"]
}

variable "private_subnets" {
  description = "Private subnet CIDR blocks"
  type        = list(string)
  default     = ["10.2.1.0/24", "10.2.2.0/24", "10.2.3.0/24"]
}

variable "kubernetes_version" {
  description = "Kubernetes version"
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

variable "gpu_spot_max" {
  description = "Maximum number of Spot GPU nodes"
  type        = number
  default     = 10
}

