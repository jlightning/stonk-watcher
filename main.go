package main

import (
	"log"
	"os"
	"stonk-watcher/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	err := os.Mkdir("data", 0600)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/stock", handlers.StockHandler)
	router.DELETE("/stock", handlers.TruncateStockInfo)
	router.GET("/stock/price", handlers.StockPriceHandler)
	router.GET("/watchlist", handlers.WatchListHandler)
	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
