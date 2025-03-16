package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sgomeza13/stock-recommender/api/controller"
)

func RegisterRoutes(router *gin.Engine) {
	helloRoutes(router)
	RegisterStockRoutes(router)
}

func helloRoutes(router *gin.Engine) {
	router.GET("/hello", controller.HelloHandler)
}
