-- Migration: Add gacha system fields to eggs table
-- File: 024_add_gacha_fields_to_eggs.sql

-- Add mint information
ALTER TABLE eggs ADD COLUMN IF NOT EXISTS mint_cost BIGINT NOT NULL DEFAULT 1;
ALTER TABLE eggs ADD COLUMN IF NOT EXISTS mint_tx_hash VARCHAR(66);

-- Add predetermined traits (visible on mint)
ALTER TABLE eggs ADD COLUMN IF NOT EXISTS class VARCHAR(20);

-- Add predetermined stats (hidden until hatch or scanned)
ALTER TABLE eggs ADD COLUMN IF NOT EXISTS predetermined_stats JSONB;
ALTER TABLE eggs ADD COLUMN IF NOT EXISTS predetermined_abilities JSONB;

-- Add stats reveal tracking
ALTER TABLE eggs ADD COLUMN IF NOT EXISTS is_stats_revealed BOOLEAN DEFAULT false;
ALTER TABLE eggs ADD COLUMN IF NOT EXISTS revealed_at TIMESTAMP;

-- Add accelerator tracking
ALTER TABLE eggs ADD COLUMN IF NOT EXISTS accelerators_applied JSONB DEFAULT '[]'::jsonb;
ALTER TABLE eggs ADD COLUMN IF NOT EXISTS effective_incubation_time INT;

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_eggs_user_id_created ON eggs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_eggs_hatched_at ON eggs(hatched_at) WHERE hatched_at IS NULL;
