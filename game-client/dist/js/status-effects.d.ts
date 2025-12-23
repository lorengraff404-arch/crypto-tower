export interface StatusEffect {
    id: number;
    effect_name: string;
    effect_type: 'BUFF' | 'DEBUFF';
    stacks: number;
    turns_remaining: number;
    icon: string;
    description: string;
    stat_modifier: number;
    damage_per_turn: number;
}
export interface EffectDefinition {
    Name: string;
    Type: 'BUFF' | 'DEBUFF';
    Icon: string;
    Description: string;
    Modifier: number;
    DamagePerTurn: number;
    DefaultDuration: number;
}
export declare const statusEffectsAPI: {
    getCharacterEffects(characterId: number): Promise<unknown>;
    getEffectDefinitions(): Promise<unknown>;
};
export declare function createEffectBadge(effect: StatusEffect): string;
export declare function renderActiveEffects(containerId: string, effects: StatusEffect[]): void;
//# sourceMappingURL=status-effects.d.ts.map