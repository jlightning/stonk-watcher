package util

import "math"

func CalculateAnnualCompoundInterest(orig, current float64, duration int) float64 {
	isNegative := 1.0
	if current-orig < 0 {
		isNegative = -1.0
	}

	percentage := math.Abs((current - orig) / orig)
	return isNegative * (math.Pow(1+percentage, 1/float64(duration-1)) - 1)
}
