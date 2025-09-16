package service

type currencyClient interface {
}

type CurrencyService struct {
	currencyClient currencyClient
}

func NewCurrencyService(currencyClient currencyClient) *CurrencyService {
	return &CurrencyService{
		currencyClient: currencyClient,
	}
}
