package main

import (
	"log"
	"stonk-watcher/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/stock", handlers.StockHandler)
	router.GET("/stock/price", handlers.StockPriceHandler)
	err := router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
