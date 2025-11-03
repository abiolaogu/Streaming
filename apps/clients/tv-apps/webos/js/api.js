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
        // Store in webOS preference
        if (typeof webOS !== 'undefined' && webOS.service) {
            webOS.service.request('luna://com.webos.service.preferences', {
                method: 'set',
                parameters: {
                    key: 'streamverse_token',
                    value: token
                }
            });
        } else {
            // Fallback to localStorage
            localStorage.setItem('access_token', token);
        }
    }

    getToken() {
        if (!this.token) {
            if (typeof webOS !== 'undefined' && webOS.service) {
                try {
                    this.token = localStorage.getItem('access_token'); // Simplified
                } catch (e) {
                    console.warn('Could not retrieve token from webOS service');
                }
            } else {
                this.token = localStorage.getItem('access_token');
            }
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

        if (options.body && typeof options.body === 'object') {
            config.body = JSON.stringify(options.body);
        }

        try {
            const response = await fetch(url, config);
            
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
            body: { email, password }
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

