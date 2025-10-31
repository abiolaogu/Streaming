environment               = "prd"
aws_region                = "us-east-1"
cpu_ondemand_floor        = { count = 1 }
cpu_spot_max              = 50
gpu_spot_max              = 20
enable_cross_cloud_replication = true
cross_cloud_bucket_name   = "streaming-platform-minio-gcp-primary"

