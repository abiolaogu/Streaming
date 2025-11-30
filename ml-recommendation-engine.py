"""
============================================================================
AI/ML Recommendation Engine - Global Streaming Platform
============================================================================
"""

import numpy as np
import pandas as pd
from typing import List, Dict, Tuple
import torch
import torch.nn as nn
from torch.utils.data import Dataset, DataLoader
from sklearn.preprocessing import LabelEncoder
from redis import Redis
import psycopg2
from datetime import datetime, timedelta
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


# ============================================================================
# Neural Collaborative Filtering Model
# ============================================================================

class NCFModel(nn.Module):
    """Neural Collaborative Filtering for recommendations"""
    
    def __init__(self, num_users: int, num_items: int, 
                 embedding_dim: int = 64, hidden_layers: List[int] = [128, 64, 32]):
        super(NCFModel, self).__init__()
        
        # User and item embeddings
        self.user_embedding = nn.Embedding(num_users, embedding_dim)
        self.item_embedding = nn.Embedding(num_items, embedding_dim)
        
        # MLP layers
        mlp_layers = []
        input_dim = embedding_dim * 2
        
        for hidden_dim in hidden_layers:
            mlp_layers.extend([
                nn.Linear(input_dim, hidden_dim),
                nn.ReLU(),
                nn.BatchNorm1d(hidden_dim),
                nn.Dropout(0.2)
            ])
            input_dim = hidden_dim
        
        # Output layer
        mlp_layers.append(nn.Linear(input_dim, 1))
        mlp_layers.append(nn.Sigmoid())
        
        self.mlp = nn.Sequential(*mlp_layers)
        
        # Initialize weights
        self._init_weights()
    
    def _init_weights(self):
        for module in self.modules():
            if isinstance(module, nn.Embedding):
                nn.init.normal_(module.weight, std=0.01)
            elif isinstance(module, nn.Linear):
                nn.init.xavier_uniform_(module.weight)
                if module.bias is not None:
                    nn.init.constant_(module.bias, 0)
    
    def forward(self, user_ids: torch.Tensor, item_ids: torch.Tensor) -> torch.Tensor:
        # Get embeddings
        user_emb = self.user_embedding(user_ids)
        item_emb = self.item_embedding(item_ids)
        
        # Concatenate embeddings
        x = torch.cat([user_emb, item_emb], dim=-1)
        
        # Pass through MLP
        output = self.mlp(x)
        
        return output.squeeze()


class InteractionDataset(Dataset):
    """Dataset for user-item interactions"""
    
    def __init__(self, user_ids: np.ndarray, item_ids: np.ndarray, 
                 ratings: np.ndarray):
        self.user_ids = torch.LongTensor(user_ids)
        self.item_ids = torch.LongTensor(item_ids)
        self.ratings = torch.FloatTensor(ratings)
    
    def __len__(self):
        return len(self.user_ids)
    
    def __getitem__(self, idx):
        return (self.user_ids[idx], self.item_ids[idx], self.ratings[idx])


# ============================================================================
# Recommendation Service
# ============================================================================

from psycopg2 import pool

class RecommendationService:
    """Main recommendation service with multiple algorithms"""
    
    def __init__(self, db_config: Dict, redis_host: str = 'localhost', 
                 redis_port: int = 6379):
        # Database connection pool
        self.db_pool = pool.SimpleConnectionPool(
            1, 20,
            **db_config
        )
        
        # Redis for caching
        self.redis_client = Redis(host=redis_host, port=redis_port, 
                                   decode_responses=True)
        
        # Model storage
        self.ncf_model = None
        self.user_encoder = LabelEncoder()
        self.item_encoder = LabelEncoder()
        
        # Device
        self.device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
        logger.info(f"Using device: {self.device}")

    def get_db_connection(self):
        return self.db_pool.getconn()

    def release_db_connection(self, conn):
        self.db_pool.putconn(conn)
    
    def load_interaction_data(self, days: int = 90) -> pd.DataFrame:
        """Load user-item interaction data from database"""
        
        query = """
        SELECT 
            wh.profile_id,
            wh.content_id,
            CASE 
                WHEN wh.completed THEN 1.0
                WHEN wh.position::float / NULLIF(wh.duration, 0) > 0.7 THEN 0.8
                WHEN wh.position::float / NULLIF(wh.duration, 0) > 0.3 THEN 0.5
                ELSE 0.2
            END as implicit_rating,
            wh.watched_at
        FROM watch_history wh
        WHERE wh.watched_at > NOW() - INTERVAL '%s days'
        ORDER BY wh.watched_at DESC
        """
        
        conn = self.get_db_connection()
        try:
            df = pd.read_sql(query, conn, params=(days,))
            logger.info(f"Loaded {len(df)} interactions from last {days} days")
            return df
        finally:
            self.release_db_connection(conn)
    
    def train_ncf_model(self, df: pd.DataFrame, epochs: int = 10, 
                        batch_size: int = 1024, lr: float = 0.001):
        """Train Neural Collaborative Filtering model"""
        
        # Encode user and item IDs
        df['user_idx'] = self.user_encoder.fit_transform(df['profile_id'])
        df['item_idx'] = self.item_encoder.fit_transform(df['content_id'])
        
        num_users = len(self.user_encoder.classes_)
        num_items = len(self.item_encoder.classes_)
        
        logger.info(f"Training NCF: {num_users} users, {num_items} items")
        
        # Create dataset and dataloader
        dataset = InteractionDataset(
            df['user_idx'].values,
            df['item_idx'].values,
            df['implicit_rating'].values
        )
        
        dataloader = DataLoader(dataset, batch_size=batch_size, 
                                shuffle=True, num_workers=4)
        
        # Initialize model
        self.ncf_model = NCFModel(num_users, num_items).to(self.device)
        
        # Loss and optimizer
        criterion = nn.BCELoss()
        optimizer = torch.optim.Adam(self.ncf_model.parameters(), lr=lr)
        
        # Training loop
        for epoch in range(epochs):
            self.ncf_model.train()
            total_loss = 0
            
            for batch_users, batch_items, batch_ratings in dataloader:
                batch_users = batch_users.to(self.device)
                batch_items = batch_items.to(self.device)
                batch_ratings = batch_ratings.to(self.device)
                
                # Forward pass
                predictions = self.ncf_model(batch_users, batch_items)
                loss = criterion(predictions, batch_ratings)
                
                # Backward pass
                optimizer.zero_grad()
                loss.backward()
                optimizer.step()
                
                total_loss += loss.item()
            
            avg_loss = total_loss / len(dataloader)
            logger.info(f"Epoch {epoch+1}/{epochs}, Loss: {avg_loss:.4f}")
        
        logger.info("NCF training completed")
    
    def get_recommendations_ncf(self, profile_id: str, top_k: int = 20) -> List[str]:
        """Get recommendations using NCF model"""
        
        # Check cache first
        cache_key = f"recs:ncf:{profile_id}:{top_k}"
        cached = self.redis_client.get(cache_key)
        if cached:
            logger.info(f"Cache hit for profile {profile_id}")
            return eval(cached)
        
        if self.ncf_model is None:
            logger.warning("NCF model not trained, falling back to popular")
            return self.get_popular_content(top_k)
        
        try:
            # Encode profile ID
            user_idx = self.user_encoder.transform([profile_id])[0]
            user_tensor = torch.LongTensor([user_idx]).to(self.device)
            
            # Get all item indices
            all_items = torch.LongTensor(range(len(self.item_encoder.classes_))).to(self.device)
            
            # Predict scores for all items
            self.ncf_model.eval()
            with torch.no_grad():
                user_expanded = user_tensor.repeat(len(all_items))
                scores = self.ncf_model(user_expanded, all_items)
            
            # Get top-k items
            top_indices = torch.topk(scores, k=top_k).indices.cpu().numpy()
            recommended_content_ids = self.item_encoder.inverse_transform(top_indices).tolist()
            
            # Cache for 1 hour
            self.redis_client.setex(cache_key, 3600, str(recommended_content_ids))
            
            logger.info(f"Generated {len(recommended_content_ids)} NCF recommendations for {profile_id}")
            return recommended_content_ids
            
        except Exception as e:
            logger.error(f"NCF recommendation failed: {e}")
            return self.get_popular_content(top_k)
    
    def get_collaborative_filtering_recs(self, profile_id: str, top_k: int = 20) -> List[str]:
        """Item-based collaborative filtering using cosine similarity"""
        
        cache_key = f"recs:cf:{profile_id}:{top_k}"
        cached = self.redis_client.get(cache_key)
        if cached:
            return eval(cached)
        
        # Get user's watch history
        query = """
        SELECT content_id, implicit_rating
        FROM (
            SELECT 
                content_id,
                CASE 
                    WHEN completed THEN 1.0
                    WHEN position::float / NULLIF(duration, 0) > 0.7 THEN 0.8
                    ELSE 0.5
                END as implicit_rating
            FROM watch_history
            WHERE profile_id = %s
            ORDER BY watched_at DESC
            LIMIT 50
        ) recent_watches
        """
        
        conn = self.get_db_connection()
        try:
            cursor = conn.cursor()
            cursor.execute(query, (profile_id,))
            watched_items = cursor.fetchall()
            cursor.close()
        finally:
            self.release_db_connection(conn)
        
        if not watched_items:
            return self.get_popular_content(top_k)
        
        # Find similar items based on co-occurrence
        watched_ids = [item[0] for item in watched_items]
        placeholders = ','.join(['%s'] * len(watched_ids))
        
        query = f"""
        SELECT 
            wh2.content_id,
            COUNT(DISTINCT wh2.profile_id) as co_occurrence,
            AVG(CASE WHEN wh2.completed THEN 1.0 ELSE 0.5 END) as avg_rating
        FROM watch_history wh1
        JOIN watch_history wh2 ON wh1.profile_id = wh2.profile_id
        WHERE wh1.content_id IN ({placeholders})
          AND wh2.content_id NOT IN ({placeholders})
          AND wh2.watched_at > NOW() - INTERVAL '90 days'
        GROUP BY wh2.content_id
        ORDER BY co_occurrence DESC, avg_rating DESC
        LIMIT %s
        """
        
        conn = self.get_db_connection()
        try:
            cursor = conn.cursor()
            cursor.execute(query, watched_ids + watched_ids + [top_k])
            recommendations = [row[0] for row in cursor.fetchall()]
            cursor.close()
        finally:
            self.release_db_connection(conn)
        
        # Cache for 30 minutes
        self.redis_client.setex(cache_key, 1800, str(recommendations))
        
        logger.info(f"Generated {len(recommendations)} CF recommendations for {profile_id}")
        return recommendations
    
    def get_popular_content(self, top_k: int = 20, region: str = None) -> List[str]:
        """Get trending/popular content"""
        
        cache_key = f"popular:{region or 'global'}:{top_k}"
        cached = self.redis_client.get(cache_key)
        if cached:
            return eval(cached)
        
        # Get popular content from last 7 days
        query = """
        SELECT 
            content_id,
            COUNT(DISTINCT profile_id) as unique_viewers,
            AVG(CASE WHEN completed THEN 1.0 ELSE 0.5 END) as engagement_score
        FROM watch_history
        WHERE watched_at > NOW() - INTERVAL '7 days'
        GROUP BY content_id
        ORDER BY unique_viewers DESC, engagement_score DESC
        LIMIT %s
        """
        
        conn = self.get_db_connection()
        try:
            cursor = conn.cursor()
            cursor.execute(query, (top_k,))
            popular = [row[0] for row in cursor.fetchall()]
            cursor.close()
        finally:
            self.release_db_connection(conn)
        
        # Cache for 10 minutes
        self.redis_client.setex(cache_key, 600, str(popular))
        
        return popular
    
    def get_content_based_recs(self, profile_id: str, top_k: int = 20) -> List[str]:
        """Content-based recommendations using genre/metadata similarity"""
        
        # Get user's genre preferences
        query = """
        WITH user_genres AS (
            SELECT 
                UNNEST(c.genres) as genre,
                COUNT(*) as watch_count
            FROM watch_history wh
            JOIN content c ON wh.content_id = c.id
            WHERE wh.profile_id = %s
              AND wh.watched_at > NOW() - INTERVAL '90 days'
            GROUP BY genre
            ORDER BY watch_count DESC
            LIMIT 5
        )
        SELECT 
            c.id,
            COUNT(*) as genre_match_count,
            c.release_year
        FROM content c
        CROSS JOIN UNNEST(c.genres) AS content_genre
        JOIN user_genres ug ON content_genre = ug.genre
        WHERE c.is_published = TRUE
          AND c.id NOT IN (
              SELECT content_id FROM watch_history WHERE profile_id = %s
          )
        GROUP BY c.id, c.release_year
        ORDER BY genre_match_count DESC, c.release_year DESC
        LIMIT %s
        """
        
        conn = self.get_db_connection()
        try:
            cursor = conn.cursor()
            cursor.execute(query, (profile_id, profile_id, top_k))
            recommendations = [row[0] for row in cursor.fetchall()]
            cursor.close()
        finally:
            self.release_db_connection(conn)
        
        return recommendations
    
    def get_hybrid_recommendations(self, profile_id: str, top_k: int = 20) -> List[Dict]:
        """Hybrid recommendations combining multiple strategies"""
        
        cache_key = f"recs:hybrid:{profile_id}:{top_k}"
        cached = self.redis_client.get(cache_key)
        if cached:
            return eval(cached)
        
        # Get recommendations from different sources
        ncf_recs = self.get_recommendations_ncf(profile_id, top_k=30)
        cf_recs = self.get_collaborative_filtering_recs(profile_id, top_k=30)
        cb_recs = self.get_content_based_recs(profile_id, top_k=30)
        popular = self.get_popular_content(top_k=20)
        
        # Score each recommendation
        scores = {}
        
        # NCF weight: 0.4
        for i, content_id in enumerate(ncf_recs):
            scores[content_id] = scores.get(content_id, 0) + 0.4 * (1 - i/len(ncf_recs))
        
        # CF weight: 0.3
        for i, content_id in enumerate(cf_recs):
            scores[content_id] = scores.get(content_id, 0) + 0.3 * (1 - i/len(cf_recs))
        
        # Content-based weight: 0.2
        for i, content_id in enumerate(cb_recs):
            scores[content_id] = scores.get(content_id, 0) + 0.2 * (1 - i/len(cb_recs))
        
        # Popular weight: 0.1
        for i, content_id in enumerate(popular):
            scores[content_id] = scores.get(content_id, 0) + 0.1 * (1 - i/len(popular))
        
        # Sort by score and take top-k
        sorted_recs = sorted(scores.items(), key=lambda x: x[1], reverse=True)[:top_k]
        
        # Fetch content details
        content_ids = [rec[0] for rec in sorted_recs]
        placeholders = ','.join(['%s'] * len(content_ids))
        
        query = f"""
        SELECT id, title, type, poster_url, rating, genres
        FROM content
        WHERE id IN ({placeholders})
        """
        
        conn = self.get_db_connection()
        try:
            cursor = conn.cursor()
            cursor.execute(query, tuple(content_ids))
            
            results = []
            for row in cursor.fetchall():
                results.append({
                    'id': str(row[0]),
                    'title': row[1],
                    'type': row[2],
                    'poster_url': row[3],
                    'rating': row[4],
                    'genres': row[5]
                })
            cursor.close()
        finally:
            self.release_db_connection(conn)
        
        # Cache for 15 minutes
        self.redis_client.setex(cache_key, 900, str(results))
        
        logger.info(f"Generated {len(results)} hybrid recommendations for {profile_id}")
        return results


# ============================================================================
# CDN Selection ML Model
# ============================================================================

class CDNSelectorModel(nn.Module):
    """ML model for optimal CDN edge selection"""
    
    def __init__(self, num_edges: int, feature_dim: int = 20):
        super(CDNSelectorModel, self).__init__()
        
        self.network = nn.Sequential(
            nn.Linear(feature_dim, 64),
            nn.ReLU(),
            nn.BatchNorm1d(64),
            nn.Dropout(0.2),
            nn.Linear(64, 32),
            nn.ReLU(),
            nn.Linear(32, num_edges),
            nn.Softmax(dim=-1)
        )
    
    def forward(self, features: torch.Tensor) -> torch.Tensor:
        return self.network(features)


class CDNSelector:
    """CDN edge selector using ML for optimal performance"""
    
    def __init__(self):
        self.model = None
        self.edge_map = {
            0: 'eu-west-edge.cdn.example.com',
            1: 'eu-central-edge.cdn.example.com',
            2: 'us-east-edge.cdn.example.com',
            3: 'us-west-edge.cdn.example.com',
            4: 'af-west-edge.cdn.example.com',
            5: 'ap-southeast-edge.cdn.example.com',
            6: 'sa-east-edge.cdn.example.com',
        }
    
    def extract_features(self, user_region: str, content_id: str, 
                         time_of_day: int, day_of_week: int) -> np.ndarray:
        """Extract features for CDN selection"""
        
        features = []
        
        # Region encoding (one-hot)
        regions = ['eu-west', 'eu-central', 'us-east', 'us-west', 
                   'af-west', 'ap-southeast', 'sa-east']
        region_encoding = [1 if user_region == r else 0 for r in regions]
        features.extend(region_encoding)
        
        # Time features
        features.append(time_of_day / 24.0)  # Normalized hour
        features.append(day_of_week / 7.0)    # Normalized day
        
        # Content hotness (from Aerospike/cache)
        # Simplified - in production, query actual hotness score
        features.append(0.5)
        
        # Historical performance metrics (simplified)
        features.extend([0.9, 0.85, 0.92, 0.88, 0.87, 0.90, 0.89])  # Hit ratios
        features.extend([0.15, 0.20, 0.12, 0.18, 0.25, 0.14, 0.22]) # Avg latency
        
        return np.array(features, dtype=np.float32)
    
    def select_edge(self, user_region: str, content_id: str) -> str:
        """Select optimal CDN edge for user request"""
        
        now = datetime.now()
        features = self.extract_features(
            user_region, content_id,
            now.hour, now.weekday()
        )
        
        if self.model is None:
            # Fallback to rule-based
            return self.edge_map.get(hash(user_region) % len(self.edge_map))
        
        # ML-based selection
        with torch.no_grad():
            features_tensor = torch.FloatTensor(features).unsqueeze(0)
            probabilities = self.model(features_tensor)
            edge_idx = torch.argmax(probabilities).item()
        
        return self.edge_map[edge_idx]


# ============================================================================
# Main Training Script
# ============================================================================

if __name__ == "__main__":
    # Database configuration
    db_config = {
        'host': 'localhost',
        'database': 'streaming',
        'user': 'postgres',
        'password': 'password'
    }
    
    # Initialize recommendation service
    rec_service = RecommendationService(db_config)
    
    # Load data
    df = rec_service.load_interaction_data(days=90)
    
    # Train NCF model
    rec_service.train_ncf_model(df, epochs=10, batch_size=1024)
    
    # Test recommendations
    test_profile_id = "test-profile-uuid"
    recommendations = rec_service.get_hybrid_recommendations(test_profile_id, top_k=20)
    
    print(f"Top 20 recommendations for profile {test_profile_id}:")
    for i, rec in enumerate(recommendations, 1):
        print(f"{i}. {rec['title']} ({rec['type']}) - {rec['genres']}")
