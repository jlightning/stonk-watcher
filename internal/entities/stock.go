package entities

type StockInfoDTO struct {
	Ticker                    string
	FinvizStockInfoDTO        *FinvizStockInfoDTO
	MarketWatchInfoDTO        *MarketWatchInfoDTO
	MorningStarPerformanceDTO *MorningStarPerformanceDTO
}

type FinvizStockInfoDTO struct {
	CompanyName   string
	Index         string
	MarketCap     Money
	Income        Money
	PE            Percentage
	DividendYield Percentage
	DebtOnEquity  Percentage
	GrossMargin   Percentage
	TargetPrice   Money
	Price         Money
	ShortRatio    Percentage
}

type MarketWatchInfoDTO struct {
	Years                    []int
	Sales                    ListMoney
	SalesGrowth              ListPercentage
	SalesGrowth5Years        Percentage
	SalesGrowth3Years        Percentage
	GrossIncome              ListMoney
	GrossIncomeMargin        ListPercentage
	PretaxIncome             ListMoney
	NetIncome                ListMoney
	EPS                      ListMoney
	EPSGrowth                ListPercentage
	EPSGrowth5Years          Percentage
	EPSGrowth3Years          Percentage
	TotalAssets              ListMoney
	TotalAssetsGrowth        ListPercentage
	ShortTermDebt            ListMoney
	LongTermDebt             ListMoney
	TotalLiabilities         ListMoney
	Equity                   ListMoney
	EquiryGrowth5Years       Percentage
	EquiryGrowth3Years       Percentage
	FreeCashFlow             ListMoney
	FreeCashFlowGrowth       ListPercentage
	FreeCashFlowGrowth5Years Percentage
	FreeCashFlowGrowth3Years Percentage
}

type MorningStarPerformanceDTO struct {
	ROI10Years   Percentage
	ROI5Years    Percentage
	ROILastYears Percentage
	ROITTM       Percentage
}
