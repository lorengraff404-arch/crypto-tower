-- Migration: Create ability and status effect tables
-- File: 026_create_ability_status_system.sql

-- Ability definitions
CREATE TABLE IF NOT EXISTS abilities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    class VARCHAR(20) NOT NULL,
    element VARCHAR(20),
    unlock_level INT NOT NULL DEFAULT 1,
    power INT NOT NULL DEFAULT 50,
    cooldown INT DEFAULT 0,
    max_pp INT DEFAULT 10,
    effect_type VARCHAR(20), -- damage, heal, buff, debuff
    effect_value INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Character ability usage (PP tracking)
CREATE TABLE IF NOT EXISTS ability_usage (
    id SERIAL PRIMARY KEY,
    character_id INT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    ability_id INT NOT NULL REFERENCES abilities(id) ON DELETE CASCADE,
    current_pp INT NOT NULL,
    max_pp INT NOT NULL,
    last_used_at TIMESTAMP,
    last_recovered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(character_id, ability_id)
);

-- Active status effects on characters
CREATE TABLE IF NOT EXISTS character_status_effects (
    id SERIAL PRIMARY KEY,
    character_id INT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    status_type VARCHAR(20) NOT NULL, -- poison, burn, freeze, paralysis, sleep
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    damage_per_turn INT DEFAULT 0
);

-- Active buffs on characters
CREATE TABLE IF NOT EXISTS character_buffs (
    id SERIAL PRIMARY KEY,
    character_id INT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    buff_type VARCHAR(20) NOT NULL, -- x_attack, x_defense, x_speed, etc
    multiplier DECIMAL(3,2) NOT NULL, -- 1.5 for +50%
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    turns_remaining INT NOT NULL
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_abilities_class ON abilities(class);
CREATE INDEX IF NOT EXISTS idx_abilities_unlock_level ON abilities(unlock_level);
CREATE INDEX IF NOT EXISTS idx_ability_usage_character ON ability_usage(character_id);
CREATE INDEX IF NOT EXISTS idx_status_effects_character ON character_status_effects(character_id);
CREATE INDEX IF NOT EXISTS idx_buffs_character ON character_buffs(character_id);
