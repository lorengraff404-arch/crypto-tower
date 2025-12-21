-- Phase 16: Rewards & Loot System

-- Items table (master list of all items)
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    item_type VARCHAR(30) NOT NULL CHECK (item_type IN ('equipment', 'consumable', 'material', 'egg', 'token')),
    rarity VARCHAR(20) NOT NULL CHECK (rarity IN ('common', 'uncommon', 'rare', 'epic', 'legendary')),
    icon_emoji VARCHAR(10),
    icon_url VARCHAR(255),
    stackable BOOLEAN DEFAULT FALSE,
    max_stack INT DEFAULT 1,
    sell_price INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- User inventory
CREATE TABLE IF NOT EXISTS user_inventory (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id INT NOT NULL REFERENCES items(id),
    quantity INT DEFAULT 1 CHECK (quantity >= 0),
    acquired_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, item_id)
);

-- Loot tables (what can drop from missions)
CREATE TABLE IF NOT EXISTS loot_tables (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    mission_id INT REFERENCES missions(id),
    drop_type VARCHAR(20) DEFAULT 'random' CHECK (drop_type IN ('guaranteed', 'random', 'bonus')),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Loot entries (individual items in loot tables)
CREATE TABLE IF NOT EXISTS loot_entries (
    id SERIAL PRIMARY KEY,
    loot_table_id INT NOT NULL REFERENCES loot_tables(id) ON DELETE CASCADE,
    item_id INT REFERENCES items(id),
    drop_chance DECIMAL(5,2) NOT NULL CHECK (drop_chance >= 0 AND drop_chance <= 100),
    min_quantity INT DEFAULT 1,
    max_quantity INT DEFAULT 1,
    rarity_weight INT DEFAULT 1
);

-- Battle rewards (track what players earned)
CREATE TABLE IF NOT EXISTS battle_rewards (
    id SERIAL PRIMARY KEY,
    raid_session_id INT NOT NULL REFERENCES raid_sessions(id),
    user_id INT NOT NULL REFERENCES users(id),
    tokens_earned INT DEFAULT 0,
    xp_earned INT DEFAULT 0,
    items_json TEXT, -- JSON array of {item_id, quantity}
    performance_grade VARCHAR(2),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_inventory_user ON user_inventory(user_id);
CREATE INDEX IF NOT EXISTS idx_user_inventory_item ON user_inventory(item_id);
CREATE INDEX IF NOT EXISTS idx_loot_tables_mission ON loot_tables(mission_id);
CREATE INDEX IF NOT EXISTS idx_battle_rewards_user ON battle_rewards(user_id);
CREATE INDEX IF NOT EXISTS idx_battle_rewards_session ON battle_rewards(raid_session_id);

-- Insert base items
INSERT INTO items (name, description, item_type, rarity, icon_emoji, stackable, max_stack, sell_price) VALUES
-- Consumables
('Health Potion', 'Restores 50% HP', 'consumable', 'common', '🧪', true, 99, 50),
('Super Potion', 'Restores 100% HP', 'consumable', 'uncommon', '💊', true, 99, 150),
('Revive', 'Revives fainted character with 50% HP', 'consumable', 'rare', '💫', true, 20, 300),
('Max Revive', 'Revives fainted character with 100% HP', 'consumable', 'epic', '✨', true, 10, 1000),

-- Materials
('Iron Ore', 'Basic crafting material', 'material', 'common', '⛏️', true, 999, 10),
('Crystal Shard', 'Rare crafting material', 'material', 'uncommon', '💎', true, 999, 50),
('Dragon Scale', 'Epic crafting material', 'material', 'rare', '🐉', true, 999, 200),
('Stardust', 'Legendary crafting material', 'material', 'legendary', '⭐', true, 999, 1000),

-- Token items
('Token Bag (Small)', 'Contains 100 tokens', 'token', 'common', '💰', true, 99, 0),
('Token Bag (Medium)', 'Contains 500 tokens', 'token', 'uncommon', '💰', true, 99, 0),
('Token Bag (Large)', 'Contains 2000 tokens', 'token', 'rare', '💰', true, 99, 0),

-- Eggs
('Common Egg', 'Hatches a common character', 'egg', 'common', '🥚', false, 1, 0),
('Rare Egg', 'Hatches a rare character', 'egg', 'rare', '🥚', false, 1, 0),
('Epic Egg', 'Hatches an epic character', 'egg', 'epic', '🥚', false, 1, 0),
('Legendary Egg', 'Hatches a legendary character', 'egg', 'legendary', '🥚', false, 1, 0)
ON CONFLICT (name) DO NOTHING;

-- Create default loot tables for existing missions
INSERT INTO loot_tables (name, mission_id, drop_type) 
SELECT 
    CONCAT('Loot Table - ', m.name),
    m.id,
    'random'
FROM missions m
WHERE NOT EXISTS (SELECT 1 FROM loot_tables WHERE mission_id = m.id)
ON CONFLICT DO NOTHING;

-- Add loot entries for each loot table
INSERT INTO loot_entries (loot_table_id, item_id, drop_chance, min_quantity, max_quantity, rarity_weight)
SELECT 
    lt.id,
    i.id,
    CASE i.rarity
        WHEN 'common' THEN 60.0
        WHEN 'uncommon' THEN 25.0
        WHEN 'rare' THEN 10.0
        WHEN 'epic' THEN 4.0
        WHEN 'legendary' THEN 1.0
    END,
    1,
    CASE i.rarity
        WHEN 'common' THEN 3
        WHEN 'uncommon' THEN 2
        ELSE 1
    END,
    CASE i.rarity
        WHEN 'common' THEN 100
        WHEN 'uncommon' THEN 50
        WHEN 'rare' THEN 20
        WHEN 'epic' THEN 5
        WHEN 'legendary' THEN 1
    END
FROM loot_tables lt
CROSS JOIN items i
WHERE i.item_type IN ('consumable', 'material', 'token')
ON CONFLICT DO NOTHING;
