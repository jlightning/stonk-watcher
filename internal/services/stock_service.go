package services

import (
	"stonk-watcher/internal/entities"
	"stonk-watcher/internal/repositories"
	"stonk-watcher/internal/util"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var stockInformation = util.NewProcessingQueue(5)

func GetStockInformation(ticker string) (*entities.StockInfoDTO, error) {
	result, err := stockInformation.Trigger(func() (interface{}, error) {
		return getStockInformationWithRetry(ticker)
	})
	if err != nil {
		return nil, err
	}

	return result.(*entities.StockInfoDTO), nil
}

func getStockInformationWithRetry(ticker string) (stockInfo *entities.StockInfoDTO, err error) {
	for i := 0; i < 5; i++ {
		if i > 0 {
			time.Sleep(1 * time.Second)
		}
		stockInfo, err = getStockInformation(ticker)
		if err == nil {
			return
		}
	}
	return
}

func getStockInformation(ticker string) (*entities.StockInfoDTO, error) {
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
