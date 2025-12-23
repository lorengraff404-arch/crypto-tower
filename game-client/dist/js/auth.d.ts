declare namespace AUTH_CONFIG {
    let SESSION_DURATION: number;
    let CHECK_INTERVAL: number;
    namespace STORAGE_KEYS {
        let WALLET: string;
        let SESSION_START: string;
        let LAST_ACTIVITY: string;
    }
    namespace URLS {
        let LOGIN: string;
        let GAME: string;
    }
}
declare class AuthManager {
    checkInterval: NodeJS.Timeout | null;
    walletCheckInterval: NodeJS.Timeout | null;
    currentWallet: string | null;
    /**
     * Initialize authentication system
     */
    init(): void;
    /**
     * Check if user is authenticated
     * NOTE: Solo verifica wallet. Token validation se hace en requireAuth()
     * @returns {boolean}
     */
    isAuthenticated(): boolean;
    /**
     * Check if session is valid
     * @returns {boolean}
     */
    isSessionValid(): boolean;
    /**
     * Get current wallet address
     * @returns {string|null}
     */
    getWallet(): string | null;
    /**
     * Set wallet and start session
     * @param {string} wallet
     */
    setWallet(wallet: string): void;
    /**
     * Update last activity timestamp
     */
    updateActivity(): void;
    /**
     * Logout user and clear session
     * @param {string} reason - Reason for logout
     */
    logout(reason?: string): void;
    /**
     * Require authentication for current page (ASYNC)
     * Redirects to login if not authenticated
     * @returns {Promise<boolean>}
     */
    requireAuth(): Promise<boolean>;
    /**
     * Start session monitoring
     */
    startSessionMonitoring(): void;
    /**
     * Stop session monitoring
     */
    stopSessionMonitoring(): void;
    /**
     * Clean expired session silently (no notification)
     * Used on page load to clear old sessions
     */
    cleanExpiredSession(): void;
    /**
     * Check session validity (with notification)
     * Used during active session monitoring
     */
    checkSession(): void;
    /**
     * Start wallet change monitoring
     */
    startWalletMonitoring(): void;
    /**
     * Stop wallet monitoring
     */
    stopWalletMonitoring(): void;
    /**
     * Handle wallet change event
     * @param {Array} accounts
     */
    handleWalletChange(accounts: any[]): void;
    /**
     * Check if wallet changed (polling method)
     */
    checkWalletChange(): Promise<void>;
    /**
     * Setup activity tracking
     */
    setupActivityTracking(): void;
    /**
     * Show notification to user
     * @param {string} message
     * @param {string} type - 'success', 'error', 'warning', 'info'
     */
    showNotification(message: string, type?: string): void;
    /**
     * Get session info
     * @returns {Object}
     */
    getSessionInfo(): Object;
}
declare const authManager: AuthManager;
//# sourceMappingURL=auth.d.ts.map