package dto

import "time"

const DefaultBaseCurrency = "RUB"

type CurrencyRequestDTO struct {
	BaseCurrency   string
	TargetCurrency string
	DateFrom       time.Time
	DateTo         time.Time
}
