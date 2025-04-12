package services

import (
	"log"
	"time"

	"auction-app/config"
	"auction-app/models"
)

// AuctionWatcher periodically checks for ended auctions and assigns winners.
func AuctionWatcher() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		var items []models.Item
		now := time.Now()
		// Find items with auctions ended and no winner assigned.
		if err := config.DB.Where("end_time <= ? AND winner_id IS NULL", now).Find(&items).Error; err != nil {
			log.Println("Auction watcher error:", err)
			continue
		}
		for _, item := range items {
			var highestBid models.Bid
			err := config.DB.Where("item_id = ?", item.ID).Order("bid_amount desc").First(&highestBid).Error
			if err != nil {
				log.Printf("No bids for item %d\n", item.ID)
				continue
			}
			item.WinnerID = &highestBid.UserID
			if err := config.DB.Save(&item).Error; err != nil {
				log.Println("Failed to update item winner:", err)
			} else {
				log.Printf("Auction ended for item %d, winner: %d with bid: %.2f\n", item.ID, highestBid.UserID, highestBid.BidAmount)
				// Insert notification logic if required.
			}
		}
	}
}
