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
	Years                      []int          `json:"years"`
	Sales                      ListMoney      `json:"sales"`
	SalesGrowth                ListPercentage `json:"sales_growth"`
	SalesGrowth5Years          Percentage     `json:"sales_growth_5_years"`
	SalesGrowth3Years          Percentage     `json:"sales_growth_3_years"`
	SalesGrowthLastYear        Percentage     `json:"sales_growth_last_year"`
	GrossIncome                ListMoney      `json:"gross_income"`
	GrossIncomeMargin          ListPercentage `json:"gross_income_margin"`
	PretaxIncome               ListMoney      `json:"pretax_income"`
	NetIncome                  ListMoney      `json:"net_income"`
	EPS                        ListMoney      `json:"eps"`
	EPSGrowth                  ListPercentage `json:"eps_growth"`
	EPSGrowth5Years            Percentage     `json:"eps_growth_5_years"`
	EPSGrowth3Years            Percentage     `json:"eps_growth_3_years"`
	EPSGrowthLastYear          Percentage     `json:"eps_growth_last_year"`
	TotalAssets                ListMoney      `json:"total_assets"`
	TotalAssetsGrowth          ListPercentage `json:"total_assets_growth"`
	ShortTermDebt              ListMoney      `json:"short_term_debt"`
	LongTermDebt               ListMoney      `json:"long_term_debt"`
	TotalLiabilities           ListMoney      `json:"total_liabilities"`
	Equity                     ListMoney      `json:"equity"`
	EquityGrowth5Years         Percentage     `json:"equity_growth_5_years"`
	EquityGrowth3Years         Percentage     `json:"equity_growth_3_years"`
	EquityGrowthLastYear       Percentage     `json:"equity_growth_last_year"`
	FreeCashFlow               ListMoney      `json:"free_cash_flow"`
	FreeCashFlowGrowth         ListPercentage `json:"free_cash_flow_growth"`
	FreeCashFlowGrowth5Years   Percentage     `json:"free_cash_flow_growth_5_years"`
	FreeCashFlowGrowth3Years   Percentage     `json:"free_cash_flow_growth_3_years"`
	FreeCashFlowGrowthLastYear Percentage     `json:"free_cash_flow_growth_last_year"`
	Url                        string         `json:"url"`
}

type MorningStarFinancialData struct {
	Revenues        []YearAmount `json:"revenues"`
	RevenueGrowths  []YearAmount `json:"revenue_growths"`
	EPS             []YearAmount `json:"eps"`
	EPSGrowths      []YearAmount `json:"eps_growths"`
	Equities        []YearAmount `json:"equities"`
	EquityGrowths   []YearAmount `json:"equity_growths"`
	CashFlows       []YearAmount `json:"cash_flows"`
	CashFlowGrowths []YearAmount `json:"cash_flow_growths"`
}

type MorningStarPerformanceDTO struct {
	ROI10Years      Percentage               `json:"roi_10_years"`
	ROI5Years       Percentage               `json:"roi_5_years"`
	ROILastYears    Percentage               `json:"roi_last_year"`
	ROITTM          Percentage               `json:"roittm"`
	LatestFairPrice float64                  `json:"latest_fair_price"`
	Url             string                   `json:"url"`
	FinancialData   MorningStarFinancialData `json:"financial_data"`
}
