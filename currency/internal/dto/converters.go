package dto

import (
	"github.com/BernsteinMondy/currency-service/pkg/currency"
	"time"
)

func CurrencyRequestFromPbToDTO(req *currency.GetRateRequest, baseCurrency string) *CurrencyRequestDTO {
	return &CurrencyRequestDTO{
		BaseCurrency:   baseCurrency,
		TargetCurrency: req.Currency,
		DateFrom:       req.DateFrom.AsTime(),
		DateTo:         req.DateTo.AsTime(),
	}
}
