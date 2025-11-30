// Ingestion Manager - Handles multi-protocol video ingestion

use dashmap::DashMap;
use std::sync::Arc;
use uuid::Uuid;

pub struct IngestConfig {
    pub max_concurrent_streams: usize,
    pub buffer_size_mb: usize,
    pub storage_endpoint: String,
    pub kafka_brokers: Vec<String>,
}

impl IngestConfig {
    pub fn from_env() -> anyhow::Result<Self> {
        Ok(Self {
            max_concurrent_streams: 10000,
            buffer_size_mb: 512,
            storage_endpoint: std::env::var("STORAGE_ENDPOINT")
                .unwrap_or_else(|_| "http://minio:9000".to_string()),
            kafka_brokers: std::env::var("KAFKA_BROKERS")
                .unwrap_or_else(|_| "kafka:9092".to_string())
                .split(',')
                .map(|s| s.to_string())
                .collect(),
        })
    }
}

pub struct Stream {
    pub id: String,
    pub protocol: String,
    pub input_url: String,
    pub status: StreamStatus,
    pub bitrate: u32,
    pub start_time: chrono::DateTime<chrono::Utc>,
}

#[derive(Debug, Clone)]
pub enum StreamStatus {
    Initializing,
    Active,
    Buffering,
    Error,
    Stopped,
}

pub struct IngestManager {
    config: IngestConfig,
    active_streams: Arc<DashMap<String, Stream>>,
}

impl IngestManager {
    pub fn new(config: IngestConfig) -> Self {
        Self {
            config,
            active_streams: Arc::new(DashMap::new()),
        }
    }

    pub async fn start_stream(
        &self,
        protocol: &str,
        input_url: &str,
    ) -> anyhow::Result<String> {
        let stream_id = Uuid::new_v4().to_string();

        let stream = Stream {
            id: stream_id.clone(),
            protocol: protocol.to_string(),
            input_url: input_url.to_string(),
            status: StreamStatus::Initializing,
            bitrate: 0,
            start_time: chrono::Utc::now(),
        };

        self.active_streams.insert(stream_id.clone(), stream);

        // Spawn ingestion task
        let stream_id_clone = stream_id.clone();
        tokio::spawn(async move {
            // Actual ingestion logic would go here
            tracing::info!("Ingestion task started for stream {}", stream_id_clone);
        });

        Ok(stream_id)
    }

    pub async fn stop_stream(&self, stream_id: &str) -> anyhow::Result<()> {
        self.active_streams.remove(stream_id);
        Ok(())
    }

    pub fn get_active_count(&self) -> usize {
        self.active_streams.len()
    }
}
