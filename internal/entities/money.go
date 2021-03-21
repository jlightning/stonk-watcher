package entities

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

type FloatSetter interface {
	Set(float64)
}

type FloatGetter interface {
	Get() float64
}

type ListFloatSetter interface {
	Set([]float64)
}

type ListFloatGetter interface {
	GetListFloat() []float64
}

type Money float64
type ListMoney []Money

func (m *Money) Set(i float64) {
	*m = Money(i)
}

func (m Money) Get() float64 {
	return float64(m)
}

func (m Money) IsNaN() bool {
	return math.IsNaN(float64(m))
}

func (m Money) MarshalJSON() ([]byte, error) {
	if m.IsNaN() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%.2f", m)), nil
}

func (m *Money) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		*m = Money(math.NaN())
	} else {
		f, err := strconv.ParseFloat(string(bytes), 64)
		if err != nil {
			return err
		}

		*m = Money(f)
	}
	return nil
}

func (lm *ListMoney) Set(list []float64) {
	*lm = make(ListMoney, 0, len(list))
	for _, i := range list {
		*lm = append(*lm, Money(i))
	}
}

func (lm ListMoney) GetListFloat() (r []float64) {
	for _, i := range lm {
		r = append(r, float64(i))
	}
	return
}

type Percentage float64
type ListPercentage []Percentage

func (p *Percentage) Set(i float64) {
	*p = Percentage(i)
}

func (p Percentage) Get() float64 {
	return float64(p)
}

func (p Percentage) IsNaN() bool {
	return math.IsNaN(float64(p))
}

func (lp *ListPercentage) Set(list []float64) {
	*lp = make(ListPercentage, 0, len(list))
	for _, i := range list {
		*lp = append(*lp, Percentage(i))
	}
}

func (lp ListPercentage) GetListFloat() (r []float64) {
	for _, i := range lp {
		r = append(r, float64(i))
	}
	return
}

func (m *Percentage) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		*m = Percentage(math.NaN())
	} else {
		var parseDest struct {
			Amount float64 `json:"amount"`
		}
		err := json.Unmarshal(bytes, &parseDest)
		if err != nil {
			return err
		}
		*m = Percentage(parseDest.Amount)
	}
	return nil
}

func (m Percentage) MarshalJSON() ([]byte, error) {
	if m.IsNaN() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`{"amount": %.2f, "percent": "%.2f%%"}`, m, m*100)), nil
}
