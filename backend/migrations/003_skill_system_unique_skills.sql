-- ============================================
-- PHASE 2B: 53 Unique Skills (Avoiding Duplicates)
-- ============================================
-- Current: 47 skills
-- Target: 100 skills
-- New unique skills: 53

-- ============================================
-- C RANK SKILLS (15 new)
-- ============================================

INSERT INTO abilities (name, description, class, unlock_level, ability_type, target_type, cooldown, mana_cost, base_damage, base_heal, duration_secs, rarity, damage_type, aoe_radius, synergy_tags, animation_name) VALUES
-- Warrior
('Sword Thrust', 'Quick piercing attack', 'Warrior', 1, 'ACTIVE', 'SINGLE_ENEMY', 2, 10, 55, 0, 0, 'C', 'physical', 0, ARRAY['warrior', 'physical'], 'sword_thrust'),
('Shield Slam', 'Bash with shield for damage', 'Warrior', 2, 'ACTIVE', 'SINGLE_ENEMY', 3, 12, 50, 0, 0, 'C', 'physical', 0, ARRAY['warrior', 'tank'], 'shield_slam'),
('Rage', 'Temporary attack boost', 'Warrior', 4, 'ACTIVE', 'SELF', 4, 15, 0, 0, 2, 'C', 'physical', 0, ARRAY['warrior', 'buff'], 'rage'),

-- Mage
('Spark', 'Small electric jolt', 'Mage', 1, 'ACTIVE', 'SINGLE_ENEMY', 2, 8, 40, 0, 0, 'C', 'magical', 0, ARRAY['mage', 'thunder'], 'spark'),
('Frost Touch', 'Minor ice damage', 'Mage', 2, 'ACTIVE', 'SINGLE_ENEMY', 3, 10, 45, 0, 0, 'C', 'magical', 0, ARRAY['mage', 'ice'], 'frost_touch'),
('Mana Tap', 'Restore small amount of mana', 'Mage', 4, 'ACTIVE', 'SELF', 6, 0, 0, 0, 0, 'C', 'magical', 0, ARRAY['mage', 'utility'], 'mana_tap'),

-- Tank
('Defensive Stance', 'Reduce damage taken', 'Tank', 1, 'ACTIVE', 'SELF', 4, 10, 0, 0, 2, 'C', 'physical', 0, ARRAY['tank', 'defense'], 'defensive_stance'),
('Bodyguard', 'Protect an ally', 'Tank', 3, 'ACTIVE', 'SINGLE_ALLY', 5, 15, 0, 0, 2, 'C', 'physical', 0, ARRAY['tank', 'support'], 'bodyguard'),

-- Archer
('Aimed Shot', 'Careful aimed arrow', 'Archer', 1, 'ACTIVE', 'SINGLE_ENEMY', 3, 12, 60, 0, 0, 'C', 'physical', 0, ARRAY['archer', 'physical'], 'aimed_shot'),
('Quick Draw', 'Fast arrow shot', 'Archer', 2, 'ACTIVE', 'SINGLE_ENEMY', 2, 10, 50, 0, 0, 'C', 'physical', 0, ARRAY['archer', 'speed'], 'quick_draw'),
('Camouflage', 'Increase evasion', 'Archer', 4, 'ACTIVE', 'SELF', 6, 12, 0, 0, 3, 'C', 'physical', 0, ARRAY['archer', 'stealth'], 'camouflage'),

-- Assassin
('Dagger Strike', 'Quick dagger attack', 'Assassin', 1, 'ACTIVE', 'SINGLE_ENEMY', 2, 10, 55, 0, 0, 'C', 'physical', 0, ARRAY['assassin', 'physical'], 'dagger_strike'),
('Poison Vial', 'Throw poison', 'Assassin', 3, 'ACTIVE', 'SINGLE_ENEMY', 4, 15, 40, 0, 0, 'C', 'physical', 0, ARRAY['assassin', 'poison'], 'poison_vial'),
('Stealth', 'Become harder to hit', 'Assassin', 5, 'ACTIVE', 'SELF', 8, 18, 0, 0, 2, 'C', 'physical', 0, ARRAY['assassin', 'stealth'], 'stealth'),

-- Healer
('Minor Heal', 'Small heal', 'Healer', 1, 'ACTIVE', 'SINGLE_ALLY', 3, 12, 0, 80, 0, 'C', 'magical', 0, ARRAY['healer', 'heal'], 'minor_heal');

-- ============================================
-- B RANK SKILLS (15 new)
-- ============================================

INSERT INTO abilities (name, description, class, unlock_level, ability_type, target_type, cooldown, mana_cost, base_damage, base_heal, duration_secs, rarity, damage_type, aoe_radius, synergy_tags, applies_buff, applies_debuff, animation_name) VALUES
-- Warrior
('Crushing Blow', 'Heavy overhead strike', 'Warrior', 12, 'ACTIVE', 'SINGLE_ENEMY', 5, 35, 130, 0, 0, 'B', 'physical', 0, ARRAY['warrior', 'physical'], NULL, 'armor_break', 'crushing_blow'),
('Battle Roar', 'Intimidate enemies', 'Warrior', 14, 'ACTIVE', 'ALL_ENEMIES', 8, 30, 0, 0, 2, 'B', 'physical', 0, ARRAY['warrior', 'debuff'], NULL, 'fear', 'battle_roar'),
('Sword Dance', 'Multiple quick strikes', 'Warrior', 16, 'ACTIVE', 'SINGLE_ENEMY', 6, 38, 140, 0, 0, 'B', 'physical', 0, ARRAY['warrior', 'combo'], NULL, NULL, 'sword_dance'),

-- Mage
('Fireball', 'Classic fire projectile', 'Mage', 12, 'ACTIVE', 'SINGLE_ENEMY', 5, 35, 125, 0, 0, 'B', 'magical', 0, ARRAY['mage', 'fire'], NULL, 'burn', 'fireball_classic'),
('Ice Lance', 'Sharp ice projectile', 'Mage', 14, 'ACTIVE', 'SINGLE_ENEMY', 6, 33, 115, 0, 0, 'B', 'magical', 0, ARRAY['mage', 'ice'], NULL, 'slow', 'ice_lance'),
('Arcane Barrier', 'Magic shield for team', 'Mage', 16, 'ACTIVE', 'ALL_ALLIES', 10, 40, 0, 0, 3, 'B', 'magical', 0, ARRAY['mage', 'support'], 'shield', NULL, 'arcane_barrier'),
('Spell Steal', 'Copy enemy buff', 'Mage', 18, 'ACTIVE', 'SINGLE_ENEMY', 12, 45, 0, 0, 0, 'B', 'magical', 0, ARRAY['mage', 'utility'], NULL, NULL, 'spell_steal'),

-- Tank
('Iron Skin', 'Harden skin', 'Tank', 12, 'ACTIVE', 'SELF', 10, 35, 0, 0, 3, 'B', 'physical', 0, ARRAY['tank', 'defense'], 'x_defense', NULL, 'iron_skin'),
('Revenge', 'Counter after taking damage', 'Tank', 14, 'PASSIVE', 'SELF', 0, 0, 100, 0, 0, 'B', 'physical', 0, ARRAY['tank', 'counter'], NULL, NULL, 'revenge'),
('Rallying Cry', 'Boost team morale', 'Tank', 16, 'ACTIVE', 'ALL_ALLIES', 12, 40, 0, 50, 0, 'B', 'physical', 0, ARRAY['tank', 'support'], 'morale', NULL, 'rallying_cry'),

-- Archer
('Poison Arrow', 'Arrow with poison tip', 'Archer', 12, 'ACTIVE', 'SINGLE_ENEMY', 5, 32, 85, 0, 0, 'B', 'physical', 0, ARRAY['archer', 'poison'], NULL, 'poison', 'poison_arrow_b'),
('Volley', 'Rain of arrows', 'Archer', 14, 'ACTIVE', 'AOE', 7, 40, 100, 0, 0, 'B', 'physical', 3, ARRAY['archer', 'aoe'], NULL, NULL, 'volley_rain'),
('Hunter Instinct', 'Increase crit chance', 'Archer', 16, 'PASSIVE', 'SELF', 0, 0, 0, 0, 0, 'B', 'physical', 0, ARRAY['archer', 'passive'], 'crit', NULL, 'hunter_instinct'),

-- Assassin
('Shadow Strike', 'Attack from shadows', 'Assassin', 12, 'ACTIVE', 'SINGLE_ENEMY', 6, 35, 140, 0, 0, 'B', 'physical', 0, ARRAY['assassin', 'stealth'], NULL, NULL, 'shadow_strike'),
('Toxic Blade', 'Poisoned weapon', 'Assassin', 14, 'ACTIVE', 'SINGLE_ENEMY', 5, 32, 95, 0, 0, 'B', 'physical', 0, ARRAY['assassin', 'poison'], NULL, 'poison', 'toxic_blade');

-- ============================================
-- A RANK SKILLS (15 new)
-- ============================================

INSERT INTO abilities (name, description, class, unlock_level, ability_type, target_type, cooldown, mana_cost, base_damage, base_heal, duration_secs, rarity, damage_type, aoe_radius, status_effect_chance, synergy_tags, applies_buff, applies_debuff, animation_name) VALUES
-- Warrior
('Berserker Fury', 'Massive damage boost, lose defense', 'Warrior', 22, 'ACTIVE', 'SELF', 15, 65, 0, 0, 4, 'A', 'physical', 0, 0, ARRAY['warrior', 'berserker'], 'x_attack', 'x_defense', 'berserker_fury'),
('Execution', 'Finish low HP enemies', 'Warrior', 24, 'ACTIVE', 'SINGLE_ENEMY', 10, 60, 300, 0, 0, 'A', 'physical', 0, 0, ARRAY['warrior', 'execute'], NULL, NULL, 'execution'),
('War Machine', 'Become unstoppable', 'Warrior', 26, 'ACTIVE', 'SELF', 20, 70, 0, 0, 3, 'A', 'physical', 0, 0, ARRAY['warrior', 'buff'], 'unstoppable', NULL, 'war_machine'),
('Seismic Slam', 'Ground pound AOE', 'Warrior', 28, 'ACTIVE', 'AOE', 12, 65, 190, 0, 0, 'A', 'physical', 4, 50, ARRAY['warrior', 'earth'], NULL, 'stun', 'seismic_slam'),

-- Mage
('Pyroblast', 'Massive fireball', 'Mage', 22, 'ACTIVE', 'SINGLE_ENEMY', 10, 70, 250, 0, 0, 'A', 'magical', 0, 0, ARRAY['mage', 'fire'], NULL, 'burn', 'pyroblast'),
('Glacial Spike', 'Huge ice spike', 'Mage', 24, 'ACTIVE', 'SINGLE_ENEMY', 11, 72, 230, 0, 0, 'A', 'magical', 0, 60, ARRAY['mage', 'ice'], NULL, 'freeze', 'glacial_spike'),
('Void Bolt', 'Dark energy blast', 'Mage', 26, 'ACTIVE', 'SINGLE_ENEMY', 9, 68, 210, 0, 0, 'A', 'magical', 0, 0, ARRAY['mage', 'dark'], NULL, 'curse', 'void_bolt'),
('Elemental Fury', 'Summon elemental storm', 'Mage', 28, 'ACTIVE', 'AOE', 15, 80, 200, 0, 0, 'A', 'magical', 5, 0, ARRAY['mage', 'elemental'], NULL, NULL, 'elemental_fury'),
('Mana Burn', 'Destroy enemy mana', 'Mage', 30, 'ACTIVE', 'SINGLE_ENEMY', 10, 50, 80, 0, 0, 'A', 'magical', 0, 0, ARRAY['mage', 'utility'], NULL, 'mana_burn', 'mana_burn_skill'),

-- Tank
('Immovable', 'Cannot be moved or stunned', 'Tank', 22, 'ACTIVE', 'SELF', 18, 65, 0, 0, 3, 'A', 'physical', 0, 0, ARRAY['tank', 'defense'], 'immovable', NULL, 'immovable_skill'),
('Guardian Shield', 'Protect entire team', 'Tank', 24, 'ACTIVE', 'ALL_ALLIES', 16, 70, 0, 0, 2, 'A', 'physical', 0, 0, ARRAY['tank', 'support'], 'shield', NULL, 'guardian_shield'),
('Retaliation', 'Damage scales with damage taken', 'Tank', 26, 'ACTIVE', 'SINGLE_ENEMY', 8, 55, 200, 0, 0, 'A', 'physical', 0, 0, ARRAY['tank', 'counter'], NULL, NULL, 'retaliation'),

-- Archer
('Headshot', 'Massive single target damage', 'Archer', 22, 'ACTIVE', 'SINGLE_ENEMY', 12, 70, 350, 0, 0, 'A', 'physical', 0, 0, ARRAY['archer', 'sniper'], NULL, NULL, 'headshot_skill'),
('Arrow Storm', 'Massive AOE arrow rain', 'Archer', 24, 'ACTIVE', 'AOE', 14, 75, 180, 0, 0, 'A', 'physical', 5, 0, ARRAY['archer', 'aoe'], NULL, NULL, 'arrow_storm'),
('Perfect Aim', 'Next attack cannot miss', 'Archer', 26, 'ACTIVE', 'SELF', 10, 50, 0, 0, 1, 'A', 'physical', 0, 0, ARRAY['archer', 'buff'], 'perfect_aim', NULL, 'perfect_aim_skill');

-- ============================================
-- S RANK SKILLS (8 new)
-- ============================================

INSERT INTO abilities (name, description, class, unlock_level, ability_type, target_type, cooldown, mana_cost, base_damage, base_heal, duration_secs, rarity, damage_type, aoe_radius, status_effect_chance, synergy_tags, is_ultimate, animation_name) VALUES
-- Warrior
('Cataclysm', 'Devastating earth attack', 'Warrior', 38, 'ACTIVE', 'AOE', 20, 105, 380, 0, 0, 'S', 'physical', 5, 50, ARRAY['warrior', 'earth'], false, 'cataclysm_skill'),
('Blade Storm', 'Become a whirlwind of blades', 'Warrior', 42, 'ACTIVE', 'AOE', 18, 100, 350, 0, 3, 'S', 'physical', 4, 0, ARRAY['warrior', 'aoe'], false, 'blade_storm'),

-- Mage
('Supernova', 'Massive explosion', 'Mage', 38, 'ACTIVE', 'AOE', 22, 120, 480, 0, 0, 'S', 'magical', 6, 0, ARRAY['mage', 'fire'], false, 'supernova'),
('Dimensional Rift', 'Tear reality for damage', 'Mage', 42, 'ACTIVE', 'AOE', 20, 115, 420, 0, 0, 'S', 'magical', 4, 0, ARRAY['mage', 'void'], false, 'dimensional_rift'),

-- Tank
('Titan Shield', 'Become invincible tank', 'Tank', 38, 'ACTIVE', 'SELF', 30, 110, 0, 0, 3, 'S', 'physical', 0, 0, ARRAY['tank', 'defense'], false, 'titan_shield'),
('Martyr', 'Sacrifice HP to heal team', 'Tank', 42, 'ACTIVE', 'ALL_ALLIES', 25, 100, 0, 500, 0, 'S', 'physical', 0, 0, ARRAY['tank', 'support'], false, 'martyr'),

-- Archer
('Judgment Shot', 'Execute low HP enemies', 'Archer', 38, 'ACTIVE', 'SINGLE_ENEMY', 20, 115, 600, 0, 0, 'S', 'physical', 0, 0, ARRAY['archer', 'execute'], false, 'judgment_shot'),
('Phoenix Arrow', 'Revive self if killed', 'Archer', 42, 'PASSIVE', 'SELF', 999, 120, 0, 0, 0, 'S', 'physical', 0, 0, ARRAY['archer', 'revive'], false, 'phoenix_arrow');

-- Verify final count
SELECT 
    'Total Skills After Update' as metric,
    COUNT(*) as count 
FROM abilities;

SELECT 
    'Skills by Rarity' as metric,
    rarity,
    COUNT(*) as count 
FROM abilities 
GROUP BY rarity 
ORDER BY 
    CASE rarity
        WHEN 'C' THEN 1
        WHEN 'B' THEN 2
        WHEN 'A' THEN 3
        WHEN 'S' THEN 4
        WHEN 'SS' THEN 5
        WHEN 'SSS' THEN 6
    END;
