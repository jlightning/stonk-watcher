package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"stonk-watcher/internal/entities"
	"stonk-watcher/internal/util"
	"strconv"
	"strings"

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
	} `json:"reported"`
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

type morningStarFinancialDataRowDTO struct {
	Label    string                          `json:"label"`
	Datum    []*float64                      `json:"datum"`
	SubLevel morningStarFinancialDataRowsDTO `json:"subLevel"`
}

type morningStarFinancialDataRowsDTO []*morningStarFinancialDataRowDTO

type morningStarFinancialDataDTO struct {
	Columns []string                        `json:"columnDefs"`
	Rows    morningStarFinancialDataRowsDTO `json:"rows"`
}

func (d morningStarFinancialDataRowsDTO) Find(label []string) (morningStarFinancialDataRowDTO, bool) {
	if len(label) == 0 {
		return morningStarFinancialDataRowDTO{}, false
	}
	for idx := range d {
		row := d[idx]
		if row.Label == label[0] {
			if len(label) == 1 {
				return *row, true
			} else {
				return row.SubLevel.Find(label[1:])
			}
		}
	}

	return morningStarFinancialDataRowDTO{}, false
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

	err = getMorningStarFinancialData(stockMSID, headerData)
	if err != nil {
		return nil, err
	}

	var roi10, roi5, roi1, roittm float64
	for _, row := range performanceDTO.Reported.Collapsed.Rows {
		if strings.Contains(strings.ToLower(row.Label), "invested capital") {
			for idx, col := range performanceDTO.Reported.Columns {
				roiStr := row.Datum[idx]
				roi, err := strconv.ParseFloat(roiStr, 64)
				if err != nil {
					continue
				}
				if regexp.MustCompile("^[0-9]+$").MatchString(col) {
					roi10 += roi

					if idx > 4 {
						roi5 += roi
					}
					if idx == 9 {
						roi1 = roi
					}
				}
				if strings.ToLower(col) == "ttm" {
					roittm = roi
				}
			}
		}
	}
	roi10 /= 10
	roi5 /= 10

	latestFairPrice, _ := util.ParseMoney(fairPriceDTO.Chart.ChartDatums.Recent.LatestFairValue)

	response := entities.MorningStarPerformanceDTO{
		ROI10Years:      entities.Percentage(roi10 / 100),
		ROI5Years:       entities.Percentage(roi5 / 100),
		ROILastYears:    entities.Percentage(roi1 / 100),
		ROITTM:          entities.Percentage(roittm / 100),
		LatestFairPrice: latestFairPrice,
		Url:             url,
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
	})

	err := c.Visit(apiURL)
	if err != nil {
		return nil, err
	}

	return &responseDTO, nil
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

func getMorningStarFinancialData(stockMSID string, headerData map[string]string) error {
	sub := func(stmType string, responseDTO *morningStarFinancialDataDTO) error {
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

			fmt.Println(string(pretty.Color(pretty.PrettyOptions(response.Body, &pretty.Options{
				Width:  180,
				Prefix: "",
				Indent: "  ",
			}), nil)))
		})

		err := c.Visit(apiURL)
		if err != nil {
			return err
		}

		//fmt.Println(util.MustJSONStringify(responseDTO, true))

		return nil
	}

	var incomeStmResp morningStarFinancialDataDTO

	if err := sub("incomeStatement", &incomeStmResp); err != nil {
		return err
	}

	dilutedEPS, _ := incomeStmResp.Rows.Find([]string{"WasoAndEpsData", "Diluted EPS"})
	fmt.Println(util.MustJSONStringify(dilutedEPS, true))
	fmt.Println(util.MustJSONStringify(incomeStmResp.Columns, true))

	//if err := sub("balanceSheet"); err != nil {
	//	return err
	//}
	//if err := sub("cashFlow"); err != nil {
	//	return err
	//}

	return nil
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
