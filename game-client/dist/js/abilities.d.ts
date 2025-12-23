export interface Ability {
    id: number;
    name: string;
    description: string;
    class: string;
    unlock_level: number;
    ability_type: string;
    target_type: string;
    cooldown: number;
    mana_cost: number;
    base_damage: number;
    icon_url: string;
    animation_name: string;
    applies_buff?: string;
    applies_debuff?: string;
}
export declare function getAbilitiesByClass(className: string): Promise<Ability[]>;
export declare function getCharacterAbilities(characterId: number): Promise<{
    learned: Ability[];
    available: Ability[];
}>;
export declare function getAbilityTypeColor(type: string): string;
export declare function getTargetTypeIcon(targetType: string): string;
export declare function formatCooldown(seconds: number): string;
export declare function createAbilityCard(ability: Ability, isLearned?: boolean): string;
//# sourceMappingURL=abilities.d.ts.map