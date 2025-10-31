terraform {
  required_version = ">= 1.5.0"
  
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.22.0"
    }
  }
  
  backend "azurerm" {
    resource_group_name  = "${var.project_name}-terraform-state"
    storage_account_name = "${var.project_name}tfstate${var.environment}"
    container_name       = "tfstate"
    key                  = "azure/${var.environment}/terraform.tfstate"
  }
}

provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

# Resource Group
resource "azurerm_resource_group" "main" {
  name     = "${var.project_name}-${var.environment}-rg"
  location = var.azure_region
  
  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

# Virtual Network
resource "azurerm_virtual_network" "main" {
  name                = "${var.project_name}-${var.environment}-vnet"
  address_space       = [var.vnet_address_space]
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
}

# Subnets
resource "azurerm_subnet" "public" {
  count = length(var.public_subnets)
  
  name                 = "${var.project_name}-public-subnet-${count.index + 1}"
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = [var.public_subnets[count.index]]
}

resource "azurerm_subnet" "private" {
  count = length(var.private_subnets)
  
  name                 = "${var.project_name}-private-subnet-${count.index + 1}"
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = [var.private_subnets[count.index]]
}

# AKS Cluster
resource "azurerm_kubernetes_cluster" "main" {
  name                = "${var.project_name}-${var.environment}"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  dns_prefix          = "${var.project_name}-${var.environment}"
  kubernetes_version  = var.kubernetes_version
  
  node_resource_group = "${var.project_name}-${var.environment}-nodes-rg"
  
  default_node_pool {
    name                = "nodepool1"
    node_count          = var.cpu_ondemand_floor.count
    vm_size             = "Standard_B2s"
    type                = "VirtualMachineScaleSets"
    enable_auto_scaling = false
    os_disk_size_gb     = 50
    zones                = [1, 2, 3]
  }
  
  # Identity
  identity {
    type = "SystemAssigned"
  }
  
  # Network
  network_profile {
    network_plugin     = "azure"
    network_policy     = "calico"
    service_cidr       = "10.0.100.0/24"
    dns_service_ip     = "10.0.100.10"
    docker_bridge_cidr = "172.17.0.1/16"
  }
  
  # RBAC
  role_based_access_control_enabled = true
  
  # Azure AD integration
  azure_active_directory_role_based_access_control {
    managed                = true
    azure_rbac_enabled     = true
    admin_group_object_ids = []
  }
  
  # Add-ons
  oms_agent {
    enabled                    = true
    log_analytics_workspace_id = azurerm_log_analytics_workspace.main.id
  }
  
  # Security
  auto_scaler_profile {
    scale_down_delay_after_add    = "5m"
    scan_interval                 = "30s"
    scale_down_unneeded           = "10m"
    scale_down_utilization_threshold = 0.5
  }
}

# Log Analytics
resource "azurerm_log_analytics_workspace" "main" {
  name                = "${var.project_name}-logs-${var.environment}"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

# Key Vault
resource "azurerm_key_vault" "main" {
  name                = "${var.project_name}-kv-${var.environment}-${random_id.kv_suffix.hex}"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
  
  purge_protection_enabled = var.environment == "prd"
  soft_delete_retention_days = 7
}

resource "random_id" "kv_suffix" {
  byte_length = 4
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault_access_policy" "cluster" {
  key_vault_id = azurerm_key_vault.main.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_kubernetes_cluster.main.identity[0].principal_id
  
  secret_permissions = [
    "Get",
    "List"
  ]
}

resource "azurerm_key_vault_key" "main" {
  name         = "${var.project_name}-encryption-key"
  key_vault_id = azurerm_key_vault.main.id
  key_type     = "RSA"
  key_size     = 2048
  
  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey"
  ]
  
  rotation_policy {
    automatic {
      time_before_expiry = "P30D"
    }
    
    expire_after         = "P365D"
    notify_before_expiry = "P30D"
  }
}

# Container Registry
resource "azurerm_container_registry" "main" {
  name                = "${var.project_name}${var.environment}${random_id.acr_suffix.hex}"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  sku                 = "Standard"
  admin_enabled       = false
  
  network_rule_set {
    default_action = "Deny"
    
    ip_rule {
      action   = "Allow"
      ip_range = azurerm_virtual_network.main.address_space[0]
    }
  }
}

resource "random_id" "acr_suffix" {
  byte_length = 4
}

# Storage Account for MinIO
resource "azurerm_storage_account" "minio" {
  name                     = "${var.project_name}minio${var.environment}${random_id.storage_suffix.hex}"
  resource_group_name      = azurerm_resource_group.main.name
  location                 = azurerm_resource_group.main.location
  account_tier             = "Standard"
  account_replication_type = "ZRS"
  account_kind             = "StorageV2"
  
  network_rules {
    default_action             = "Deny"
    virtual_network_subnet_ids = azurerm_subnet.private[*].id
  }
  
  blob_properties {
    versioning_enabled = true
    
    delete_retention_policy {
      days = 7
    }
  }
  
  tags = {
    Environment = var.environment
  }
}

resource "random_id" "storage_suffix" {
  byte_length = 4
}

resource "azurerm_storage_container" "minio" {
  name                  = "data"
  storage_account_name  = azurerm_storage_account.minio.name
  container_access_type = "private"
}

# Spot CPU Node Pool
resource "azurerm_kubernetes_cluster_node_pool" "spot_cpu" {
  name                  = "spotcpu"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.main.id
  vm_size               = "Standard_D2s_v3"
  node_count            = 0
  
  priority        = "Spot"
  eviction_policy = "Delete"
  spot_max_price  = 0.10
  
  enable_auto_scaling = true
  min_count           = 0
  max_count           = var.cpu_spot_max
  
  os_disk_size_gb = 50
  
  node_labels = {
    workload-type = "general"
    spot-enabled  = "true"
  }
  
  node_taints = [
    "spot=true:NoSchedule"
  ]
  
  zones = [1, 2, 3]
}

# Outputs
output "cluster_id" {
  value = azurerm_kubernetes_cluster.main.id
}

output "cluster_fqdn" {
  value = azurerm_kubernetes_cluster.main.fqdn
}

output "cluster_host" {
  value     = azurerm_kubernetes_cluster.main.kube_config[0].host
  sensitive = true
}

output "cluster_client_certificate" {
  value     = azurerm_kubernetes_cluster.main.kube_config[0].client_certificate
  sensitive = true
}

output "cluster_client_key" {
  value     = azurerm_kubernetes_cluster.main.kube_config[0].client_key
  sensitive = true
}

output "cluster_cluster_ca_certificate" {
  value     = azurerm_kubernetes_cluster.main.kube_config[0].cluster_ca_certificate
  sensitive = true
}

output "key_vault_name" {
  value = azurerm_key_vault.main.name
}

output "container_registry_name" {
  value = azurerm_container_registry.main.name
}

output "storage_account_name" {
  value = azurerm_storage_account.minio.name
}

