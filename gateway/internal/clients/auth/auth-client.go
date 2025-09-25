package auth

import (
	"context"
	"errors"
	"fmt"
	errors2 "github.com/BernsteinMondy/currency-service/gateway/internal/clients/auth/errors"
	"github.com/BernsteinMondy/currency-service/gateway/internal/config"
	"github.com/BernsteinMondy/currency-service/gateway/internal/service"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	pingPath     = "/ping"
	generatePath = "/generate"
	validatePath = "/validate"

	authorizationHeader = "Authorization"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

var _ service.AuthClient = (*Client)(nil)

func NewClient(cfg config.AuthConfig) (*Client, error) {
	parsedURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Client{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Duration(cfg.TimeoutSeconds) * time.Second,
		},
	}, nil
}

func (c *Client) CloseIdleConnections() {
	c.httpClient.CloseIdleConnections()
}

func (c *Client) Ping() (_ string, err error) {
	relativePingPath, _ := url.Parse(pingPath)
	fullURL := *c.baseURL.ResolveReference(relativePingPath)

	resp, err := c.httpClient.Get(fullURL.String())
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close response body: %w", closeErr))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	return string(body), nil
}

func (c *Client) GenerateToken(ctx context.Context, login string) (_ string, err error) {
	relativeGeneratePath, _ := url.Parse(generatePath)
	fullURL := *c.baseURL.ResolveReference(relativeGeneratePath)

	query := fullURL.Query()
	query.Set("login", login)
	fullURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL.String(), http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("httpClient.Do: %w", err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close response body: %w", closeErr))
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %w", err)
		}
		return string(bodyBytes), nil
	case http.StatusBadRequest:
		return "", fmt.Errorf("%w: bad request", errors2.ErrClientTokenGeneration)
	case http.StatusUnauthorized:
		return "", fmt.Errorf("%w: unauthorized", errors2.ErrClientInvalidCredentials)
	default:
		return "", fmt.Errorf("%w: %d", errors2.ErrClientUnexpectedStatusCode, resp.StatusCode)
	}
}

func (c *Client) ValidateToken(ctx context.Context, token string) (err error) {
	relativeValidatePath, _ := url.Parse(validatePath)
	fullURL := *c.baseURL.ResolveReference(relativeValidatePath)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL.String(), http.NoBody)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Set(authorizationHeader, "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpClient.Do: %w", err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to close response body: %w", closeErr))
		}
	}()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	errorMessage := string(body)

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", errors2.ErrClientTokenNotFound, errorMessage)
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", errors2.ErrClientInvalidOrExpiredToken, errorMessage)
	default:
		return fmt.Errorf("%w %d: %s", errors2.ErrClientUnexpectedStatusCode, resp.StatusCode, errorMessage)
	}
}
