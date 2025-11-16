export interface User {
  id: string
  email: string
  username: string
  profile_picture?: string
  subscription_tier: 'free' | 'basic' | 'standard' | 'premium'
  created_at: string
}

export interface Content {
  id: string
  title: string
  description: string
  thumbnail_url: string
  video_url: string
  duration: number
  type: 'movie' | 'series' | 'documentary' | 'short'
  genres: string[]
  release_date: string
  rating: number
  views: number
  created_at: string
}

export interface Episode {
  id: string
  series_id: string
  season_number: number
  episode_number: number
  title: string
  description: string
  thumbnail_url: string
  video_url: string
  duration: number
  release_date: string
}

export interface Profile {
  id: string
  user_id: string
  name: string
  avatar_url: string
  is_kids: boolean
  language: string
}

export interface WatchProgress {
  content_id: string
  profile_id: string
  progress: number
  last_watched: string
}

export interface AuthResponse {
  user: User
  token: string
  refresh_token: string
}
