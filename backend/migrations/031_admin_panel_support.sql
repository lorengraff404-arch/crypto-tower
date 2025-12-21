-- Migration: Admin Panel Support Tables
-- Created: 2025-12-20
-- Description: Creates tables for revenue tracking, admin actions, anti-cheat, and game configuration

-- Revenue Transactions Table
CREATE TABLE IF NOT EXISTS revenue_transactions (
    id SERIAL PRIMARY KEY,
    source VARCHAR(50) NOT NULL,
    amount_gtk DECIMAL(20,2) NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_revenue_source (source),
    INDEX idx_revenue_user (user_id),
    INDEX idx_revenue_created (created_at)
);

COMMENT ON TABLE revenue_transactions IS 'Tracks all GTK revenue from various sources';
COMMENT ON COLUMN revenue_transactions.source IS 'Source of revenue: shop, gacha, marketplace, breeding, battle_wager';
COMMENT ON COLUMN revenue_transactions.metadata IS 'Additional data: item_id, battle_id, etc.';

-- Admin Actions Log Table
CREATE TABLE IF NOT EXISTS admin_actions (
    id SERIAL PRIMARY KEY,
    admin_wallet VARCHAR(42) NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    target_user_id INT REFERENCES users(id) ON DELETE SET NULL,
    details JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_admin_wallet (admin_wallet),
    INDEX idx_admin_action_type (action_type),
    INDEX idx_admin_created (created_at)
);

COMMENT ON TABLE admin_actions IS 'Logs all admin actions for audit trail';
COMMENT ON COLUMN admin_actions.action_type IS 'Type: ban_user, adjust_balance, resolve_flag, update_config, etc.';

-- Anti-Cheat Flags Table
CREATE TABLE IF NOT EXISTS anti_cheat_flags (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    flag_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    battle_id INT REFERENCES battles(id) ON DELETE SET NULL,
    details JSONB NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'reviewing', 'resolved', 'false_positive')),
    resolved_at TIMESTAMP,
    resolved_by VARCHAR(42),
    resolution_notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_flag_user (user_id),
    INDEX idx_flag_type (flag_type),
    INDEX idx_flag_severity (severity),
    INDEX idx_flag_status (status),
    INDEX idx_flag_created (created_at)
);

COMMENT ON TABLE anti_cheat_flags IS 'Anti-cheat detection flags for suspicious activity';
COMMENT ON COLUMN anti_cheat_flags.flag_type IS 'Type: impossible_stats, rapid_completion, suspicious_winrate, invalid_transaction';
COMMENT ON COLUMN anti_cheat_flags.details IS 'Detection details: expected vs actual values, timestamps, etc.';

-- Game Configuration Table
CREATE TABLE IF NOT EXISTS game_config (
    id SERIAL PRIMARY KEY,
    config_key VARCHAR(100) UNIQUE NOT NULL,
    config_value JSONB NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT,
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by VARCHAR(42),
    
    INDEX idx_config_category (category),
    INDEX idx_config_key (config_key)
);

COMMENT ON TABLE game_config IS 'Centralized game configuration managed by admin';
COMMENT ON COLUMN game_config.category IS 'Category: gacha, shop, battle, breeding, sprite, etc.';
COMMENT ON COLUMN game_config.config_value IS 'JSON value: rates, prices, timers, etc.';

-- Insert default game configurations
INSERT INTO game_config (config_key, config_value, category, description) VALUES
('gacha.mint_cost', '{"gtk": 100}', 'gacha', 'Cost to mint a new egg'),
('gacha.incubation_time', '{"hours": 24}', 'gacha', 'Default incubation time in hours'),
('breeding.cost', '{"gtk": 500}', 'breeding', 'Cost to breed two characters'),
('breeding.incubation_time', '{"hours": 48}', 'breeding', 'Breeding incubation time'),
('sprite.provider', '{"active": "mock", "auto_generate": true}', 'sprite', 'Sprite generation provider and settings'),
('battle.wager_fee', '{"percentage": 5}', 'battle', 'Platform fee percentage on wager battles'),
('shop.discount_multiplier', '{"value": 1.0}', 'shop', 'Global shop discount multiplier')
ON CONFLICT (config_key) DO NOTHING;
