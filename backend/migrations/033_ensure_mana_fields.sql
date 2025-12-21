-- Migration: Add and initialize mana system fields
-- File: backend/migrations/033_ensure_mana_fields.sql

-- Add columns if they don't exist (PostgreSQL will skip if exists)
DO $$ 
BEGIN
    -- Add current_mana column
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'characters' AND column_name = 'current_mana'
    ) THEN
        ALTER TABLE characters ADD COLUMN current_mana INTEGER DEFAULT 100 NOT NULL;
    END IF;

    -- Add max_mana column
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'characters' AND column_name = 'max_mana'
    ) THEN
        ALTER TABLE characters ADD COLUMN max_mana INTEGER DEFAULT 100 NOT NULL;
    END IF;

    -- Add mana_regen_rate column
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'characters' AND column_name = 'mana_regen_rate'
    ) THEN
        ALTER TABLE characters ADD COLUMN mana_regen_rate INTEGER DEFAULT 10 NOT NULL;
    END IF;
END $$;

-- Update existing characters with NULL or 0 values
UPDATE characters 
SET current_mana = 100
WHERE current_mana IS NULL OR current_mana = 0;

UPDATE characters 
SET max_mana = 100
WHERE max_mana IS NULL OR max_mana = 0;

UPDATE characters 
SET mana_regen_rate = 10
WHERE mana_regen_rate IS NULL OR mana_regen_rate = 0;

-- Verify the update
SELECT 
    COUNT(*) as total_characters,
    COUNT(CASE WHEN current_mana > 0 THEN 1 END) as with_current_mana,
    COUNT(CASE WHEN max_mana > 0 THEN 1 END) as with_max_mana,
    COUNT(CASE WHEN mana_regen_rate > 0 THEN 1 END) as with_mana_regen
FROM characters;
