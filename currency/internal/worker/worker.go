package worker

import (
	"context"
	"fmt"
	"github.com/BernsteinMondy/currency-service/currency/internal/config"
	"github.com/BernsteinMondy/currency-service/currency/internal/service"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

type Currency struct {
	currencyService service.Currency
	cron            *cron.Cron
	schedule        string
	baseCurrency    string
	targetCurrency  string
	logger          *zap.Logger
}

func NewCurrency(
	cfg config.WorkerConfig,
	service service.Currency,
	cron *cron.Cron,
	logger *zap.Logger,
) *Currency {
	return &Currency{
		currencyService: service,
		cron:            cron,
		schedule:        cfg.Schedule,
		baseCurrency:    cfg.CurrencyPair.BaseCurrency,
		targetCurrency:  cfg.CurrencyPair.TargetCurrency,
		logger:          logger,
	}
}

func (w *Currency) StartFetchingCurrencyRates() error {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // todo move to config
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
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // todo move to config
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
		return fmt.Errorf("Cron.AddFunc: %w", err)
	}

	w.cron.Start()

	return nil
}

func (w *Currency) Stop() {
	w.cron.Stop()
}
