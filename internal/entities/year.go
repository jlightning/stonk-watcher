package entities

import (
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

	if strings.ToLower(yearStr) == "current" {
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

func (y Year) Equal(b Year) bool {
	return y.Year == b.Year && y.IsTTM == b.IsTTM && y.IsPeriod == b.IsPeriod
}
