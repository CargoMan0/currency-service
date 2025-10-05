package currency

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CargoMan0/currency-service/currency/internal/config"
	"github.com/CargoMan0/currency-service/currency/internal/dto"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Currency struct {
	baseURL    *url.URL
	httpClient *http.Client
	logger     *zap.Logger
}

func NewClient(cfg config.CurrencyAPIConfig, logger *zap.Logger) (*Currency, error) {
	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return &Currency{}, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Currency{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		logger: logger,
	}, nil
}

func (c *Currency) FetchCurrentRates(ctx context.Context, currency string) (_ dto.RatesResponse, err error) {
	relativeCurrencyPath, _ := url.Parse(fmt.Sprintf("/v1/currencies/%s.json", strings.ToLower(currency)))
	fullURL := *c.baseURL.ResolveReference(relativeCurrencyPath)

	fullURLStr := fullURL.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURLStr, http.NoBody)
	if err != nil {
		return dto.RatesResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return dto.RatesResponse{}, fmt.Errorf("failed to make request to currency API: %w", err)
	}

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			c.logger.Error("failed to close response body", zap.Error(closeErr))
			err = errors.Join(err, closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return dto.RatesResponse{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.RatesResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var rateResponse dto.RatesResponse
	if err := json.Unmarshal(body, &rateResponse); err != nil {
		return dto.RatesResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return rateResponse, nil
}
