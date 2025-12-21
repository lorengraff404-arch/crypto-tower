// progression.ts - Client-side progression display logic
import { apiClient } from './api.js';

export interface ProgressionInfo {
    level: number;
    total_xp: number;
    current_level_xp: number;
    xp_for_next_level: number;
    progress_percent: number;
    rarity: string;
    evolution_stage: number;
    can_level_up: boolean;
}

// Get XP formula (matches backend)
export function getXPForLevel(level: number): number {
    if (level <= 1) return 0;
    return Math.floor(100 * Math.pow(level, 2.5));
}

// Get rarity for level
export function getRarityForLevel(level: number): string {
    if (level >= 95) return 'SSS';
    if (level >= 80) return 'SS';
    if (level >= 60) return 'S';
    if (level >= 40) return 'A';
    if (level >= 20) return 'B';
    return 'C';
}

// Get evolution stage for level
export function getEvolutionStage(level: number): number {
    if (level >= 75) return 3;
    if (level >= 50) return 2;
    if (level >= 25) return 1;
    return 0;
}

// Get evolution name
export function getEvolutionName(stage: number): string {
    const names = ['Base', 'Growth', 'Mature', 'Ultimate'];
    return names[stage] || 'Base';
}

// Get rarity color
export function getRarityColor(rarity: string): string {
    const colors: Record<string, string> = {
        'SSS': '#FFD700', // Gold
        'SS': '#E535AB',  // Pink
        'S': '#9333EA',   // Purple
        'A': '#3B82F6',   // Blue
        'B': '#10B981',   // Green
        'C': '#6B7280'    // Gray
    };
    return (colors[rarity] || colors['C']) as string;
}

// Create XP bar HTML
export function createXPBar(currentXP: number, requiredXP: number, percent: number): string {
    return `
        <div class="xp-bar-container">
            <div class="xp-bar-fill" style="width: ${percent}%"></div>
            <div class="xp-bar-text">${currentXP} / ${requiredXP} XP (${Math.floor(percent)}%)</div>
        </div>
    `;
}

// Create stat display with level scaling indicator
export function createStatDisplay(statName: string, current: number, baseValue: number): string {
    const bonus = current - baseValue;
    const bonusText = bonus > 0 ? ` (+${bonus})` : '';
    const color = bonus > 0 ? '#10B981' : '#6B7280';

    return `
        <div class="stat-item">
            <span class="stat-name">${statName}:</span>
            <span class="stat-value" style="color: ${color}">${current}${bonusText}</span>
        </div>
    `;
}

// Format large numbers
export function formatNumber(num: number): string {
    if (num >= 1000000) {
        return (num / 1000000).toFixed(1) + 'M';
    }
    if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'K';
    }
    return num.toString();
}

// Fetch progression info from API
export async function fetchProgressionInfo(characterId: number): Promise<ProgressionInfo> {
    return await apiClient.get<ProgressionInfo>(`/characters/${characterId}/progression`);
}

// Grant XP (called after battles, quests, etc)
export async function grantXP(characterId: number, xpAmount: number, source: string, difficulty: string = ''): Promise<any> {
    return await apiClient.post(`/characters/${characterId}/gain-xp`, {
        xp_gained: xpAmount,
        source: source,
        difficulty: difficulty
    });
}
