package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/roothash-pay/wallet-services/services/market/model"
)

type DefiLlamaProvider struct {
	client *http.Client
}

func NewDefiLlamaProvider() *DefiLlamaProvider {
	return &DefiLlamaProvider{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *DefiLlamaProvider) Name() string {
	return "defillama"
}

/*
GET https://coins.llama.fi/prices/current/ethereum:0x...
*/

type llamaResp struct {
	Coins map[string]struct {
		Price     float64 `json:"price"`
		Timestamp int64   `json:"timestamp"`
	} `json:"coins"`
}

func (p *DefiLlamaProvider) FetchQuotes(ctx context.Context) ([]model.Quote, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://coins.llama.fi/prices/current/coingecko:ethereum",
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result llamaResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	eth, ok := result.Coins["coingecko:ethereum"]
	if !ok {
		return nil, nil
	}

	return []model.Quote{
		{
			BaseAsset:  "ETH",
			QuoteAsset: "USDT",
			Price:      eth.Price,
			Volume24h:  0,
			Source:     p.Name(),
			Timestamp:  time.Unix(eth.Timestamp, 0),
		},
	}, nil
}
