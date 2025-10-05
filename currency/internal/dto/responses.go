package dto

import (
	"github.com/CargoMan0/currency-service/pkg/currency"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type CurrencyResponseDTO struct {
	Currency string
	Rates    []RateRecordDTO
}

type RateRecordDTO struct {
	Date time.Time
	Rate float32
}

func (dto *CurrencyResponseDTO) ToProtobuf() *currency.GetRateResponse {
	rateRecords := make([]*currency.RateRecord, 0, len(dto.Rates))
	for _, record := range dto.Rates {
		rateRecords = append(
			rateRecords, &currency.RateRecord{
				Date: timestamppb.New(record.Date),
				Rate: record.Rate,
			},
		)
	}

	return &currency.GetRateResponse{
		Currency: dto.Currency,
		Rates:    rateRecords,
	}
}

type RatesResponse struct {
	Date string             `json:"date"`
	Rub  map[string]float64 `json:"rub"`
}
