import { apiClient } from './api.js';
export const statusEffectsAPI = {
    async getCharacterEffects(characterId) {
        try {
            const response = await apiClient.get(`/characters/${characterId}/effects`);
            return response;
        }
        catch (error) {
            console.error('Failed to fetch character effects:', error);
            throw error;
        }
    },
    async getEffectDefinitions() {
        try {
            const response = await apiClient.get('/effects/definitions');
            return response;
        }
        catch (error) {
            console.error('Failed to fetch effect definitions:', error);
            throw error;
        }
    }
};
export function createEffectBadge(effect) {
    const isBuff = effect.effect_type === 'BUFF';
    const bgColor = isBuff ? 'rgba(16, 185, 129, 0.2)' : 'rgba(239, 68, 68, 0.2)';
    const borderColor = isBuff ? '#10b981' : '#ef4444';
    const textColor = isBuff ? '#6ee7b7' : '#fca5a5';
    let details = effect.description;
    if (effect.stacks > 1) {
        details += ` (${effect.stacks}x stacks)`;
    }
    return `
        <div class="effect-badge" style="
            display: inline-flex;
            align-items: center;
            gap: 6px;
            padding: 6px 12px;
            background: ${bgColor};
            border: 1px solid ${borderColor};
            border-radius: 20px;
            color: ${textColor};
            font-size: 13px;
            font-weight: 600;
            cursor: help;
            position: relative;
        " title="${details} • ${effect.turns_remaining} turns left">
            <span style="font-size: 16px;">${effect.icon}</span>
            <span>${effect.effect_name}</span>
            ${effect.stacks > 1 ? `<span style="background: ${borderColor}; color: white; padding: 1px 6px; border-radius: 10px; font-size: 10px;">${effect.stacks}</span>` : ''}
            <span style="opacity: 0.7; font-size: 11px; margin-left: 4px;">⏳ ${effect.turns_remaining}</span>
        </div>
    `;
}
export function renderActiveEffects(containerId, effects) {
    const container = document.getElementById(containerId);
    if (!container)
        return;
    if (effects.length === 0) {
        container.innerHTML = `<span style="color: #64748b; font-size: 14px; font-style: italic;">No active effects</span>`;
        return;
    }
    container.innerHTML = `<div style="display: flex; gap: 8px; flex-wrap: wrap;">
        ${effects.map(effect => createEffectBadge(effect)).join('')}
    </div>`;
}
//# sourceMappingURL=status-effects.js.map