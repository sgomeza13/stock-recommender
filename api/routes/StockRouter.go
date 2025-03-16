package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sgomeza13/stock-recommender/api/controller"
	"github.com/sgomeza13/stock-recommender/api/repository"
	"github.com/sgomeza13/stock-recommender/api/service"
)

func RegisterStockRoutes(router *gin.Engine) {
	stockRepo := repository.NewStockRepository()
	stockService := service.NewStockService(stockRepo)
	stockController := controller.NewStockController(stockService)

	// ✅ Define route for getting all stocks
	router.GET("/stocks", stockController.GetAllStocks)

	// ✅ Define route for pagination
	router.GET("/stocksByPage", stockController.GetStocksPaginated)

	// ✅ Define route for creating stocks
	router.POST("/stocks", stockController.CreateStocks)

	// ✅ Define route for creating stock
	router.POST("/stock", stockController.CreateStock)

	// ✅ Define route for getting stock by id
	router.GET("/stock/:id", stockController.GetStockByID)

	// ✅ Define route for deleting stock by id
	router.DELETE("/stock/:id", stockController.DeleteStockByID)

	// ✅ Define route for updating stock by id
	router.PUT("/stock/:id", stockController.UpdateStockByID)

}
