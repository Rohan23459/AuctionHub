package models

import "time"

// User represents a platform user.
type User struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Username      string `gorm:"unique" json:"username"`
	Password      string `json:"password"` // In production, store a hashed password.
	ContactNumber string `gorm:"unique" json:"contactNumber"`
	EmailID       string `gorm:"unique" json:"emailID"`
}

// Item represents an auction item.
type Item struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartPrice  float64   `json:"start_price"`
	CreatedAt   time.Time `json:"created_at"`
	EndTime     time.Time `json:"end_time"`
	WinnerID    *uint     `json:"winner_id"`
}

// Bid represents a bid placed by a user.
type Bid struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	ItemID    uint      `json:"item_id"`
	BidAmount float64   `json:"bid_amount"`
	CreatedAt time.Time `json:"created_at"`
}
