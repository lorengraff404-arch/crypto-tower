-- Migration 030: Fix inventory table name
-- Created: 2024-12-18

-- Rename table
ALTER TABLE IF EXISTS user_inventory RENAME TO user_inventories;

-- Rename indexes if they exist
ALTER INDEX IF EXISTS idx_user_inventory_user RENAME TO idx_user_inventories_user;
ALTER INDEX IF EXISTS idx_user_inventory_item RENAME TO idx_user_inventories_item;
ALTER INDEX IF EXISTS user_inventory_user_id_item_id_key RENAME TO user_inventories_user_id_item_id_key;
ALTER INDEX IF EXISTS user_inventory_pkey RENAME TO user_inventories_pkey;
