package currency

import (
	"context"
	"fmt"
	"github.com/BernsteinMondy/currency-service/gateway/internal/service"
	"github.com/BernsteinMondy/currency-service/pkg/currency"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	currencyGRPCClient currency.CurrencyServiceClient
}

func NewClient(currencyGRPCClient currency.CurrencyServiceClient) *Client {
	return &Client{
		currencyGRPCClient: currencyGRPCClient,
	}
}

func (c *Client) GetCurrencyRates(ctx context.Context, request service.CurrencyRequest) (*service.CurrencyResponse, error) {
	pbResp, err := c.currencyGRPCClient.GetRate(
		ctx, &currency.GetRateRequest{
			Currency: request.Currency,
			DateFrom: timestamppb.New(request.DateFrom),
			DateTo:   timestamppb.New(request.DateTo),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("currency grpc client: get currency rate: %s", err)
	}

	resp := &service.CurrencyResponse{
		Currency: pbResp.GetCurrency(),
		Rates:    make([]service.CurrencyRate, 0, len(pbResp.Rates)),
	}

	for _, rate := range pbResp.Rates {
		resp.Rates = append(
			resp.Rates, service.CurrencyRate{
				Rate: rate.Rate,
				Date: rate.Date.AsTime(),
			},
		)
	}

	return resp, nil
}
