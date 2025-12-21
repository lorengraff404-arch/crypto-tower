-- Remove duplicate marketplace listings
-- Keep only the oldest listing for each character

-- Step 1: Cancel duplicate listings (keep oldest one per character)
UPDATE marketplace_listings
SET status = 'CANCELLED'
WHERE id NOT IN (
    SELECT MIN(id)
    FROM marketplace_listings
    WHERE status = 'ACTIVE' AND character_id IS NOT NULL
    GROUP BY character_id
)
AND status = 'ACTIVE'
AND character_id IS NOT NULL;

-- Step 2: Add unique constraint to prevent future duplicates
-- This ensures only ONE active listing per character at a time
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_active_character_listing 
ON marketplace_listings (character_id) 
WHERE status = 'ACTIVE' AND character_id IS NOT NULL;

-- Step 3: Add unique constraint for items too
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_active_item_listing 
ON marketplace_listings (item_id) 
WHERE status = 'ACTIVE' AND item_id IS NOT NULL;

-- Verify cleanup
SELECT character_id, COUNT(*) as listing_count
FROM marketplace_listings
WHERE status = 'ACTIVE' AND character_id IS NOT NULL
GROUP BY character_id
HAVING COUNT(*) > 1;
-- Should return 0 rows after cleanup
