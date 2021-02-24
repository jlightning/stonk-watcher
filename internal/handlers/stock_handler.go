package handlers

import (
	"net/http"
	"stonk-watcher/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func StockHandler(c *gin.Context) {
	ticker := c.Query("ticker")
	if len(ticker) == 0 {
		c.String(http.StatusNotFound, "ticker not found")
		return
	}

	resp, err := services.GetStockInformation(strings.ToUpper(ticker))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func StockPriceHandler(c *gin.Context) {
	ticker := c.Query("ticker")
	if len(ticker) == 0 {
		c.String(http.StatusNotFound, "ticker not found")
		return
	}

	resp, err := services.GetDataFromFinviz(strings.ToUpper(ticker))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}