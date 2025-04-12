package controllers

import (
	"net/http"
	"strconv"
	"time"

	"auction-app/config"
	"auction-app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// PlaceBid allows a user to place a bid on an auction item.
func PlaceBid(c *gin.Context) {
	var input struct {
		UserID    uint    `json:"user_id"`
		ItemID    uint    `json:"item_id"`
		BidAmount float64 `json:"bid_amount"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the auction item.
	var item models.Item
	if err := config.DB.First(&item, input.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	// Ensure the auction hasn't ended.
	if time.Now().After(item.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Auction has ended"})
		return
	}

	// Get current highest bid from Redis (or fallback to DB).
	highestBidKey := "item:" + strconv.Itoa(int(item.ID)) + ":highest_bid"
	highestBidStr, err := config.Rdb.Get(config.Ctx, highestBidKey).Result()
	var highestBid float64
	if err == redis.Nil {
		var bid models.Bid
		if err := config.DB.Where("item_id = ?", item.ID).Order("bid_amount desc").First(&bid).Error; err == nil {
			highestBid = bid.BidAmount
		} else {
			highestBid = item.StartPrice
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
		return
	} else {
		highestBid, err = strconv.ParseFloat(highestBidStr, 64)
		if err != nil {
			highestBid = item.StartPrice
		}
	}

	if input.BidAmount <= highestBid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bid amount must be higher than the current highest bid"})
		return
	}

	// Create new bid record.
	newBid := models.Bid{
		UserID:    input.UserID,
		ItemID:    input.ItemID,
		BidAmount: input.BidAmount,
		CreatedAt: time.Now(),
	}
	if err := config.DB.Create(&newBid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the highest bid in Redis.
	if err := config.Rdb.Set(config.Ctx, highestBidKey, input.BidAmount, time.Minute*10).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Redis"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bid placed successfully", "bid": newBid})
}

// FulfillOrder simulates order fulfillment for a completed auction.
func FulfillOrder(c *gin.Context) {
	var input struct {
		ItemID uint `json:"item_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var item models.Item
	if err := config.DB.First(&item, input.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	if item.WinnerID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Auction not ended or no winner"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order fulfilled", "item": item})
}
