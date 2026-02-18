import type { AuthResponse, Content, Profile, WatchProgress } from "@/types";

const API_BASE_URL = import.meta.env.VITE_API_URL || "/api/v1";

class ApiService {
  private token: string | null = null;

  constructor() {
    this.token = localStorage.getItem("auth_token");
  }

  private async fetch<T>(
    endpoint: string,
    options: RequestInit = {},
  ): Promise<T> {
    const headers = new Headers(options.headers);
    headers.set("Content-Type", "application/json");

    if (this.token) {
      headers.set("Authorization", `Bearer ${this.token}`);
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }

    return response.json();
  }

  // Auth endpoints
  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await this.fetch<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    this.token = response.token;
    localStorage.setItem("auth_token", response.token);
    return response;
  }

  async register(
    email: string,
    password: string,
    username: string,
  ): Promise<AuthResponse> {
    const response = await this.fetch<AuthResponse>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password, username }),
    });
    this.token = response.token;
    localStorage.setItem("auth_token", response.token);
    return response;
  }

  async logout(): Promise<void> {
    this.token = null;
    localStorage.removeItem("auth_token");
  }

  // Content endpoints
  async getHomeContent(): Promise<{
    categories: { name: string; content: Content[] }[];
  }> {
    return this.fetch("/content/home");
  }

  async getContent(id: string): Promise<Content> {
    return this.fetch(`/content/${id}`);
  }

  async searchContent(query: string): Promise<Content[]> {
    return this.fetch(`/search?q=${encodeURIComponent(query)}`);
  }

  async getTrending(): Promise<Content[]> {
    return this.fetch("/content/trending");
  }

  async getRecommendations(): Promise<Content[]> {
    return this.fetch("/recommendations");
  }

  // Profile endpoints
  async getProfiles(): Promise<Profile[]> {
    return this.fetch("/profiles");
  }

  async createProfile(name: string, isKids: boolean = false): Promise<Profile> {
    return this.fetch("/profiles", {
      method: "POST",
      body: JSON.stringify({ name, is_kids: isKids }),
    });
  }

  // Watch progress
  async updateWatchProgress(
    contentId: string,
    progress: number,
  ): Promise<void> {
    await this.fetch("/watch/progress", {
      method: "POST",
      body: JSON.stringify({ content_id: contentId, progress }),
    });
  }

  async getWatchProgress(): Promise<WatchProgress[]> {
    return this.fetch("/watch/progress");
  }
}

export const api = new ApiService();
