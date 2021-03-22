package services

import (
	"fmt"
	"regexp"
	"stonk-watcher/internal/entities"
	"stonk-watcher/internal/util"
	"strings"

	"github.com/gocolly/colly"
)

func GetDataFromFinviz(ticker string) (*entities.FinvizStockInfoDTO, error) {
	url := fmt.Sprintf("https://finviz.com/quote.ashx?t=%s", ticker)

	c := colly.NewCollector()

	stockInfo := entities.FinvizStockInfoDTO{
		Url: url,
	}

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
				if i == 0 && regexp.MustCompile(fmt.Sprintf("%s \\[[A-Z]+\\]", strings.ToUpper(ticker))).MatchString(element.Text) {
					found = true
				}

				if found && i == 1 {
					stockInfo.CompanyName = element.Text
				}

				if found && i == 2 {
					stockInfo.Sector = strings.TrimSpace(strings.Split(element.Text, "|")[0])
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
		dest   entities.FloatSetter
		title  string
		parser func(string) (float64, error)
	}{
		{dest: &stockInfo.MarketCap, title: "Market Cap", parser: util.ParseMoney},
		{dest: &stockInfo.Income, title: "Income", parser: util.ParseMoney},
		{dest: &stockInfo.PE, title: "P/E", parser: util.ParseMoney},
		{dest: &stockInfo.PB, title: "P/B", parser: util.ParseMoney},
		{dest: &stockInfo.DividendYield, title: "Dividend %", parser: util.ParsePercentage},
		{dest: &stockInfo.DebtOnEquity, title: "Debt/Eq", parser: util.ParseMoney},
		{dest: &stockInfo.GrossMargin, title: "Gross Margin", parser: util.ParsePercentage},
		{dest: &stockInfo.TargetPrice, title: "Target Price", parser: util.ParseMoney},
		{dest: &stockInfo.Price, title: "Price", parser: util.ParseMoney},
		{dest: &stockInfo.ShortFloat, title: "Short Float", parser: util.ParsePercentage},
		{dest: &stockInfo.RSI, title: "RSI (14)", parser: util.ParseMoney},
		{dest: &stockInfo.EPSNextYear, title: "EPS next Y", parser: util.ParsePercentage},
		{dest: &stockInfo.EPSNext5Years, title: "EPS next 5Y", parser: util.ParsePercentage},
		{dest: &stockInfo.EPSTTM, title: "EPS (ttm)", parser: util.ParseMoney},
		{dest: &stockInfo.ShareOutstanding, title: "Shs Outstand", parser: util.ParseMoney},
	}

	for _, pair := range floatPairs {
		value, err := pair.parser(data[pair.title])
		if err != nil {
			return nil, err
		}

		pair.dest.Set(value)
	}

	return &stockInfo, nil
}
