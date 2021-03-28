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

func NewMoney(i float64) *Money {
	m := Money(i)
	return &m
}

func (m *Money) Set(i float64) {
	*m = Money(i)
}

func (m Money) Get() float64 {
	return float64(m)
}

func (m Money) IsNaN() bool {
	return math.IsNaN(float64(m))
}

func (m *Money) Add(b Amount) Amount {
	return NewMoney(float64(*m) + b.Get())
}

func (m *Money) Multiply(b Amount) Amount {
	return NewMoney(float64(*m) * b.Get())
}

func (m *Money) FlipSign() Amount {
	return NewPercentage(-float64(*m))
}

func (m Money) MarshalJSON() ([]byte, error) {
	if m.IsNaN() {
		return []byte("null"), nil
	}
	if math.IsInf(m.Get(), 0) {
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

type Percentage float64

func NewPercentage(i float64) *Percentage {
	p := Percentage(i)
	return &p
}

func (p *Percentage) Set(i float64) {
	*p = Percentage(i)
}

func (p Percentage) Get() float64 {
	return float64(p)
}

func (p Percentage) IsNaN() bool {
	return math.IsNaN(float64(p))
}

func (p *Percentage) Add(b Amount) Amount {
	return NewPercentage(float64(*p) + b.Get())
}

func (p *Percentage) Multiply(b Amount) Amount {
	return NewPercentage(float64(*p) * b.Get())
}

func (p *Percentage) FlipSign() Amount {
	return NewPercentage(-float64(*p))
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
	if math.IsInf(m.Get(), 0) {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`{"amount": %.4f, "percent": "%.2f%%"}`, m, m*100)), nil
}
