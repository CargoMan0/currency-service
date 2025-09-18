package dto

import "time"

type CurrencyResponse struct {
	Currency string
	Rates    []CurrencyRate
}

type CurrencyRate struct {
	Rate float32
	Date time.Time
}
