package handlers

import (
	"io/ioutil"
	"net/http"
	"stonk-watcher/internal/repositories"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

var defaultWatchlist = []string{
	"MSFT",
	"ADBE",
	"AMZN",
	"GOOG",
	"FB",
	"NFLX",
}

func GetWatchlistHandler(c *gin.Context) {
	watchlist, err := repositories.GetWatchlist()
	if err != nil {
		logrus.WithError(err).Error("error while getting watch list from repository")
		watchlist = defaultWatchlist
	}

	c.JSON(http.StatusOK, watchlist)
}

func UpdateWatchlistHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Cannot read body")
		return
	}

	arr := strings.Split(string(body), ",")
	var arr2 []string
	for _, item := range arr {
		item = strings.ToUpper(strings.TrimSpace(item))
		if len(item) > 0 {
			arr2 = append(arr2, item)
		}
	}

	arr2 = funk.UniqString(arr2)

	err = repositories.PersistWatchlist(arr2)
	if err != nil {
		logrus.WithError(err).Error("error while persisting watchlist")
	}

	c.JSON(http.StatusOK, arr2)
}
