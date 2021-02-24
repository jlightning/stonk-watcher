package handlers

import (
	"net/http"
	"stonk-watcher/internal/services"

	"github.com/gin-gonic/gin"
)

func StockHandler(c *gin.Context) {
	err := services.GetStockInformation("MSFT")
	if err != nil {
		panic(err)
	}
	c.String(http.StatusOK, "test")
}
