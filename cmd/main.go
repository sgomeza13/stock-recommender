package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
	"github.com/sgomeza13/stock-recommender/api/routes"
	"github.com/sgomeza13/stock-recommender/config"
	"github.com/sgomeza13/stock-recommender/db"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	defer config.CloseDB() // Close the database connection when the server exits

	db.RunMigrations()

	router := gin.Default()
	// Apply CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Change to your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.SetTrustedProxies(nil)

	routes.RegisterRoutes(router)

	port := config.GetPort()

	log.Printf("Server is running on port %s", port)

	if err := router.Run(port); err != nil {
		log.Fatal("Server failed to start", err)
	}

}
