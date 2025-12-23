// Common UI utilities - Loading, toasts, error handling
// ============= Loading States =============
export function showLoadingSpinner(containerId = 'main-content') {
    const container = document.getElementById(containerId);
    if (!container)
        return;
    const spinner = document.createElement('div');
    spinner.id = 'loading-spinner';
    spinner.className = 'loading-spinner-overlay';
    spinner.innerHTML = `
    <div class="loading-spinner">
      <div class="spinner-ring"></div>
      <p>Loading...</p>
    </div>
  `;
    container.appendChild(spinner);
}
export function hideLoadingSpinner() {
    const spinner = document.getElementById('loading-spinner');
    if (spinner) {
        spinner.remove();
    }
}
export function showToast(message, type = 'info', duration = 3000) {
    // Create toast container if it doesn't exist
    let container = document.getElementById('toast-container');
    if (!container) {
        container = document.createElement('div');
        container.id = 'toast-container';
        container.className = 'toast-container';
        document.body.appendChild(container);
    }
    // Create toast element
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    const icon = getToastIcon(type);
    toast.innerHTML = `
    <span class="toast-icon">${icon}</span>
    <span class="toast-message">${message}</span>
    <button class="toast-close">&times;</button>
  `;
    // Add to container
    container.appendChild(toast);
    // Animate in
    setTimeout(() => toast.classList.add('toast-show'), 10);
    // Close button functionality
    const closeBtn = toast.querySelector('.toast-close');
    if (closeBtn) {
        closeBtn.addEventListener('click', () => removeToast(toast));
    }
    // Auto-remove after duration
    setTimeout(() => removeToast(toast), duration);
}
function removeToast(toast) {
    toast.classList.remove('toast-show');
    setTimeout(() => toast.remove(), 300);
}
function getToastIcon(type) {
    const icons = {
        success: '✓',
        error: '✕',
        warning: '⚠',
        info: 'ℹ'
    };
    return icons[type];
}
// ============= Error Handling =============
export function handleAPIError(error) {
    console.error('API Error:', error);
    if (error instanceof Error) {
        if (error.message.includes('401') || error.message.includes('Unauthorized')) {
            showToast('Session expired. Please reconnect your wallet.', 'error', 5000);
            setTimeout(() => {
                window.location.href = 'index.html';
            }, 2000);
            return;
        }
        if (error.message.includes('429')) {
            showToast('Too many requests. Please wait a moment.', 'warning', 5000);
            return;
        }
        if (error.message.includes('404')) {
            showToast('Resource not found', 'error');
            return;
        }
        if (error.message.includes('500')) {
            showToast('Server error. Please try again later.', 'error', 5000);
            return;
        }
        // Generic error
        showToast(error.message || 'Something went wrong', 'error');
    }
    else {
        showToast('Unknown error occurred', 'error');
    }
}
// ============= Offline Detection =============
export function setupOfflineDetection() {
    window.addEventListener('online', () => {
        showToast('Connection restored', 'success');
        // Optionally reload page
        setTimeout(() => location.reload(), 1000);
    });
    window.addEventListener('offline', () => {
        showToast('No internet connection', 'error', 0); // Stay visible
    });
    // Check initial state
    if (!navigator.onLine) {
        showToast('You are offline', 'error', 0);
    }
}
// ============= Global Error Handler =============
export function setupGlobalErrorHandlers() {
    // Unhandled promise rejections
    window.addEventListener('unhandledrejection', (event) => {
        console.error('Unhandled promise rejection:', event.reason);
        handleAPIError(event.reason);
    });
    // General errors
    window.addEventListener('error', (event) => {
        console.error('Global error:', event.error);
        showToast('An unexpected error occurred', 'error');
    });
}
// ============= Debouncing =============
export function debounce(func, wait) {
    let timeout = null;
    return function executedFunction(...args) {
        const later = () => {
            timeout = null;
            func(...args);
        };
        if (timeout)
            clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}
// ============= Request Debouncing =============
const requestLocks = new Map();
export async function withRequestLock(key, fn) {
    if (requestLocks.get(key)) {
        console.warn(`Request ${key} already in progress`);
        return null;
    }
    requestLocks.set(key, true);
    try {
        return await fn();
    }
    finally {
        requestLocks.delete(key);
    }
}
// ============= Empty State Display =============
export function showEmptyState(containerId, message, icon = '📭') {
    const container = document.getElementById(containerId);
    if (!container)
        return;
    container.innerHTML = `
    <div class="empty-state">
      <div class="empty-state-icon">${icon}</div>
      <p class="empty-state-message">${message}</p>
    </div>
  `;
}
// ============= Confirmation Dialogs =============
export async function confirmAction(message) {
    // For now, use native confirm, but could be replaced with custom modal
    return window.confirm(message);
}
// Export all utilities
export default {
    showLoadingSpinner,
    hideLoadingSpinner,
    showToast,
    handleAPIError,
    setupOfflineDetection,
    setupGlobalErrorHandlers,
    debounce,
    withRequestLock,
    showEmptyState,
    confirmAction
};
//# sourceMappingURL=ui-utils.js.map