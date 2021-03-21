package services

import (
	"fmt"
	"math"
	"regexp"
	"sort"
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

	stockInfo := entities.MarketWatchInfoDTO{
		Url: financialUrl,
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
		{dest: &stockInfo.EPS, growthDest: &stockInfo.EPSGrowths, title: "EPS (Diluted)", source: incomeStmData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.TotalAssets, title: "Total Assets", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.ShortTermDebt, title: "Short Term Debt", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
		{dest: &stockInfo.LongTermDebt, title: "Long-Term Debt", source: balanceSheetData, parser: util.ParseMultipleFloat(util.ParseMoney)},
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
			*pair.growthDest = calculateGrowth(*pair.dest, []int{5, 3, 2})
		}
	}

	for i, sale := range stockInfo.Sales {
		if len(stockInfo.GrossIncome) > i {
			grossIncome := stockInfo.GrossIncome[i]

			if !sale.Amount.IsNaN() && !grossIncome.Amount.IsNaN() {
				amount := grossIncome.Amount.Get() / sale.Amount.Get()
				percentage := entities.Percentage(amount)
				stockInfo.GrossIncomeMargin = append(stockInfo.GrossIncomeMargin, entities.YearAmount{Year: entities.Year{Year: uint(years[i])}, Amount: &percentage})
			} else {
				percentage := entities.Percentage(math.NaN())
				stockInfo.GrossIncomeMargin = append(stockInfo.GrossIncomeMargin, entities.YearAmount{Year: entities.Year{Year: uint(years[i])}, Amount: &percentage})
			}
		}
	}

	sort.Sort(stockInfo.GrossIncomeMargin)

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
