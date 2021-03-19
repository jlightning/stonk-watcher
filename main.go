package main

import (
	"embed"
	"stonk-watcher/internal/services"
)

//go:embed assets/src/build/**/* assets/src/build/*
var staticFileFS embed.FS

func main() {
	//router := gin.Default()
	//router.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"*"},
	//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	//	AllowHeaders:     []string{"Origin"},
	//	ExposeHeaders:    []string{"Content-Length"},
	//	AllowCredentials: true,
	//	MaxAge:           12 * time.Hour,
	//}))
	//router.GET("/stock", handlers.StockHandler)
	//router.DELETE("/stock", handlers.TruncateStockInfo)
	//router.GET("/stock/price", handlers.StockPriceHandler)
	//router.GET("/watchlist", handlers.GetWatchlistHandler)
	//router.POST("/watchlist", handlers.UpdateWatchlistHandler)
	//router.GET("/", handlers.StaticHandler(staticFileFS))
	//router.GET("/static/*file", handlers.StaticHandler(staticFileFS))
	//err := router.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
	_, err := services.GetDataFromMorningstar("ADBE")
	if err != nil {
		panic(err)
	}
}
