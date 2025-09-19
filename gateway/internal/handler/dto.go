package handler

import (
	"github.com/BernsteinMondy/currency-service/gateway/internal/service"
	"time"
)

type registerRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type loginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type currencyRequest struct {
	Currency string `form:"currency" binding:"required"`
	DateFrom string `form:"date_from" binding:"required,datetime=2006-01-02"`
	DateTo   string `form:"date_to" binding:"required,datetime=2006-01-02"`
}

func currencyResponseFromServiceToDTO(currencyResp *service.CurrencyResponse) *currencyResponse {
	rates := make([]currencyRate, 0, len(currencyResp.Rates))

	for _, rate := range currencyResp.Rates {
		rates = append(rates, currencyRate{
			Rate: rate.Rate,
			Date: rate.Date.Format(time.RFC3339),
		})
	}

	return &currencyResponse{
		Currency: currencyResp.Currency,
		Rates:    rates,
	}
}

type currencyResponse struct {
	Currency string
	Rates    []currencyRate
}

type currencyRate struct {
	Rate float32
	Date string
}
