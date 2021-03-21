package entities

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Year struct {
	Year     uint `json:"year"`
	IsTTM    bool `json:"is_ttm"`
	IsPeriod bool `json:"is_period"`
}

func NewYear(yearStr string) (Year, error) {
	if strings.ToLower(yearStr) == "ttm" {
		return Year{IsTTM: true, Year: uint(time.Now().Year())}, nil
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return Year{}, err
	}

	return Year{Year: uint(year)}, nil
}

func NewYearPeriod(year int) Year {
	return Year{Year: uint(year), IsPeriod: true}
}

func (y *Year) PeriodFrom(orig Year) Year {
	return Year{Year: y.Year - orig.Year, IsPeriod: true}
}

type Amount interface {
	FloatSetter
	FloatGetter
	IsNaN() bool
	UnmarshalJSON(b []byte) error
}

type YearAmount struct {
	Year   Year   `json:"year"`
	Amount Amount `json:"amount"`
}

type yearAmountAlias YearAmount

func (ymo *YearAmount) UnmarshalJSON(bytes []byte) error {
	var ym yearAmountAlias

	money := Money(0)
	percentage := Percentage(0)
	ym.Amount = &money

	err := json.Unmarshal(bytes, &ym)
	if err == nil {
		goto end
	}

	ym.Amount = &percentage
	err = json.Unmarshal(bytes, &ym)
	if err != nil {
		return err
	}

end:
	ymo.Amount = ym.Amount
	ymo.Year = ym.Year

	return nil
}

func NewYearAmount(year Year, amount Amount) YearAmount {
	return YearAmount{Year: year, Amount: amount}
}
