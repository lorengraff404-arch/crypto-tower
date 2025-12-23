-- Crypto Tower Defense - Complete Database Schema
-- Generated: 2025-12-21
-- Version: 1.1.0 (Audit & Audit Fix)
-- 
-- This migration creates ALL tables from scratch
-- Run this on a FRESH database to avoid conflicts

-- ============================================================
-- USERS & AUTHENTICATION
-- ============================================================

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    wallet_address VARCHAR(42) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE,
    email VARCHAR(255),
    is_admin BOOLEAN DEFAULT FALSE,
    gtk_balance DECIMAL(20,2) DEFAULT 0,
    tower_balance DECIMAL(20,2) DEFAULT 0,
    total_spent DECIMAL(20,2) DEFAULT 0,
    referral_code VARCHAR(20) UNIQUE,
    referred_by INTEGER REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'active',
    is_banned BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_wallet ON users(wallet_address);
CREATE INDEX idx_users_referral ON users(referral_code);

-- ============================================================
-- CHARACTERS
-- ============================================================

CREATE TABLE characters (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    class VARCHAR(20) NOT NULL CHECK (class IN ('Warrior', 'Mage', 'Tank')),
    element VARCHAR(20) CHECK (element IN ('Fire', 'Water', 'Earth', 'Air', 'Lightning', 'Dark', 'Light', 'Neutral')),
    rarity VARCHAR(3) CHECK (rarity IN ('C', 'B', 'A', 'S', 'SS', 'SSS')),
    level INTEGER DEFAULT 1,
    experience INTEGER DEFAULT 0,
    total_experience INTEGER DEFAULT 0,
    max_hp INTEGER DEFAULT 100,
    current_hp INTEGER DEFAULT 100,
    max_mana INTEGER DEFAULT 50,
    current_mana INTEGER DEFAULT 50,
    attack INTEGER DEFAULT 10,
    defense INTEGER DEFAULT 5,
    speed INTEGER DEFAULT 10,
    critical_rate DECIMAL(5,2) DEFAULT 5.00,
    critical_damage DECIMAL(5,2) DEFAULT 150.00,
    fatigue INTEGER DEFAULT 0,
    durability INTEGER DEFAULT 100,
    personality VARCHAR(50),
    background_story TEXT,
    sprite_url TEXT,
    token_id INTEGER,
    is_egg BOOLEAN DEFAULT FALSE,
    owner_id INTEGER REFERENCES users(id),
    nft_token_id VARCHAR(100),
    is_minted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_characters_user ON characters(user_id);
CREATE INDEX idx_characters_class ON characters(class);
CREATE INDEX idx_characters_rarity ON characters(rarity);

-- ============================================================
-- ABILITIES SYSTEM
-- ============================================================

CREATE TABLE abilities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    class VARCHAR(20) NOT NULL CHECK (class IN ('Warrior', 'Mage', 'Tank', 'Universal')),
    rarity VARCHAR(3) CHECK (rarity IN ('C', 'B', 'A', 'S', 'SS', 'SSS')),
    unlock_level INTEGER DEFAULT 1,
    base_damage INTEGER DEFAULT 0,
    mana_cost INTEGER DEFAULT 10,
    cooldown INTEGER DEFAULT 0,
    element VARCHAR(20),
    target_type VARCHAR(20) DEFAULT 'single' CHECK (target_type IN ('single', 'aoe', 'self', 'ally')),
    status_effect VARCHAR(50),
    effect_duration INTEGER DEFAULT 0,
    effect_power INTEGER DEFAULT 0,
    icon_url TEXT,
    animation_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_abilities_class ON abilities(class);
CREATE INDEX idx_abilities_rarity ON abilities(rarity);
CREATE INDEX idx_abilities_level ON abilities(unlock_level);

-- Character Abilities (Learned - Unlimited)
CREATE TABLE character_abilities (
    id SERIAL PRIMARY KEY,
    character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    ability_id INTEGER NOT NULL REFERENCES abilities(id) ON DELETE CASCADE,
    learned_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(character_id, ability_id)
);

CREATE INDEX idx_character_abilities_character ON character_abilities(character_id);

-- Equipped Abilities (Battle Loadout - Slot Limited)
CREATE TABLE equipped_abilities (
    id SERIAL PRIMARY KEY,
    character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    ability_id INTEGER NOT NULL REFERENCES abilities(id) ON DELETE CASCADE,
    slot_position INTEGER NOT NULL CHECK (slot_position BETWEEN 1 AND 16),
    equipped_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(character_id, slot_position),
    UNIQUE(character_id, ability_id)
);

CREATE INDEX idx_equipped_abilities_character ON equipped_abilities(character_id);

-- ============================================================
-- BATTLE SYSTEM
-- ============================================================

CREATE TABLE battles (
    id SERIAL PRIMARY KEY,
    battle_type VARCHAR(20) NOT NULL CHECK (battle_type IN ('pve', 'pvp', 'raid', 'ranked', 'wager')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'completed', 'cancelled')),
    winner_id INTEGER REFERENCES users(id),
    wager_amount DECIMAL(20,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE battle_participants (
    id SERIAL PRIMARY KEY,
    battle_id INTEGER NOT NULL REFERENCES battles(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id),
    character_id INTEGER REFERENCES characters(id),
    team_side VARCHAR(10) CHECK (team_side IN ('team_a', 'team_b')),
    is_ai BOOLEAN DEFAULT FALSE,
    final_hp INTEGER DEFAULT 0,
    damage_dealt INTEGER DEFAULT 0,
    damage_taken INTEGER DEFAULT 0,
    abilities_used INTEGER DEFAULT 0
);

CREATE INDEX idx_battle_participants_battle ON battle_participants(battle_id);

CREATE TABLE battle_states (
    id SERIAL PRIMARY KEY,
    battle_id INTEGER NOT NULL REFERENCES battles(id) ON DELETE CASCADE,
    turn_number INTEGER DEFAULT 1,
    current_turn VARCHAR(10),
    state_data JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_battle_states_battle ON battle_states(battle_id);

CREATE TABLE status_effects (
    id SERIAL PRIMARY KEY,
    character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    effect_type VARCHAR(20) NOT NULL,
    effect_name VARCHAR(30) NOT NULL,
    stacks INTEGER DEFAULT 1 NOT NULL,
    duration INTEGER NOT NULL,
    max_duration INTEGER NOT NULL,
    turns_remaining INTEGER DEFAULT 0,
    expires_at TIMESTAMP,
    source_ability_id INTEGER,
    caster_id INTEGER,
    battle_id INTEGER,
    stat_modifier DECIMAL(5,2),
    damage_per_turn INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_status_effects_character ON status_effects(character_id);

-- ============================================================
-- ECONOMY & ITEMS
-- ============================================================

-- Items (NFTs / Instances)
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nft_token_id VARCHAR(100) UNIQUE,
    is_minted BOOLEAN DEFAULT FALSE,
    item_type VARCHAR(20) NOT NULL,
    name VARCHAR(50) NOT NULL,
    rarity VARCHAR(5) NOT NULL,
    icon_url VARCHAR(255),
    attack_bonus INTEGER DEFAULT 0,
    defense_bonus INTEGER DEFAULT 0,
    hp_bonus INTEGER DEFAULT 0,
    speed_bonus INTEGER DEFAULT 0,
    special_effects TEXT,
    durability INTEGER DEFAULT 100,
    is_broken BOOLEAN DEFAULT FALSE,
    is_consumable BOOLEAN DEFAULT FALSE,
    consume_effect VARCHAR(50),
    is_crafting_material BOOLEAN DEFAULT FALSE,
    is_stackable BOOLEAN DEFAULT FALSE,
    quantity INTEGER DEFAULT 1,
    is_listed BOOLEAN DEFAULT FALSE,
    list_price BIGINT,
    listed_at TIMESTAMP,
    is_equipped BOOLEAN DEFAULT FALSE,
    equipped_by_id INTEGER REFERENCES characters(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_items_owner ON items(owner_id);
CREATE INDEX idx_items_nft ON items(nft_token_id);

-- Shop Items (Catalog / Definitions)
CREATE TABLE shop_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    category VARCHAR(20) NOT NULL,
    effect_type VARCHAR(20) NOT NULL,
    effect_value BIGINT,
    gtk_cost BIGINT NOT NULL,
    is_consumable BOOLEAN DEFAULT TRUE,
    max_stack BIGINT DEFAULT 99,
    is_available BOOLEAN DEFAULT TRUE,
    icon_url VARCHAR(200),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- User Inventory (For Shop/Stackable Items)
CREATE TABLE user_inventories (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES shop_items(id),
    quantity INTEGER DEFAULT 1 CHECK (quantity >= 0),
    acquired_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_inventories_user ON user_inventories(user_id);
CREATE UNIQUE INDEX idx_user_inventories_unique ON user_inventories(user_id, item_id);

-- Transactions
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_type VARCHAR(30) NOT NULL,
    token_type VARCHAR(10) NOT NULL,
    amount BIGINT NOT NULL,
    balance_before BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    battle_id INTEGER REFERENCES battles(id),
    character_id INTEGER REFERENCES characters(id),
    item_id INTEGER REFERENCES items(id),
    description VARCHAR(255),
    metadata TEXT,
    blockchain_tx_hash VARCHAR(66),
    chain_id BIGINT,
    is_on_chain BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_transactions_user ON transactions(user_id);

-- ============================================================
-- CAMPAIGN & MISSIONS
-- ============================================================

CREATE TABLE islands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    difficulty VARCHAR(20),
    required_level INTEGER DEFAULT 1,
    reward_gtk DECIMAL(20,2) DEFAULT 0,
    reward_xp INTEGER DEFAULT 0,
    island_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT FALSE, -- Default false in audits? No, keep logic.
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE raid_bosses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    max_hp INTEGER NOT NULL,
    attack INTEGER DEFAULT 10,
    defense INTEGER DEFAULT 5,
    speed INTEGER DEFAULT 5,
    element VARCHAR(20),
    abilities JSONB,
    loot_table JSONB,
    difficulty INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE raid_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    island_id INTEGER REFERENCES islands(id),
    boss_id INTEGER REFERENCES raid_bosses(id),
    status VARCHAR(20) DEFAULT 'active',
    current_turn INTEGER DEFAULT 1,
    player_hp INTEGER,
    boss_hp INTEGER,
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

CREATE INDEX idx_raid_sessions_user ON raid_sessions(user_id);

CREATE TABLE missions (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    mission_type VARCHAR(50),
    objective_data JSONB,
    reward_gtk DECIMAL(20,2) DEFAULT 0,
    reward_xp INTEGER DEFAULT 0,
    required_level INTEGER DEFAULT 1,
    is_daily BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- TEAMS
-- ============================================================

CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100),
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE team_members (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    character_id INTEGER NOT NULL REFERENCES characters(id),
    position INTEGER CHECK (position BETWEEN 1 AND 5),
    UNIQUE(team_id, position),
    UNIQUE(team_id, character_id)
);

-- ============================================================
-- SYSTEM SETTINGS
-- ============================================================

CREATE TABLE system_settings (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert default ability slot limits
INSERT INTO system_settings (key, value, description) VALUES
('ability_slots_c', '4', 'Maximum ability slots for C-rank characters'),
('ability_slots_b', '6', 'Maximum ability slots for B-rank characters'),
('ability_slots_a', '8', 'Maximum ability slots for A-rank characters'),
('ability_slots_s', '10', 'Maximum ability slots for S-rank characters'),
('ability_slots_ss', '12', 'Maximum ability slots for SS-rank characters'),
('ability_slots_sss', '16', 'Maximum ability slots for SSS-rank characters')
ON CONFLICT DO NOTHING;

-- ============================================================
-- EGGS & BREEDING
-- ============================================================

CREATE TABLE eggs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent1_id INTEGER REFERENCES characters(id),
    parent2_id INTEGER REFERENCES characters(id),
    
    -- Mint info
    mint_cost BIGINT DEFAULT 0,
    mint_tx_hash VARCHAR(66),
    
    -- Properties
    rarity VARCHAR(20) NOT NULL,
    element VARCHAR(20),
    character_type VARCHAR(20),
    class VARCHAR(20),
    
    -- Hidden traits
    predetermined_stats JSONB,
    predetermined_abilities JSONB,
    
    -- Reveal status
    is_stats_revealed BOOLEAN DEFAULT FALSE,
    revealed_at TIMESTAMP,
    
    -- Incubation
    incubation_time INTEGER NOT NULL DEFAULT 0,
    effective_incubation_time INTEGER DEFAULT 0,
    incubation_started_at TIMESTAMP,
    accelerators_applied TEXT,
    hatched_at TIMESTAMP,
    character_id INTEGER REFERENCES characters(id),
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_eggs_user ON eggs(user_id);

CREATE TABLE breeding_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    parent1_id INTEGER NOT NULL REFERENCES characters(id),
    parent2_id INTEGER NOT NULL REFERENCES characters(id),
    status VARCHAR(20) DEFAULT 'in_progress',
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    egg_id INTEGER REFERENCES eggs(id),
    tokens_spent INTEGER DEFAULT 0
);

CREATE INDEX idx_breeding_sessions_user ON breeding_sessions(user_id);

-- ============================================================
-- DAILY QUESTS
-- ============================================================

CREATE TABLE daily_quests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quest_type VARCHAR(50),
    target_value INTEGER,
    current_progress INTEGER DEFAULT 0,
    reward_gtk DECIMAL(20,2) DEFAULT 0,
    reward_xp INTEGER DEFAULT 0,
    is_completed BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_daily_quests_user ON daily_quests(user_id);

-- ============================================================
-- STORY & DIALOGUES
-- ============================================================

CREATE TABLE story_dialogues (
    id SERIAL PRIMARY KEY,
    chapter VARCHAR(50),
    scene_id VARCHAR(50),
    character_name VARCHAR(100),
    dialogue_text TEXT,
    order_index INTEGER,
    choices JSONB,
    next_scene VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- ADMIN & AUDIT
-- ============================================================

CREATE TABLE admin_audit_logs (
    id SERIAL PRIMARY KEY,
    admin_id INTEGER REFERENCES users(id),
    action VARCHAR(100),
    target_type VARCHAR(50),
    target_id INTEGER,
    changes JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_admin ON admin_audit_logs(admin_id);
CREATE INDEX idx_audit_logs_created ON admin_audit_logs(created_at);

-- ============================================================
-- END OF MIGRATION
-- ============================================================
