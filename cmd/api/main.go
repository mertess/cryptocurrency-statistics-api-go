package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mertess/cryptocurrency-statistics-api-go/internal/api/handlers"
)

func main() {

	router := gin.Default()

	router.GET("/last-price", handlers.GetLastPrice)

	router.Run()
}
