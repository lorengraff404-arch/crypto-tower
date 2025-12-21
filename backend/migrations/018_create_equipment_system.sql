-- Phase 18: Equipment System
CREATE TABLE IF NOT EXISTS equipment (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL REFERENCES items(id),
    slot VARCHAR(20) NOT NULL CHECK (slot IN ('weapon', 'armor', 'accessory')),
    required_level INT DEFAULT 1,
    required_class VARCHAR(20),
    bonus_attack INT DEFAULT 0,
    bonus_defense INT DEFAULT 0,
    bonus_hp INT DEFAULT 0,
    bonus_speed INT DEFAULT 0,
    special_effect VARCHAR(50),
    effect_value INT DEFAULT 0,
    upgrade_level INT DEFAULT 0,
    max_upgrade_level INT DEFAULT 5,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS character_equipment (
    id SERIAL PRIMARY KEY,
    character_id INT NOT NULL UNIQUE REFERENCES characters(id) ON DELETE CASCADE,
    weapon_id INT REFERENCES equipment(id),
    armor_id INT REFERENCES equipment(id),
    accessory_id INT REFERENCES equipment(id)
);

CREATE INDEX IF NOT EXISTS idx_equipment_item ON equipment(item_id);
CREATE INDEX IF NOT EXISTS idx_equipment_slot ON equipment(slot);
CREATE INDEX IF NOT EXISTS idx_char_equipment_char ON character_equipment(character_id);

-- Insert base equipment items
INSERT INTO items (name, description, item_type, rarity, icon_emoji, sell_price) VALUES
('Iron Sword', '+10 ATK', 'equipment', 'common', '⚔️', 100),
('Steel Armor', '+15 DEF', 'equipment', 'uncommon', '🛡️', 250),
('Speed Boots', '+5 SPD', 'equipment', 'common', '👟', 150),
('Dragon Blade', '+30 ATK, Fire damage', 'equipment', 'epic', '🗡️', 2000),
('Mage Staff', '+25 ATK, Magic boost', 'equipment', 'rare', '🪄', 800)
ON CONFLICT (name) DO NOTHING;

-- Create equipment entries
INSERT INTO equipment (item_id, slot, required_level, bonus_attack, bonus_defense, bonus_hp, bonus_speed)
SELECT i.id, 'weapon', 1, 10, 0, 0, 0 FROM items i WHERE i.name = 'Iron Sword'
ON CONFLICT DO NOTHING;

INSERT INTO equipment (item_id, slot, required_level, bonus_defense)
SELECT i.id, 'armor', 5, 15 FROM items i WHERE i.name = 'Steel Armor'
ON CONFLICT DO NOTHING;

INSERT INTO equipment (item_id, slot, required_level, bonus_speed)
SELECT i.id, 'accessory', 1, 5 FROM items i WHERE i.name = 'Speed Boots'
ON CONFLICT DO NOTHING;
