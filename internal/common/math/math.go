package math

import "github.com/shopspring/decimal"

func Percent(percent, total float64) float64 {
	d := decimal.NewFromFloat(percent).Mul(decimal.NewFromFloat(total))
	d = d.Div(decimal.NewFromFloat(100))
	f, _ := d.Float64()

	return f
}

func PercentOf(part, total float64) float64 {
	d := decimal.NewFromFloat(part).Mul(decimal.NewFromFloat(100))
	d = d.Div(decimal.NewFromFloat(total))
	f, _ := d.Float64()

	return f
}

func Round(total float64, places int32) float64 {
	d := decimal.NewFromFloat(total).Round(places)
	f, _ := d.Float64()

	return f
}
