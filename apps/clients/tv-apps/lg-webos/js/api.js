/**
 * API Service for StreamVerse webOS App
 */

const API_BASE_URL = 'https://api.streamverse.com/';

class ApiService {
    constructor() {
        this.token = null;
    }

    setToken(token) {
        this.token = token;
        if (typeof webOS !== 'undefined' && webOS.service) {
            webOS.service.request('luna://com.palm.systemservice', {
                method: 'getPreferenceValues',
                parameters: { keys: ['access_token'] },
                onSuccess: (response) => {
                    // Store token using webOS preferences
                },
                onFailure: (error) => {
                    console.error('Failed to store token:', error);
                }
            });
        } else {
            localStorage.setItem('access_token', token);
        }
    }

    getToken() {
        if (!this.token) {
            this.token = localStorage.getItem('access_token');
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

        try {
            const response = await fetch(url, {
                ...options,
                headers
            });
            
            if (!response.ok) {
                if (response.status === 401) {
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
        localStorage.removeItem('access_token');
    }

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

    async getHomeContent() {
        return this.request('api/v1/content/home');
    }

    async getContentById(id) {
        return this.request(`api/v1/content/${id}`);
    }

    async searchContent(query) {
        return this.request(`api/v1/content/search?q=${encodeURIComponent(query)}`);
    }
}

const apiService = new ApiService();

