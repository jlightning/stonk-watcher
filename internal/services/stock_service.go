package services

import (
	"stonk-watcher/internal/entities"
	"stonk-watcher/internal/repositories"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetStockInformation(ticker string) (*entities.StockInfoDTO, error) {
	stockInfo, err := repositories.GetStockInfo(ticker)
	if err != nil {
		logrus.WithError(err).Warnf("error while getting data from repository for ticker: %s", ticker)
	}
	if stockInfo != nil {
		return stockInfo, nil
	}

	finvizInfo, err := GetDataFromFinviz(ticker)
	if err != nil {
		return nil, err
	}

	marketWatchInfo, err := GetDataFromMarketWatch(ticker)
	if err != nil {
		return nil, err
	}

	morningStarInfo, err := GetDataFromMorningstar(ticker)
	if err != nil {
		return nil, err
	}

	dto := &entities.StockInfoDTO{
		Ticker:                    strings.ToUpper(ticker),
		FinvizStockInfoDTO:        finvizInfo,
		MarketWatchInfoDTO:        marketWatchInfo,
		MorningStarPerformanceDTO: morningStarInfo,
	}

	err = repositories.PersistStockInfo(ticker, dto)
	if err != nil {
		logrus.WithError(err).Warnf("error while persisting data for ticker: %s", ticker)
	}
	return dto, nil
}
