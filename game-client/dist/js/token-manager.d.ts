/**
 * Validate wallet consistency between localStorage and MetaMask
 * @returns {Promise<boolean>}
 */
declare function validateWalletConsistency(): Promise<boolean>;
/**
 * Ensure token is valid, auto-renew if necessary
 * @returns {Promise<string|null>} - Valid token or null if failed
 */
declare function ensureValidToken(): Promise<string | null>;
/**
 * Start background validation (polling every 30s)
 */
declare function startBackgroundValidation(): void;
/**
 * Handle page visibility change (refresh on focus)
 */
declare function setupVisibilityHandling(): void;
/**
 * Handle pageshow event (detect back/forward navigation)
 */
declare function setupPageShowHandling(): void;
/**
 * TOKEN MANAGER - Gestión robusta de tokens JWT
 * Maneja expiración, renovación automática y validación
 */
declare class TokenManager {
    TOKEN_EXPIRY_MS: number;
    /**
     * Get current token if valid, null if expired
     * @returns {string|null}
     */
    getToken(): string | null;
    /**
     * Set token with expiry tracking
     * @param {string} token - JWT token
     * @param {number} expiresInMs - Expiration time in milliseconds
     */
    setToken(token: string, expiresInMs?: number): void;
    /**
     * Check if token is expired
     * @returns {boolean}
     */
    isTokenExpired(): boolean;
    /**
     * Check if token is about to expire (within threshold)
     * @param {number} thresholdMs - Time threshold in milliseconds
     * @returns {boolean}
     */
    isTokenExpiringSoon(thresholdMs?: number): boolean;
    /**
     * Clear all token data
     */
    clearToken(): void;
    /**
     * Get token info for debugging
     * @returns {Object}
     */
    getTokenInfo(): Object;
}
declare const tokenManager: TokenManager;
//# sourceMappingURL=token-manager.d.ts.map