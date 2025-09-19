package service

import "time"

type User struct {
	Login    string
	Password string
}

type CurrencyRequest struct {
	Currency string
	DateFrom time.Time
	DateTo   time.Time
}

type CurrencyResponse struct {
	Currency string
	Rates    []CurrencyRate
}

type CurrencyRate struct {
	Rate float32
	Date time.Time
}
