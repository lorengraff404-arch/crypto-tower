-- Phase 21: Daily Systems
CREATE TABLE IF NOT EXISTS daily_quests (
    id SERIAL PRIMARY KEY,
    quest_type VARCHAR(30) NOT NULL,
    description TEXT NOT NULL,
    requirement_value INT NOT NULL,
    reward_tokens INT DEFAULT 0,
    reward_xp INT DEFAULT 0,
    reward_items TEXT,
    active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS user_daily_progress (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quest_id INT NOT NULL REFERENCES daily_quests(id),
    progress INT DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    date DATE DEFAULT CURRENT_DATE,
    UNIQUE(user_id, quest_id, date)
);

CREATE TABLE IF NOT EXISTS login_rewards (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    day_number INT NOT NULL,
    claimed_at TIMESTAMP DEFAULT NOW(),
    reward_tokens INT,
    reward_items TEXT
);

CREATE INDEX IF NOT EXISTS idx_daily_progress_user ON user_daily_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_daily_progress_date ON user_daily_progress(date);
CREATE INDEX IF NOT EXISTS idx_login_rewards_user ON login_rewards(user_id);

-- Insert base daily quests
INSERT INTO daily_quests (quest_type, description, requirement_value, reward_tokens, reward_xp) VALUES
('win_battles', 'Win 3 battles', 3, 100, 50),
('defeat_enemies', 'Defeat 10 enemies', 10, 50, 25),
('use_characters', 'Use 5 different characters', 5, 150, 75),
('breed_character', 'Breed 1 character', 1, 200, 100),
('upgrade_equipment', 'Upgrade 1 equipment', 1, 100, 50)
ON CONFLICT DO NOTHING;
