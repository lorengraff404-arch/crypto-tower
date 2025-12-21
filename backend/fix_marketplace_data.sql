-- FIX MARKETPLACE DATA INTEGRITY

-- 1. Identify Orphaned Listings (Listings pointing to deleted/non-existent characters)
SELECT 'ORPHAN_CHECK' as check_type, count(*) as count 
FROM marketplace_listings 
WHERE character_id IS NOT NULL 
AND character_id NOT IN (SELECT id FROM characters WHERE deleted_at IS NULL);

-- 2. Delete Orphaned Listings
DELETE FROM marketplace_listings 
WHERE character_id IS NOT NULL 
AND character_id NOT IN (SELECT id FROM characters WHERE deleted_at IS NULL);

-- 3. Identify Duplicates (Same character listed multiple times)
SELECT 'DUPLICATE_CHECK' as check_type, character_id, count(*) 
FROM marketplace_listings 
WHERE status = 'ACTIVE' 
AND character_id IS NOT NULL 
GROUP BY character_id 
HAVING count(*) > 1;

-- 4. Delete Duplicates (Keep only the oldest/first listing)
DELETE FROM marketplace_listings
WHERE status = 'ACTIVE'
AND character_id IS NOT NULL
AND id NOT IN (
    SELECT MIN(id)
    FROM marketplace_listings
    WHERE status = 'ACTIVE' 
    AND character_id IS NOT NULL
    GROUP BY character_id
);

-- 5. Verify Clean State
SELECT * FROM marketplace_listings WHERE status = 'ACTIVE';
