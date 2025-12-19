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

type CMCProvider struct {
	client  *http.Client
	apiKey  string
	symbols []string // ["BTC","ETH",...]
	baseURL string   // https://pro-api.coinmarketcap.com
	quoteAs string   // "USDT" or "USD"
}

func NewCMCProvider(apiKey string, symbols []string) *CMCProvider {
	ss := make([]string, 0, len(symbols))
	for _, s := range symbols {
		s = strings.TrimSpace(strings.ToUpper(s))
		if s != "" {
			ss = append(ss, s)
		}
	}
	return &CMCProvider{
		client:  &http.Client{Timeout: 10 * time.Second},
		apiKey:  apiKey,
		symbols: ss,
		baseURL: "https://pro-api.coinmarketcap.com",
		quoteAs: "USDT",
	}
}

func (p *CMCProvider) Name() string { return "coinmarketcap" }

type cmcResp struct {
	Status struct {
		Timestamp string `json:"timestamp"`
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_message"`
	} `json:"status"`
	Data map[string]struct {
		Symbol string `json:"symbol"`
		Quote  map[string]struct {
			Price            float64 `json:"price"`
			Volume24h        float64 `json:"volume_24h"`
			PercentChange24h float64 `json:"percent_change_24h"`
			LastUpdated      string  `json:"last_updated"`
		} `json:"quote"`
	} `json:"data"`
}

func (p *CMCProvider) FetchQuotes(ctx context.Context) ([]model.Quote, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("cmc api key is empty")
	}
	if len(p.symbols) == 0 {
		return nil, nil
	}

	endpoint := p.baseURL + "/v2/cryptocurrency/quotes/latest"
	q := url.Values{}
	q.Set("symbol", strings.Join(p.symbols, ","))
	q.Set("convert", "USD")

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("cmc http status %d", resp.StatusCode)
	}

	var out cmcResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if out.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("cmc error %d: %s", out.Status.ErrorCode, out.Status.ErrorMsg)
	}

	quotes := make([]model.Quote, 0, len(out.Data))
	for sym, item := range out.Data {
		qusd, ok := item.Quote["USD"]
		if !ok || qusd.Price <= 0 {
			continue
		}

		ts := time.Now()
		if qusd.LastUpdated != "" {
			if t, err := time.Parse(time.RFC3339, qusd.LastUpdated); err == nil {
				ts = t
			}
		}

		quotes = append(quotes, model.Quote{
			BaseAsset:  strings.ToUpper(sym),
			QuoteAsset: p.quoteAs,
			Price:      qusd.Price,
			Volume24h:  qusd.Volume24h,
			Source:     p.Name(),
			Timestamp:  ts,
		})
	}
	return quotes, nil
}
