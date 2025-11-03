terraform {
  required_version = ">= 1.5.0"
  
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }

  backend "azurerm" {
    resource_group_name  = "streamverse-terraform-state"
    storage_account_name = "streamversetfstate"
    container_name       = "terraform-state"
    key                  = "azure/terraform.tfstate"
  }
}

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = var.environment == "production"
    }
  }
}

# Resource Group
resource "azurerm_resource_group" "main" {
  name     = "${var.project_name}-rg"
  location = var.azure_region

  tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
  }
}

# Virtual Network
resource "azurerm_virtual_network" "main" {
  name                = "${var.project_name}-vnet"
  address_space       = [var.vpc_cidr]
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name

  tags = {
    Name = "${var.project_name}-vnet"
  }
}

# Subnets
resource "azurerm_subnet" "private" {
  count = length(var.availability_zones)

  name                 = "${var.project_name}-private-${count.index + 1}"
  resource_group_name  = azurerm_resource_group.main.name
  virtual_network_name = azurerm_virtual_network.main.name
  address_prefixes     = [cidrsubnet(var.vpc_cidr, 8, count.index)]
}

# AKS Cluster
resource "azurerm_kubernetes_cluster" "main" {
  name                = "${var.project_name}-aks"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  dns_prefix          = var.project_name
  kubernetes_version  = var.k8s_version

  default_node_pool {
    name                = "default"
    node_count          = var.node_desired_size
    vm_size             = var.node_vm_size
    enable_auto_scaling = true
    min_count           = var.node_min_size
    max_count           = var.node_max_size
    vnet_subnet_id      = azurerm_subnet.private[0].id
    os_disk_size_gb     = 50
    os_disk_type        = "Ephemeral"
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin     = "azure"
    network_policy     = "calico"
    service_cidr       = "10.1.0.0/16"
    dns_service_ip     = "10.1.0.10"
    docker_bridge_cidr = "172.17.0.1/16"
  }

  role_based_access_control {
    enabled = true
    azure_active_directory {
      managed                = true
      admin_group_object_ids  = var.admin_group_object_ids
      azure_rbac_enabled      = true
    }
  }

  azure_policy_enabled = true

  oms_agent {
    enabled                    = true
    log_analytics_workspace_id = azurerm_log_analytics_workspace.main.id
  }

  key_vault_secrets_provider {
    secret_rotation_enabled  = true
    secret_rotation_interval = "2m"
  }

  tags = {
    Name = "${var.project_name}-aks"
  }
}

# PostgreSQL Flexible Server
resource "azurerm_postgresql_flexible_server" "main" {
  name                   = "${var.project_name}-db"
  resource_group_name    = azurerm_resource_group.main.name
  location               = azurerm_resource_group.main.location
  version                = "15"
  delegated_subnet_id    = azurerm_subnet.private[0].id
  private_dns_zone_id    = azurerm_private_dns_zone.postgres.id
  administrator_login    = var.db_username
  administrator_password = var.db_password

  backup_retention_days        = 7
  geo_redundant_backup_enabled = var.environment == "production"

  sku_name   = var.db_sku_name
  storage_mb = var.db_storage_mb

  maintenance_window {
    day_of_week  = 0
    start_hour   = 3
    start_minute = 0
  }

  high_availability {
    mode = var.environment == "production" ? "ZoneRedundant" : "Disabled"
  }

  depends_on = [azurerm_private_dns_zone_virtual_network_link.postgres]
}

resource "azurerm_postgresql_flexible_server_database" "main" {
  name      = var.db_name
  server_id = azurerm_postgresql_flexible_server.main.id
  charset   = "UTF8"
  collation = "en_US.utf8"
}

# Private DNS Zone for PostgreSQL
resource "azurerm_private_dns_zone" "postgres" {
  name                = "${var.project_name}.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.main.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "postgres" {
  name                  = "${var.project_name}-postgres-link"
  resource_group_name   = azurerm_resource_group.main.name
  private_dns_zone_name = azurerm_private_dns_zone.postgres.name
  virtual_network_id    = azurerm_virtual_network.main.id
}

# Azure Cache for Redis
resource "azurerm_redis_cache" "main" {
  name                = "${var.project_name}-redis"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  capacity            = var.redis_capacity
  family              = var.redis_family
  sku_name            = var.redis_sku_name
  enable_non_ssl_port = false
  minimum_tls_version = "1.2"

  subnet_id = azurerm_subnet.private[0].id

  redis_configuration {
    maxmemory_reserved = 2
    maxmemory_delta   = 2
    maxmemory_policy  = "allkeys-lru"
  }
}

# Azure Service Bus (Kafka alternative)
resource "azurerm_servicebus_namespace" "main" {
  name                = "${var.project_name}-sb"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  sku                 = var.service_bus_sku

  tags = {
    Name = "${var.project_name}-servicebus"
  }
}

resource "azurerm_servicebus_topic" "events" {
  name         = "${var.project_name}-events"
  namespace_id = azurerm_servicebus_namespace.main.id
}

resource "azurerm_servicebus_subscription" "events" {
  name     = "${var.project_name}-events-sub"
  topic_id = azurerm_servicebus_topic.events.id
}

# Storage Account (Media files)
resource "azurerm_storage_account" "media" {
  name                     = "${var.project_name}media${var.environment}"
  resource_group_name      = azurerm_resource_group.main.name
  location                 = azurerm_resource_group.main.location
  account_tier             = "Standard"
  account_replication_type = var.environment == "production" ? "GRS" : "LRS"
  account_kind              = "StorageV2"

  blob_properties {
    versioning_enabled       = true
    change_feed_enabled      = true
    delete_retention_policy {
      days = 90
    }
    container_delete_retention_policy {
      days = 30
    }
  }

  tags = {
    Name = "${var.project_name}-media"
  }
}

resource "azurerm_storage_container" "media" {
  name                  = "media"
  storage_account_name  = azurerm_storage_account.media.name
  container_access_type = "private"
}

# Container Registry
resource "azurerm_container_registry" "main" {
  name                = "${var.project_name}acr"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  sku                 = var.environment == "production" ? "Premium" : "Standard"
  admin_enabled       = false

  network_rule_set {
    default_action = "Deny"
    virtual_network {
      subnet_id = azurerm_subnet.private[0].id
    }
  }

  georeplications {
    location = var.azure_secondary_region
    tags     = {}
  }

  tags = {
    Name = "${var.project_name}-acr"
  }
}

# Log Analytics Workspace
resource "azurerm_log_analytics_workspace" "main" {
  name                = "${var.project_name}-logs"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  sku                 = "PerGB2018"
  retention_in_days   = 30

  tags = {
    Name = "${var.project_name}-logs"
  }
}

