package util

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateAnnualCompoundInterest(t *testing.T) {
	a := assert.New(t)

	a.Equal(10.0, math.Round(CalculateAnnualCompoundInterest(10, 12.1, 3)*100))
}
