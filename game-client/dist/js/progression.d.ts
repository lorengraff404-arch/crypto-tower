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
export declare function getXPForLevel(level: number): number;
export declare function getRarityForLevel(level: number): string;
export declare function getEvolutionStage(level: number): number;
export declare function getEvolutionName(stage: number): string;
export declare function getRarityColor(rarity: string): string;
export declare function createXPBar(currentXP: number, requiredXP: number, percent: number): string;
export declare function createStatDisplay(statName: string, current: number, baseValue: number): string;
export declare function formatNumber(num: number): string;
export declare function fetchProgressionInfo(characterId: number): Promise<ProgressionInfo>;
export declare function grantXP(characterId: number, xpAmount: number, source: string, difficulty?: string): Promise<any>;
//# sourceMappingURL=progression.d.ts.map