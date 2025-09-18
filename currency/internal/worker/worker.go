package worker

import (
	"context"
	"fmt"
	"github.com/BernsteinMondy/currency-service/currency/internal/config"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

type CurrencyService interface {
	FetchAndSaveCurrencyRates(ctx context.Context, baseCurrency string) error
}

type Currency struct {
	currencyService CurrencyService
	cron            *cron.Cron
	schedule        string
	baseCurrency    string
	targetCurrency  string
	timeoutSeconds  int
	logger          *zap.Logger
}

func NewCurrency(
	cfg config.WorkerConfig,
	service CurrencyService,
	cron *cron.Cron,
	logger *zap.Logger,
) *Currency {
	return &Currency{
		currencyService: service,
		cron:            cron,
		schedule:        cfg.Schedule,
		baseCurrency:    cfg.CurrencyPair.BaseCurrency,
		targetCurrency:  cfg.CurrencyPair.TargetCurrency,
		timeoutSeconds:  cfg.TimoutSeconds,
		logger:          logger,
	}
}

func (w *Currency) StartFetchingCurrencyRates() error {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(w.timeoutSeconds)*time.Second)
		defer cancel()

		err := w.currencyService.FetchAndSaveCurrencyRates(ctx, w.baseCurrency)
		if err != nil {
			w.logger.Error(
				"Failed to fetch currency rate immediately on startup",
				zap.Time("timestamp", time.Now()),
				zap.Error(err),
			)
		}
	}()

	_, err := w.cron.AddFunc(
		w.schedule, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(w.timeoutSeconds)*time.Second)
			defer cancel()

			err := w.currencyService.FetchAndSaveCurrencyRates(ctx, w.baseCurrency)
			if err != nil {
				w.logger.Error(
					"Failed to fetch currency rate on scheduled run",
					zap.Time("timestamp", time.Now()),
					zap.Error(err),
					zap.String("schedule", w.schedule),
				)
			}
		},
	)
	if err != nil {
		return fmt.Errorf("add func to cron: %w", err)
	}

	w.cron.Start()
	return nil
}

func (w *Currency) Stop() {
	w.cron.Stop()
}
