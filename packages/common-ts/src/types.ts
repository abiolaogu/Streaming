/**
 * Common TypeScript type definitions
 */

export interface User {
  id: string;
  email: string;
  name?: string;
  avatar?: string;
  roles: string[];
  createdAt: string;
  updatedAt: string;
}

export interface Content {
  id: string;
  title: string;
  description: string;
  genre: string;
  category: 'movie' | 'show' | 'live';
  posterUrl: string;
  backdropUrl: string;
  streamUrl: string;
  duration: number;
  releaseYear: number;
  rating: number;
  isDrmProtected: boolean;
  drmType?: 'widevine' | 'playready' | 'fairplay';
  thumbnailUrl?: string;
  cast: string[];
  directors: string[];
  tags: string[];
  createdAt: string;
  updatedAt: string;
}

export interface ContentRow {
  id: string;
  title: string;
  items: Content[];
}

export interface AuthResponse {
  token: string;
  refreshToken?: string;
  user?: User;
  expiresAt?: string;
  expiresIn?: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface SearchRequest {
  query: string;
  filters?: {
    genre?: string;
    year?: number;
    rating?: number;
    category?: string;
  };
  sort?: 'relevance' | 'popularity' | 'date';
  page?: number;
  pageSize?: number;
}

export interface SearchResponse {
  results: Content[];
  total: number;
  page: number;
  pageSize: number;
  aggregations?: Record<string, any>;
}

