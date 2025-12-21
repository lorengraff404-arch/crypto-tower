-- Add sprite sheet URL columns to characters table
-- Migration: add_sprite_fields
-- Created: 2025-12-20

ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_idle VARCHAR(255);
ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_walk VARCHAR(255);
ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_run VARCHAR(255);
ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_attack VARCHAR(255);
ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_skill VARCHAR(255);
ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_hit VARCHAR(255);
ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_block VARCHAR(255);
ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_dodge VARCHAR(255);
ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_death VARCHAR(255);
ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_victory VARCHAR(255);
ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_gen_status VARCHAR(20) DEFAULT 'pending';
ALTER TABLE characters ADD COLUMN IF NOT EXISTS sprite_gen_job_id INTEGER;

-- Create index on sprite_gen_job_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_characters_sprite_gen_job_id ON characters(sprite_gen_job_id);

-- Create sprite_generation_jobs table
CREATE TABLE IF NOT EXISTS sprite_generation_jobs (
    id SERIAL PRIMARY KEY,
    character_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    progress INTEGER DEFAULT 0,
    error_msg TEXT,
    provider VARCHAR(50),
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

-- Create indexes for sprite_generation_jobs
CREATE INDEX IF NOT EXISTS idx_sprite_jobs_character_id ON sprite_generation_jobs(character_id);
CREATE INDEX IF NOT EXISTS idx_sprite_jobs_status ON sprite_generation_jobs(status);
CREATE INDEX IF NOT EXISTS idx_sprite_jobs_created_at ON sprite_generation_jobs(created_at);

-- Add foreign key constraint for sprite_gen_job_id
ALTER TABLE characters ADD CONSTRAINT fk_characters_sprite_gen_job 
    FOREIGN KEY (sprite_gen_job_id) REFERENCES sprite_generation_jobs(id) ON DELETE SET NULL;
