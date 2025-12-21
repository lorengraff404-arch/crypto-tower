-- ============================================
-- PHASE 2: Seed 61 New Skills (Total: 100)
-- ============================================
-- Current: 39 skills
-- Target: 100 skills
-- New skills: 61

-- Distribution:
-- C Rank: 15 new (5 → 20)
-- B Rank: 15 new (10 → 25)
-- A Rank: 15 new (10 → 25)
-- S Rank: 8 new (7 → 15)
-- SS Rank: 5 new (5 → 10)
-- SSS Rank: 3 new (2 → 5)

-- ============================================
-- C RANK SKILLS (15 new - Basic/Common)
-- ============================================

INSERT INTO abilities (name, description, class, unlock_level, ability_type, target_type, cooldown, mana_cost, base_damage, base_heal, duration_secs, rarity, damage_type, aoe_radius, animation_name) VALUES
-- Warrior
('Quick Slash', 'A fast melee attack', 'Warrior', 1, 'ACTIVE', 'SINGLE_ENEMY', 2, 10, 50, 0, 0, 'C', 'physical', 0, 'slash_quick'),
('Power Strike', 'A powerful overhead strike', 'Warrior', 3, 'ACTIVE', 'SINGLE_ENEMY', 3, 15, 70, 0, 0, 'C', 'physical', 0, 'strike_power'),
('Guard Stance', 'Increase defense temporarily', 'Warrior', 5, 'ACTIVE', 'SELF', 4, 12, 0, 0, 2, 'C', 'physical', 0, 'guard_stance'),
('Battle Cry', 'Boost team attack', 'Warrior', 7, 'ACTIVE', 'ALL_ALLIES', 5, 15, 0, 0, 2, 'C', 'physical', 0, 'battle_cry'),

-- Mage
('Magic Dart', 'Quick magical projectile', 'Mage', 1, 'ACTIVE', 'SINGLE_ENEMY', 2, 8, 45, 0, 0, 'C', 'magical', 0, 'magic_dart'),
('Mana Bolt', 'Concentrated mana attack', 'Mage', 3, 'ACTIVE', 'SINGLE_ENEMY', 3, 12, 60, 0, 0, 'C', 'magical', 0, 'mana_bolt'),
('Arcane Shield', 'Create protective barrier', 'Mage', 5, 'ACTIVE', 'SELF', 5, 15, 0, 50, 0, 'C', 'magical', 0, 'arcane_shield'),
('Focus', 'Increase critical chance', 'Mage', 7, 'ACTIVE', 'SELF', 4, 10, 0, 0, 3, 'C', 'magical', 0, 'focus'),

-- Tank
('Shield Bash', 'Bash with shield', 'Tank', 1, 'ACTIVE', 'SINGLE_ENEMY', 4, 15, 40, 0, 0, 'C', 'physical', 0, 'shield_bash'),
('Fortify', 'Increase defense', 'Tank', 3, 'ACTIVE', 'SELF', 5, 12, 0, 0, 2, 'C', 'physical', 0, 'fortify'),
('Endure', 'Survive fatal blow', 'Tank', 7, 'ACTIVE', 'SELF', 30, 20, 0, 0, 0, 'C', 'physical', 0, 'endure'),

-- Archer
('Rapid Shot', 'Quick arrow', 'Archer', 1, 'ACTIVE', 'SINGLE_ENEMY', 2, 10, 55, 0, 0, 'C', 'physical', 0, 'rapid_shot'),
('Piercing Arrow', 'Ignore some defense', 'Archer', 3, 'ACTIVE', 'SINGLE_ENEMY', 3, 15, 65, 0, 0, 'C', 'physical', 0, 'piercing_arrow'),
('Evasion', 'Increase dodge chance', 'Archer', 5, 'ACTIVE', 'SELF', 5, 12, 0, 0, 2, 'C', 'physical', 0, 'evasion'),
('Mark Target', 'Mark enemy for bonus damage', 'Archer', 7, 'ACTIVE', 'SINGLE_ENEMY', 4, 10, 0, 0, 3, 'C', 'physical', 0, 'mark_target');

-- ============================================
-- B RANK SKILLS (15 new - Improved/Uncommon)
-- ============================================

INSERT INTO abilities (name, description, class, unlock_level, ability_type, target_type, cooldown, mana_cost, base_damage, base_heal, duration_secs, rarity, damage_type, aoe_radius, applies_buff, applies_debuff, animation_name) VALUES
-- Warrior
('Whirlwind', 'Spin attack hitting all nearby enemies', 'Warrior', 10, 'ACTIVE', 'AOE', 5, 30, 90, 0, 0, 'B', 'physical', 3, NULL, NULL, 'whirlwind'),
('Berserker Rage', 'Increase attack, reduce defense', 'Warrior', 12, 'ACTIVE', 'SELF', 8, 25, 0, 0, 3, 'B', 'physical', 0, 'x_attack', 'x_defense', 'berserker_rage'),
('Cleave', 'Strike two enemies', 'Warrior', 14, 'ACTIVE', 'SINGLE_ENEMY', 6, 35, 110, 0, 0, 'B', 'physical', 0, NULL, NULL, 'cleave'),
('Iron Will', 'Immune to debuffs', 'Warrior', 16, 'ACTIVE', 'SELF', 12, 30, 0, 0, 2, 'B', 'physical', 0, NULL, NULL, 'iron_will'),

-- Mage
('Flame Burst', 'Fire explosion', 'Mage', 10, 'ACTIVE', 'SINGLE_ENEMY', 5, 35, 120, 0, 0, 'B', 'magical', 0, NULL, 'burn', 'flame_burst'),
('Frost Bite', 'Ice attack that slows', 'Mage', 12, 'ACTIVE', 'SINGLE_ENEMY', 6, 32, 100, 0, 0, 'B', 'magical', 0, NULL, 'slow', 'frost_bite'),
('Chain Lightning', 'Lightning that chains', 'Mage', 14, 'ACTIVE', 'SINGLE_ENEMY', 7, 40, 90, 0, 0, 'B', 'magical', 0, NULL, NULL, 'chain_lightning'),
('Mana Shield', 'Absorb damage with mana', 'Mage', 16, 'ACTIVE', 'SELF', 8, 30, 0, 150, 0, 'B', 'magical', 0, NULL, NULL, 'mana_shield'),

-- Tank
('Shield Wall', 'Boost team defense', 'Tank', 10, 'ACTIVE', 'ALL_ALLIES', 10, 35, 0, 0, 2, 'B', 'physical', 0, 'x_defense', NULL, 'shield_wall'),
('Counter Strike', 'Reflect damage', 'Tank', 12, 'ACTIVE', 'SELF', 8, 30, 0, 0, 2, 'B', 'physical', 0, NULL, NULL, 'counter_strike'),
('Provoke', 'Force all enemies to attack you', 'Tank', 14, 'ACTIVE', 'ALL_ENEMIES', 10, 25, 0, 0, 1, 'B', 'physical', 0, NULL, NULL, 'provoke'),
('Last Stand', 'Boost stats when low HP', 'Tank', 16, 'PASSIVE', 'SELF', 20, 40, 0, 0, 0, 'B', 'physical', 0, NULL, NULL, 'last_stand'),

-- Archer
('Multi-Shot', 'Hit three enemies', 'Archer', 10, 'ACTIVE', 'AOE', 5, 32, 70, 0, 0, 'B', 'physical', 0, NULL, NULL, 'multi_shot'),
('Explosive Arrow', 'AOE explosion', 'Archer', 12, 'ACTIVE', 'AOE', 7, 40, 130, 0, 0, 'B', 'physical', 2, NULL, NULL, 'explosive_arrow'),
('Eagle Eye', 'Ignore evasion', 'Archer', 14, 'ACTIVE', 'SELF', 6, 25, 0, 0, 3, 'B', 'physical', 0, 'accuracy', NULL, 'eagle_eye');

-- ============================================
-- A RANK SKILLS (15 new - Advanced/Rare)
-- ============================================

INSERT INTO abilities (name, description, class, unlock_level, ability_type, target_type, cooldown, mana_cost, base_damage, base_heal, duration_secs, rarity, damage_type, aoe_radius, status_effect_chance, applies_buff, applies_debuff, animation_name) VALUES
-- Warrior
('Earthquake', 'Earth-shattering strike', 'Warrior', 20, 'ACTIVE', 'AOE', 10, 60, 180, 0, 0, 'A', 'physical', 4, 40, NULL, 'stun', 'earthquake'),
('Titan Strength', 'Massive attack boost', 'Warrior', 22, 'ACTIVE', 'SELF', 12, 55, 0, 0, 3, 'A', 'physical', 0, 0, 'x_attack', NULL, 'titan_strength'),
('Reckless Assault', 'High damage, take recoil', 'Warrior', 24, 'ACTIVE', 'SINGLE_ENEMY', 8, 65, 250, 0, 0, 'A', 'physical', 0, 0, NULL, NULL, 'reckless_assault'),
('War Banner', 'Team attack boost', 'Warrior', 26, 'ACTIVE', 'ALL_ALLIES', 15, 60, 0, 0, 4, 'A', 'physical', 0, 0, 'x_attack', NULL, 'war_banner'),
('Unstoppable', 'Immune to crowd control', 'Warrior', 28, 'ACTIVE', 'SELF', 18, 50, 0, 0, 3, 'A', 'physical', 0, 0, NULL, NULL, 'unstoppable'),

-- Mage
('Inferno', 'Massive fire damage with burn', 'Mage', 20, 'ACTIVE', 'SINGLE_ENEMY', 10, 70, 220, 0, 0, 'A', 'magical', 0, 0, NULL, 'burn', 'inferno'),
('Blizzard', 'Ice storm with freeze', 'Mage', 22, 'ACTIVE', 'AOE', 12, 75, 180, 0, 0, 'A', 'magical', 3, 50, NULL, 'freeze', 'blizzard'),
('Thunder Storm', 'Lightning storm with paralyze', 'Mage', 24, 'ACTIVE', 'AOE', 11, 72, 200, 0, 0, 'A', 'magical', 3, 40, NULL, 'paralyze', 'thunder_storm'),
('Arcane Explosion', 'Pure arcane damage AOE', 'Mage', 26, 'ACTIVE', 'AOE', 13, 80, 240, 0, 0, 'A', 'magical', 4, 0, NULL, NULL, 'arcane_explosion'),
('Mana Drain', 'Steal mana and deal damage', 'Mage', 28, 'ACTIVE', 'SINGLE_ENEMY', 9, 45, 100, 0, 0, 'A', 'magical', 0, 0, NULL, NULL, 'mana_drain'),

-- Tank
('Fortress', 'Massive defense boost', 'Tank', 20, 'ACTIVE', 'SELF', 15, 65, 0, 0, 3, 'A', 'physical', 0, 0, 'x_defense', NULL, 'fortress'),
('Thorns', 'Reflect damage to attackers', 'Tank', 22, 'ACTIVE', 'SELF', 14, 60, 0, 0, 3, 'A', 'physical', 0, 0, NULL, NULL, 'thorns'),
('Sacrifice', 'Take damage for ally', 'Tank', 24, 'ACTIVE', 'SINGLE_ALLY', 12, 55, 0, 0, 0, 'A', 'physical', 0, 0, NULL, NULL, 'sacrifice'),
('Unbreakable', 'Cannot go below 1 HP', 'Tank', 26, 'ACTIVE', 'SELF', 25, 70, 0, 0, 2, 'A', 'physical', 0, 0, NULL, NULL, 'unbreakable'),
('Rallying Shout', 'Heal and boost team defense', 'Tank', 28, 'ACTIVE', 'ALL_ALLIES', 16, 65, 0, 100, 0, 'A', 'physical', 0, 0, 'x_defense', NULL, 'rallying_shout');

-- ============================================
-- S RANK SKILLS (8 new - Elite/Epic)
-- ============================================

INSERT INTO abilities (name, description, class, unlock_level, ability_type, target_type, cooldown, mana_cost, base_damage, base_heal, duration_secs, rarity, damage_type, aoe_radius, status_effect_chance, is_ultimate, animation_name) VALUES
-- Warrior
('Ragnarok Strike', 'Devastating fire strike', 'Warrior', 35, 'ACTIVE', 'SINGLE_ENEMY', 18, 100, 400, 0, 0, 'S', 'physical', 0, 0, false, 'ragnarok_strike'),
('God of War', 'Massive stat boost', 'Warrior', 40, 'ACTIVE', 'SELF', 25, 110, 0, 0, 4, 'S', 'physical', 0, 0, false, 'god_of_war'),

-- Mage
('Meteor Storm', 'Rain of fire meteors', 'Mage', 35, 'ACTIVE', 'AOE', 22, 120, 450, 0, 0, 'S', 'magical', 5, 0, false, 'meteor_storm'),
('Absolute Zero', 'Freeze everything', 'Mage', 40, 'ACTIVE', 'AOE', 20, 115, 380, 0, 0, 'S', 'magical', 4, 70, false, 'absolute_zero'),

-- Tank
('Immortal Bastion', 'Invulnerable for 2 turns', 'Tank', 35, 'ACTIVE', 'SELF', 35, 110, 0, 0, 2, 'S', 'physical', 0, 0, false, 'immortal_bastion'),
('Phoenix Shield', 'Revive if killed', 'Tank', 40, 'ACTIVE', 'SELF', 999, 120, 0, 0, 0, 'S', 'physical', 0, 0, false, 'phoenix_shield'),

-- Archer
('Divine Arrow', 'Ignore all defense', 'Archer', 35, 'ACTIVE', 'SINGLE_ENEMY', 25, 115, 550, 0, 0, 'S', 'physical', 0, 0, false, 'divine_arrow'),
('Barrage', 'Rapid fire 20 arrows', 'Archer', 40, 'ACTIVE', 'AOE', 20, 110, 600, 0, 0, 'S', 'physical', 0, 0, false, 'barrage');

-- ============================================
-- SS RANK SKILLS (5 new - Master/Legendary)
-- ============================================

INSERT INTO abilities (name, description, class, unlock_level, ability_type, target_type, cooldown, mana_cost, base_damage, base_heal, duration_secs, rarity, damage_type, aoe_radius, is_ultimate, animation_name) VALUES
-- Warrior
('Apocalypse', 'Dark destruction reducing all stats', 'Warrior', 50, 'ACTIVE', 'AOE', 35, 150, 600, 0, 0, 'SS', 'true', 5, false, 'apocalypse'),
('Berserker Soul', 'Immune to death, massive attack', 'Warrior', 55, 'ACTIVE', 'SELF', 45, 160, 0, 0, 3, 'SS', 'physical', 0, false, 'berserker_soul'),

-- Mage
('Time Stop', 'Freeze all enemies for 1 turn', 'Mage', 50, 'ACTIVE', 'ALL_ENEMIES', 60, 180, 0, 0, 1, 'SS', 'magical', 0, false, 'time_stop'),

-- Tank
('Titan Form', 'Transform into titan', 'Tank', 50, 'ACTIVE', 'SELF', 50, 170, 0, 0, 5, 'SS', 'physical', 0, false, 'titan_form'),

-- Archer
('Celestial Arrow', 'Holy arrow that heals team', 'Archer', 50, 'ACTIVE', 'SINGLE_ENEMY', 45, 175, 800, 400, 0, 'SS', 'magical', 0, false, 'celestial_arrow');

-- ============================================
-- SSS RANK SKILLS (3 new - Mythic/Unique)
-- ============================================

INSERT INTO abilities (name, description, class, unlock_level, ability_type, target_type, cooldown, mana_cost, base_damage, base_heal, duration_secs, rarity, damage_type, is_ultimate, animation_name) VALUES
-- Ultimate Skills
('Genesis', 'Full heal and revive all allies', 'Mage', 60, 'ULTIMATE', 'ALL_ALLIES', 999, 250, 0, 9999, 0, 'SSS', 'magical', true, 'genesis'),
('Black Hole', 'Instant kill enemies below 40% HP', 'Mage', 65, 'ULTIMATE', 'AOE', 120, 280, 0, 0, 0, 'SSS', 'true', true, 'black_hole'),
('Immortality', 'Cannot die for 5 turns', 'Tank', 70, 'ULTIMATE', 'SELF', 999, 300, 0, 0, 5, 'SSS', 'physical', true, 'immortality');

-- Verify new skills
SELECT 
    'Total Skills' as metric,
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
