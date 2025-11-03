/**
 * API Service for StreamVerse Tizen App
 */

const API_BASE_URL = 'https://api.streamverse.com/';

class ApiService {
    constructor() {
        this.token = null;
    }

    setToken(token) {
        this.token = token;
        // Store in Tizen preference
        if (typeof tizen !== 'undefined') {
            tizen.preference.setValue('access_token', token);
        }
    }

    getToken() {
        if (!this.token && typeof tizen !== 'undefined') {
            this.token = tizen.preference.getValue('access_token');
        }
        return this.token;
    }

    async request(endpoint, options = {}) {
        const url = API_BASE_URL + endpoint;
        const headers = {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            ...options.headers
        };

        const token = this.getToken();
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const config = {
            ...options,
            headers
        };

        try {
            const response = await fetch(url, config);
            
            if (!response.ok) {
                if (response.status === 401) {
                    // Token expired, redirect to login
                    this.clearToken();
                    window.location.hash = '#login';
                    throw new Error('Unauthorized');
                }
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }

    clearToken() {
        this.token = null;
        if (typeof tizen !== 'undefined') {
            tizen.preference.remove('access_token');
        }
    }

    // Auth endpoints
    async login(email, password) {
        return this.request('api/v1/auth/login', {
            method: 'POST',
            body: JSON.stringify({ email, password })
        });
    }

    async logout() {
        try {
            await this.request('api/v1/auth/logout', { method: 'POST' });
        } finally {
            this.clearToken();
        }
    }

    // Content endpoints
    async getHomeContent() {
        return this.request('api/v1/content/home');
    }

    async getContentById(id) {
        return this.request(`api/v1/content/${id}`);
    }

    async searchContent(query) {
        return this.request(`api/v1/content/search?q=${encodeURIComponent(query)}`);
    }

    async getContentByCategory(category) {
        return this.request(`api/v1/content/category/${category}`);
    }
}

const apiService = new ApiService();

