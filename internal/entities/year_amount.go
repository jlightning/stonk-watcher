package entities

import "encoding/json"

type Amount interface {
	FloatSetter
	FloatGetter
	IsNaN() bool
	UnmarshalJSON(b []byte) error
	Add(b Amount) Amount
	Multiply(b Amount) Amount
	FlipSign() Amount
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

type ListYearAmount []YearAmount

func (x ListYearAmount) Len() int {
	return len(x)
}

func (x ListYearAmount) Less(i, j int) bool {
	return x[i].Year.Year < x[j].Year.Year
}

func (x ListYearAmount) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x ListYearAmount) Add(b ListYearAmount) ListYearAmount {
	var res ListYearAmount
	for _, xitem := range x {
		for _, bitem := range b {
			if xitem.Year.Equal(bitem.Year) && !bitem.Amount.IsNaN() {
				xitem.Amount = xitem.Amount.Add(bitem.Amount)
			}
		}

		res = append(res, xitem)
	}

	return res
}

func (x ListYearAmount) Multiply(b ListYearAmount) ListYearAmount {
	var res ListYearAmount
	for _, xitem := range x {
		for _, bitem := range b {
			if xitem.Year.Equal(bitem.Year) && !bitem.Amount.IsNaN() {
				xitem.Amount = xitem.Amount.Multiply(bitem.Amount)
			}
		}

		res = append(res, xitem)
	}

	return res
}

func (x ListYearAmount) AddToAll(addFactor float64) ListYearAmount {
	var res ListYearAmount
	for _, xitem := range x {
		xitem.Amount = xitem.Amount.Add(NewMoney(addFactor))

		res = append(res, xitem)
	}

	return res
}

func (x ListYearAmount) FlipSign() ListYearAmount {
	var res ListYearAmount
	for _, xitem := range x {
		xitem.Amount = xitem.Amount.FlipSign()
		res = append(res, xitem)
	}

	return res
}

func (x ListYearAmount) Last() YearAmount {
	if len(x) == 0 {
		return YearAmount{}
	}
	return x[len(x)-1]
}
