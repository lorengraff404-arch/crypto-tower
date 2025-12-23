"use strict";
/**
 * CRYPTO TOWER DEFENSE - MODAL UTILITIES
 * JavaScript utilities for modal/popup management
 * Version: 1.0.0
 */
class ModalManager {
    constructor() {
        this.activeModal = null;
        this.activeToast = null;
    }
    /**
     * Show a modal dialog
     * @param {Object} options - Modal configuration
     * @param {string} options.title - Modal title
     * @param {string} options.content - Modal content (HTML)
     * @param {Array} options.buttons - Array of button configs
     * @param {Function} options.onClose - Callback when modal closes
     */
    showModal(options) {
        // Remove existing modal if any
        this.closeModal();
        // Create modal overlay
        const overlay = document.createElement('div');
        overlay.className = 'modal-overlay active';
        overlay.id = 'dynamicModal';
        // Create modal content
        const modal = document.createElement('div');
        modal.className = 'modal-content';
        // Header
        const header = document.createElement('div');
        header.className = 'modal-header';
        header.innerHTML = `
            <h2 class="modal-title">${options.title || 'Notice'}</h2>
            <button class="modal-close" onclick="modalManager.closeModal()">×</button>
        `;
        // Body
        const body = document.createElement('div');
        body.className = 'modal-body';
        body.innerHTML = options.content || '';
        // Footer with buttons
        const footer = document.createElement('div');
        footer.className = 'modal-footer';
        if (options.buttons && options.buttons.length > 0) {
            options.buttons.forEach(btn => {
                const button = document.createElement('button');
                button.className = `modal-btn modal-btn-${btn.type || 'secondary'}`;
                button.textContent = btn.text;
                button.onclick = () => {
                    if (btn.onClick)
                        btn.onClick();
                    if (btn.closeOnClick !== false)
                        this.closeModal();
                };
                footer.appendChild(button);
            });
        }
        else {
            // Default close button
            const closeBtn = document.createElement('button');
            closeBtn.className = 'modal-btn modal-btn-primary';
            closeBtn.textContent = 'OK';
            closeBtn.onclick = () => this.closeModal();
            footer.appendChild(closeBtn);
        }
        // Assemble modal
        modal.appendChild(header);
        modal.appendChild(body);
        modal.appendChild(footer);
        overlay.appendChild(modal);
        // Add to document
        document.body.appendChild(overlay);
        this.activeModal = overlay;
        // Close on overlay click
        overlay.addEventListener('click', (e) => {
            if (e.target === overlay) {
                this.closeModal();
            }
        });
        // Close on ESC key
        const escHandler = (e) => {
            if (e.key === 'Escape') {
                this.closeModal();
                document.removeEventListener('keydown', escHandler);
            }
        };
        document.addEventListener('keydown', escHandler);
        // Call onClose callback if provided
        if (options.onClose) {
            overlay.dataset.onClose = 'true';
            overlay.onCloseCallback = options.onClose;
        }
    }
    /**
     * Close active modal
     */
    closeModal() {
        if (this.activeModal) {
            // Call onClose callback if exists
            if (this.activeModal.onCloseCallback) {
                this.activeModal.onCloseCallback();
            }
            this.activeModal.remove();
            this.activeModal = null;
        }
    }
    /**
     * Show confirmation dialog
     * @param {Object} options - Confirm configuration
     * @param {string} options.title - Dialog title
     * @param {string} options.message - Confirmation message
     * @param {string} options.type - 'warning', 'danger', 'info'
     * @param {Function} options.onConfirm - Callback on confirm
     * @param {Function} options.onCancel - Callback on cancel
     */
    confirm(options) {
        const iconMap = {
            warning: '⚠️',
            danger: '❌',
            info: 'ℹ️'
        };
        const icon = iconMap[options.type] || iconMap.info;
        this.showModal({
            title: options.title || 'Confirm',
            content: `
                <div class="confirm-dialog">
                    <div class="confirm-icon ${options.type || 'info'}">${icon}</div>
                    <p style="text-align: center; font-size: 16px;">${options.message}</p>
                </div>
            `,
            buttons: [
                {
                    text: options.cancelText || 'Cancel',
                    type: 'secondary',
                    onClick: options.onCancel
                },
                {
                    text: options.confirmText || 'Confirm',
                    type: options.type === 'danger' ? 'danger' : 'primary',
                    onClick: options.onConfirm
                }
            ]
        });
    }
    /**
     * Show notification toast
     * @param {string} message - Notification message
     * @param {string} type - 'success', 'error', 'warning', 'info'
     * @param {number} duration - Duration in ms (default 4000)
     */
    showNotification(message, type = 'info', duration = 4000) {
        // Remove existing toast
        if (this.activeToast) {
            this.activeToast.remove();
        }
        // Icon map
        const iconMap = {
            success: '✓',
            error: '✗',
            warning: '⚠',
            info: 'ℹ'
        };
        // Create toast
        const toast = document.createElement('div');
        toast.className = `notification-toast ${type} active`;
        toast.innerHTML = `
            <div class="notification-content">
                <div class="notification-icon">${iconMap[type] || iconMap.info}</div>
                <div class="notification-message">${message}</div>
            </div>
        `;
        document.body.appendChild(toast);
        this.activeToast = toast;
        // Auto-remove after duration
        setTimeout(() => {
            toast.classList.remove('active');
            setTimeout(() => {
                if (toast.parentNode) {
                    toast.remove();
                }
                if (this.activeToast === toast) {
                    this.activeToast = null;
                }
            }, 300);
        }, duration);
    }
    /**
     * Show alert (simple modal with OK button)
     * @param {string} title - Alert title
     * @param {string} message - Alert message
     * @param {string} type - 'success', 'error', 'warning', 'info'
     */
    alert(title, message, type = 'info') {
        const iconMap = {
            success: '✓',
            error: '✗',
            warning: '⚠',
            info: 'ℹ'
        };
        this.showModal({
            title: title,
            content: `
                <div class="confirm-dialog">
                    <div class="confirm-icon ${type}">${iconMap[type] || iconMap.info}</div>
                    <p style="text-align: center;">${message}</p>
                </div>
            `,
            buttons: [
                {
                    text: 'OK',
                    type: 'primary'
                }
            ]
        });
    }
}
// Create global instance
const modalManager = new ModalManager();
// Export for use in other scripts
window.modalManager = modalManager;
// Convenience functions
window.showModal = (options) => modalManager.showModal(options);
window.showNotification = (message, type, duration) => modalManager.showNotification(message, type, duration);
window.showConfirm = (options) => modalManager.confirm(options);
window.showAlert = (title, message, type) => modalManager.alert(title, message, type);
//# sourceMappingURL=modals.js.map