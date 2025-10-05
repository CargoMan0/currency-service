package currency

import (
	"context"
	"fmt"
	"github.com/CargoMan0/currency-service/gateway/internal/dto"
	"github.com/CargoMan0/currency-service/gateway/internal/service"
	"github.com/CargoMan0/currency-service/pkg/currency"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	currencyGRPCClient currency.CurrencyServiceClient
}

var _ service.CurrencyClient = (*Client)(nil)

func NewClient(currencyGRPCClient currency.CurrencyServiceClient) *Client {
	return &Client{
		currencyGRPCClient: currencyGRPCClient,
	}
}

func (c *Client) GetCurrencyRates(ctx context.Context, request dto.CurrencyRequest) (*dto.CurrencyResponse, error) {
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

	resp := &dto.CurrencyResponse{
		Currency: pbResp.GetCurrency(),
		Rates:    make([]dto.CurrencyRate, 0, len(pbResp.Rates)),
	}

	for _, rate := range pbResp.Rates {
		resp.Rates = append(
			resp.Rates, dto.CurrencyRate{
				Rate: rate.Rate,
				Date: rate.Date.AsTime(),
			},
		)
	}

	return resp, nil
}
