package main

import (
	"log"
	"os"
	"stonk-watcher/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func main() {
	err := os.Mkdir("data", 0600)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Use(func(context *gin.Context) {
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			Debug:            false,
		}).HandlerFunc(context.Writer, context.Request)
	})
	router.GET("/stock", handlers.StockHandler)
	router.DELETE("/stock", handlers.TruncateStockInfo)
	router.GET("/stock/price", handlers.StockPriceHandler)
	router.GET("/watchlist", handlers.WatchListHandler)
	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
