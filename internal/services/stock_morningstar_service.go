package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"stonk-watcher/internal/entities"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func GetDataFromMorningstar(ticker string) (*MorningStarPerformanceDTO, error) {
	stockMSID, err := getMorningstarStockID(ticker, nil)
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile("morningstarKey.tmp.txt")
	if err != nil {
		return nil, err
	}
	keyStr := strings.Trim(string(key), "\n")

	performanceDTO, err := getMorningstarPerformance(stockMSID, keyStr)
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
					return nil, err
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

	response := MorningStarPerformanceDTO{
		ROI10Years:   entities.Percentage(roi10 / 100),
		ROI5Years:    entities.Percentage(roi5 / 100),
		ROILastYears: entities.Percentage(roi1 / 100),
		ROITTM:       entities.Percentage(roittm / 100),
	}

	return &response, nil
}

func getMorningstarPerformance(stockMSID string, apiKey string) (*MorningStarPerformanceResponseDTO, error) {
	c := colly.NewCollector()
	apiUrl := fmt.Sprintf("https://api-global.morningstar.com/sal-service/v1/stock/operatingPerformance/v2/%s", stockMSID)

	var responseDTO MorningStarPerformanceResponseDTO

	c.OnRequest(func(request *colly.Request) {
		request.Headers.Add("ApiKey", apiKey)
	})
	c.OnResponse(func(response *colly.Response) {
		err := json.Unmarshal(response.Body, &responseDTO)
		if err != nil {
			fmt.Println("Error while decoding Morningstar response: ", err.Error())
		}
	})

	err := c.Visit(apiUrl)
	if err != nil {
		return nil, err
	}

	return &responseDTO, nil
}

func getMorningstarStockID(ticker string, prefix *string) (string, error) {
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
		return "", err
	}

	return stockMSID, nil
}
