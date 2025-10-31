# Local GPU Baseline for Tier-1 PoPs
# Deploy 1-2 GPU servers per region for ultra-low-latency media AI

variable "gpu_baseline_per_region" {
  description = "Number of GPU baseline servers per region"
  type        = number
  default     = 1
}

variable "gpu_instance_flavor" {
  description = "OpenStack flavor for GPU instances"
  type        = string
  default     = "gpu.a100.2x"
}

resource "openstack_compute_instance_v2" "local_gpu_baseline" {
  count     = var.gpu_baseline_per_region
  name      = "local-gpu-baseline-${format("%02d", count.index + 1)}"
  region    = var.openstack_region
  
  flavor_name = var.gpu_instance_flavor
  
  image_name = "ubuntu-22.04-nvidia-docker"
  
  key_pair = var.openstack_key_pair
  
  security_groups = [
    openstack_compute_secgroup_v2.gpu_sg.name
  ]
  
  network {
    uuid = var.private_network_id
  }
  
  user_data = <<-EOF
    #!/bin/bash
    set -e
    
    # Install NVIDIA drivers
    apt-get update
    apt-get install -y nvidia-driver-535 nvidia-container-toolkit
    
    # Configure Docker for GPU
    cat > /etc/docker/daemon.json <<'DOCKEREOF'
    {
      "default-runtime": "nvidia",
      "runtimes": {
        "nvidia": {
          "path": "nvidia-container-runtime",
          "runtimeArgs": []
        }
      }
    }
    DOCKEREOF
    
    systemctl restart docker
    
    # Install Kubernetes worker node components
    curl -fsSL https://get.k8s.io | bash -
    
    # Join cluster
    kubeadm join ${var.kubernetes_master} \
      --token ${var.bootstrap_token} \
      --discovery-token-unsafe-skip-ca-verification
    
    # Label as local GPU baseline
    kubectl label node $(hostname) \
      workload-type=gpu \
      gpu-type=local-baseline \
      preemptible=false \
      nvidia.com/gpu=true \
      --overwrite
    
    # Taint for dedicated use
    kubectl taint node $(hostname) \
      gpu-type=local-baseline:NoSchedule \
      --overwrite
  EOF
  
  lifecycle {
    ignore_changes = [user_data]
  }
  
  tags = {
    Name        = "local-gpu-baseline"
    Environment = var.environment
    Region      = var.openstack_region
    Tier        = "1"
  }
}

resource "openstack_compute_secgroup_v2" "gpu_sg" {
  name        = "local-gpu-baseline-sg"
  description = "Security group for local GPU baseline nodes"
  
  rule {
    ip_protocol = "tcp"
    from_port   = 22
    to_port     = 22
    cidr        = "10.0.0.0/8"
  }
  
  rule {
    ip_protocol = "tcp"
    from_port   = 6443
    to_port     = 6443
    cidr        = var.kubernetes_master_cidr
  }
  
  rule {
    ip_protocol = "tcp"
    from_port   = 10250
    to_port     = 10250
    cidr        = "0.0.0.0/0"
  }
}

# Auto-scaling policies
resource "k8s_autoscaling_v2_priority_class" "local_gpu" {
  api_version = "scheduling.k8s.io/v1"
  kind        = "PriorityClass"
  
  metadata {
    name = "local-gpu-baseline"
  }
  
  value             = 1500000000
  global_default    = false
  description       = "Priority for local GPU baseline workloads"
  preemption_policy = "Never"
}

resource "k8s_apps_v1_horizontal_pod_autoscaler" "local_gpu_checkpoint" {
  api_version = "autoscaling/v2"
  kind        = "HorizontalPodAutoscaler"
  
  metadata {
    name      = "local-gpu-checkpoint"
    namespace = "media"
  }
  
  spec {
    min_replicas = 0
    max_replicas = 2
    target_cpu_utilization_percentage = 80
    target_memory_utilization_percentage = 80
    
    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = "ome-transcoder"
    }
    
    behavior {
      scale_down {
        stabilization_window_seconds = 600  # 10 minutes
        policies {
          type          = "Percent"
          value         = 50
          period_seconds = 120
        }
      }
      
      scale_up {
        stabilization_window_seconds = 30
        policies {
          type          = "Percent"
          value         = 100
          period_seconds = 60
        }
        policies {
          type          = "Pods"
          value         = 2
          period_seconds = 60
        }
        select_policy = "Max"
      }
    }
  }
}

