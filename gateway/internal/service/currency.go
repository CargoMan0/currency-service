package service

import (
	"context"
	"fmt"
	"github.com/BernsteinMondy/currency-service/gateway/internal/dto"
)

type CurrencyService struct {
	currencyClient CurrencyClient
}

type CurrencyClient interface {
	GetCurrencyRates(ctx context.Context, request dto.CurrencyRequest) (*dto.CurrencyResponse, error)
}

func NewCurrencyService(currencyClient CurrencyClient) *CurrencyService {
	return &CurrencyService{
		currencyClient: currencyClient,
	}
}

func (svc *CurrencyService) GetCurrencyRates(ctx context.Context, request dto.CurrencyRequest) (*dto.CurrencyResponse, error) {
	resp, err := svc.currencyClient.GetCurrencyRates(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("currency client: get currency rates: %w", err)
	}

	return resp, nil
}
