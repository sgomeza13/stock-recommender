package main

import (
	"log"

	"github.com/gin-gonic/gin"

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
	router.SetTrustedProxies(nil)

	routes.RegisterRoutes(router)

	port := config.GetPort()

	log.Printf("Server is running on port %s", port)

	if err := router.Run(port); err != nil {
		log.Fatal("Server failed to start", err)
	}

}
