-- Phase 17: Breeding & Egg System

-- Eggs table
CREATE TABLE IF NOT EXISTS eggs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent1_id INT REFERENCES characters(id),
    parent2_id INT REFERENCES characters(id),
    
    -- Egg properties
    rarity VARCHAR(20) NOT NULL CHECK (rarity IN ('common', 'uncommon', 'rare', 'epic', 'legendary')),
    element VARCHAR(20),
    character_type VARCHAR(20),
    
    -- Incubation
    incubation_time INT NOT NULL, -- seconds
    incubation_started_at TIMESTAMP,
    hatched_at TIMESTAMP,
    
    -- Result
    character_id INT REFERENCES characters(id),
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- Breeding sessions
CREATE TABLE IF NOT EXISTS breeding_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent1_id INT NOT NULL REFERENCES characters(id),
    parent2_id INT NOT NULL REFERENCES characters(id),
    
    -- Status
    status VARCHAR(20) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed', 'cancelled')),
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    
    -- Result
    egg_id INT REFERENCES eggs(id),
    
    -- Cost
    tokens_spent INT DEFAULT 0
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_eggs_user ON eggs(user_id);
CREATE INDEX IF NOT EXISTS idx_eggs_hatched ON eggs(hatched_at);
CREATE INDEX IF NOT EXISTS idx_breeding_user ON breeding_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_breeding_status ON breeding_sessions(status);

-- Update characters table to track breeding count
ALTER TABLE characters ADD COLUMN IF NOT EXISTS breed_count INT DEFAULT 0;
ALTER TABLE characters ADD COLUMN IF NOT EXISTS last_bred_at TIMESTAMP;
