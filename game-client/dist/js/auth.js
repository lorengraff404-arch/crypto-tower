"use strict";
/**
 * CRYPTO TOWER DEFENSE - AUTHENTICATION UTILITY
 * Centralized authentication and session management
 * Version: 1.0.0
 */
const AUTH_CONFIG = {
    // Session expires after 24 hours
    SESSION_DURATION: 24 * 60 * 60 * 1000,
    // Check session every minute
    CHECK_INTERVAL: 60 * 1000,
    // Storage keys
    STORAGE_KEYS: {
        WALLET: 'wallet_address',
        SESSION_START: 'session_start',
        LAST_ACTIVITY: 'last_activity'
    },
    // Redirect URLs
    URLS: {
        LOGIN: '/index.html',
        GAME: '/game.html'
    }
};
class AuthManager {
    constructor() {
        this.checkInterval = null;
        this.walletCheckInterval = null;
        this.currentWallet = null;
    }
    /**
     * Initialize authentication system
     */
    init() {
        console.log('[Auth] Initializing authentication system...');
        // Start session monitoring
        this.startSessionMonitoring();
        // Start wallet change detection
        this.startWalletMonitoring();
        // Update activity on user interaction
        this.setupActivityTracking();
        console.log('[Auth] Authentication system initialized');
    }
    /**
     * Check if user is authenticated
     * NOTE: Solo verifica wallet. Token validation se hace en requireAuth()
     * @returns {boolean}
     */
    isAuthenticated() {
        const wallet = this.getWallet();
        return wallet !== null && wallet !== '';
    }
    /**
     * Check if session is valid
     * @returns {boolean}
     */
    isSessionValid() {
        const sessionStart = localStorage.getItem(AUTH_CONFIG.STORAGE_KEYS.SESSION_START);
        if (!sessionStart) {
            return false;
        }
        const now = Date.now();
        const elapsed = now - parseInt(sessionStart);
        return elapsed < AUTH_CONFIG.SESSION_DURATION;
    }
    /**
     * Get current wallet address
     * @returns {string|null}
     */
    getWallet() {
        return localStorage.getItem(AUTH_CONFIG.STORAGE_KEYS.WALLET);
    }
    /**
     * Set wallet and start session
     * @param {string} wallet
     */
    setWallet(wallet) {
        const now = Date.now();
        localStorage.setItem(AUTH_CONFIG.STORAGE_KEYS.WALLET, wallet);
        localStorage.setItem(AUTH_CONFIG.STORAGE_KEYS.SESSION_START, now.toString());
        localStorage.setItem(AUTH_CONFIG.STORAGE_KEYS.LAST_ACTIVITY, now.toString());
        this.currentWallet = wallet;
        console.log('[Auth] Session started for wallet:', wallet.substring(0, 10) + '...');
    }
    /**
     * Update last activity timestamp
     */
    updateActivity() {
        const now = Date.now();
        localStorage.setItem(AUTH_CONFIG.STORAGE_KEYS.LAST_ACTIVITY, now.toString());
    }
    /**
     * Logout user and clear session
     * @param {string} reason - Reason for logout
     */
    logout(reason = 'Manual logout') {
        console.log('[Auth] Logging out:', reason);
        // Clear all auth data (wallet AND token)
        localStorage.removeItem(AUTH_CONFIG.STORAGE_KEYS.WALLET);
        localStorage.removeItem(AUTH_CONFIG.STORAGE_KEYS.SESSION_START);
        localStorage.removeItem(AUTH_CONFIG.STORAGE_KEYS.LAST_ACTIVITY);
        localStorage.removeItem('token'); // Clear JWT token
        // Stop monitoring
        this.stopSessionMonitoring();
        this.stopWalletMonitoring();
        // Show notification
        this.showNotification('Session ended: ' + reason, 'warning');
        // Redirect to login after short delay
        setTimeout(() => {
            window.location.href = AUTH_CONFIG.URLS.LOGIN;
        }, 1500);
    }
    /**
     * Require authentication for current page (ASYNC)
     * Redirects to login if not authenticated
     * @returns {Promise<boolean>}
     */
    async requireAuth() {
        // 1. Check wallet
        if (!this.isAuthenticated()) {
            console.log('[Auth] Not authenticated, redirecting to login...');
            this.showNotification('Please connect your wallet to continue', 'error');
            setTimeout(() => {
                window.location.href = AUTH_CONFIG.URLS.LOGIN;
            }, 1500);
            return false;
        }
        // 2. Check session validity
        if (!this.isSessionValid()) {
            console.log('[Auth] Session expired');
            this.logout('Session expired');
            return false;
        }
        // 3. Validate wallet consistency (if token-manager loaded)
        if (window.validateWalletConsistency) {
            const isConsistent = await window.validateWalletConsistency();
            if (!isConsistent) {
                return false; // validateWalletConsistency handles logout
            }
        }
        // 4. Ensure token is valid (if token-manager loaded)
        if (window.ensureValidToken) {
            const token = await window.ensureValidToken();
            if (!token) {
                return false; // ensureValidToken handles logout on failure
            }
        }
        console.log('[Auth] Authentication verified ✅');
        return true;
    }
    /**
     * Start session monitoring
     */
    startSessionMonitoring() {
        // Check immediately
        this.cleanExpiredSession();
        // Then check every minute
        this.checkInterval = setInterval(() => {
            this.checkSession();
        }, AUTH_CONFIG.CHECK_INTERVAL);
    }
    /**
     * Stop session monitoring
     */
    stopSessionMonitoring() {
        if (this.checkInterval) {
            clearInterval(this.checkInterval);
            this.checkInterval = null;
        }
    }
    /**
     * Clean expired session silently (no notification)
     * Used on page load to clear old sessions
     */
    cleanExpiredSession() {
        if (!this.isAuthenticated()) {
            return;
        }
        if (!this.isSessionValid()) {
            console.log('[Auth] Cleaning expired session silently');
            // Clear session data without notification or redirect
            localStorage.removeItem(AUTH_CONFIG.STORAGE_KEYS.WALLET);
            localStorage.removeItem(AUTH_CONFIG.STORAGE_KEYS.SESSION_START);
            localStorage.removeItem(AUTH_CONFIG.STORAGE_KEYS.LAST_ACTIVITY);
            localStorage.removeItem('token'); // Clear JWT token
        }
    }
    /**
     * Check session validity (with notification)
     * Used during active session monitoring
     */
    checkSession() {
        if (!this.isAuthenticated()) {
            return;
        }
        if (!this.isSessionValid()) {
            this.logout('Session expired (24 hours)');
        }
    }
    /**
     * Start wallet change monitoring
     */
    startWalletMonitoring() {
        this.currentWallet = this.getWallet();
        // Check for MetaMask
        if (typeof window.ethereum !== 'undefined') {
            // Listen for account changes
            window.ethereum.on('accountsChanged', (accounts) => {
                this.handleWalletChange(accounts);
            });
            // Also poll every 5 seconds as backup
            this.walletCheckInterval = setInterval(() => {
                this.checkWalletChange();
            }, 5000);
        }
    }
    /**
     * Stop wallet monitoring
     */
    stopWalletMonitoring() {
        if (this.walletCheckInterval) {
            clearInterval(this.walletCheckInterval);
            this.walletCheckInterval = null;
        }
    }
    /**
     * Handle wallet change event
     * @param {Array} accounts
     */
    handleWalletChange(accounts) {
        if (accounts.length === 0) {
            // User disconnected wallet
            this.logout('Wallet disconnected');
        }
        else {
            const newWallet = accounts[0];
            const storedWallet = this.getWallet();
            if (storedWallet && newWallet.toLowerCase() !== storedWallet.toLowerCase()) {
                // Wallet changed
                this.logout('Wallet changed');
            }
        }
    }
    /**
     * Check if wallet changed (polling method)
     */
    async checkWalletChange() {
        if (typeof window.ethereum === 'undefined') {
            return;
        }
        try {
            const accounts = await window.ethereum.request({ method: 'eth_accounts' });
            this.handleWalletChange(accounts);
        }
        catch (error) {
            console.error('[Auth] Error checking wallet:', error);
        }
    }
    /**
     * Setup activity tracking
     */
    setupActivityTracking() {
        // Update activity on mouse move, click, keypress
        const events = ['mousedown', 'keydown', 'scroll', 'touchstart'];
        events.forEach(event => {
            document.addEventListener(event, () => {
                this.updateActivity();
            }, { passive: true });
        });
    }
    /**
     * Show notification to user
     * @param {string} message
     * @param {string} type - 'success', 'error', 'warning', 'info'
     */
    showNotification(message, type = 'info') {
        // Try to use existing notification system if available
        if (typeof showNotification === 'function') {
            showNotification(message, type);
            return;
        }
        // Fallback to console
        console.log(`[Auth] ${type.toUpperCase()}: ${message}`);
        // Simple alert as last resort
        if (type === 'error' || type === 'warning') {
            alert(message);
        }
    }
    /**
     * Get session info
     * @returns {Object}
     */
    getSessionInfo() {
        const wallet = this.getWallet();
        const sessionStart = localStorage.getItem(AUTH_CONFIG.STORAGE_KEYS.SESSION_START);
        const lastActivity = localStorage.getItem(AUTH_CONFIG.STORAGE_KEYS.LAST_ACTIVITY);
        if (!wallet || !sessionStart) {
            return null;
        }
        const now = Date.now();
        const sessionAge = now - parseInt(sessionStart);
        const timeRemaining = AUTH_CONFIG.SESSION_DURATION - sessionAge;
        return {
            wallet: wallet,
            sessionStart: new Date(parseInt(sessionStart)),
            lastActivity: new Date(parseInt(lastActivity)),
            sessionAge: sessionAge,
            timeRemaining: timeRemaining,
            isValid: this.isSessionValid()
        };
    }
}
// Create global instance
const authManager = new AuthManager();
// Auto-initialize on page load
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        authManager.init();
    });
}
else {
    authManager.init();
}
// Export for use in other scripts
window.authManager = authManager;
// Convenience functions (async wrapper for requireAuth)
window.requireAuth = () => authManager.requireAuth();
window.logout = (reason) => authManager.logout(reason);
window.getSessionInfo = () => authManager.getSessionInfo();
//# sourceMappingURL=auth.js.map