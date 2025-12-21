-- Migration: Add gacha fields to characters and users
-- File: 027_add_gacha_fields_to_characters_users.sql

-- Add to characters table
ALTER TABLE characters ADD COLUMN IF NOT EXISTS unique_name VARCHAR(100) UNIQUE;
ALTER TABLE characters ADD COLUMN IF NOT EXISTS unlocked_abilities JSONB DEFAULT '[]'::jsonb;
ALTER TABLE characters ADD COLUMN IF NOT EXISTS is_fainted BOOLEAN DEFAULT false;
ALTER TABLE characters ADD COLUMN IF NOT EXISTS fainted_at TIMESTAMP;

-- Add to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS has_minted_first_char BOOLEAN DEFAULT false;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_characters_unique_name ON characters(unique_name) WHERE unique_name IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_characters_fainted ON characters(is_fainted) WHERE is_fainted = true;
