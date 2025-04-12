// routes/routes.go
package routes

import (
	"auction-app/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter registers all routes and returns the Gin engine.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Enable CORS with default options (you can customize this if needed).
	router.Use(cors.Default())

	// Authentication routes.
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	// Auction item routes.
	router.POST("/items", controllers.CreateItem)
	router.GET("/items", controllers.ListItems)

	// Bid and order routes.
	router.POST("/bid", controllers.PlaceBid)
	router.POST("/fulfill", controllers.FulfillOrder)

	return router
}
