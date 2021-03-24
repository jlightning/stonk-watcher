package services

import (
	"stonk-watcher/internal/repositories"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type CronLogger struct {
}

func (c *CronLogger) Info(msg string, keysAndValues ...interface{}) {
	logrus.Infof(msg, keysAndValues)
}

func (c *CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	logrus.WithError(err).Errorf(msg, keysAndValues)
}

func init() {
	cronGroup := cron.New(cron.WithChain(
		cron.Recover(&CronLogger{}), // or use cron.DefaultLogger
		cron.DelayIfStillRunning(&CronLogger{}),
	))

	_, err := cronGroup.AddFunc("@every 1h", func() {
		logrus.Info("CRON: reload stock price")

		watchlist, err := repositories.GetWatchlist()
		if err != nil {
			panic(err)
		}

		for _, ticker := range watchlist {
			finvizData, err := GetDataFromFinviz(ticker)
			if err != nil {
				panic(err)
			}

			savedStockData, err := repositories.GetStockInfo(ticker)
			if err == nil && savedStockData != nil {
				savedStockData.FinvizStockInfoDTO = finvizData

				err = repositories.PersistStockInfo(ticker, savedStockData)
				if err != nil {
					panic(err)
				}

				logrus.Infof("CRON: reloaded stock price for %s", ticker)
			}
		}
	})
	if err != nil {
		panic(err)
	}

	cronGroup.Start()
}
