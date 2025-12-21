-- Migration 025: Create shop system tables
-- Created: 2024-12-18

-- Shop items table
CREATE TABLE IF NOT EXISTS shop_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    category VARCHAR(20) NOT NULL,
    effect_type VARCHAR(20) NOT NULL,
    effect_value INTEGER,
    gtk_cost BIGINT NOT NULL,
    is_consumable BOOLEAN DEFAULT true,
    max_stack INTEGER DEFAULT 99,
    is_available BOOLEAN DEFAULT true,
    icon_url VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_shop_items_category ON shop_items(category);
CREATE INDEX idx_shop_items_available ON shop_items(is_available);

-- User inventory table
CREATE TABLE IF NOT EXISTS user_inventories (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES shop_items(id) ON DELETE CASCADE,
    quantity INTEGER DEFAULT 1 CHECK (quantity >= 0),
    acquired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, item_id)
);

CREATE INDEX idx_user_inventories_user ON user_inventories(user_id);
CREATE INDEX idx_user_inventories_item ON user_inventories(item_id);
