package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/roothash-pay/wallet-services/services/market/model"
)

type CryptoCompareProvider struct {
	client *http.Client
	apiKey string
}

func NewCryptoCompareProvider(apiKey string) *CryptoCompareProvider {
	return &CryptoCompareProvider{
		client: &http.Client{Timeout: 10 * time.Second},
		apiKey: apiKey,
	}
}

func (p *CryptoCompareProvider) Name() string {
	return "cryptocompare"
}

/*
GET https://min-api.cryptocompare.com/data/pricemultifull
?fsyms=BTC,ETH
&tsyms=USDT
*/

type ccResp struct {
	RAW map[string]map[string]struct {
		PRICE      float64 `json:"PRICE"`
		VOLUME24H  float64 `json:"VOLUME24H"`
		LASTUPDATE int64   `json:"LASTUPDATE"`
	} `json:"RAW"`
}

func (p *CryptoCompareProvider) FetchQuotes(ctx context.Context) ([]model.Quote, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://min-api.cryptocompare.com/data/pricemultifull?fsyms=BTC,ETH,BNB,SOL&tsyms=USDT",
		nil,
	)
	if err != nil {
		return nil, err
	}

	if p.apiKey != "" {
		req.Header.Set("Authorization", "Apikey "+p.apiKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ccResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	quotes := make([]model.Quote, 0, 8)

	for base, m := range result.RAW {
		usdt, ok := m["USDT"]
		if !ok {
			continue
		}

		quotes = append(quotes, model.Quote{
			BaseAsset:  base,
			QuoteAsset: "USDT",
			Price:      usdt.PRICE,
			Volume24h:  usdt.VOLUME24H,
			Source:     p.Name(),
			Timestamp:  time.Unix(usdt.LASTUPDATE, 0),
		})
	}

	return quotes, nil
}
