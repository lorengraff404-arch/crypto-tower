export declare function showLoadingSpinner(containerId?: string): void;
export declare function hideLoadingSpinner(): void;
type ToastType = 'success' | 'error' | 'warning' | 'info';
export declare function showToast(message: string, type?: ToastType, duration?: number): void;
export declare function handleAPIError(error: unknown): void;
export declare function setupOfflineDetection(): void;
export declare function setupGlobalErrorHandlers(): void;
export declare function debounce<T extends (...args: any[]) => any>(func: T, wait: number): (...args: Parameters<T>) => void;
export declare function withRequestLock<T>(key: string, fn: () => Promise<T>): Promise<T | null>;
export declare function showEmptyState(containerId: string, message: string, icon?: string): void;
export declare function confirmAction(message: string): Promise<boolean>;
declare const _default: {
    showLoadingSpinner: typeof showLoadingSpinner;
    hideLoadingSpinner: typeof hideLoadingSpinner;
    showToast: typeof showToast;
    handleAPIError: typeof handleAPIError;
    setupOfflineDetection: typeof setupOfflineDetection;
    setupGlobalErrorHandlers: typeof setupGlobalErrorHandlers;
    debounce: typeof debounce;
    withRequestLock: typeof withRequestLock;
    showEmptyState: typeof showEmptyState;
    confirmAction: typeof confirmAction;
};
export default _default;
//# sourceMappingURL=ui-utils.d.ts.map