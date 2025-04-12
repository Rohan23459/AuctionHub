package controllers

import (
	"net/http"
	"time"

	"auction-app/config"
	"auction-app/models"

	"github.com/gin-gonic/gin"
)

// CreateItem creates a new auction item.
func CreateItem(c *gin.Context) {
	var input struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		StartPrice  float64 `json:"start_price"`
		EndTime     string  `json:"end_time"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	istLoc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load IST timezone"})
		return
	}

	parsedTime, err := time.Parse(time.RFC3339, input.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format. Use RFC3339."})
		return
	}

	// Convert the parsed time to IST.
	endTimeIST := parsedTime.In(istLoc)

	item := models.Item{
		Title:       input.Title,
		Description: input.Description,
		StartPrice:  input.StartPrice,
		CreatedAt:   time.Now().In(istLoc),
		EndTime:     endTimeIST,
	}
	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item created", "item": item})
}

// ListItems returns all auction items.
func ListItems(c *gin.Context) {
	var items []models.Item
	if err := config.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
