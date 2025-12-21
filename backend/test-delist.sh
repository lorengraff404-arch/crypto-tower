#!/bin/bash

echo "🧪 Testing Marketplace Delist Functionality"
echo "==========================================="
echo ""

# Database connection
DB="tower_defense_dev"
USER="lorengraff"

echo "📊 Current State of Listings:"
psql -U $USER -d $DB -c "
SELECT 
    ml.id as listing_id,
    ml.status,
    c.id as char_id,
    c.name,
    c.is_listed,
    ml.seller_id
FROM marketplace_listings ml
JOIN characters c ON ml.character_id = c.id
WHERE ml.id IN (9, 10, 11)
ORDER BY ml.id;
"

echo ""
echo "🔓 Simulating Delist for Listing #9 (Warrior Alpha)..."
echo "   This would normally be done via DELETE /api/v1/marketplace/9"
echo ""

# Simulate what the CancelListing function does
psql -U $USER -d $DB -c "
BEGIN;
-- Update listing status
UPDATE marketplace_listings SET status = 'CANCELLED' WHERE id = 9;
-- Unlock character
UPDATE characters SET is_listed = false WHERE id = 6;
COMMIT;
"

echo ""
echo "✅ After Delist - Listing #9:"
psql -U $USER -d $DB -c "
SELECT 
    ml.id as listing_id,
    ml.status,
    c.id as char_id,
    c.name,
    c.is_listed as 'Available for Teams?'
FROM marketplace_listings ml
JOIN characters c ON ml.character_id = c.id
WHERE ml.id = 9;
"

echo ""
echo "📝 Summary:"
echo "   ✓ Listing status changed to CANCELLED"
echo "   ✓ Character is_listed set to FALSE"
echo "   ✓ Character is now available for teams, battles, and raids"
echo ""
echo "🔄 To restore for testing, run:"
echo "   psql -U $USER -d $DB -c \"UPDATE marketplace_listings SET status = 'ACTIVE' WHERE id = 9; UPDATE characters SET is_listed = true WHERE id = 6;\""
