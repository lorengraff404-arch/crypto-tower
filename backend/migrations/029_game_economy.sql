-- Migration: Add revenue distribution and game economy tables
-- Created: 2024-12-18

-- Revenue Distribution Table
CREATE TABLE IF NOT EXISTS revenue_distributions (
    id SERIAL PRIMARY KEY,
    tx_hash VARCHAR(66) NOT NULL UNIQUE,
    source VARCHAR(50) NOT NULL, -- 'shop', 'breeding', 'battle', etc.
    total_amount DECIMAL(20,8) NOT NULL,
    growth_fund DECIMAL(20,8) NOT NULL,     -- 10%
    security_fund DECIMAL(20,8) NOT NULL,   -- 10%
    operations DECIMAL(20,8) NOT NULL,      -- 5%
    rewards_pool DECIMAL(20,8) NOT NULL,    -- 30%
    dev_team DECIMAL(20,8) NOT NULL,        -- 20%
    tower_liquidity DECIMAL(20,8) NOT NULL, -- 25%
    created_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT check_percentages CHECK (
        ABS((growth_fund + security_fund + operations + rewards_pool + dev_team + tower_liquidity) - total_amount) < 0.00000001
    )
);

CREATE INDEX idx_revenue_source ON revenue_distributions(source);
CREATE INDEX idx_revenue_created ON revenue_distributions(created_at DESC);

-- Game Mode Configuration
CREATE TABLE IF NOT EXISTS battle_modes (
    id SERIAL PRIMARY KEY,
    mode VARCHAR(20) NOT NULL UNIQUE, -- 'free', 'ranked', 'wager'
    min_bet DECIMAL(20,8) DEFAULT 0,
    max_bet DECIMAL(20,8),
    gtk_reward DECIMAL(20,8) DEFAULT 0,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert default game modes
INSERT INTO battle_modes (mode, min_bet, max_bet, gtk_reward, enabled) VALUES
('free', 0, 0, 10, true),
('ranked', 0, 0, 50, true),
('wager', 1, 1000, 0, true)
ON CONFLICT (mode) DO NOTHING;

-- Battles Table
CREATE TABLE IF NOT EXISTS battles (
    id SERIAL PRIMARY KEY,
    mode VARCHAR(20) NOT NULL REFERENCES battle_modes(mode),
    player1_id INTEGER NOT NULL REFERENCES users(id),
    player2_id INTEGER REFERENCES users(id), -- NULL for PvE
    wager_amount DECIMAL(20,8) DEFAULT 0,
    escrow_tx VARCHAR(66),
    winner_id INTEGER REFERENCES users(id),
    replay_data JSONB,
    replay_checksum VARCHAR(64),
    anti_cheat_flags JSONB DEFAULT '[]'::jsonb,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'in_progress', 'completed', 'disputed', 'cancelled'
    result_data JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    
    CONSTRAINT check_wager_mode CHECK (
        (mode != 'wager') OR (wager_amount > 0 AND escrow_tx IS NOT NULL)
    ),
    CONSTRAINT check_completion CHECK (
        (status != 'completed') OR (winner_id IS NOT NULL AND completed_at IS NOT NULL)
    )
);

CREATE INDEX idx_battles_player1 ON battles(player1_id);
CREATE INDEX idx_battles_player2 ON battles(player2_id);
CREATE INDEX idx_battles_mode ON battles(mode);
CREATE INDEX idx_battles_status ON battles(status);
CREATE INDEX idx_battles_created ON battles(created_at DESC);

-- Anti-Cheat Flags Table
CREATE TABLE IF NOT EXISTS anti_cheat_flags (
    id SERIAL PRIMARY KEY,
    battle_id INTEGER NOT NULL REFERENCES battles(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    flag_type VARCHAR(50) NOT NULL, -- 'bot_detection', 'collusion', 'replay_invalid'
    severity VARCHAR(20) NOT NULL, -- 'low', 'medium', 'high', 'critical'
    details JSONB NOT NULL,
    auto_flagged BOOLEAN DEFAULT true,
    reviewed BOOLEAN DEFAULT false,
    reviewer_id INTEGER REFERENCES users(id),
    reviewed_at TIMESTAMP,
    action_taken VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_anti_cheat_battle ON anti_cheat_flags(battle_id);
CREATE INDEX idx_anti_cheat_user ON anti_cheat_flags(user_id);
CREATE INDEX idx_anti_cheat_severity ON anti_cheat_flags(severity);
CREATE INDEX idx_anti_cheat_reviewed ON anti_cheat_flags(reviewed);

-- Player ELO Ratings (for Ranked mode)
CREATE TABLE IF NOT EXISTS player_elo (
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    elo_rating INTEGER DEFAULT 1200,
    games_played INTEGER DEFAULT 0,
    wins INTEGER DEFAULT 0,
    losses INTEGER DEFAULT 0,
    win_streak INTEGER DEFAULT 0,
    highest_elo INTEGER DEFAULT 1200,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_elo_rating ON player_elo(elo_rating DESC);

-- Treasury Wallet Addresses Configuration
CREATE TABLE IF NOT EXISTS treasury_wallets (
    id SERIAL PRIMARY KEY,
    wallet_type VARCHAR(50) NOT NULL UNIQUE,
    address VARCHAR(42) NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert treasury wallet placeholders (to be updated with real addresses)
INSERT INTO treasury_wallets (wallet_type, address, description) VALUES
('growth_fund', '0x02e4bb10328eB9f38137c142BD3CAD3d5cC5a9BD', 'Growth and marketing fund'),
('security_fund', '0x87Af9Dd54C4358D3B60c9568ea0e943E9147b464', 'Security and audit fund'),
('operations', '0x105FfE090e3cc6Da33067fAF72c98127dB3daE91', 'Operations and infrastructure'),
('rewards_pool', '0x3A405d7374D2fF47baeA4d8798412736c6BC7bDe', 'Top 10 players rewards pool'),
('dev_team', '0xd210925940D236951F0Bfa632A89d76fA1C8883d', 'Development team'),
('tower_liquidity', '0x67E29D8d8f6F69d10D0965CBB59C750119C3b8f8', 'TOWER token liquidity pool')
ON CONFLICT (wallet_type) DO NOTHING;

-- Battle History for analytics
CREATE TABLE IF NOT EXISTS battle_history (
    id SERIAL PRIMARY KEY,
    battle_id INTEGER NOT NULL REFERENCES battles(id),
    player_id INTEGER NOT NULL REFERENCES users(id),
    elo_before INTEGER,
    elo_after INTEGER,
    elo_change INTEGER,
    gtk_reward DECIMAL(20,8) DEFAULT 0,
    tower_reward DECIMAL(20,8) DEFAULT 0,
    is_winner BOOLEAN,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_battle_history_player ON battle_history(player_id);
CREATE INDEX idx_battle_history_battle ON battle_history(battle_id);
