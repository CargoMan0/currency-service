package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/currency-service/currency/internal/dto"
	"github.com/BernsteinMondy/currency-service/currency/internal/service"
	"time"
)

type Repository struct {
	db *sql.DB
}

var _ service.Repository = (*Repository)(nil)

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type CurrencyRate struct {
	Date time.Time
	Rate float32
}

func (r *Repository) SaveCurrencyExchangeRates(ctx context.Context, date time.Time, baseCurrency string, rates map[string]float64) error {
	ratesJSON, err := json.Marshal(rates)
	if err != nil {
		return fmt.Errorf("failed to marshal currency rates: %w", err)
	}

	const query = `INSERT INTO exchange_rates (date, base_currency, currency_rates) VALUES ($1, $2, $3)`

	_, err = r.db.ExecContext(ctx, query, date, baseCurrency, ratesJSON)
	if err != nil {
		return fmt.Errorf("failed to save exchange rates: %w", err)
	}

	return nil
}

func (r *Repository) GetCurrencyExchangeRatesInInterval(ctx context.Context, dto *dto.CurrencyRequestDTO) (_ []CurrencyRate, err error) {
	const query = `
		SELECT date, (currency_rates ->> $1)::float 
		FROM exchange_rates
		WHERE date::date BETWEEN $2 AND $3 AND base_currency = $4
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		dto.TargetCurrency,
		dto.DateFrom.Format("2006-01-02"),
		dto.DateTo.Format("2006-01-02"),
		dto.BaseCurrency,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query exchange rates: %w", err)
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close rows: %w", closeErr))
		}
	}()

	var rates []CurrencyRate
	for rows.Next() {
		var rate CurrencyRate
		if err := rows.Scan(&rate.Date, &rate.Rate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		rates = append(rates, rate)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return rates, nil
}
