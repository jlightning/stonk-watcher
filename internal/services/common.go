package services

import "stonk-watcher/internal/entities"

type StockInfoDTO struct {
	FinvizStockInfoDTO        *FinvizStockInfoDTO
	MarketWatchInfoDTO        *MarketWatchInfoDTO
	MorningStarPerformanceDTO *MorningStarPerformanceDTO
}

type FinvizStockInfoDTO struct {
	CompanyName   string
	Index         string
	MarketCap     entities.Money
	Income        entities.Money
	PE            entities.Percentage
	DividendYield entities.Percentage
	DebtOnEquity  entities.Percentage
	GrossMargin   entities.Percentage
	TargetPrice   entities.Money
	Price         entities.Money
	ShortRatio    entities.Percentage
}

type MarketWatchInfoDTO struct {
	Years                    []int
	Sales                    entities.ListMoney
	SalesGrowth              entities.ListPercentage
	SalesGrowth5Years        entities.Percentage
	SalesGrowth3Years        entities.Percentage
	GrossIncome              entities.ListMoney
	GrossIncomeMargin        entities.ListPercentage
	PretaxIncome             entities.ListMoney
	NetIncome                entities.ListMoney
	EPS                      entities.ListMoney
	EPSGrowth                entities.ListPercentage
	EPSGrowth5Years          entities.Percentage
	EPSGrowth3Years          entities.Percentage
	TotalAssets              entities.ListMoney
	TotalAssetsGrowth        entities.ListPercentage
	ShortTermDebt            entities.ListMoney
	LongTermDebt             entities.ListMoney
	TotalLiabilities         entities.ListMoney
	Equity                   entities.ListMoney
	EquiryGrowth5Years       entities.Percentage
	EquiryGrowth3Years       entities.Percentage
	FreeCashFlow             entities.ListMoney
	FreeCashFlowGrowth       entities.ListPercentage
	FreeCashFlowGrowth5Years entities.Percentage
	FreeCashFlowGrowth3Years entities.Percentage
}

type MorningStarPerformanceResponseDTO struct {
	Reported struct {
		Columns   []string `json:"columnDefs"`
		Collapsed struct {
			Rows []struct {
				Label      string   `json:"label"`
				Datum      []string `json:"datum"`
				Percentage bool     `json:"percentage"`
			} `json:"rows"`
		} `json:"Collapsed"`
	} `json:"reported"`
}

type MorningStarPerformanceDTO struct {
	ROI10Years   entities.Percentage
	ROI5Years    entities.Percentage
	ROILastYears entities.Percentage
	ROITTM       entities.Percentage
}
