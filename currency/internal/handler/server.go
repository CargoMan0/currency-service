package handler

import (
	"context"
	"github.com/CargoMan0/currency-service/currency/internal/dto"
	"github.com/CargoMan0/currency-service/currency/internal/middleware"
	"github.com/CargoMan0/currency-service/currency/internal/repository"
	"github.com/CargoMan0/currency-service/pkg/currency"
	"go.uber.org/zap"
)

type CurrencyService interface {
	GetCurrencyRatesInInterval(ctx context.Context, reqDTO *dto.CurrencyRequestDTO) ([]repository.CurrencyRate, error)
	FetchAndSaveCurrencyRates(ctx context.Context, baseCurrency string) error
}

type CurrencyServer struct {
	currency.UnimplementedCurrencyServiceServer
	mw      *middleware.Middleware
	service CurrencyService
	logger  *zap.Logger
}

func NewCurrencyServer(svc CurrencyService, logger *zap.Logger) CurrencyServer {
	return CurrencyServer{
		service: svc,
		logger:  logger,
	}
}
