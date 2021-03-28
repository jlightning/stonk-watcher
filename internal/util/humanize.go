package util

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func ParseMultipleFloat(parser func(string) (float64, error)) func([]string) ([]float64, error) {
	return func(input []string) (res []float64, err error) {
		for _, i := range input {
			if i == "-" {
				res = append(res, math.NaN())
			} else {
				amount, err := parser(i)
				if err != nil {
					return nil, err
				}

				res = append(res, amount)
			}
		}

		return
	}
}

func ParseMoney(str string) (float64, error) {
	if len(str) == 0 {
		return math.NaN(), nil
	}
	if str == "-" {
		return math.NaN(), nil
	}
	if strings.HasPrefix(str, "$") {
		str = str[1:]
	}

	orgStr := str

	if regexp.MustCompile("^\\(.*\\)$").MatchString(str) {
		str = "-" + str[1:len(str)-1]
	}

	unit := ""
	if str[len(str)-1] >= 'A' && str[len(str)-1] <= 'z' {
		unit = string(str[len(str)-1])
		str = str[:len(str)-1]
	}

	str = strings.ReplaceAll(str, ",", "")

	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	switch unit {
	case "T":
		res *= 1000000000000
	case "B":
		res *= 1000000000
	case "M":
		res *= 1000000
	case "K":
		res *= 1000
	case "":
	default:
		return 0, errors.New(fmt.Sprintf("invalid amount: %s", orgStr))
	}

	return res, nil
}

func ParsePercentage(str string) (float64, error) {
	if !strings.HasSuffix(str, "%") {
		return math.NaN(), nil
	}
	if str == "-" {
		return math.NaN(), nil
	}

	str = str[:len(str)-1]

	str = strings.ReplaceAll(str, ",", "")

	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	res /= 100

	return res, nil
}
