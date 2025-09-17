package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "admin"
	dbname   = "currency_db"
)

func main() {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}
	defer db.Close()

	if err := generateTestData(db); err != nil {
		log.Fatalf("Failed to generate test data: %v", err)
	}

	fmt.Println("Test data generated successfully.")
}

func generateTestData(db *sql.DB) error {
	ctx := context.Background()
	baseCurrency := "RUB"
	targetCurrencies := []string{"usd", "eur", "gbp", "jpy", "cny"}

	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	insertStmt := `INSERT INTO exchange_rates (date, base_currency, currency_rates) VALUES ($1, $2, $3)`
	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		ratesMap := make(map[string]float64)
		for _, targetCurrency := range targetCurrencies {
			ratesMap[targetCurrency] = generateRandomRate()
		}
		ratesJSON, err := json.Marshal(ratesMap)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		_, err = tx.ExecContext(ctx, insertStmt, date, baseCurrency, ratesJSON)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert test data: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func generateRandomRate() float64 {
	return 0.01 + (0.1-0.01)*rand.Float64()
}
