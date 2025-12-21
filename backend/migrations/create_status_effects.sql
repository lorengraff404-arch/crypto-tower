-- Phase 15.3: Status Effects System - Database Schema

-- Table: status_effects (master list of all possible effects)
CREATE TABLE IF NOT EXISTS status_effects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('buff', 'debuff')),
    effect_type VARCHAR(30) NOT NULL,
    magnitude INT NOT NULL,
    base_duration INT NOT NULL,
    description TEXT,
    icon_emoji VARCHAR(10),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Table: active_status_effects (effects currently active in battles)
CREATE TABLE IF NOT EXISTS active_status_effects (
    id SERIAL PRIMARY KEY,
    raid_session_id INT NOT NULL REFERENCES raid_sessions(id) ON DELETE CASCADE,
    character_id INT REFERENCES characters(id) ON DELETE CASCADE,
    is_enemy BOOLEAN DEFAULT FALSE,
    status_effect_id INT NOT NULL REFERENCES status_effects(id),
    turns_remaining INT NOT NULL,
    applied_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(raid_session_id, character_id, status_effect_id)
);

-- Insert base status effects
INSERT INTO status_effects (name, type, effect_type, magnitude, base_duration, description, icon_emoji) VALUES
-- Buffs
('ATK Up', 'buff', 'atk_up', 30, 2, '+30% Attack for 2 turns', '🔺'),
('DEF Up', 'buff', 'def_up', 30, 2, '+30% Defense for 2 turns', '🛡️'),
('SPD Up', 'buff', 'spd_up', 50, 2, '+50% Speed for 2 turns', '⚡'),
('Regeneration', 'buff', 'regen', 10, 3, 'Restore 10% HP per turn for 3 turns', '💚'),
('Critical Boost', 'buff', 'crit_boost', 30, 2, '+30% Critical chance for 2 turns', '⭐'),

-- Debuffs
('Burn', 'debuff', 'burn', 5, 3, '5% max HP damage per turn for 3 turns', '🔥'),
('Freeze', 'debuff', 'freeze', 50, 2, '50% chance to skip turn for 2 turns', '❄️'),
('Poison', 'debuff', 'poison', 7, 4, '7% max HP damage per turn for 4 turns', '☠️'),
('Stun', 'debuff', 'stun', 100, 1, 'Skip next turn', '⚡'),
('ATK Down', 'debuff', 'atk_down', 30, 2, '-30% Attack for 2 turns', '🔻'),
('DEF Down', 'debuff', 'def_down', 30, 2, '-30% Defense for 2 turns', '💔')
ON CONFLICT (name) DO NOTHING;

-- Index for performance
CREATE INDEX IF NOT EXISTS idx_active_status_session ON active_status_effects(raid_session_id);
CREATE INDEX IF NOT EXISTS idx_active_status_char ON active_status_effects(character_id);
