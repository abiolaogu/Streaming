# Apache Traffic Control/Server CDN Infrastructure

terraform {
  required_version = ">= 1.5.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Traffic Server L1 Edge Cache nodes
resource "aws_instance" "traffic_server_edge" {
  count         = var.edge_count
  ami           = var.ami_id
  instance_type = var.instance_type
  subnet_id     = var.subnet_ids[count.index % length(var.subnet_ids)]

  tags = {
    Name = "traffic-server-edge-${count.index + 1}"
    Role = "cdn-edge"
  }
}

# Traffic Control components
resource "aws_instance" "traffic_ops" {
  ami           = var.ami_id
  instance_type = "t3.medium"
  
  tags = {
    Name = "traffic-ops"
    Role = "cdn-control"
  }
}

resource "aws_instance" "traffic_router" {
  ami           = var.ami_id
  instance_type = "t3.small"
  
  tags = {
    Name = "traffic-router"
    Role = "cdn-router"
  }
}

