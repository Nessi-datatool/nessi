// Authentication functions for the Nessi Monitoring Dashboard

// Check if user is authenticated
function isAuthenticated() {
    return localStorage.getItem('auth_token') !== null;
}

// Get the current user's username
function getCurrentUser() {
    return localStorage.getItem('username');
}

// Get the current user's role
function getUserRole() {
    return localStorage.getItem('role');
}

// Get the authentication token
function getAuthToken() {
    return localStorage.getItem('auth_token');
}

// Logout the current user
function logout() {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('username');
    localStorage.removeItem('role');
    window.location.href = '/login';
}

// Add authorization headers to fetch requests
function fetchWithAuth(url, options = {}) {
    // Default options
    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json'
        }
    };

    // Merge options
    const mergedOptions = {
        ...defaultOptions,
        ...options,
        headers: {
            ...defaultOptions.headers,
            ...options.headers
        }
    };

    // Add authorization header if authenticated
    if (isAuthenticated()) {
        mergedOptions.headers['Authorization'] = `Bearer ${getAuthToken()}`;
    }

    // Make the request
    return fetch(url, mergedOptions);
}

// Initialize authentication UI
function initAuth() {
    const userMenu = document.getElementById('user-menu');
    const loginMenu = document.getElementById('login-menu');
    const usernameDisplay = document.getElementById('username-display');
    const logoutLink = document.getElementById('logout-link');

    if (isAuthenticated()) {
        // User is authenticated
        userMenu.style.display = 'block';
        loginMenu.style.display = 'none';
        usernameDisplay.textContent = getCurrentUser();
        
        // Add logout handler
        logoutLink.addEventListener('click', function(e) {
            e.preventDefault();
            logout();
        });
    } else {
        // User is not authenticated
        userMenu.style.display = 'none';
        loginMenu.style.display = 'block';
    }
}

// Check if a user has a specific role
function hasRole(role) {
    return getUserRole() === role;
}

// Check if the current user is an admin
function isAdmin() {
    return hasRole('admin');
}

// Initialize admin features
function initAdminFeatures() {
    // Only show admin features if the user is an admin
    if (isAdmin()) {
        // Show admin elements
        const adminElements = document.querySelectorAll('.admin-only');
        adminElements.forEach(element => {
            element.style.display = 'block';
        });
    }
}

// Initialize authentication when the page loads
document.addEventListener('DOMContentLoaded', function() {
    initAuth();
    initAdminFeatures();
});
