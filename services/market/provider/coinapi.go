package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/roothash-pay/wallet-services/services/market/model"
)

type CoinAPIProvider struct {
	client *http.Client
	apiKey string
}

func NewCoinAPIProvider(apiKey string) *CoinAPIProvider {
	return &CoinAPIProvider{
		client: &http.Client{Timeout: 10 * time.Second},
		apiKey: apiKey,
	}
}

func (p *CoinAPIProvider) Name() string {
	return "coinapi"
}

/*
GET https://rest.coinapi.io/v1/assets
*/

type coinAPIAsset struct {
	AssetID   string  `json:"asset_id"`
	PriceUSDT float64 `json:"price_usd"`
}

func (p *CoinAPIProvider) FetchQuotes(ctx context.Context) ([]model.Quote, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://rest.coinapi.io/v1/assets",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-CoinAPI-Key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var assets []coinAPIAsset
	if err := json.NewDecoder(resp.Body).Decode(&assets); err != nil {
		return nil, err
	}

	quotes := make([]model.Quote, 0, 16)
	now := time.Now()

	for _, a := range assets {
		if a.PriceUSDT <= 0 {
			continue
		}

		quotes = append(quotes, model.Quote{
			BaseAsset:  a.AssetID,
			QuoteAsset: "USDT",
			Price:      a.PriceUSDT,
			Volume24h:  0,
			Source:     p.Name(),
			Timestamp:  now,
		})
	}

	return quotes, nil
}
