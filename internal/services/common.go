package services

type FinvizStockInfoDTO struct {
	CompanyName   string
	Index         string
	MarketCap     float64
	Income        float64
	PE            float64
	DividendYield float64
	DebtOnEquity  float64
	GrossMargin   float64
	TargetPrice   float64
	Price         float64
	ShortRatio    float64
}

type MarketWatchInfoDTO struct {
	Years                    []int
	Sales                    []*float64
	SalesGrowth              []*float64
	SalesGrowth5Years        *float64
	SalesGrowth3Years        *float64
	GrossIncome              []*float64
	GrossIncomeMargin        []*float64
	PretaxIncome             []*float64
	NetIncome                []*float64
	EPS                      []*float64
	EPSGrowth                []*float64
	EPSGrowth5Years          *float64
	EPSGrowth3Years          *float64
	TotalAssets              []*float64
	TotalAssetsGrowth        []*float64
	ShortTermDebt            []*float64
	LongTermDebt             []*float64
	TotalLiabilities         []*float64
	Equity                   []*float64
	EquiryGrowth5Years       *float64
	EquiryGrowth3Years       *float64
	FreeCashFlow             []*float64
	FreeCashFlowGrowth       []*float64
	FreeCashFlowGrowth5Years *float64
	FreeCashFlowGrowth3Years *float64
}
