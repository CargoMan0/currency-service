package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/CargoMan0/currency-service/currency/internal/dto"
	"github.com/CargoMan0/currency-service/pkg/currency"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

func (s CurrencyServer) GetRate(ctx context.Context, request *currency.GetRateRequest) (*currency.GetRateResponse, error) {
	reqDTO := dto.CurrencyRequestFromPbToDTO(request, dto.DefaultBaseCurrency)

	if strings.ToUpper(reqDTO.TargetCurrency) == dto.DefaultBaseCurrency {
		return nil, errors.New("target currency can not be equal to base currency")
	}

	rates, err := s.service.GetCurrencyRatesInInterval(ctx, reqDTO)
	if err != nil {
		return nil, fmt.Errorf("service: get currenct rates in interval: %w", err)
	}

	rateRecords := make([]*currency.RateRecord, len(rates))
	for i, rate := range rates {
		rateRecords[i] = &currency.RateRecord{
			Date: timestamppb.New(rate.Date),
			Rate: rate.Rate,
		}
	}

	return &currency.GetRateResponse{
		Currency: reqDTO.TargetCurrency,
		Rates:    rateRecords,
	}, nil
}
