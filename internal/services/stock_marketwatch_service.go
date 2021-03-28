package services

import (
	"fmt"
	"regexp"
	"sort"
	"stonk-watcher/internal/entities"
	"stonk-watcher/internal/util"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/gocolly/colly"
)

func GetDataFromMarketWatch(ticker string) (*entities.MarketWatchInfoDTO, error) {
	return getFinancialDataFromMarketWatch(ticker)
}

func getFinancialDataFromMarketWatch(ticker string) (*entities.MarketWatchInfoDTO, error) {
	financialUrl := fmt.Sprintf("https://www.marketwatch.com/investing/stock/%s/financials", strings.ToLower(ticker))
	cashFlowUrl := fmt.Sprintf("https://www.marketwatch.com/investing/stock/%s/financials/cash-flow", strings.ToLower(ticker))
	balanceSheetUrl := fmt.Sprintf("https://www.marketwatch.com/investing/stock/%s/financials/balance-sheet", strings.ToLower(ticker))

	stockInfo := entities.MarketWatchInfoDTO{
		Url: financialUrl,
	}

	beta, marketCap, err := getMarketWatchOverviewData(ticker)
	if err != nil {
		return nil, err
	}

	incomeStmData, years, err := getMarketwatchTableData(financialUrl)
	if err != nil {
		return nil, err
	}

	balanceSheetData, _, err := getMarketwatchTableData(balanceSheetUrl)
	if err != nil {
		return nil, err
	}

	cashFlowData, _, err := getMarketwatchTableData(cashFlowUrl)
	if err != nil {
		return nil, err
	}

	stockInfo.Years = years

	floatPairs := []struct {
		dest       *entities.ListYearAmount
		growthDest *entities.ListYearAmount
		title      string
		source     map[string][]string
		parser     func([]string) ([]float64, error)
	}{
		{dest: &stockInfo.Sales, growthDest: &stockInfo.SalesGrowth, title: "Sales/Revenue", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.GrossIncome, title: "Gross Income", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.PretaxIncome, title: "Pretax Income", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.NetIncome, title: "Net Income", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.IncomeTax, title: "Income Tax", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.InterestExpense, title: "Interest Expense", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.EPS, growthDest: &stockInfo.EPSGrowths, title: "EPS (Diluted)", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.TotalAssets, title: "Total Assets", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.ShortTermDebt, title: "Short Term Debt", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.CurrentPortionOfLongTermDebt, title: "Current Portion of Long Term Debt", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.LongTermDebt, title: "Long-Term Debt", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.TotalCurrentLiabilities, title: "Total Current Liabilities", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.TotalLiabilities, title: "Total Liabilities", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.Equities, growthDest: &stockInfo.EquityGrowths, title: "Total Shareholders' Equity", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.FreeCashFlow, growthDest: &stockInfo.FreeCashFlowGrowths, title: "Free Cash Flow", source: cashFlowData, parser: util.ParseMultipleFloat(util.ParseMoney)},
	}

	for _, pair := range floatPairs {
		value, err := pair.parser(pair.source[pair.title])
		if err != nil {
			return nil, err
		}

		for i, v := range value {
			*pair.dest = append(*pair.dest, entities.YearAmount{
				Year:   entities.Year{Year: uint(years[i])},
				Amount: entities.NewMoney(v),
			})
		}

		sort.Sort(*pair.dest)

		if pair.growthDest != nil {
			*pair.growthDest = calculateGrowth(*pair.dest, []int{5, 3, 1})
		}
	}

	stockInfo.GrossIncomeMargin = calculateMargin(stockInfo.Sales, stockInfo.GrossIncome)
	stockInfo.NetIncomeMargins = calculateMargin(stockInfo.Sales, stockInfo.NetIncome)

	func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Error("error while calculating WACC", r)
			}
		}()
		interestExpenseOnLiabilities := calculateMargin(stockInfo.LongTermDebt.Add(stockInfo.CurrentPortionOfLongTermDebt), stockInfo.InterestExpense)

		taxRate := calculateMargin(stockInfo.PretaxIncome, stockInfo.IncomeTax)

		costOfDebts := interestExpenseOnLiabilities.Multiply(taxRate.FlipSign().AddToAll(1))

		const riskFreeRate = 0.02
		const expectedReturnOfMarket = 0.1

		costOfEquity := riskFreeRate + beta*(expectedReturnOfMarket-riskFreeRate)

		debtWeight := stockInfo.TotalCurrentLiabilities.Last().Amount.Get() / marketCap

		equityWeight := 1 - debtWeight

		wacc := debtWeight*costOfDebts.Last().Amount.Get() + equityWeight*costOfEquity
		stockInfo.WACC = entities.NewPercentage(wacc)
	}()

	return &stockInfo, nil
}

func getMarketWatchOverviewData(ticker string) (beta, marketCap float64, err error) {
	overviewUrl := fmt.Sprintf("https://www.marketwatch.com/investing/stock/%s", strings.ToLower(ticker))

	c := colly.NewCollector()

	c.OnHTML("li", func(e *colly.HTMLElement) {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(e.Text)), "beta") {
			arr := strings.Split(e.Text, "\n")
			var newArr []string
			for _, item := range arr {
				item = strings.Trim(item, " \n")
				if len(item) > 0 {
					newArr = append(newArr, item)
				}
			}
			if len(newArr) == 2 {
				beta, _ = strconv.ParseFloat(newArr[1], 64)
			}
		}

		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(e.Text)), "market cap") {
			arr := strings.Split(e.Text, "\n")
			var newArr []string
			for _, item := range arr {
				item = strings.Trim(item, " \n")
				if len(item) > 0 {
					newArr = append(newArr, item)
				}
			}
			if len(newArr) == 2 {
				marketCap, _ = util.ParseMoney(newArr[1])
			}
		}
	})

	err = c.Visit(overviewUrl)
	if err != nil {
		return 0, 0, err
	}

	return beta, marketCap, nil
}

func getMarketwatchTableData(url string) (map[string][]string, []int, error) {
	c := colly.NewCollector()

	data := map[string][]string{}
	var years []int

	c.OnHTML("table", func(e *colly.HTMLElement) {
		found := false
		e.ForEach("th", func(i int, th *colly.HTMLElement) {
			if i == 0 && parseMarketWatchRowTitle(th) == "Item" {
				found = true
			}
		})

		if found {
			e.ForEach("th", func(i int, th *colly.HTMLElement) {
				if i > 0 {
					parsed := parseMarketWatchRowTitle(th)
					if regexp.MustCompile("^[0-9]+$").Match([]byte(parsed)) {
						year, _ := strconv.Atoi(parsed)
						years = append(years, year)
					}
				}
			})

			e.ForEach("tr", func(i int, tr *colly.HTMLElement) {
				rowKey := ""
				tr.ForEach("td", func(i int, td *colly.HTMLElement) {
					if i == 0 {
						rowKey = parseMarketWatchRowTitle(td)
					} else if len(strings.TrimSpace(td.Text)) > 0 {
						data[rowKey] = append(data[rowKey], td.Text)
					}
				})
			})
		}
	})

	err := c.Visit(url)
	if err != nil {
		return nil, nil, err
	}
	return data, years, nil
}

func parseMarketWatchRowTitle(e *colly.HTMLElement) string {
	text := ""
	e.ForEach("div", func(i int, element *colly.HTMLElement) {
		if i == 0 {
			text = element.Text
		}
	})
	return text
}
