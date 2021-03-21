package entities

type StockInfoDTO struct {
	Ticker                    string                     `json:"ticker"`
	FinvizStockInfoDTO        *FinvizStockInfoDTO        `json:"finviz_info"`
	MarketWatchInfoDTO        *MarketWatchInfoDTO        `json:"marketwatch_info"`
	MorningStarPerformanceDTO *MorningStarPerformanceDTO `json:"morningstar_info"`
}

type FinvizStockInfoDTO struct {
	CompanyName      string     `json:"company_name"`
	Index            string     `json:"index"`
	MarketCap        Money      `json:"market_cap"`
	Income           Money      `json:"income"`
	PE               Percentage `json:"pe"`
	PB               Percentage `json:"pb"`
	DividendYield    Percentage `json:"dividend_yield"`
	DebtOnEquity     Percentage `json:"debt_on_equity"`
	GrossMargin      Percentage `json:"gross_margin"`
	TargetPrice      Money      `json:"target_price"`
	Price            Money      `json:"price"`
	ShortFloat       Percentage `json:"short_float"`
	RSI              Percentage `json:"rsi"`
	EPSNextYear      Percentage `json:"eps_next_year"`
	EPSNext5Years    Percentage `json:"eps_next_5_years"`
	EPSTTM           Money      `json:"epsttm"`
	Url              string     `json:"url"`
	ShareOutstanding Money      `json:"share_outstanding"`
}

type MarketWatchInfoDTO struct {
	Years               []int          `json:"years"`
	Sales               ListYearAmount `json:"sales"`
	SalesGrowth         ListYearAmount `json:"sales_growth"`
	GrossIncome         ListYearAmount `json:"gross_income"`
	GrossIncomeMargin   ListYearAmount `json:"gross_income_margin"`
	PretaxIncome        ListYearAmount `json:"pretax_income"`
	NetIncome           ListYearAmount `json:"net_income"`
	EPS                 ListYearAmount `json:"eps"`
	EPSGrowths          ListYearAmount `json:"eps_growths"`
	TotalAssets         ListYearAmount `json:"total_assets"`
	ShortTermDebt       ListYearAmount `json:"short_term_debt"`
	LongTermDebt        ListYearAmount `json:"long_term_debt"`
	TotalLiabilities    ListYearAmount `json:"total_liabilities"`
	Equities            ListYearAmount `json:"equities"`
	EquityGrowths       ListYearAmount `json:"equity_growths"`
	FreeCashFlow        ListYearAmount `json:"free_cash_flow"`
	FreeCashFlowGrowths ListYearAmount `json:"free_cash_flow_growths"`
	Url                 string         `json:"url"`
}

type MorningStarFinancialData struct {
	Revenues        ListYearAmount `json:"revenues"`
	RevenueGrowths  ListYearAmount `json:"revenue_growths"`
	EPS             ListYearAmount `json:"eps"`
	EPSGrowths      ListYearAmount `json:"eps_growths"`
	Equities        ListYearAmount `json:"equities"`
	EquityGrowths   ListYearAmount `json:"equity_growths"`
	CashFlows       ListYearAmount `json:"cash_flows"`
	CashFlowGrowths ListYearAmount `json:"cash_flow_growths"`
}

type MorningStarPerformanceDTO struct {
	ROIs            ListYearAmount           `json:"rois"`
	ROIGrowths      ListYearAmount           `json:"roi_growths"`
	LatestFairPrice float64                  `json:"latest_fair_price"`
	Url             string                   `json:"url"`
	FinancialData   MorningStarFinancialData `json:"financial_data"`
}
