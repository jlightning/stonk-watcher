package main

import (
	"embed"
	"log"
	"os"
	"stonk-watcher/internal/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed assets/src/build/**/* assets/src/build/*
var staticFileFS embed.FS

func main() {
	err := os.Mkdir("data", 0600)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/stock", handlers.StockHandler)
	router.DELETE("/stock", handlers.TruncateStockInfo)
	router.GET("/stock/price", handlers.StockPriceHandler)
	router.GET("/watchlist", handlers.GetWatchlistHandler)
	router.POST("/watchlist", handlers.UpdateWatchlistHandler)
	router.GET("/", handlers.StaticHandler(staticFileFS))
	router.GET("/static/*file", handlers.StaticHandler(staticFileFS))
	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
