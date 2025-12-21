// API Client -// Type definitions for API responses and models
import type { User } from "./types";
import type {
    AuthResponse,
    NonceResponse,
    HealthResponse,
    APIError
} from './types';

const API_BASE_URL = 'http://localhost:8080/api/v1';

// Type-safe API client class
export class APIClient {
    private baseURL: string;

    constructor(baseURL: string = API_BASE_URL) {
        this.baseURL = baseURL;
    }

    // Set authentication token
    setToken(token: string): void {
        localStorage.setItem('token', token);
    }

    // Clear authentication token
    clearToken(): void {
        localStorage.removeItem('token');
    }

    // Get authorization headers
    private getHeaders(includeAuth = true): HeadersInit {
        const headers: HeadersInit = {
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
    private async fetch<T>(
        endpoint: string,
        options: RequestInit = {},
        requiresAuth = true
    ): Promise<T> {
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
                    const error: APIError = await response.json();
                    if (error.error) errorMessage = error.error;
                } catch (e) {
                    // Response was not JSON, use default message
                }
                throw new Error(errorMessage);
            }

            return await response.json();
        } catch (error) {
            if (error instanceof Error) {
                throw error;
            }
            throw new Error('Unknown error occurred');
        }
    }

    // GET request
    async get<T>(endpoint: string, requiresAuth = true): Promise<T> {
        return this.fetch<T>(endpoint, { method: 'GET' }, requiresAuth);
    }

    // POST request
    async post<T>(
        endpoint: string,
        data?: any,
        requiresAuth = true
    ): Promise<T> {
        return this.fetch<T>(
            endpoint,
            {
                method: 'POST',
                body: data ? JSON.stringify(data) : undefined,
            },
            requiresAuth
        );
    }

    // PUT request
    async put<T>(
        endpoint: string,
        data?: any,
        requiresAuth = true
    ): Promise<T> {
        return this.fetch<T>(
            endpoint,
            {
                method: 'PUT',
                body: data ? JSON.stringify(data) : undefined,
            },
            requiresAuth
        );
    }

    // DELETE request
    async delete<T>(endpoint: string, requiresAuth = true): Promise<T> {
        return this.fetch<T>(endpoint, { method: 'DELETE' }, requiresAuth);
    }

    // Health check (public)
    async health(): Promise<HealthResponse> {
        return this.get<HealthResponse>('/health', false);
    }

    // Auth: Get nonce
    async getNonce(walletAddress: string): Promise<NonceResponse> {
        return this.post<NonceResponse>('/auth/nonce', { wallet_address: walletAddress }, false);
    }

    // Auth: Verify signature
    async verifySignature(
        walletAddress: string,
        signature: string
    ): Promise<AuthResponse> {
        const response = await this.post<AuthResponse>(
            '/auth/verify',
            {
                wallet_address: walletAddress,
                signature,
            },
            false
        );

        // Store token
        if (response.token) {
            this.setToken(response.token);
        }

        return response;
    }

    // Get user profile (requires auth)
    async getProfile(): Promise<User> {
        return this.get<User>('/auth/profile');
    }

    // Get user's characters (requires auth)
    async getMyCharacters(): Promise<{ characters: any[] }> {
        return this.get<{ characters: any[] }>('/characters');
    }

    // Get single character details
    async getCharacter(id: number): Promise<any> {
        return this.get<any>(`/characters/${id}`);
    }

    // Find a battle match (requires auth)
    async findMatch(characterIds: number[], betAmount: number): Promise<any> {
        return this.post<any>('/battles/matchmaking', { character_ids: characterIds, bet_amount: betAmount });
    }

    // Get battle status (requires auth)
    async getBattleStatus(battleId: string): Promise<any> {
        return this.get<any>(`/battles/${battleId}`);
    }

    // Logout (clear token on client side)
    async logout(): Promise<void> {
        this.clearToken();
    }
}


// Export singleton instance
export const apiClient = new APIClient();

// Export for convenience
export default apiClient;
