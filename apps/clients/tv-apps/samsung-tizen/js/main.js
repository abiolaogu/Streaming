/**
 * Main Application Entry Point
 */

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    initializeApp();
});

function initializeApp() {
    // Check if running on Tizen
    if (typeof tizen === 'undefined') {
        console.warn('Tizen APIs not available - running in browser mode');
    }

    // Initialize player
    tizenPlayer.init();

    // Check authentication
    checkAuth();

    // Setup event listeners
    setupEventListeners();

    // Handle hash routing
    window.addEventListener('hashchange', handleRoute);
    handleRoute();

    // Update navigation when screen changes
    const observer = new MutationObserver(() => {
        tvNavigation.updateFocusableElements();
    });
    observer.observe(document.body, { childList: true, subtree: true });
}

function checkAuth() {
    const token = apiService.getToken();
    if (token) {
        // User is logged in
        window.location.hash = '#home';
    } else {
        // Show login screen
        window.location.hash = '#login';
    }
}

function setupEventListeners() {
    // Login form
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }

    // Navigation buttons
    document.getElementById('search-button')?.addEventListener('click', () => {
        window.location.hash = '#search';
    });

    document.getElementById('settings-button')?.addEventListener('click', () => {
        window.location.hash = '#settings';
    });

    document.getElementById('back-button')?.addEventListener('click', () => {
        window.location.hash = '#home';
    });

    document.getElementById('search-back-button')?.addEventListener('click', () => {
        window.location.hash = '#home';
    });

    document.getElementById('settings-back-button')?.addEventListener('click', () => {
        window.location.hash = '#home';
    });

    document.getElementById('logout-button')?.addEventListener('click', handleLogout);

    document.getElementById('player-close-button')?.addEventListener('click', () => {
        tizenPlayer.stop();
        window.location.hash = '#detail';
    });

    // Search input
    const searchInput = document.getElementById('search-input');
    if (searchInput) {
        searchInput.addEventListener('input', debounce(handleSearch, 500));
    }
}

function handleRoute() {
    const hash = window.location.hash || '#login';
    const route = hash.substring(1);

    // Hide all screens
    document.querySelectorAll('.screen').forEach(screen => {
        screen.classList.remove('active');
    });

    // Show active screen
    const activeScreen = document.getElementById(`${route}-screen`);
    if (activeScreen) {
        activeScreen.classList.add('active');
        
        // Load screen-specific content
        switch(route) {
            case 'home':
                loadHomeContent();
                break;
            case 'login':
                // Already handled
                break;
        }
    } else {
        // Default to login if route not found
        document.getElementById('login-screen').classList.add('active');
    }

    // Update navigation
    setTimeout(() => {
        tvNavigation.updateFocusableElements();
        tvNavigation.focusFirst();
    }, 100);
}

async function handleLogin(e) {
    e.preventDefault();
    
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const errorDiv = document.getElementById('login-error');

    errorDiv.textContent = '';

    try {
        const response = await apiService.login(email, password);
        apiService.setToken(response.token);
        window.location.hash = '#home';
    } catch (error) {
        errorDiv.textContent = 'Login failed. Please check your credentials.';
        console.error('Login error:', error);
    }
}

async function handleLogout() {
    try {
        await apiService.logout();
        window.location.hash = '#login';
    } catch (error) {
        console.error('Logout error:', error);
        // Still redirect to login
        window.location.hash = '#login';
    }
}

async function loadHomeContent() {
    const container = document.getElementById('content-container');
    if (!container) return;

    container.innerHTML = '<div class="loading-container"><div class="spinner"></div><p>Loading content...</p></div>';

    try {
        const rows = await apiService.getHomeContent();
        renderContentRows(rows);
    } catch (error) {
        container.innerHTML = `<div class="error-message">Error loading content: ${error.message}</div>`;
        console.error('Error loading home content:', error);
    }
}

function renderContentRows(rows) {
    const container = document.getElementById('content-container');
    container.innerHTML = '';

    rows.forEach(row => {
        const rowDiv = document.createElement('div');
        rowDiv.className = 'content-row';

        const title = document.createElement('h2');
        title.className = 'row-title';
        title.textContent = row.title;

        const list = document.createElement('div');
        list.className = 'content-list';

        row.items.forEach(content => {
            const card = createContentCard(content);
            list.appendChild(card);
        });

        rowDiv.appendChild(title);
        rowDiv.appendChild(list);
        container.appendChild(rowDiv);
    });

    tvNavigation.updateFocusableElements();
}

function createContentCard(content) {
    const card = document.createElement('div');
    card.className = 'content-card';
    card.tabIndex = 0;

    card.innerHTML = `
        <img src="${content.posterUrl}" alt="${content.title}" onerror="this.src='placeholder.jpg'">
        <div class="content-card-info">
            <div class="content-card-title">${content.title}</div>
            <div class="content-card-meta">${content.releaseYear} • ${content.rating} ⭐</div>
        </div>
    `;

    card.addEventListener('click', () => {
        showContentDetail(content);
    });

    return card;
}

async function showContentDetail(content) {
    try {
        const fullContent = await apiService.getContentById(content.id);
        renderContentDetail(fullContent);
        window.location.hash = '#detail';
    } catch (error) {
        console.error('Error loading content detail:', error);
        alert('Failed to load content details');
    }
}

function renderContentDetail(content) {
    const container = document.getElementById('detail-content');
    container.innerHTML = `
        <img src="${content.backdropUrl}" alt="${content.title}" class="detail-backdrop">
        <div class="detail-info">
            <h2>${content.title}</h2>
            <div class="detail-meta">${content.releaseYear} • ${content.rating} ⭐ • ${content.genre}</div>
            <p class="detail-description">${content.description}</p>
            <button class="tv-button primary" id="play-button">Play</button>
        </div>
    `;

    document.getElementById('play-button').addEventListener('click', () => {
        playContent(content);
    });
}

function playContent(content) {
    window.location.hash = '#player';
    
    setTimeout(() => {
        tizenPlayer.play(content).catch(error => {
            console.error('Playback error:', error);
            alert('Failed to play content. Please try again.');
        });
    }, 100);
}

async function handleSearch() {
    const query = document.getElementById('search-input').value;
    const resultsDiv = document.getElementById('search-results');

    if (!query || query.length < 2) {
        resultsDiv.innerHTML = '';
        return;
    }

    try {
        const results = await apiService.searchContent(query);
        renderSearchResults(results);
    } catch (error) {
        console.error('Search error:', error);
        resultsDiv.innerHTML = '<div class="error-message">Search failed</div>';
    }
}

function renderSearchResults(results) {
    const resultsDiv = document.getElementById('search-results');
    
    if (results.length === 0) {
        resultsDiv.innerHTML = '<p>No results found</p>';
        return;
    }

    resultsDiv.innerHTML = results.map(content => `
        <div class="content-card" tabindex="0">
            <img src="${content.posterUrl}" alt="${content.title}">
            <div class="content-card-info">
                <div class="content-card-title">${content.title}</div>
                <div class="content-card-meta">${content.releaseYear}</div>
            </div>
        </div>
    `).join('');

    // Add click handlers
    resultsDiv.querySelectorAll('.content-card').forEach((card, index) => {
        card.addEventListener('click', () => {
            showContentDetail(results[index]);
        });
    });
}

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

