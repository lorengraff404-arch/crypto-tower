-- Phase 19: Marketplace & Economy
CREATE TABLE IF NOT EXISTS marketplace_listings (
    id SERIAL PRIMARY KEY,
    seller_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_type VARCHAR(30) NOT NULL CHECK (item_type IN ('character', 'equipment', 'egg', 'item')),
    item_id INT NOT NULL,
    price INT NOT NULL CHECK (price > 0),
    currency VARCHAR(20) DEFAULT 'TOWER' CHECK (currency IN ('TOWER', 'GTK')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'sold', 'cancelled')),
    listed_at TIMESTAMP DEFAULT NOW(),
    sold_at TIMESTAMP,
    buyer_id INT REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS trade_history (
    id SERIAL PRIMARY KEY,
    listing_id INT REFERENCES marketplace_listings(id),
    seller_id INT NOT NULL,
    buyer_id INT NOT NULL,
    item_type VARCHAR(30) NOT NULL,
    item_id INT NOT NULL,
    price INT NOT NULL,
    currency VARCHAR(20) NOT NULL,
    completed_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_marketplace_seller ON marketplace_listings(seller_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_status ON marketplace_listings(status);
CREATE INDEX IF NOT EXISTS idx_marketplace_type ON marketplace_listings(item_type);
CREATE INDEX IF NOT EXISTS idx_trade_history_seller ON trade_history(seller_id);
CREATE INDEX IF NOT EXISTS idx_trade_history_buyer ON trade_history(buyer_id);
