package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WatchListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, []string{
		"MSFT",
		"ADBE",
		"AMZN",
		"GOOG",
		"TSLA",
		"VMW",
		"XOM",
		"PFE",
		"NFLX",
		"FTI",
		"GME",
	})
}
