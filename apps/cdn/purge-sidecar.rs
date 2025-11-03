use axum::{
    extract::{Path, Query},
    http::StatusCode,
    routing::{get, post},
    Router,
};
use quiche::h3;
use serde::{Deserialize, Serialize};
use std::{collections::HashMap, sync::Arc};
use tokio::sync::RwLock;
use tower::ServiceBuilder;
use tower_http::cors::CorsLayer;

#[derive(Clone)]
struct PurgeState {
    nodes: Arc<RwLock<HashMap<String, Vec<String>>>>,
    token_bucket: Arc<RwLock<TokenBucket>>,
}

#[derive(Clone, Copy)]
struct TokenBucket {
    tokens: usize,
    rate: usize,
    max_tokens: usize,
}

impl TokenBucket {
    fn new(rate: usize, max_tokens: usize) -> Self {
        Self {
            tokens: max_tokens,
            rate,
            max_tokens,
        }
    }

    async fn consume(&mut self, n: usize) -> bool {
        self.tokens = (self.tokens + self.rate).min(self.max_tokens);
        if self.tokens >= n {
            self.tokens -= n;
            true
        } else {
            false
        }
    }
}

#[derive(Deserialize)]
struct PurgeRequest {
    url: String,
    pattern: Option<String>,
}

#[derive(Serialize)]
struct PurgeResponse {
    success: bool,
    message: String,
    purged_nodes: Vec<String>,
}

async fn health_check() -> StatusCode {
    StatusCode::OK
}

async fn purge_url(
    Path(cache_tier): Path<String>,
    Query(params): Query<PurgeRequest>,
    state: axum::extract::State<PurgeState>,
) -> Result<impl axum::response::IntoResponse, StatusCode> {
    if !state
        .token_bucket
        .write()
        .await
        .consume(1)
        .await
    {
        return Err(StatusCode::TOO_MANY_REQUESTS);
    }

    let nodes = state.nodes.read().await;
    let cache_nodes = nodes
        .get(&cache_tier)
        .ok_or_else(|| StatusCode::NOT_FOUND)?;

    let mut purged = Vec::new();
    let client = reqwest::Client::new();

    for node in cache_nodes {
        let url = format!("http://{}/cache/purge?url={}", node, params.url);
        match client.get(&url).send().await {
            Ok(_) => purged.push(node.clone()),
            Err(_) => continue,
        }
    }

    Ok(axum::Json(PurgeResponse {
        success: !purged.is_empty(),
        message: format!("Purged on {} nodes", purged.len()),
        purged_nodes: purged,
    }))
}

async fn ban_pattern(
    Path(cache_tier): Path<String>,
    Query(params): Query<PurgeRequest>,
    state: axum::extract::State<PurgeState>,
) -> Result<impl axum::response::IntoResponse, StatusCode> {
    if !state
        .token_bucket
        .write()
        .await
        .consume(1)
        .await
    {
        return Err(StatusCode::TOO_MANY_REQUESTS);
    }

    let nodes = state.nodes.read().await;
    let cache_nodes = nodes
        .get(&cache_tier)
        .ok_or_else(|| StatusCode::NOT_FOUND)?;

    let pattern = params.pattern.unwrap_or_else(|| params.url);

    let mut banned = Vec::new();
    let client = reqwest::Client::new();

    for node in cache_nodes {
        let url = format!("http://{}/cache/ban?pattern={}", node, pattern);
        match client.get(&url).send().await {
            Ok(_) => banned.push(node.clone()),
            Err(_) => continue,
        }
    }

    Ok(axum::Json(PurgeResponse {
        success: !banned.is_empty(),
        message: format!("Banned on {} nodes", banned.len()),
        purged_nodes: banned,
    }))
}

async fn fanout_purge(
    body: axum::Json<PurgeRequest>,
    state: axum::extract::State<PurgeState>,
) -> Result<impl axum::response::IntoResponse, StatusCode> {
    if !state
        .token_bucket
        .write()
        .await
        .consume(10)
        .await
    {
        return Err(StatusCode::TOO_MANY_REQUESTS);
    }

    let nodes = state.nodes.read().await;
    let mut all_purged = Vec::new();

    for (tier, cache_nodes) in nodes.iter() {
        let client = reqwest::Client::new();
        for node in cache_nodes {
            let url = format!("http://{}/cache/purge?url={}", node, body.url);
            if client.get(&url).send().await.is_ok() {
                all_purged.push(format!("{}/{}", tier, node));
            }
        }
    }

    Ok(axum::Json(PurgeResponse {
        success: !all_purged.is_empty(),
        message: format!("Fanout purge on {} nodes", all_purged.len()),
        purged_nodes: all_purged,
    }))
}

#[tokio::main]
async fn main() {
    let mut nodes = HashMap::new();
    nodes.insert(
        "shield".to_string(),
        vec![
            "varnish-shield-1:80".to_string(),
            "varnish-shield-2:80".to_string(),
            "varnish-shield-3:80".to_string(),
        ],
    );
    nodes.insert(
        "edge".to_string(),
        vec![
            "ats-edge-1:8080".to_string(),
            "ats-edge-2:8080".to_string(),
            "ats-edge-3:8080".to_string(),
        ],
    );

    let state = PurgeState {
        nodes: Arc::new(RwLock::new(nodes)),
        token_bucket: Arc::new(RwLock::new(TokenBucket::new(100, 1000))),
    };

    let app = Router::new()
        .route("/health", get(health_check))
        .route("/purge/:tier", get(purge_url))
        .route("/ban/:tier", get(ban_pattern))
        .route("/fanout", post(fanout_purge))
        .layer(
            ServiceBuilder::new()
                .layer(CorsLayer::permissive())
                .into_inner(),
        )
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:8080").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

