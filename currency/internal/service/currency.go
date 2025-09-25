package service

import (
	"context"
	"fmt"
	"github.com/BernsteinMondy/currency-service/currency/internal/dto"
	"github.com/BernsteinMondy/currency-service/currency/internal/repository"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Repository interface {
	SaveCurrencyExchangeRates(ctx context.Context, date time.Time, baseCurrency string, rates map[string]float64) error
	GetCurrencyExchangeRatesInInterval(ctx context.Context, dto *dto.CurrencyRequestDTO) ([]repository.CurrencyRate, error)
}

type CurrencyAPIClient interface {
	FetchCurrentRates(ctx context.Context, currency string) (dto.RatesResponse, error)
}

type Currency struct {
	currencyRepo Repository
	client       CurrencyAPIClient
	logger       *zap.Logger
}

func NewCurrency(
	repo Repository,
	client CurrencyAPIClient,
	logger *zap.Logger,
) *Currency {
	return &Currency{
		currencyRepo: repo,
		client:       client,
		logger:       logger,
	}
}

func (s *Currency) GetCurrencyRatesInInterval(ctx context.Context, reqDTO *dto.CurrencyRequestDTO) ([]repository.CurrencyRate, error) {
	reqDTO.TargetCurrency = strings.ToLower(reqDTO.TargetCurrency)
	rates, err := s.currencyRepo.GetCurrencyExchangeRatesInInterval(ctx, reqDTO)
	if err != nil {
		return nil, fmt.Errorf("currency repo: get currency exchange rates in interval: %w", err)
	}

	return rates, nil
}

func (s *Currency) FetchAndSaveCurrencyRates(ctx context.Context, baseCurrency string) error {
	rates, err := s.client.FetchCurrentRates(ctx, baseCurrency)
	if err != nil {
		return fmt.Errorf("client: fetch current rates: %s", err)
	}

	date, err := time.Parse("2006-01-02", rates.Date)
	if err != nil {
		return fmt.Errorf("failed to parse currency date: %v ", err)
	}

	err = s.currencyRepo.SaveCurrencyExchangeRates(ctx, date, baseCurrency, rates.Rub)
	if err != nil {
		return fmt.Errorf("Failed to save currency rates: %v ", err)
	}

	s.logger.Info("Currency rates fetched and saved", zap.Any("rates", rates))
	return nil
}
