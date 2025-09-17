package dto

import (
	"github.com/BernsteinMondy/currency-service/pkg/currency"
	"time"
)

const DefaultBaseCurrency = "RUB"

type CurrencyRequestDTO struct {
	BaseCurrency   string
	TargetCurrency string
	DateFrom       time.Time
	DateTo         time.Time
}

func CurrencyRequestFromPbToDTO(req *currency.GetRateRequest, baseCurrency string) *CurrencyRequestDTO {
	return &CurrencyRequestDTO{
		BaseCurrency:   baseCurrency,
		TargetCurrency: req.Currency,
		DateFrom:       req.DateFrom.AsTime(),
		DateTo:         req.DateTo.AsTime(),
	}
}
