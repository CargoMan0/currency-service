package dto

import "time"

type ParsedCurrencyRequest struct {
	Currency string
	DateFrom time.Time
	DateTo   time.Time
}
