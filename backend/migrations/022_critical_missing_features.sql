-- Critical Missing Features & Security

-- Add soft delete support (best practice)
ALTER TABLE characters ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE items ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE eggs ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- Add audit trails (security & debugging)
CREATE TABLE IF NOT EXISTS audit_log (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INT,
    old_values TEXT,
    new_values TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Add rate limiting table (prevent abuse)
CREATE TABLE IF NOT EXISTS rate_limits (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    action_type VARCHAR(50) NOT NULL,
    count INT DEFAULT 1,
    window_start TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, action_type)
);

-- Add character evolution tracking
ALTER TABLE characters ADD COLUMN IF NOT EXISTS evolution_history TEXT; -- JSON array

-- Add equipment durability system
ALTER TABLE items ADD COLUMN IF NOT EXISTS current_durability INT DEFAULT 100;
ALTER TABLE items ADD COLUMN IF NOT EXISTS max_durability INT DEFAULT 100;

-- Add battle replay system
CREATE TABLE IF NOT EXISTS battle_replays (
    id SERIAL PRIMARY KEY,
    raid_session_id INT REFERENCES raid_sessions(id),
    replay_data TEXT NOT NULL, -- JSON battle log
    created_at TIMESTAMP DEFAULT NOW()
);

-- Add leaderboards
CREATE TABLE IF NOT EXISTS leaderboards (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    category VARCHAR(30) NOT NULL, -- 'pvp_rating', 'total_wins', 'breeding_count'
    score INT NOT NULL,
    rank INT,
    season INT DEFAULT 1,
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, category, season)
);

-- Add notifications system
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(30) NOT NULL,
    title VARCHAR(100),
    message TEXT,
    data TEXT, -- JSON
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Add friend system
CREATE TABLE IF NOT EXISTS friendships (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    friend_id INT NOT NULL REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'pending', -- pending, accepted, blocked
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, friend_id)
);

-- Add referral system
CREATE TABLE IF NOT EXISTS referrals (
    id SERIAL PRIMARY KEY,
    referrer_id INT NOT NULL REFERENCES users(id),
    referred_id INT NOT NULL REFERENCES users(id),
    reward_claimed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Add character favorites
CREATE TABLE IF NOT EXISTS character_favorites (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    character_id INT NOT NULL REFERENCES characters(id),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, character_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_audit_log_user ON audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_created ON audit_log(created_at);
CREATE INDEX IF NOT EXISTS idx_rate_limits_user ON rate_limits(user_id);
CREATE INDEX IF NOT EXISTS idx_leaderboards_category ON leaderboards(category, season);
CREATE INDEX IF NOT EXISTS idx_leaderboards_rank ON leaderboards(rank);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(read);
CREATE INDEX IF NOT EXISTS idx_friendships_user ON friendships(user_id);
CREATE INDEX IF NOT EXISTS idx_friendships_friend ON friendships(friend_id);
CREATE INDEX IF NOT EXISTS idx_battle_replays_session ON battle_replays(raid_session_id);
