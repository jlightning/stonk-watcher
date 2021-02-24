package services

func GetStockInformation(ticker string) (*StockInfoDTO, error) {
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

	return &StockInfoDTO{
		FinvizStockInfoDTO:        finvizInfo,
		MarketWatchInfoDTO:        marketWatchInfo,
		MorningStarPerformanceDTO: morningStarInfo,
	}, nil
}
