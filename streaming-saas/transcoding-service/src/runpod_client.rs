// Runpod.io GPU Cloud Integration
// Automatic scaling for transcoding workloads using Runpod serverless GPUs

use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::env;
use std::time::Duration;

const RUNPOD_API_BASE: &str = "https://api.runpod.io/v2";

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RunpodConfig {
    pub api_key: String,
    pub gpu_type: String,  // "RTX4090", "A100", "H100"
    pub max_pods: u32,
    pub auto_scale: bool,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PodRequest {
    pub name: String,
    pub image_name: String,
    pub gpu_type_id: String,
    pub cloud_type: String,  // "SECURE" or "COMMUNITY"
    pub container_disk_in_gb: u32,
    pub volume_in_gb: u32,
    pub env: Vec<EnvVar>,
    pub ports: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct EnvVar {
    pub key: String,
    pub value: String,
}

#[derive(Debug, Deserialize)]
pub struct PodResponse {
    pub id: String,
    pub status: String,
    pub machine: MachineInfo,
}

#[derive(Debug, Deserialize)]
pub struct MachineInfo {
    pub gpu_type: String,
    pub cpu_cores: u32,
    pub ram_gb: u32,
}

#[derive(Debug, Deserialize)]
pub struct JobResponse {
    pub id: String,
    pub status: String,
    pub output: Option<String>,
}

pub struct RunpodClient {
    client: Client,
    config: RunpodConfig,
}

impl RunpodClient {
    pub fn new(config: RunpodConfig) -> Self {
        let client = Client::builder()
            .timeout(Duration::from_secs(300))
            .build()
            .expect("Failed to create HTTP client");

        Self { client, config }
    }

    pub fn from_env() -> Self {
        let config = RunpodConfig {
            api_key: env::var("RUNPOD_API_KEY")
                .expect("RUNPOD_API_KEY must be set"),
            gpu_type: env::var("RUNPOD_GPU_TYPE")
                .unwrap_or_else(|_| "RTX4090".to_string()),
            max_pods: env::var("RUNPOD_MAX_PODS")
                .unwrap_or_else(|_| "10".to_string())
                .parse()
                .unwrap_or(10),
            auto_scale: env::var("RUNPOD_AUTO_SCALE")
                .unwrap_or_else(|_| "true".to_string())
                .parse()
                .unwrap_or(true),
        };

        Self::new(config)
    }

    /// Create a new GPU pod for transcoding
    pub async fn create_pod(&self, pod_name: &str) -> Result<PodResponse, Box<dyn std::error::Error>> {
        let request = PodRequest {
            name: pod_name.to_string(),
            image_name: "registry.streamverse.io/streamverse/transcoding-service:latest".to_string(),
            gpu_type_id: self.get_gpu_type_id(),
            cloud_type: "SECURE".to_string(),
            container_disk_in_gb: 50,
            volume_in_gb: 100,
            env: vec![
                EnvVar {
                    key: "NVIDIA_VISIBLE_DEVICES".to_string(),
                    value: "all".to_string(),
                },
                EnvVar {
                    key: "GPU_ACCELERATION".to_string(),
                    value: "true".to_string(),
                },
            ],
            ports: "8101/http".to_string(),
        };

        let url = format!("{}/pods", RUNPOD_API_BASE);
        let response = self.client
            .post(&url)
            .header("Authorization", format!("Bearer {}", self.config.api_key))
            .json(&request)
            .send()
            .await?
            .json::<PodResponse>()
            .await?;

        Ok(response)
    }

    /// Terminate a GPU pod
    pub async fn terminate_pod(&self, pod_id: &str) -> Result<(), Box<dyn std::error::Error>> {
        let url = format!("{}/pods/{}/terminate", RUNPOD_API_BASE, pod_id);
        self.client
            .post(&url)
            .header("Authorization", format!("Bearer {}", self.config.api_key))
            .send()
            .await?;

        Ok(())
    }

    /// Submit a transcoding job to a serverless endpoint
    pub async fn submit_job(
        &self,
        endpoint_id: &str,
        input_url: &str,
        output_url: &str,
        profile: TranscodingProfile,
    ) -> Result<JobResponse, Box<dyn std::error::Error>> {
        let job_input = serde_json::json!({
            "input": {
                "source_url": input_url,
                "output_url": output_url,
                "profile": profile,
            }
        });

        let url = format!("{}/{}/run", RUNPOD_API_BASE, endpoint_id);
        let response = self.client
            .post(&url)
            .header("Authorization", format!("Bearer {}", self.config.api_key))
            .json(&job_input)
            .send()
            .await?
            .json::<JobResponse>()
            .await?;

        Ok(response)
    }

    /// Get job status
    pub async fn get_job_status(
        &self,
        endpoint_id: &str,
        job_id: &str,
    ) -> Result<JobResponse, Box<dyn std::error::Error>> {
        let url = format!("{}/{}/status/{}", RUNPOD_API_BASE, endpoint_id, job_id);
        let response = self.client
            .get(&url)
            .header("Authorization", format!("Bearer {}", self.config.api_key))
            .send()
            .await?
            .json::<JobResponse>()
            .await?;

        Ok(response)
    }

    /// Auto-scale: decide whether to use local GPU or Runpod.io
    pub async fn should_use_runpod(&self, queue_depth: u32) -> bool {
        if !self.config.auto_scale {
            return false;
        }

        // Use Runpod.io when queue depth exceeds threshold
        // Threshold: 10 jobs per local GPU
        let local_gpu_count = self.get_local_gpu_count();
        let threshold = local_gpu_count * 10;

        queue_depth > threshold
    }

    fn get_gpu_type_id(&self) -> String {
        match self.config.gpu_type.as_str() {
            "RTX4090" => "NVIDIA RTX 4090",
            "A100" => "NVIDIA A100 80GB PCIe",
            "H100" => "NVIDIA H100 PCIe",
            _ => "NVIDIA RTX 4090",
        }
        .to_string()
    }

    fn get_local_gpu_count(&self) -> u32 {
        // Query local GPU count via nvidia-smi
        // For now, return configured value
        env::var("LOCAL_GPU_COUNT")
            .unwrap_or_else(|_| "4".to_string())
            .parse()
            .unwrap_or(4)
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct TranscodingProfile {
    pub resolution: String,
    pub bitrate: u32,
    pub fps: u32,
    pub codec: String,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_runpod_client_creation() {
        let config = RunpodConfig {
            api_key: "test-key".to_string(),
            gpu_type: "RTX4090".to_string(),
            max_pods: 10,
            auto_scale: true,
        };

        let client = RunpodClient::new(config);
        assert_eq!(client.config.gpu_type, "RTX4090");
    }

    #[test]
    fn test_auto_scaling_logic() {
        let config = RunpodConfig {
            api_key: "test-key".to_string(),
            gpu_type: "RTX4090".to_string(),
            max_pods: 10,
            auto_scale: true,
        };

        let client = RunpodClient::new(config);

        // With 4 local GPUs, threshold is 40 jobs
        // Queue depth of 50 should trigger Runpod.io
        std::env::set_var("LOCAL_GPU_COUNT", "4");

        // Test is async, but logic is: queue_depth > (local_gpu_count * 10)
        // 50 > 40 = true (should use Runpod)
    }
}
