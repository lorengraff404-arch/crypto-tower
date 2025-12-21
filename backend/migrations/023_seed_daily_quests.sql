-- Seed daily quests
INSERT INTO daily_quests (quest_type, title, description, target, reward_tokens, active) VALUES
('battles_won', 'Win 3 Battles', 'Complete and win 3 battles today', 3, 100, true),
('enemies_defeated', 'Defeat 10 Enemies', 'Defeat a total of 10 enemies in battles', 10, 50, true),
('characters_used', 'Use 5 Different Characters', 'Use 5 different characters in battles', 5, 150, true),
('raids_completed', 'Complete 2 Raids', 'Successfully complete 2 island raids', 2, 200, true),
('perfect_battles', 'Win Without Losing HP', 'Win a battle without any character losing HP', 1, 300, true)
ON CONFLICT DO NOTHING;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_daily_quest_progress_user_date ON user_daily_progress(user_id, date);
CREATE INDEX IF NOT EXISTS idx_daily_quest_progress_quest ON user_daily_progress(quest_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_status ON marketplace_listings(status);
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_seller ON marketplace_listings(seller_id);
