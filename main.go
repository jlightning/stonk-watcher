package main

import "stonk-watcher/internal/services"

func main() {
	//router := gin.Default()
	//router.GET("/stock", handlers.StockHandler)
	//err := router.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}

	err := services.GetStockInformation("ADBE")
	if err != nil {
		panic(err)
	}
}
