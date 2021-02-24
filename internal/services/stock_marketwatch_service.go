package services

import (
	"fmt"
	"math"
	"regexp"
	"stonk-watcher/internal/entities"
	"stonk-watcher/internal/util"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func GetDataFromMarketWatch(ticker string) (*entities.MarketWatchInfoDTO, error) {
	return getFinancialDataFromMarketWatch(ticker)
}

func getFinancialDataFromMarketWatch(ticker string) (*entities.MarketWatchInfoDTO, error) {
	financialUrl := fmt.Sprintf("https://www.marketwatch.com/investing/stock/%s/financials", strings.ToLower(ticker))
	cashFlowUrl := fmt.Sprintf("https://www.marketwatch.com/investing/stock/%s/financials/cash-flow", strings.ToLower(ticker))
	balanceSheetUrl := fmt.Sprintf("https://www.marketwatch.com/investing/stock/%s/financials/balance-sheet", strings.ToLower(ticker))

	var stockInfo entities.MarketWatchInfoDTO

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
		dest   entities.ListFloatSetter
		title  string
		source map[string][]string
		parser func([]string) ([]float64, error)
	}{
		{dest: &stockInfo.Sales, title: "Sales/Revenue", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.SalesGrowth, title: "Sales Growth", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParsePercentage)},
		{dest: &stockInfo.GrossIncome, title: "Gross Income", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.PretaxIncome, title: "Pretax Income", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.NetIncome, title: "Net Income", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.EPS, title: "EPS (Diluted)", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.EPSGrowth, title: "EPS (Diluted) Growth", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParsePercentage)},
		{dest: &stockInfo.TotalAssets, title: "Total Assets", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.TotalAssetsGrowth, title: "Total Assets Growth", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParsePercentage)},
		{dest: &stockInfo.ShortTermDebt, title: "Short Term Debt", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.LongTermDebt, title: "Long-Term Debt", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.TotalLiabilities, title: "Total Liabilities", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.Equity, title: "Total Shareholders' Equity", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.FreeCashFlow, title: "Free Cash Flow", source: cashFlowData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.FreeCashFlowGrowth, title: "Free Cash Flow Growth", source: cashFlowData, parser: util.ParseMultipleFloat(util.ParsePercentage)},
	}

	for _, pair := range floatPairs {
		value, err := pair.parser(pair.source[pair.title])
		if err != nil {
			return nil, err
		}

		pair.dest.Set(value)
	}

	for i, sale := range stockInfo.Sales {
		if len(stockInfo.GrossIncome) > i {
			grossIncome := stockInfo.GrossIncome[i]

			if !sale.IsNaN() && !grossIncome.IsNaN() {
				amount := grossIncome / sale
				stockInfo.GrossIncomeMargin = append(stockInfo.GrossIncomeMargin, entities.Percentage(amount))
			} else {
				stockInfo.GrossIncomeMargin = append(stockInfo.GrossIncomeMargin, entities.Percentage(math.NaN()))
			}
		}
	}

	compoundInterestPairs2 := []struct {
		dest     entities.FloatSetter
		source   entities.ListFloatGetter
		duration int
	}{
		{dest: &stockInfo.SalesGrowth5Years, source: stockInfo.Sales, duration: 5},
		{dest: &stockInfo.SalesGrowth3Years, source: stockInfo.Sales, duration: 3},
		{dest: &stockInfo.SalesGrowthLastYear, source: stockInfo.Sales, duration: 2},
		{dest: &stockInfo.EPSGrowth5Years, source: stockInfo.EPS, duration: 5},
		{dest: &stockInfo.EPSGrowth3Years, source: stockInfo.EPS, duration: 3},
		{dest: &stockInfo.EPSGrowthLastYear, source: stockInfo.EPS, duration: 2},
		{dest: &stockInfo.EquityGrowth5Years, source: stockInfo.Equity, duration: 5},
		{dest: &stockInfo.EquityGrowth3Years, source: stockInfo.Equity, duration: 3},
		{dest: &stockInfo.EquityGrowthLastYear, source: stockInfo.Equity, duration: 2},
		{dest: &stockInfo.FreeCashFlowGrowth5Years, source: stockInfo.FreeCashFlow, duration: 5},
		{dest: &stockInfo.FreeCashFlowGrowth3Years, source: stockInfo.FreeCashFlow, duration: 3},
		{dest: &stockInfo.FreeCashFlowGrowthLastYear, source: stockInfo.FreeCashFlow, duration: 2},
	}

	for _, pair := range compoundInterestPairs2 {
		source := pair.source.GetListFloat()

		if len(source)-pair.duration < 0 {
			continue
		}

		if !math.IsNaN(source[len(source)-pair.duration]) && !math.IsNaN(source[len(source)-1]) {
			pair.dest.Set(util.CalculateAnnualCompoundInterest(source[len(source)-pair.duration], source[len(source)-1], pair.duration))
		}
	}

	return &stockInfo, nil
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
