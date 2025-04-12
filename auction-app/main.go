package main

import (
	"log"

	"auction-app/config"
	"auction-app/models"
	"auction-app/routes"
	"auction-app/services"
)

func main() {
	// Initialize PostgreSQL and Redis connections.
	config.InitDB()
	config.InitRedis()

	// Auto-migrate models.
	if err := config.DB.AutoMigrate(&models.User{}, &models.Item{}, &models.Bid{}); err != nil {
		log.Fatal("Failed to migrate models:", err)
	}

	// Start background auction watcher.
	go services.AuctionWatcher()

	// Setup Gin router with defined routes.
	router := routes.SetupRouter()

	// Start the server.
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server Run Failed:", err)
	}
}
