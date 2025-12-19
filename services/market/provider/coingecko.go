package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/roothash-pay/wallet-services/services/market/model"
)

type CoinGeckoProvider struct {
	client        *http.Client
	symbolToID    map[string]string // BTC -> bitcoin
	vsCurrency    string            // "usd"
	stableAsQuote string            // "USDT" or "USD" (我们统一输出 QuoteAsset)
}

func NewCoinGeckoProvider(symbolToID map[string]string) *CoinGeckoProvider {
	defaultMap := map[string]string{
		"BTC":   "bitcoin",
		"ETH":   "ethereum",
		"SOL":   "solana",
		"BNB":   "binancecoin",
		"ARB":   "arbitrum",
		"OP":    "optimism",
		"MATIC": "matic-network",
	}

	for k, v := range symbolToID {
		defaultMap[strings.ToUpper(k)] = v
	}

	return &CoinGeckoProvider{
		client:     &http.Client{Timeout: 10 * time.Second},
		symbolToID: defaultMap,
		vsCurrency: "usd",
		// 你 Quote 体系现在偏 USDT，我就统一写 USDT（本质还是美元价）
		stableAsQuote: "USDT",
	}
}

func (p *CoinGeckoProvider) Name() string { return "coingecko" }

type cgResp map[string]struct {
	USD           float64 `json:"usd"`
	USD24hVol     float64 `json:"usd_24h_vol"`
	USD24hChange  float64 `json:"usd_24h_change"`
	LastUpdatedAt int64   `json:"last_updated_at"`
}

func (p *CoinGeckoProvider) FetchQuotes(ctx context.Context) ([]model.Quote, error) {
	ids := make([]string, 0, len(p.symbolToID))
	idToSymbol := make(map[string]string, len(p.symbolToID))
	for sym, id := range p.symbolToID {
		ids = append(ids, id)
		idToSymbol[id] = sym
	}

	endpoint := "https://api.coingecko.com/api/v3/simple/price"
	q := url.Values{}
	q.Set("ids", strings.Join(ids, ","))
	q.Set("vs_currencies", p.vsCurrency)
	q.Set("include_24hr_vol", "true")
	q.Set("include_24hr_change", "true")
	q.Set("include_last_updated_at", "true")

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("coingecko http status %d", resp.StatusCode)
	}

	var out cgResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	quotes := make([]model.Quote, 0, len(out))
	for id, v := range out {
		sym := idToSymbol[id]
		if sym == "" || v.USD <= 0 {
			continue
		}
		ts := time.Now()
		if v.LastUpdatedAt > 0 {
			ts = time.Unix(v.LastUpdatedAt, 0)
		}

		quotes = append(quotes, model.Quote{
			BaseAsset:  sym,
			QuoteAsset: p.stableAsQuote,
			Price:      v.USD,
			Volume24h:  v.USD24hVol,
			// Quote 里如果有 PriceChange24h 字段就填；没有就先不填
			Source:    p.Name(),
			Timestamp: ts,
		})
	}

	return quotes, nil
}
