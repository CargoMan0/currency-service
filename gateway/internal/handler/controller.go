package handler

import "github.com/BernsteinMondy/currency-service/gateway/internal/service"

type Controller struct {
	authService     service.AuthService
	currencyService service.CurrencyService
}
