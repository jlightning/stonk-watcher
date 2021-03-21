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
	Years              []int          `json:"years"`
	Sales              ListMoney      `json:"sales"`
	SalesGrowth        ListPercentage `json:"sales_growth"`
	GrossIncome        ListMoney      `json:"gross_income"`
	GrossIncomeMargin  ListYearAmount `json:"gross_income_margin"`
	PretaxIncome       ListMoney      `json:"pretax_income"`
	NetIncome          ListMoney      `json:"net_income"`
	EPS                ListMoney      `json:"eps"`
	EPSGrowth          ListPercentage `json:"eps_growth"`
	TotalAssets        ListMoney      `json:"total_assets"`
	TotalAssetsGrowth  ListPercentage `json:"total_assets_growth"`
	ShortTermDebt      ListMoney      `json:"short_term_debt"`
	LongTermDebt       ListMoney      `json:"long_term_debt"`
	TotalLiabilities   ListMoney      `json:"total_liabilities"`
	Equity             ListMoney      `json:"equity"`
	FreeCashFlow       ListMoney      `json:"free_cash_flow"`
	FreeCashFlowGrowth ListPercentage `json:"free_cash_flow_growth"`
	Url                string         `json:"url"`
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
