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

variable "azure_region" {
  description = "Azure region"
  type        = string
  default     = "eastus"
}

variable "azure_secondary_region" {
  description = "Azure secondary region for geo-replication"
  type        = string
  default     = "westus"
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
  default     = ["1", "2", "3"]
}

variable "k8s_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.28"
}

variable "node_vm_size" {
  description = "AKS node VM size"
  type        = string
  default     = "Standard_D2s_v3"
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

variable "admin_group_object_ids" {
  description = "Azure AD admin group object IDs"
  type        = list(string)
  default     = []
}

variable "db_sku_name" {
  description = "PostgreSQL SKU name"
  type        = string
  default     = "GP_Standard_D2s_v3"
}

variable "db_storage_mb" {
  description = "PostgreSQL storage (MB)"
  type        = number
  default     = 32768
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

variable "redis_capacity" {
  description = "Redis cache capacity"
  type        = number
  default     = 1
}

variable "redis_family" {
  description = "Redis cache family"
  type        = string
  default     = "C"
}

variable "redis_sku_name" {
  description = "Redis SKU name"
  type        = string
  default     = "Standard"
}

variable "service_bus_sku" {
  description = "Service Bus SKU"
  type        = string
  default     = "Standard"
}

