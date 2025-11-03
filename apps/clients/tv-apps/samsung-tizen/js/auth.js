/**
 * Authentication Module
 */

// Auth functions are handled in api.js
// This file can be extended with additional auth-related utilities

function isAuthenticated() {
    return apiService.getToken() !== null;
}

function requireAuth() {
    if (!isAuthenticated()) {
        window.location.hash = '#login';
        return false;
    }
    return true;
}

