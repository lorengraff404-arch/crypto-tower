package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/lorengraff/crypto-tower-defense/internal/db"
	"github.com/lorengraff/crypto-tower-defense/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MarketplaceService handles marketplace operations
type MarketplaceService struct {
	ledger *LedgerService // Ledger Integration
}

// NewMarketplaceService creates a new service (helper internal)
// Ideally this should be a proper constructor called in main.go, but for now we follow the pattern
func (s *MarketplaceService) init() {
	if s.ledger == nil {
		s.ledger = NewLedgerService()
	}
}

// CreateListing creates a new marketplace listing
func (s *MarketplaceService) CreateListing(userID uint, itemType string, itemID uint, price int, currency string) (*models.MarketplaceListing, error) {
	// Verify ownership
	if err := s.verifyOwnership(userID, itemType, itemID); err != nil {
		return nil, err
	}

	// Check for existing active listing to prevent duplicates
	var existingListing models.MarketplaceListing
	duplicateQuery := db.DB.Where("status = ?", "ACTIVE")

	switch itemType {
	case "character":
		duplicateQuery = duplicateQuery.Where("character_id = ?", itemID)
	case "item", "equipment":
		duplicateQuery = duplicateQuery.Where("item_id = ?", itemID)
	}

	if err := duplicateQuery.First(&existingListing).Error; err == nil {
		return nil, errors.New("this item is already listed on the marketplace")
	}

	// Create listing
	listing := models.MarketplaceListing{
		SellerID:  userID,
		AssetType: itemType,
		Price:     int64(price),
		Status:    "ACTIVE",
	}

	// Set appropriate ID field based on type
	switch itemType {
	case "character":
		listing.CharacterID = &itemID
	case "item", "equipment":
		listing.ItemID = &itemID
	default:
		return nil, errors.New("invalid item type")
	}

	if err := db.DB.Create(&listing).Error; err != nil {
		return nil, err
	}

	// Mark item as listed (prevent double listing)
	s.markItemAsListed(itemType, itemID, true)

	return &listing, nil
}

// BuyListing purchases an item from marketplace
func (s *MarketplaceService) BuyListing(buyerID, listingID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var listing models.MarketplaceListing
		// 1. Lock listing row to prevent double-buy race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&listing, listingID).Error; err != nil {
			return err
		}

		if listing.Status != "ACTIVE" {
			return errors.New("listing is no longer active")
		}

		if listing.SellerID == buyerID {
			return errors.New("cannot buy your own listing")
		}

		// 2. Check buyer funds
		var buyer models.User
		if err := tx.First(&buyer, buyerID).Error; err != nil {
			return err
		}

		if int64(buyer.GTKBalance) < listing.Price {
			return errors.New("insufficient tokens")
		}

		// 3. Ledger Transfer (Atomic)
		feePercent := 0.05
		feeAmount := int64(float64(listing.Price) * feePercent)
		sellerAmount := listing.Price - feeAmount

		if s.ledger == nil {
			s.ledger = NewLedgerService()
		}

		// Fetch accounts using global DB (Locking not strictly needed for fetch, only updates)
		buyerAcc, err := s.ledger.GetOrCreateAccount(&buyerID, models.AccountTypeWallet, "GTK")
		if err != nil {
			return err
		}
		sellerAcc, err := s.ledger.GetOrCreateAccount(&listing.SellerID, models.AccountTypeWallet, "GTK")
		if err != nil {
			return err
		}
		treasuryAcc, err := s.ledger.GetOrCreateAccount(nil, models.AccountTypeTreasury, "GTK")
		if err != nil {
			return err
		}

		entries := []models.LedgerEntry{
			{AccountID: buyerAcc.ID, Amount: -listing.Price, Type: "DEBIT"},
			{AccountID: sellerAcc.ID, Amount: sellerAmount, Type: "CREDIT"},
			{AccountID: treasuryAcc.ID, Amount: feeAmount, Type: "CREDIT"},
		}

		// Use the transaction-aware ledger method
		if err := s.ledger.CreateTransactionWithTx(tx, models.TxTypeMarketBuy, fmt.Sprintf("market_buy_%d", listing.ID), fmt.Sprintf("Market purchase: Listing #%d", listing.ID), entries); err != nil {
			return fmt.Errorf("ledger transaction failed: %v", err)
		}

		// 4. Legacy Balance Updates (DB Sync - keep in same tx)
		if err := tx.Exec("UPDATE users SET tokens = tokens - ? WHERE id = ?", listing.Price, buyerID).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET tokens = tokens + ? WHERE id = ?", sellerAmount, listing.SellerID).Error; err != nil {
			return err
		}

		// 5. Transfer Ownership & Update Status
		itemID := uint(0)
		if listing.CharacterID != nil {
			itemID = *listing.CharacterID
		} else if listing.ItemID != nil {
			itemID = *listing.ItemID
		}

		// Update Ownership via TX
		switch listing.AssetType {
		case "character":
			if err := tx.Exec("UPDATE characters SET owner_id = ?, is_listed = false WHERE id = ?", buyerID, itemID).Error; err != nil {
				return err
			}
		case "equipment", "item":
			if err := tx.Exec("UPDATE items SET owner_id = ?, is_listed = false WHERE id = ?", buyerID, itemID).Error; err != nil {
				return err
			}
		case "egg":
			if err := tx.Exec("UPDATE eggs SET user_id = ? WHERE id = ?", buyerID, itemID).Error; err != nil {
				return err
			}
		}

		// Update Listing Status
		now := time.Now()
		listing.Status = "SOLD"
		listing.BuyerID = &buyerID
		listing.SoldAt = &now
		if err := tx.Save(&listing).Error; err != nil {
			return err
		}

		// 6. Create Trade History
		history := models.TradeHistory{
			ListingID: &listing.ID,
			SellerID:  listing.SellerID,
			BuyerID:   buyerID,
			ItemType:  listing.AssetType,
			ItemID:    itemID,
			Price:     listing.Price,
			Currency:  "TOWER",
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	})
}

// CancelListing cancels an active listing
func (s *MarketplaceService) CancelListing(userID, listingID uint) error {
	var listing models.MarketplaceListing
	if err := db.DB.First(&listing, listingID).Error; err != nil {
		return err
	}

	if listing.SellerID != userID {
		return errors.New("not your listing")
	}

	if listing.Status != "ACTIVE" {
		return errors.New("listing is not active")
	}

	listing.Status = "CANCELLED"
	db.DB.Save(&listing)

	// Get item ID
	var itemID uint
	if listing.CharacterID != nil {
		itemID = *listing.CharacterID
	} else if listing.ItemID != nil {
		itemID = *listing.ItemID
	}

	// Unmark item as listed
	s.markItemAsListed(listing.AssetType, itemID, false)

	return nil
}

// GetActiveListings returns all active marketplace listings with full character data
func (s *MarketplaceService) GetActiveListings(itemType string, limit, offset int) ([]models.MarketplaceListing, error) {
	var listings []models.MarketplaceListing
	query := db.DB.Preload("Character", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("Item").Preload("Seller", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Where("status = ?", "ACTIVE")

	// Simple query: get all active listings
	// No hidden filtering, no deduplication. Shows exactly what is in the DB.
	if itemType != "" {
		query = query.Where("asset_type = ?", itemType)
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&listings).Error; err != nil {
		return nil, err
	}

	// MANUAL RECOVERY for Soft-Deleted Data (Fix for 'Unknown' assets)
	for i := range listings {
		// Recover Character
		if (listings[i].AssetType == "character" || listings[i].AssetType == "Character") && listings[i].Character == nil && listings[i].CharacterID != nil {
			var char models.Character
			if err := db.DB.Unscoped().First(&char, *listings[i].CharacterID).Error; err == nil {
				listings[i].Character = &char
				fmt.Printf("♻️ RECOVERED character for Listing %d: %s\n", listings[i].ID, char.Name)
			} else {
				fmt.Printf("❌ FAILED to recover character %d: %v\n", *listings[i].CharacterID, err)
			}
		}

		// Recover Seller
		if listings[i].Seller.ID == 0 && listings[i].SellerID > 0 {
			var seller models.User
			if err := db.DB.Unscoped().First(&seller, listings[i].SellerID).Error; err == nil {
				listings[i].Seller = seller
			}
		}
	}

	return listings, nil
}

// verifyOwnership checks if user owns the item
func (s *MarketplaceService) verifyOwnership(userID uint, itemType string, itemID uint) error {
	switch itemType {
	case "character":
		var char models.Character
		if err := db.DB.First(&char, itemID).Error; err != nil {
			return err
		}
		if char.OwnerID != userID {
			return errors.New("you don't own this character")
		}
	case "equipment", "item":
		var item models.Item
		if err := db.DB.First(&item, itemID).Error; err != nil {
			return err
		}
		if item.OwnerID != userID {
			return errors.New("you don't own this item")
		}
	case "egg":
		var egg models.Egg
		if err := db.DB.First(&egg, itemID).Error; err != nil {
			return err
		}
		if egg.UserID != userID {
			return errors.New("you don't own this egg")
		}
	}
	return nil
}

// transferOwnership transfers item to new owner
func (s *MarketplaceService) transferOwnership(itemType string, itemID, fromUserID, toUserID uint) {
	// Verify ownership before transfer (security check)
	if err := s.verifyOwnership(fromUserID, itemType, itemID); err != nil {
		// Log error but don't fail the transaction (already verified earlier)
		// This is a double-check for security
		return
	}

	switch itemType {
	case "character":
		db.DB.Exec("UPDATE characters SET owner_id = ? WHERE id = ? AND owner_id = ?", toUserID, itemID, fromUserID)
	case "equipment", "item":
		db.DB.Exec("UPDATE items SET owner_id = ? WHERE id = ? AND owner_id = ?", toUserID, itemID, fromUserID)
	case "egg":
		db.DB.Exec("UPDATE eggs SET user_id = ? WHERE id = ? AND user_id = ?", toUserID, itemID, fromUserID)
	}
}

// markItemAsListed marks item as listed/unlisted
func (s *MarketplaceService) markItemAsListed(itemType string, itemID uint, listed bool) {
	switch itemType {
	case "character":
		db.DB.Exec("UPDATE characters SET is_listed = ? WHERE id = ?", listed, itemID)
	case "equipment", "item":
		db.DB.Exec("UPDATE items SET is_listed = ? WHERE id = ?", listed, itemID)
	}
}
