/**
 * TOKEN MANAGER - Gestión robusta de tokens JWT
 * Maneja expiración, renovación automática y validación
 */

class TokenManager {
    constructor() {
        this.TOKEN_EXPIRY_MS = 24 * 60 * 60 * 1000; // 24 horas
    }

    /**
     * Get current token if valid, null if expired
     * @returns {string|null}
     */
    getToken() {
        const token = localStorage.getItem('token');
        if (!token) return null;

        const expiresAt = localStorage.getItem('token_expires_at');

        // Check expiry
        if (expiresAt && Date.now() > parseInt(expiresAt)) {
            console.log('[TokenManager] Token expired, clearing...');
            this.clearToken();
            return null;
        }

        return token;
    }

    /**
     * Set token with expiry tracking
     * @param {string} token - JWT token
     * @param {number} expiresInMs - Expiration time in milliseconds
     */
    setToken(token, expiresInMs = null) {
        const now = Date.now();
        const expiry = expiresInMs || this.TOKEN_EXPIRY_MS;

        localStorage.setItem('token', token);
        localStorage.setItem('token_issued_at', now.toString());
        localStorage.setItem('token_expires_at', (now + expiry).toString());

        console.log('[TokenManager] Token set, expires in', Math.floor(expiry / 1000 / 60), 'minutes');
    }

    /**
     * Check if token is expired
     * @returns {boolean}
     */
    isTokenExpired() {
        const expiresAt = localStorage.getItem('token_expires_at');
        if (!expiresAt) return true;

        return Date.now() > parseInt(expiresAt);
    }

    /**
     * Check if token is about to expire (within threshold)
     * @param {number} thresholdMs - Time threshold in milliseconds
     * @returns {boolean}
     */
    isTokenExpiringSoon(thresholdMs = 5 * 60 * 1000) {
        const expiresAt = localStorage.getItem('token_expires_at');
        if (!expiresAt) return true;

        const timeUntilExpiry = parseInt(expiresAt) - Date.now();
        return timeUntilExpiry < thresholdMs;
    }

    /**
     * Clear all token data
     */
    clearToken() {
        localStorage.removeItem('token');
        localStorage.removeItem('token_issued_at');
        localStorage.removeItem('token_expires_at');
        console.log('[TokenManager] Token cleared');
    }

    /**
     * Get token info for debugging
     * @returns {Object}
     */
    getTokenInfo() {
        const token = localStorage.getItem('token');
        const issuedAt = localStorage.getItem('token_issued_at');
        const expiresAt = localStorage.getItem('token_expires_at');

        if (!token) return { exists: false };

        const now = Date.now();
        const timeUntilExpiry = parseInt(expiresAt) - now;

        return {
            exists: true,
            isExpired: this.isTokenExpired(),
            issuedAt: new Date(parseInt(issuedAt)),
            expiresAt: new Date(parseInt(expiresAt)),
            timeUntilExpiryMinutes: Math.floor(timeUntilExpiry / 1000 / 60)
        };
    }
}

// Global instance
const tokenManager = new TokenManager();

/**
 * Validate wallet consistency between localStorage and MetaMask
 * @returns {Promise<boolean>}
 */
async function validateWalletConsistency() {
    const storedWallet = localStorage.getItem('wallet_address');

    if (!storedWallet) {
        console.log('[Wallet] No stored wallet');
        return false;
    }

    // Check MetaMask
    if (!window.ethereum) {
        console.warn('[Wallet] MetaMask not detected');
        return true; // Don't block if MetaMask not available
    }

    try {
        const accounts = await window.ethereum.request({
            method: 'eth_accounts'
        });

        if (accounts.length === 0) {
            console.warn('[Wallet] MetaMask disconnected');
            if (window.authManager) {
                window.authManager.logout('Wallet disconnected');
            }
            return false;
        }

        const currentWallet = accounts[0].toLowerCase();
        const storedWalletLower = storedWallet.toLowerCase();

        if (currentWallet !== storedWalletLower) {
            console.error('[Wallet] Wallet mismatch!', {
                stored: storedWalletLower,
                current: currentWallet
            });
            if (window.authManager) {
                window.authManager.logout('Wallet changed - please reconnect');
            }
            return false;
        }

        return true;
    } catch (error) {
        console.error('[Wallet] Validation error:', error);
        return false;
    }
}

/**
 * Ensure token is valid, auto-renew if necessary
 * @returns {Promise<string|null>} - Valid token or null if failed
 */
async function ensureValidToken() {
    const wallet = localStorage.getItem('wallet_address');

    // No wallet = force login
    if (!wallet) {
        console.log('[Auth] No wallet found, redirecting to login...');
        window.location.href = '/index.html';
        return null;
    }

    // Check current token
    const token = tokenManager.getToken();

    if (token && !tokenManager.isTokenExpired()) {
        // Token is valid
        return token;
    }

    // Token expired or missing - need to renew
    console.log('[Auth] Token expired/missing, attempting to renew...');

    try {
        // Import apiClient if not already available
        if (!window.gameAPI) {
            console.error('[Auth] gameAPI not available for token renewal');
            return null;
        }

        // Get new nonce
        const nonceResponse = await window.gameAPI.getNonce(wallet);

        // Request user signature
        if (!window.ethereum) {
            throw new Error('MetaMask not available');
        }

        const signature = await window.ethereum.request({
            method: 'personal_sign',
            params: [nonceResponse.message, wallet]
        });

        // Verify signature and get new token
        const authResponse = await window.gameAPI.verifySignature(wallet, signature);

        if (authResponse.token) {
            tokenManager.setToken(authResponse.token);
            console.log('[Auth] Token renewed successfully ✅');
            return authResponse.token;
        }

        throw new Error('No token in response');

    } catch (error) {
        console.error('[Auth] Token renewal failed:', error);

        // If renewal fails, logout
        if (window.authManager) {
            window.authManager.logout('Failed to renew authentication');
        }

        return null;
    }
}

/**
 * Start background validation (polling every 30s)
 */
function startBackgroundValidation() {
    // Clear any existing interval
    if (window.authValidationInterval) {
        clearInterval(window.authValidationInterval);
    }

    window.authValidationInterval = setInterval(async () => {
        // Only run if authenticated
        if (!window.authManager || !window.authManager.isAuthenticated()) {
            return;
        }

        console.log('[Auth] Background validation check...');

        // Validate wallet consistency
        const isValid = await validateWalletConsistency();
        if (!isValid) {
            return; // validateWalletConsistency already handles logout
        }

        // Pre-renew token if expiring soon (5 min threshold)
        if (tokenManager.isTokenExpiringSoon(5 * 60 * 1000)) {
            console.log('[Auth] Token expiring soon, pre-renewing...');
            await ensureValidToken();
        }
    }, 30000); // Every 30 seconds

    console.log('[Auth] Background validation started (30s interval)');
}

/**
 * Handle page visibility change (refresh on focus)
 */
function setupVisibilityHandling() {
    document.addEventListener('visibilitychange', async () => {
        if (!document.hidden) {
            // Page became visible
            if (window.authManager && window.authManager.isAuthenticated()) {
                console.log('[Auth] Tab focused, validating session...');

                await validateWalletConsistency();
                await ensureValidToken();
            }
        }
    });

    console.log('[Auth] Visibility change handler registered');
}

/**
 * Handle pageshow event (detect back/forward navigation)
 */
function setupPageShowHandling() {
    window.addEventListener('pageshow', (event) => {
        if (event.persisted) {
            // Page loaded from cache (back/forward navigation)
            console.log('[Auth] Page loaded from cache, forcing reload...');
            location.reload(true);
        }
    });

    console.log('[Auth] Pageshow handler registered');
}

// Auto-initialize on load
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        startBackgroundValidation();
        setupVisibilityHandling();
        setupPageShowHandling();
    });
} else {
    startBackgroundValidation();
    setupVisibilityHandling();
    setupPageShowHandling();
}

// Export to global scope
window.tokenManager = tokenManager;
window.validateWalletConsistency = validateWalletConsistency;
window.ensureValidToken = ensureValidToken;
