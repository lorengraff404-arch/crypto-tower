// abilities.ts - Frontend ability management
import { apiClient } from './api.js';

export interface Ability {
    id: number;
    name: string;
    description: string;
    class: string;
    unlock_level: number;
    ability_type: string; // ACTIVE, PASSIVE, ULTIMATE
    target_type: string;
    cooldown: number;
    mana_cost: number;
    base_damage: number;
    icon_url: string;
    animation_name: string;
    applies_buff?: string;
    applies_debuff?: string;
}

// Fetch all abilities for a class
export async function getAbilitiesByClass(className: string): Promise<Ability[]> {
    const response = await apiClient.get<{ abilities: Ability[] }>(`/abilities?class=${className}`);
    return response.abilities;
}

// Fetch learned abilities for a character
export async function getCharacterAbilities(characterId: number): Promise<{ learned: Ability[], available: Ability[] }> {
    const response = await apiClient.get<{ learned: Ability[], available: Ability[] }>(`/characters/${characterId}/abilities`);
    return { learned: response.learned, available: response.available };
}

// Get type color for ability type
export function getAbilityTypeColor(type: string): string {
    const colors: Record<string, string> = {
        'ACTIVE': '#3B82F6',    // Blue
        'PASSIVE': '#10B981',   // Green
        'ULTIMATE': '#EF4444'   // Red
    };
    return colors[type] ?? '#6B7280';
}

// Get icon for target type
export function getTargetTypeIcon(targetType: string): string {
    const icons: Record<string, string> = {
        'SELF': '👤',
        'SINGLE_ENEMY': '🎯',
        'AOE': '💥',
        'ALL_ALLIES': '👥',
        'CHAIN': '⚡'
    };
    return icons[targetType] ?? '❓';
}

// Format cooldown display
export function formatCooldown(seconds: number): string {
    if (seconds < 60) return `${seconds}s`;
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return secs > 0 ? `${mins}m ${secs}s` : `${mins}m`;
}

// Create ability card HTML (compact version)
export function createAbilityCard(ability: Ability, isLearned: boolean = false): string {
    const typeColor = getAbilityTypeColor(ability.ability_type);
    const targetIcon = getTargetTypeIcon(ability.target_type);
    const learnedBadge = isLearned ? '<span class="learned-badge">✓ Learned</span>' : '<span class="locked-badge">🔒 Locked</span>';

    return `
        <div class="ability-card" data-ability-id="${ability.id}">
            <div class="ability-header">
                <span class="ability-icon">${targetIcon}</span>
                <span class="ability-name">${ability.name}</span>
                ${learnedBadge}
            </div>
            <div class="ability-type" style="background: ${typeColor}">
                ${ability.ability_type}
            </div>
            <p class="ability-desc">${ability.description}</p>
            <div class="ability-stats">
                <span class="stat-item">⚡ ${ability.mana_cost} Mana</span>
                <span class="stat-item">⏱️ ${formatCooldown(ability.cooldown)}</span>
                ${ability.base_damage > 0 ? `<span class="stat-item">⚔️ ${ability.base_damage} Dmg</span>` : ''}
            </div>
            <div class="ability-unlock">Unlock at Level ${ability.unlock_level}</div>
        </div>
    `;
}
