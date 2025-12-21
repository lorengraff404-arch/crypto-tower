-- Daily Quest System Migration
-- Creates tables for daily quests and progress tracking

-- Main daily quests table
CREATE TABLE IF NOT EXISTS daily_quests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quest_type VARCHAR(50) NOT NULL, -- combat, collection, progression, special
    quest_name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    target_value INTEGER NOT NULL,
    current_progress INTEGER DEFAULT 0,
    difficulty VARCHAR(20) NOT NULL, -- common, uncommon, rare, epic
    reward_item_id INTEGER, -- NULL if no item reward
    reward_gtk INTEGER DEFAULT 0,
    reward_tower DECIMAL(18, 6) DEFAULT 0,
    is_completed BOOLEAN DEFAULT FALSE,
    is_claimed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    claimed_at TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Quest progress tracking for analytics
CREATE TABLE IF NOT EXISTS quest_progress_tracking (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quest_id INTEGER NOT NULL REFERENCES daily_quests(id) ON DELETE CASCADE,
    action_type VARCHAR(50) NOT NULL, -- battle_won, egg_minted, item_purchased, etc
    progress_increment INTEGER DEFAULT 1,
    metadata JSONB, -- Additional context (battle_id, character_id, etc)
    tracked_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_daily_quests_user ON daily_quests(user_id);
CREATE INDEX IF NOT EXISTS idx_daily_quests_expires ON daily_quests(expires_at);
CREATE INDEX IF NOT EXISTS idx_daily_quests_active ON daily_quests(user_id, is_completed, expires_at);
CREATE INDEX IF NOT EXISTS idx_quest_tracking_user ON quest_progress_tracking(user_id);
CREATE INDEX IF NOT EXISTS idx_quest_tracking_quest ON quest_progress_tracking(quest_id);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_daily_quest_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for auto-updating timestamp
CREATE TRIGGER update_daily_quest_timestamp
    BEFORE UPDATE ON daily_quests
    FOR EACH ROW
    EXECUTE FUNCTION update_daily_quest_timestamp();

-- Function to clean up expired quests (run daily via cron)
CREATE OR REPLACE FUNCTION cleanup_expired_quests()
RETURNS void AS $$
BEGIN
    DELETE FROM daily_quests
    WHERE expires_at < NOW() - INTERVAL '7 days'
    AND is_claimed = TRUE;
END;
$$ LANGUAGE plpgsql;
