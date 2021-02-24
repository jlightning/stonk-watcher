package services

import (
	"fmt"
	"stonk-watcher/internal/util"
	"strings"

	"github.com/gocolly/colly"
)

func GetDataFromFinviz(ticker string) (*FinvizStockInfoDTO, error) {
	url := fmt.Sprintf("https://finviz.com/quote.ashx?t=%s", ticker)

	c := colly.NewCollector()

	var stockInfo FinvizStockInfoDTO

	data := map[string]string{}
	// getting stock table
	c.OnHTML("body", func(body *colly.HTMLElement) {
		body.ForEach("table", func(i int, e *colly.HTMLElement) {
			matchDetailTable := 0
			e.ForEach("td", func(i int, element *colly.HTMLElement) {
				if i == 0 && element.Text == "Index" {
					matchDetailTable++
				}

				if i == 2 && element.Text == "P/E" {
					matchDetailTable++
				}
			})

			if matchDetailTable == 2 {
				var allText []string
				e.ForEach("td", func(i int, element *colly.HTMLElement) {
					allText = append(allText, element.Text)
				})

				for i := 0; i < len(allText); i += 2 {
					data[allText[i]] = allText[i+1]
				}
			}
		})

		body.ForEach("table", func(i int, e *colly.HTMLElement) {
			found := false
			e.ForEach("td", func(i int, element *colly.HTMLElement) {
				if i == 0 && strings.HasPrefix(element.Text, strings.ToUpper(ticker)) {
					found = true
				}

				if found && i == 1 {
					stockInfo.CompanyName = element.Text
				}
			})
		})
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	stockInfo.Index = data["Index"]

	floatPairs := []struct {
		dest   *float64
		title  string
		parser func(string) (float64, error)
	}{
		{dest: &stockInfo.MarketCap, title: "Market Cap", parser: util.ParseMoney},
		{dest: &stockInfo.Income, title: "Income", parser: util.ParseMoney},
		{dest: &stockInfo.PE, title: "P/E", parser: util.ParseMoney},
		{dest: &stockInfo.DividendYield, title: "Dividend %", parser: util.ParsePercentage},
		{dest: &stockInfo.DebtOnEquity, title: "Debt/Eq", parser: util.ParseMoney},
		{dest: &stockInfo.GrossMargin, title: "Gross Margin", parser: util.ParsePercentage},
		{dest: &stockInfo.TargetPrice, title: "Target Price", parser: util.ParseMoney},
		{dest: &stockInfo.Price, title: "Price", parser: util.ParseMoney},
		{dest: &stockInfo.ShortRatio, title: "Short Ratio", parser: util.ParseMoney},
	}

	for _, pair := range floatPairs {
		value, err := pair.parser(data[pair.title])
		if err != nil {
			return nil, err
		}

		*pair.dest = value
	}

	return &stockInfo, nil
}
