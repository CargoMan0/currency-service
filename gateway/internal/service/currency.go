package service

import (
	"context"
	"fmt"
)

type CurrencyService struct {
	currencyClient CurrencyClient
}

type CurrencyClient interface {
	GetCurrencyRates(ctx context.Context, request CurrencyRequest) (*CurrencyResponse, error)
}

func NewCurrencyService(currencyClient CurrencyClient) *CurrencyService {
	return &CurrencyService{
		currencyClient: currencyClient,
	}
}

func (svc *CurrencyService) GetCurrencyRates(ctx context.Context, request CurrencyRequest) (*CurrencyResponse, error) {
	resp, err := svc.currencyClient.GetCurrencyRates(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("currency client: get currency rates: %w", err)
	}

	return resp, nil
}
