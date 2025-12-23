/**
 * CRYPTO TOWER DEFENSE - MODAL UTILITIES
 * JavaScript utilities for modal/popup management
 * Version: 1.0.0
 */
declare class ModalManager {
    activeModal: HTMLDivElement | null;
    activeToast: HTMLDivElement | null;
    /**
     * Show a modal dialog
     * @param {Object} options - Modal configuration
     * @param {string} options.title - Modal title
     * @param {string} options.content - Modal content (HTML)
     * @param {Array} options.buttons - Array of button configs
     * @param {Function} options.onClose - Callback when modal closes
     */
    showModal(options: {
        title: string;
        content: string;
        buttons: any[];
        onClose: Function;
    }): void;
    /**
     * Close active modal
     */
    closeModal(): void;
    /**
     * Show confirmation dialog
     * @param {Object} options - Confirm configuration
     * @param {string} options.title - Dialog title
     * @param {string} options.message - Confirmation message
     * @param {string} options.type - 'warning', 'danger', 'info'
     * @param {Function} options.onConfirm - Callback on confirm
     * @param {Function} options.onCancel - Callback on cancel
     */
    confirm(options: {
        title: string;
        message: string;
        type: string;
        onConfirm: Function;
        onCancel: Function;
    }): void;
    /**
     * Show notification toast
     * @param {string} message - Notification message
     * @param {string} type - 'success', 'error', 'warning', 'info'
     * @param {number} duration - Duration in ms (default 4000)
     */
    showNotification(message: string, type?: string, duration?: number): void;
    /**
     * Show alert (simple modal with OK button)
     * @param {string} title - Alert title
     * @param {string} message - Alert message
     * @param {string} type - 'success', 'error', 'warning', 'info'
     */
    alert(title: string, message: string, type?: string): void;
}
declare const modalManager: ModalManager;
//# sourceMappingURL=modals.d.ts.map