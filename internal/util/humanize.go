package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseMultipleFloat(parser func(string) (float64, error)) func([]string) ([]*float64, error) {
	return func(input []string) (res []*float64, err error) {
		for _, i := range input {
			if i == "-" {
				res = append(res, nil)
			} else {
				amount, err := parser(i)
				if err != nil {
					return nil, err
				}

				res = append(res, &amount)
			}
		}

		return
	}
}

func ParseMoney(str string) (float64, error) {
	orgStr := str

	unit := ""
	if str[len(str)-1] >= 'A' && str[len(str)-1] <= 'z' {
		unit = string(str[len(str)-1])
		str = str[:len(str)-1]
	}

	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	switch unit {
	case "B":
		res *= 1000000000
	case "M":
		res *= 1000000
	case "":
	default:
		return 0, errors.New(fmt.Sprintf("invalid amount: %s", orgStr))
	}

	return res, nil
}

func ParsePercentage(str string) (float64, error) {
	if !strings.HasSuffix(str, "%") {
		return 0, errors.New(fmt.Sprintf("invalid percentage: %s", str))
	}

	str = str[:len(str)-1]

	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	res /= 100

	return res, nil
}
