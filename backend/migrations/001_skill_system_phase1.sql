-- ============================================
-- PHASE 1: Mana System & Skill Slots Migration
-- ============================================

-- 1. Add Mana System to Characters
ALTER TABLE characters 
ADD COLUMN IF NOT EXISTS current_mana INT DEFAULT 100,
ADD COLUMN IF NOT EXISTS max_mana INT DEFAULT 100,
ADD COLUMN IF NOT EXISTS mana_regen_rate INT DEFAULT 10;

-- Update existing characters with rarity-based mana
UPDATE characters SET 
    max_mana = CASE rarity
        WHEN 'C' THEN 80
        WHEN 'B' THEN 100
        WHEN 'A' THEN 120
        WHEN 'S' THEN 150
        WHEN 'SS' THEN 200
        WHEN 'SSS' THEN 300
        ELSE 100
    END,
    mana_regen_rate = CASE rarity
        WHEN 'C' THEN 8
        WHEN 'B' THEN 10
        WHEN 'A' THEN 12
        WHEN 'S' THEN 15
        WHEN 'SS' THEN 20
        WHEN 'SSS' THEN 30
        ELSE 10
    END,
    current_mana = max_mana
WHERE current_mana IS NULL OR max_mana IS NULL;

-- 2. Create Active Skill Slots Table
CREATE TABLE IF NOT EXISTS character_active_skills (
    id SERIAL PRIMARY KEY,
    character_id INT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    ability_id INT NOT NULL REFERENCES abilities(id) ON DELETE CASCADE,
    slot_position INT NOT NULL CHECK (slot_position BETWEEN 1 AND 7),
    is_locked BOOLEAN DEFAULT true,
    unlock_level INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(character_id, slot_position),
    UNIQUE(character_id, ability_id)
);

CREATE INDEX IF NOT EXISTS idx_active_skills_char ON character_active_skills(character_id);
CREATE INDEX IF NOT EXISTS idx_active_skills_ability ON character_active_skills(ability_id);

-- 3. Update Abilities Table with Rarity and Synergy
ALTER TABLE abilities
ADD COLUMN IF NOT EXISTS rarity VARCHAR(10) DEFAULT 'C',
ADD COLUMN IF NOT EXISTS max_pp INT DEFAULT 10,
ADD COLUMN IF NOT EXISTS is_ultimate BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS synergy_tags TEXT[],
ADD COLUMN IF NOT EXISTS required_element TEXT[],
ADD COLUMN IF NOT EXISTS required_class TEXT[];

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_abilities_rarity ON abilities(rarity);
CREATE INDEX IF NOT EXISTS idx_abilities_synergy ON abilities USING GIN(synergy_tags);
CREATE INDEX IF NOT EXISTS idx_abilities_class_array ON abilities USING GIN(required_class);
CREATE INDEX IF NOT EXISTS idx_abilities_element_array ON abilities USING GIN(required_element);

-- 4. Update existing abilities with rarity based on mana cost
UPDATE abilities SET 
    rarity = CASE 
        WHEN mana_cost <= 20 THEN 'C'
        WHEN mana_cost <= 40 THEN 'B'
        WHEN mana_cost <= 70 THEN 'A'
        WHEN mana_cost <= 100 THEN 'S'
        WHEN mana_cost <= 150 THEN 'SS'
        ELSE 'SSS'
    END
WHERE rarity IS NULL OR rarity = 'C';

-- 5. Create Skill Cooldown Tracking Table
CREATE TABLE IF NOT EXISTS character_skill_cooldowns (
    id SERIAL PRIMARY KEY,
    character_id INT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    ability_id INT NOT NULL REFERENCES abilities(id) ON DELETE CASCADE,
    cooldown_remaining INT DEFAULT 0,
    last_used_at TIMESTAMP,
    
    UNIQUE(character_id, ability_id)
);

CREATE INDEX IF NOT EXISTS idx_cooldowns_char ON character_skill_cooldowns(character_id);

-- 6. Add Battle-Related Fields to Abilities
ALTER TABLE abilities
ADD COLUMN IF NOT EXISTS damage_type VARCHAR(20) DEFAULT 'physical', -- physical, magical, true
ADD COLUMN IF NOT EXISTS aoe_radius INT DEFAULT 0, -- 0 = single target
ADD COLUMN IF NOT EXISTS status_effect_chance INT DEFAULT 0, -- percentage
ADD COLUMN IF NOT EXISTS buff_duration INT DEFAULT 0, -- turns
ADD COLUMN IF NOT EXISTS debuff_duration INT DEFAULT 0; -- turns

COMMENT ON COLUMN characters.current_mana IS 'Current mana available for skills';
COMMENT ON COLUMN characters.max_mana IS 'Maximum mana capacity (scales with rarity)';
COMMENT ON COLUMN characters.mana_regen_rate IS 'Mana regenerated per turn (higher for higher rarity)';
COMMENT ON TABLE character_active_skills IS 'Tracks which skills are active in battle slots';
COMMENT ON TABLE character_skill_cooldowns IS 'Tracks cooldown state of character skills';

-- Verify migration
SELECT 
    'Characters with Mana' as check_name,
    COUNT(*) as count 
FROM characters 
WHERE current_mana IS NOT NULL AND max_mana IS NOT NULL;

SELECT 
    'Abilities with Rarity' as check_name,
    rarity,
    COUNT(*) as count 
FROM abilities 
GROUP BY rarity 
ORDER BY rarity;
