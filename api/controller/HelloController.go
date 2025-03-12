package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sgomeza13/stock-recommender/api/service"
)

func HelloHandler(c *gin.Context) {
	message := service.GetHelloMessage()

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
