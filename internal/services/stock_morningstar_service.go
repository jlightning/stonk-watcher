package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"sort"
	"stonk-watcher/internal/entities"
	"stonk-watcher/internal/util"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/pretty"

	"github.com/sirupsen/logrus"

	"github.com/gocolly/colly"
)

type morningStarPerformanceResponseDTO struct {
	Reported struct {
		Columns   []string `json:"columnDefs"`
		Collapsed struct {
			Rows []struct {
				Label      string   `json:"label"`
				Datum      []string `json:"datum"`
				Percentage bool     `json:"percentage"`
			} `json:"rows"`
		} `json:"Collapsed"`
		Expanded struct {
			Rows []struct {
				Label      string   `json:"label"`
				Datum      []string `json:"datum"`
				Percentage bool     `json:"percentage"`
			} `json:"rows"`
		} `json:"Expanded"`
	} `json:"reported"`
}

type morningStarValuationResponseDTO struct {
	Collapsed struct {
		Rows []struct {
			Label      string   `json:"label"`
			Datum      []string `json:"datum"`
			Percentage bool     `json:"percentage"`
		} `json:"rows"`
		Columns []string `json:"columnDefs"`
	} `json:"collapsed"`
}

type morningStarFairPriceDTO struct {
	Chart struct {
		ChartDatums struct {
			Recent struct {
				LatestFairValue string `json:"latestFairValue"`
			} `json:"recent"`
		} `json:"chartDatums"`
	} `json:"chart"`
}

type morningStarDividendDataDTO struct {
	Columns []string `json:"columnDefs_labels"`
	Rows    []struct {
		Label    string                          `json:"label"`
		Datum    []string                        `json:"datum"`
		SubLevel morningStarFinancialDataRowsDTO `json:"subLevel"`
	} `json:"rows"`
}

type morningStarFinancialDataRowDTO struct {
	Label    string                          `json:"label"`
	Datum    []*float64                      `json:"datum"`
	SubLevel morningStarFinancialDataRowsDTO `json:"subLevel"`
}

type morningStarFinancialDataRowsDTO []*morningStarFinancialDataRowDTO

type morningStarFinancialDataDTO struct {
	Columns []string                        `json:"columnDefs"`
	Rows    morningStarFinancialDataRowsDTO `json:"rows"`
	Footer  struct {
		OrderOfMagnitude string `json:"orderOfMagnitude"`
	} `json:"footer"`
}

func (d morningStarFinancialDataRowsDTO) find(label []string) (morningStarFinancialDataRowDTO, bool) {
	if len(label) == 0 {
		return morningStarFinancialDataRowDTO{}, false
	}
	for idx := range d {
		row := d[idx]
		if row.Label == label[0] {
			if len(label) == 1 {
				return *row, true
			} else {
				return row.SubLevel.find(label[1:])
			}
		}
	}

	return morningStarFinancialDataRowDTO{}, false
}

func (d morningStarFinancialDataDTO) getMoney(amount float64) entities.Money {
	if strings.ToLower(d.Footer.OrderOfMagnitude) == "billion" {
		return entities.Money(amount * 1000000000)
	}
	if strings.ToLower(d.Footer.OrderOfMagnitude) == "million" {
		return entities.Money(amount * 1000000)
	}
	if strings.ToLower(d.Footer.OrderOfMagnitude) == "thousand" {
		return entities.Money(amount * 1000)
	}
	return entities.Money(amount)
}

func GetDataFromMorningstar(ticker string) (*entities.MorningStarPerformanceDTO, error) {
	stockMSID, url, err := getMorningstarStockID(ticker, nil)
	if err != nil {
		return nil, err
	}

	headerStr, err := ioutil.ReadFile("morningstarKey.tmp.json")
	if err != nil {
		return nil, err
	}
	headerData := make(map[string]string)
	if err := json.Unmarshal(headerStr, &headerData); err != nil {
		return nil, err
	}

	performanceDTO, err := getMorningstarPerformance(stockMSID, headerData)
	if err != nil {
		return nil, err
	}

	fairPriceDTO, err := getMorningstarFairPrice(stockMSID, headerData)
	if err != nil {
		return nil, err
	}

	financialData, err := getMorningStarFinancialData(stockMSID, headerData)
	if err != nil {
		return nil, err
	}

	valuationDTO, err := getMorningStarValuationData(stockMSID, headerData)
	if err != nil {
		return nil, err
	}

	dividendData, err := getMorningStarDividendData(stockMSID, headerData)
	if err != nil {
		return nil, err
	}

	var rois []entities.YearAmount
	for _, row := range performanceDTO.Reported.Collapsed.Rows {
		for idx, col := range performanceDTO.Reported.Columns {
			if len(row.Datum) <= idx {
				continue
			}
			dataStr := row.Datum[idx]
			amount, err := strconv.ParseFloat(dataStr, 64)
			if err != nil {
				continue
			}

			if regexp.MustCompile("^[0-9]+$").MatchString(col) || strings.ToLower(col) == "ttm" {
				year, err := entities.NewYear(col)
				if err != nil {
					return nil, err
				}
				if strings.Contains(strings.ToLower(row.Label), "invested capital") {
					roiPercentage := entities.Percentage(amount / 100)
					rois = append(rois, entities.YearAmount{
						Year:   year,
						Amount: &roiPercentage,
					})
				}
			}
		}
	}

	for _, row := range performanceDTO.Reported.Expanded.Rows {
		for idx, col := range performanceDTO.Reported.Columns {
			if len(row.Datum) <= idx {
				continue
			}
			dataStr := row.Datum[idx]
			amount, err := strconv.ParseFloat(dataStr, 64)
			if err != nil {
				continue
			}

			if regexp.MustCompile("^[0-9]+$").MatchString(col) || strings.ToLower(col) == "ttm" {
				year, err := entities.NewYear(col)
				if err != nil {
					return nil, err
				}

				if strings.Contains(strings.ToLower(row.Label), "gross margin %") {
					financialData.GrossProfitMargins = append(financialData.GrossProfitMargins, entities.YearAmount{
						Year:   year,
						Amount: entities.NewPercentage(amount / 100),
					})
				}

				if strings.Contains(strings.ToLower(row.Label), "net margin %") {
					financialData.NetProfitMargins = append(financialData.NetProfitMargins, entities.YearAmount{
						Year:   year,
						Amount: entities.NewPercentage(amount / 100),
					})
				}
			}
		}
	}

	latestFairPrice, _ := util.ParseMoney(fairPriceDTO.Chart.ChartDatums.Recent.LatestFairValue)

	response := entities.MorningStarPerformanceDTO{
		ROIs:            rois,
		ROIGrowths:      calculateAverage(rois),
		LatestFairPrice: latestFairPrice,
		Url:             url,
		FinancialData:   *financialData,
		ValuationData:   *valuationDTO,
		DividendData:    dividendData,
	}

	return &response, nil
}

func getMorningstarPerformance(stockMSID string, headerData map[string]string) (*morningStarPerformanceResponseDTO, error) {
	c := colly.NewCollector()
	apiURL := fmt.Sprintf("https://api-global.morningstar.com/sal-service/v1/stock/operatingPerformance/v2/%s", stockMSID)

	var responseDTO morningStarPerformanceResponseDTO

	c.OnRequest(func(request *colly.Request) {
		for k, v := range headerData {
			request.Headers.Add(k, v)
		}
	})
	c.OnResponse(func(response *colly.Response) {
		err := json.Unmarshal(response.Body, &responseDTO)
		if err != nil {
			logrus.Warnf("Error while decoding Morningstar response: %s", err.Error())
		}

		//ioutil.WriteFile("performance.tmp.json", pretty.PrettyOptions(response.Body, &pretty.Options{
		//	Width:  180,
		//	Prefix: "",
		//	Indent: "  ",
		//}), 0600)
	})

	err := c.Visit(apiURL)
	if err != nil {
		return nil, err
	}

	return &responseDTO, nil
}

func getMorningStarValuationData(stockMSID string, headerData map[string]string) (*entities.MorningStarValuationData, error) {
	c := colly.NewCollector()
	apiURL := fmt.Sprintf("https://api-global.morningstar.com/sal-service/v1/stock/valuation/v3/%s?clientId=MDC&benchmarkId=category&version=3.41.0", stockMSID)

	var responseDTO morningStarValuationResponseDTO

	c.OnRequest(func(request *colly.Request) {
		for k, v := range headerData {
			request.Headers.Add(k, v)
		}
	})
	c.OnResponse(func(response *colly.Response) {
		err := json.Unmarshal(response.Body, &responseDTO)
		if err != nil {
			logrus.Warnf("Error while decoding Morningstar response: %s", err.Error())
		}

		//ioutil.WriteFile("valuation.tmp.json", pretty.PrettyOptions(response.Body, &pretty.Options{
		//	Width:  180,
		//	Prefix: "",
		//	Indent: "  ",
		//}), 0600)
	})

	err := c.Visit(apiURL)
	if err != nil {
		return nil, err
	}

	var valuationDTO entities.MorningStarValuationData

	for _, row := range responseDTO.Collapsed.Rows {
		for idx, col := range responseDTO.Collapsed.Columns {
			if idx == 0 {
				continue
			}
			if len(row.Datum) <= idx-1 {
				continue
			}
			dataStr := row.Datum[idx-1]
			amount, err := strconv.ParseFloat(dataStr, 64)
			if err != nil {
				continue
			}

			if regexp.MustCompile("^[0-9]+$").MatchString(col) || strings.ToLower(col) == "ttm" || strings.ToLower(col) == "current" {
				year, err := entities.NewYear(col)
				if err != nil {
					return nil, err
				}

				if strings.Contains(strings.ToLower(row.Label), "price/earnings") {
					valuationDTO.PriceOnEarnings = append(valuationDTO.PriceOnEarnings, entities.YearAmount{
						Year:   year,
						Amount: entities.NewPercentage(amount),
					})
				}

				if strings.Contains(strings.ToLower(row.Label), "price/book") {
					valuationDTO.PriceOnBooks = append(valuationDTO.PriceOnBooks, entities.YearAmount{
						Year:   year,
						Amount: entities.NewPercentage(amount),
					})
				}
			}
		}
	}

	return &valuationDTO, nil
}

func getMorningstarFairPrice(stockMSID string, headerData map[string]string) (*morningStarFairPriceDTO, error) {
	c := colly.NewCollector()
	apiURL := fmt.Sprintf("https://api-global.morningstar.com/sal-service/v1/stock/priceFairValue/v2/%s/data?secExchangeList=&languageId=en&locale=en&clientId=MDC&benchmarkId=category&component=sal-components-price-fairvalue&version=3.41.0", stockMSID)

	var responseDTO morningStarFairPriceDTO

	c.OnRequest(func(request *colly.Request) {
		for k, v := range headerData {
			request.Headers.Add(k, v)
		}
	})
	c.OnResponse(func(response *colly.Response) {
		err := json.Unmarshal(response.Body, &responseDTO)
		if err != nil {
			logrus.Warnf("Error while decoding Morningstar response: %s", err.Error())
		}
	})

	err := c.Visit(apiURL)
	if err != nil {
		return nil, err
	}

	return &responseDTO, nil
}

func getMorningStarFinancialData(stockMSID string, headerData map[string]string) (*entities.MorningStarFinancialData, error) {
	getStmData := func(stmType string, responseDTO *morningStarFinancialDataDTO) error {
		c := colly.NewCollector()
		apiURL := fmt.Sprintf("https://api-global.morningstar.com/sal-service/v1/stock/newfinancials/%s/%s/detail?dataType=A&reportType=A&locale=en&clientId=MDC&benchmarkId=category&version=3.41.0", stockMSID, stmType)

		c.OnRequest(func(request *colly.Request) {
			for k, v := range headerData {
				request.Headers.Add(k, v)
			}
		})
		c.OnResponse(func(response *colly.Response) {
			err := json.Unmarshal(response.Body, responseDTO)
			if err != nil {
				logrus.Warnf("Error while decoding Morningstar response: %s", err.Error())
			}

			//ioutil.WriteFile(stmType+".tmp.json", pretty.PrettyOptions(response.Body, &pretty.Options{
			//	Width:  180,
			//	Prefix: "",
			//	Indent: "  ",
			//}), 0600)
		})

		err := c.Visit(apiURL)
		if err != nil {
			return err
		}

		return nil
	}

	var res entities.MorningStarFinancialData
	var incomeStmResp, balanceSheetStmResp, cashFlowStmResp morningStarFinancialDataDTO

	if err := getStmData("incomeStatement", &incomeStmResp); err != nil {
		return nil, err
	}
	if err := getStmData("balanceSheet", &balanceSheetStmResp); err != nil {
		return nil, err
	}
	if err := getStmData("cashFlow", &cashFlowStmResp); err != nil {
		return nil, err
	}

	if err := populateAmount(incomeStmResp, []string{"IncomeStatement", "Gross Profit", "Total Revenue"}, &res.Revenues, &res.RevenueGrowths, true, true); err != nil {
		return nil, err
	}

	if err := populateAmount(incomeStmResp, []string{"IncomeStatement", "Gross Profit"}, &res.GrossProfits, &res.GrossProfitGrowths, true, true); err != nil {
		return nil, err
	}

	if err := populateAmount(incomeStmResp, []string{"IncomeStatement", "Total Operating Profit/Loss"}, &res.NetProfits, &res.NetProfitGrowths, true, true); err != nil {
		return nil, err
	}

	if err := populateAmount(incomeStmResp, []string{"WasoAndEpsData", "Diluted EPS"}, &res.EPS, &res.EPSGrowths, true, false); err != nil {
		return nil, err
	}

	if err := populateAmount(balanceSheetStmResp, []string{"BalanceSheet", "Total Equity"}, &res.Equities, &res.EquityGrowths, true, true); err != nil {
		return nil, err
	}

	if err := populateAmount(cashFlowStmResp, []string{"CashFlow", "Cash and Cash Equivalents, End of Period"}, &res.CashFlows, &res.CashFlowGrowths, true, true); err != nil {
		return nil, err
	}

	return &res, nil
}

func getMorningStarDividendData(stockMSID string, headerData map[string]string) (*entities.MorningStarDividendData, error) {
	c := colly.NewCollector()
	apiURL := fmt.Sprintf("https://api-global.morningstar.com/sal-service/v1/stock/dividends/v4/%s/data?clientId=MDC&benchmarkId=category&version=3.41.0", stockMSID)

	var responseDTO morningStarDividendDataDTO

	c.OnRequest(func(request *colly.Request) {
		for k, v := range headerData {
			request.Headers.Add(k, v)
		}
	})
	c.OnResponse(func(response *colly.Response) {
		err := json.Unmarshal(response.Body, &responseDTO)
		if err != nil {
			logrus.Warnf("Error while decoding Morningstar response: %s", err.Error())
		}

		ioutil.WriteFile("dividends.tmp.json", pretty.PrettyOptions(response.Body, &pretty.Options{
			Width:  180,
			Prefix: "",
			Indent: "  ",
		}), 0600)
	})

	err := c.Visit(apiURL)
	if err != nil {
		return nil, err
	}

	var newRows morningStarFinancialDataRowsDTO
	for _, row := range responseDTO.Rows {
		var parsedDatum []*float64
		for _, d := range row.Datum {
			if d == "" {
				parsedDatum = append(parsedDatum, nil)
			} else {
				parsedAmount, err := strconv.ParseFloat(d, 64)
				if err != nil {
					return nil, err
				}

				parsedDatum = append(parsedDatum, &parsedAmount)
			}
		}

		newRows = append(newRows, &morningStarFinancialDataRowDTO{
			Label: row.Label,
			Datum: parsedDatum,
		})
	}

	financialDTO := morningStarFinancialDataDTO{
		Columns: responseDTO.Columns[1:],
		Rows:    newRows,
	}

	var res entities.MorningStarDividendData

	if err := populateAmount(financialDTO, []string{"Dividend Per Share"}, &res.DividendPerShares, &res.DividendPerShareGrowths, true, false); err != nil {
		return nil, err
	}
	if err := populateAmount(financialDTO, []string{"Trailing Dividend Yield %"}, &res.TotalYields, nil, true, false); err != nil {
		return nil, err
	}
	if err := populateAmount(financialDTO, []string{"Total Yield %"}, &res.TotalYields, nil, true, false); err != nil {
		return nil, err
	}

	return &res, nil
}

func calculateGrowth(input []entities.YearAmount, periods []int) []entities.YearAmount {
	var currentAmount entities.YearAmount
	for i := len(input) - 1; i >= 0; i-- {
		if !input[i].Amount.IsNaN() {
			currentAmount = input[i]
			break
		}
	}

	var res []entities.YearAmount
	for _, period := range periods {
		for _, yearAmount := range input {
			periodFrom := currentAmount.Year.PeriodFrom(yearAmount.Year)
			if !yearAmount.Amount.IsNaN() && int(periodFrom.Year) <= period {
				res = append(res, entities.YearAmount{
					Year:   periodFrom,
					Amount: entities.NewPercentage(util.CalculateAnnualCompoundInterest(yearAmount.Amount.Get(), currentAmount.Amount.Get(), int(periodFrom.Year))),
				})

				break
			}
		}
	}

	return res
}

func populateAmount(stm morningStarFinancialDataDTO, find []string, amountList *entities.ListYearAmount, growthList *entities.ListYearAmount, asMoney bool, parseByOrderOfMagnitude bool) error {
	amounts, ok := stm.Rows.find(find)
	if ok {
		for idx, amount := range amounts.Datum {
			yearStr := stm.Columns[idx]
			if yearStr == "dividends.headers.current" || yearStr == "dividends.headers.fiveyear" {
				continue
			}
			if yearStr == "dividends.headers.oneyearttm" {
				yearStr = "ttm"
			}
			year, err := entities.NewYear(yearStr)
			if err != nil {
				return err
			}

			if amount != nil {
				var value entities.Amount

				if asMoney {
					value = entities.NewMoney(*amount)
					if parseByOrderOfMagnitude {
						_value := stm.getMoney(*amount)
						value = &_value
					}
				} else {
					value = entities.NewPercentage(*amount)
				}
				*amountList = append(*amountList, entities.NewYearAmount(year, value))
			} else {
				if asMoney {
					*amountList = append(*amountList, entities.NewYearAmount(year, entities.NewMoney(math.NaN())))
				} else {
					*amountList = append(*amountList, entities.NewYearAmount(year, entities.NewPercentage(math.NaN())))
				}
			}
		}

		if growthList != nil {
			*growthList = calculateGrowth(*amountList, []int{10, 5, 2})
		}
	}

	return nil
}

func calculateAverage(input []entities.YearAmount) []entities.YearAmount {
	periods := []int{10, 5, 1}
	currentYear := entities.Year{Year: uint(time.Now().Year())}

	type amountCount struct {
		amount float64
		count  int
	}
	YearAmountMap := make(map[int]amountCount, len(periods))

	var res []entities.YearAmount
	for _, yearAmount := range input {
		if yearAmount.Year.IsTTM {
			continue
		}
		periodFrom := currentYear.PeriodFrom(yearAmount.Year)
		for _, p := range periods {
			pYearAmount := YearAmountMap[p]
			if int(periodFrom.Year) <= p {
				pYearAmount.amount += yearAmount.Amount.Get()
				pYearAmount.count++

				YearAmountMap[p] = pYearAmount
			}
		}
	}

	for _, p := range periods {
		if pYearAmount, ok := YearAmountMap[p]; ok {
			res = append(res, entities.YearAmount{
				Year:   entities.NewYearPeriod(p),
				Amount: entities.NewPercentage(pYearAmount.amount / float64(pYearAmount.count)),
			})
		}
	}

	return res
}

func calculateMargin(marginOfs entities.ListYearAmount, amounts entities.ListYearAmount) (res entities.ListYearAmount) {
	for _, marginOf := range marginOfs {
		if marginOf.Amount.IsNaN() {
			continue
		}

		for _, amount := range amounts {
			if amount.Amount.IsNaN() || amount.Year.Year != marginOf.Year.Year || amount.Year.IsTTM != marginOf.Year.IsTTM {
				continue
			}

			ym := entities.YearAmount{
				Year:   marginOf.Year,
				Amount: entities.NewPercentage(amount.Amount.Get() / marginOf.Amount.Get()),
			}
			if math.IsInf(ym.Amount.Get(), 0) {
				continue
			}
			res = append(res, ym)
		}
	}

	sort.Sort(res)

	return
}

func getMorningstarStockID(ticker string, prefix *string) (string, string, error) {
	c := colly.NewCollector()

	if prefix == nil {
		_prefix := "xnas"
		prefix = &_prefix
	}

	url := fmt.Sprintf("https://www.morningstar.com/stocks/%s/%s/quote", *prefix, strings.ToLower(ticker))
	stockMSID := ""

	c.OnHTML("body", func(body *colly.HTMLElement) {
		text := body.Text
		found := regexp.MustCompile(`byId:{"([0-9a-zA-Z]+)"`).FindStringSubmatch(text)
		if len(found) > 1 {
			stockMSID = found[1]
		}
	})

	err := c.Visit(url)
	if err != nil {
		if err.Error() == "Not Found" && *prefix == "xnas" {
			_prefix := "xnys"
			return getMorningstarStockID(ticker, &_prefix)
		}
		return "", "", err
	}

	return stockMSID, url, nil
}
