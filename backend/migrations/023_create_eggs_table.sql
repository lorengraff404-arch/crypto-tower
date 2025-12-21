-- Migration 023: Create eggs table for gacha system
CREATE TABLE IF NOT EXISTS eggs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent1_id INTEGER REFERENCES characters(id),
    parent2_id INTEGER REFERENCES characters(id),
    
    -- Gacha fields
    mint_cost BIGINT DEFAULT 0,
    mint_tx_hash VARCHAR(66),
    
    -- Predetermined traits
    rarity VARCHAR(5) NOT NULL,
    element VARCHAR(20) NOT NULL,
    character_type VARCHAR(20) NOT NULL,
    class VARCHAR(20) NOT NULL,
    predetermined_stats TEXT,
    predetermined_abilities TEXT,
    
    -- Reveal status
    is_stats_revealed BOOLEAN DEFAULT FALSE,
    revealed_at TIMESTAMP,
    
    -- Incubation
    incubation_time INTEGER NOT NULL,
    effective_incubation_time INTEGER,
    incubation_started_at TIMESTAMP,
    accelerators_applied TEXT DEFAULT '[]',
    
    -- Hatching
    hatched_at TIMESTAMP,
    character_id INTEGER REFERENCES characters(id),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_eggs_user ON eggs(user_id);
CREATE INDEX idx_eggs_hatched ON eggs(hatched_at);
