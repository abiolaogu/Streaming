terraform {
  required_version = ">= 1.5.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.22.0"
    }
  }

  backend "s3" {
    bucket = "${var.project_name}-terraform-state-${var.environment}"
    key    = "aws/${var.environment}/terraform.tfstate"
    region = var.aws_region
  }
}

provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

# VPC
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"
  
  name = "${var.project_name}-${var.environment}"
  cidr = var.vpc_cidr
  
  azs                     = slice(data.aws_availability_zones.available.names, 0, 3)
  public_subnets          = var.public_subnets
  private_subnets         = var.private_subnets
  enable_nat_gateway      = true
  single_nat_gateway      = false
  enable_dns_hostnames    = true
  enable_dns_support      = true
  
  public_subnet_tags = {
    Type = "public"
  }
  
  private_subnet_tags = {
    Type = "private"
  }
}

# EKS Cluster
module "eks" {
  source = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"
  
  cluster_name    = "${var.project_name}-${var.environment}"
  cluster_version = var.kubernetes_version
  
  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets
  
  # IRSA for pod identity
  enable_irsa = true
  
  # Cluster logging
  cluster_enabled_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]
  
  cluster_addons = {
    aws-ebs-csi-driver = {
      most_recent = true
    }
    coredns = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
  }
  
  # Self-managed node groups for on-demand baseline
  self_managed_node_groups = {
    cpu_ondemand_floor = {
      desired_size = var.cpu_ondemand_floor.count
      min_size     = var.cpu_ondemand_floor.count
      max_size     = var.cpu_ondemand_floor.count
      instance_types = ["t3.medium"]
      disk_size      = 50
      
      launch_template_name = "${var.project_name}-ondemand-cpu"
      
      labels = {
        workload-type = "general"
        spot-enabled   = "false"
      }
      
      taints = []
      
      pre_bootstrap_user_data = <<-EOT
        #!/bin/bash
        /etc/eks/bootstrap.sh ${var.project_name}-${var.environment}
      EOT
    }
  }
  
  # Managed node groups with Spot
  eks_managed_node_groups = {
    cpu_spot = {
      min_size     = 0
      max_size     = var.cpu_spot_max
      desired_size = 0
      instance_types = ["t3.medium", "t3a.medium", "t3.large"]
      capacity_type = "SPOT"
      
      labels = {
        workload-type = "general"
        spot-enabled   = "true"
      }
      
      taints = []
      
      update_config = {
        max_unavailable_percentage = 50
      }
    }
    
    gpu_spot = {
      min_size     = 0
      max_size     = var.gpu_spot_max
      desired_size = 0
      instance_types = ["g4dn.xlarge", "g5.xlarge"]
      capacity_type = "SPOT"
      
      labels = {
        workload-type = "gpu"
        spot-enabled   = "true"
        nvidia.com/gpu = "true"
      }
      
      taints = [
        {
          key    = "nvidia.com/gpu"
          value  = "true"
          effect = "NO_SCHEDULE"
        }
      ]
      
      update_config = {
        max_unavailable_percentage = 50
      }
    }
    
    gpu_ondemand_fallback = {
      min_size     = 0
      max_size     = 1
      desired_size = 0
      instance_types = ["g5.xlarge"]
      capacity_type = "ON_DEMAND"
      
      labels = {
        workload-type = "gpu"
        spot-enabled   = "false"
        nvidia.com/gpu = "true"
      }
      
      taints = [
        {
          key    = "nvidia.com/gpu"
          value  = "true"
          effect = "NO_SCHEDULE"
        }
      ]
      
      update_config = {
        max_unavailable_percentage = 50
      }
    }
  }
  
  node_security_group_additional_rules = {
    ingress_self_all = {
      description = "Node to node all ports/protocols"
      protocol    = "-1"
      from_port   = 0
      to_port     = 0
      type        = "ingress"
      self        = true
    }
  }
}

# KMS for encryption
resource "aws_kms_key" "eks" {
  description             = "EKS encryption key for ${var.environment}"
  deletion_window_in_days = 10
  enable_key_rotation     = true
  
  tags = {
    Name = "${var.project_name}-eks-key-${var.environment}"
  }
}

resource "aws_kms_alias" "eks" {
  name          = "alias/${var.project_name}-eks-${var.environment}"
  target_key_id = aws_kms_key.eks.key_id
}

# ECR repositories
resource "aws_ecr_repository" "platform" {
  for_each = toset(var.ecr_repositories)
  
  name                 = "${var.project_name}/${each.value}"
  image_tag_mutability = "MUTABLE"
  
  image_scanning_configuration {
    scan_on_push = true
  }
  
  encryption_configuration {
    encryption_type = "KMS"
    kms_key         = aws_kms_key.eks.arn
  }
}

resource "aws_ecr_lifecycle_policy" "platform" {
  for_each = aws_ecr_repository.platform
  
  repository = each.value.name
  
  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 30 images"
        selection = {
          tagStatus     = "any"
          countType     = "imageCountMoreThan"
          countNumber   = 30
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# S3 buckets for MinIO
resource "aws_s3_bucket" "minio" {
  bucket = "${var.project_name}-minio-${var.environment}-${random_id.bucket_suffix.hex}"
  
  tags = {
    Name        = "MinIO Origin Storage"
    Environment = var.environment
  }
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

resource "aws_s3_bucket_versioning" "minio" {
  bucket = aws_s3_bucket.minio.id
  
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "minio" {
  bucket = aws_s3_bucket.minio.id
  
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.eks.arn
    }
  }
}

resource "aws_s3_bucket_replication_configuration" "minio" {
  count = var.enable_cross_cloud_replication ? 1 : 0
  
  role   = aws_iam_role.replication[0].arn
  bucket = aws_s3_bucket.minio.id
  
  rule {
    id     = "replicate-to-other-clouds"
    status = "Enabled"
    
    destination {
      bucket = "arn:aws:s3:::${var.cross_cloud_bucket_name}"
      
      replication_time {
        status = "Enabled"
        time {
          minutes = 15
        }
      }
      
      metrics {
        status = "Enabled"
        event_threshold {
          minutes = 15
        }
      }
    }
    
    filter {}
  }
  
  depends_on = [aws_s3_bucket_versioning.minio]
}

resource "aws_iam_role" "replication" {
  count = var.enable_cross_cloud_replication ? 1 : 0
  
  name = "${var.project_name}-s3-replication-${var.environment}"
  
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "s3.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "replication" {
  count = var.enable_cross_cloud_replication ? 1 : 0
  
  role = aws_iam_role.replication[0].id
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetReplicationConfiguration",
          "s3:ListBucket"
        ]
        Resource = [
          aws_s3_bucket.minio.arn
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObjectVersionForReplication",
          "s3:GetObjectVersionAcl",
          "s3:GetObjectVersionTagging"
        ]
        Resource = [
          "${aws_s3_bucket.minio.arn}/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "s3:ReplicateObject",
          "s3:ReplicateDelete",
          "s3:ReplicateTags"
        ]
        Resource = [
          "${aws_s3_bucket.minio.arn}/*"
        ]
        Condition = {
          StringLikeIfExists = {
            "s3:x-amz-server-side-encryption" = ["aws:kms", "AES256"]
          }
        }
      }
    ]
  })
}

# CloudWatch for logging
resource "aws_cloudwatch_log_group" "eks" {
  name              = "/aws/eks/${var.project_name}-${var.environment}/cluster"
  retention_in_days = 30
  
  kms_key_id = aws_kms_key.eks.arn
}

# Outputs
output "cluster_id" {
  value = module.eks.cluster_id
}

output "cluster_endpoint" {
  value = module.eks.cluster_endpoint
}

output "cluster_ca_certificate" {
  value     = module.eks.cluster_certificate_authority_data
  sensitive = true
}

output "vpc_id" {
  value = module.vpc.vpc_id
}

output "private_subnets" {
  value = module.vpc.private_subnets
}

output "public_subnets" {
  value = module.vpc.public_subnets
}

output "minio_bucket" {
  value = aws_s3_bucket.minio.bucket
}

output "ecr_repositories" {
  value = { for k, v in aws_ecr_repository.platform : k => v.repository_url }
}

output "kms_key_arn" {
  value = aws_kms_key.eks.arn
}

