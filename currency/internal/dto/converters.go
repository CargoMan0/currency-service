package dto

import (
	"github.com/CargoMan0/currency-service/pkg/currency"
)

func CurrencyRequestFromPbToDTO(req *currency.GetRateRequest, baseCurrency string) *CurrencyRequestDTO {
	return &CurrencyRequestDTO{
		BaseCurrency:   baseCurrency,
		TargetCurrency: req.Currency,
		DateFrom:       req.DateFrom.AsTime(),
		DateTo:         req.DateTo.AsTime(),
	}
}
