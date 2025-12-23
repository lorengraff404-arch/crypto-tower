const API_BASE_URL = 'http://localhost:8080/api/v1';
// Type-safe API client class
export class APIClient {
    baseURL;
    constructor(baseURL = API_BASE_URL) {
        this.baseURL = baseURL;
    }
    // Set authentication token
    setToken(token) {
        localStorage.setItem('token', token);
    }
    // Clear authentication token
    clearToken() {
        localStorage.removeItem('token');
    }
    // Get authorization headers
    getHeaders(includeAuth = true) {
        const headers = {
            'Content-Type': 'application/json',
        };
        if (includeAuth) {
            // Always read token from localStorage to get latest value
            const currentToken = localStorage.getItem('token');
            if (currentToken) {
                headers['Authorization'] = `Bearer ${currentToken}`;
            }
        }
        return headers;
    }
    // Generic fetch with error handling
    async fetch(endpoint, options = {}, requiresAuth = true) {
        // Cache-busting: agregar timestamp a la URL
        const cacheBuster = Date.now();
        const separator = endpoint.includes('?') ? '&' : '?';
        const url = `${this.baseURL}${endpoint}${separator}_t=${cacheBuster}`;
        const headers = this.getHeaders(requiresAuth);
        // Headers anti-cache
        const finalHeaders = {
            ...headers,
            'Cache-Control': 'no-cache, no-store, must-revalidate',
            'Pragma': 'no-cache',
            'Expires': '0',
            ...options.headers,
        };
        try {
            const response = await fetch(url, {
                ...options,
                headers: finalHeaders,
            });
            if (!response.ok) {
                let errorMessage = `HTTP ${response.status}: ${response.statusText}`;
                try {
                    const error = await response.json();
                    if (error.error)
                        errorMessage = error.error;
                }
                catch (e) {
                    // Response was not JSON, use default message
                }
                throw new Error(errorMessage);
            }
            return await response.json();
        }
        catch (error) {
            if (error instanceof Error) {
                throw error;
            }
            throw new Error('Unknown error occurred');
        }
    }
    // Compatibility wrapper
    async request(method, endpoint, data = null) {
        const options = { method };
        if (data) {
            options.body = JSON.stringify(data);
        }
        return this.fetch(endpoint, options);
    }

    // GET request
    async get(endpoint, requiresAuth = true) {
        return this.fetch(endpoint, { method: 'GET' }, requiresAuth);
    }
    // POST request
    async post(endpoint, data, requiresAuth = true) {
        return this.fetch(endpoint, {
            method: 'POST',
            body: data ? JSON.stringify(data) : undefined,
        }, requiresAuth);
    }
    // PUT request
    async put(endpoint, data, requiresAuth = true) {
        return this.fetch(endpoint, {
            method: 'PUT',
            body: data ? JSON.stringify(data) : undefined,
        }, requiresAuth);
    }
    // DELETE request
    async delete(endpoint, requiresAuth = true) {
        return this.fetch(endpoint, { method: 'DELETE' }, requiresAuth);
    }
    // Health check (public)
    async health() {
        return this.get('/health', false);
    }
    // Auth: Get nonce
    async getNonce(walletAddress) {
        return this.post('/auth/nonce', { wallet_address: walletAddress }, false);
    }
    // Auth: Verify signature
    async verifySignature(walletAddress, signature) {
        const response = await this.post('/auth/verify', {
            wallet_address: walletAddress,
            signature,
        }, false);
        // Store token
        if (response.token) {
            this.setToken(response.token);
        }
        return response;
    }
    // Get user profile (requires auth)
    async getProfile() {
        return this.get('/auth/profile');
    }
    // Get user's characters (requires auth)
    async getMyCharacters() {
        return this.get('/characters');
    }
    // Get single character details
    async getCharacter(id) {
        return this.get(`/characters/${id}`);
    }
    // Find a battle match (requires auth)
    async findMatch(characterIds, betAmount) {
        return this.post('/battles/matchmaking', { character_ids: characterIds, bet_amount: betAmount });
    }
    // Get battle status (requires auth)
    async getBattleStatus(battleId) {
        return this.get(`/battles/${battleId}`);
    }
    // Logout (clear token on client side)
    async logout() {
        this.clearToken();
    }
}
// Export singleton instance
export const apiClient = new APIClient();
// Export for convenience
export default apiClient;
//# sourceMappingURL=api.js.map