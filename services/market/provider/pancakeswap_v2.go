package provider

import (
	"context"
	"strings"
	"time"

	"github.com/roothash-pay/wallet-services/services/market/model"
	"github.com/roothash-pay/wallet-services/services/market/provider/utils"
)

type PancakeSwapV2GraphProvider struct {
	graph *utils.GraphClient
	limit int
}

func NewPancakeSwapV2GraphProvider(limit int) *PancakeSwapV2GraphProvider {
	if limit <= 0 {
		limit = 50
	}
	return &PancakeSwapV2GraphProvider{
		graph: utils.NewGraphClient(),
		limit: limit,
	}
}

func (p *PancakeSwapV2GraphProvider) Name() string { return "pancakeswap_v2_graph" }

// 注意：subgraph 可能会随时间变更（如果调用失败，再换最新 endpoint）
const pancakeV2URL = "https://api.thegraph.com/subgraphs/name/pancakeswap/exchange-v2"

func (p *PancakeSwapV2GraphProvider) FetchQuotes(ctx context.Context) ([]model.Quote, error) {
	query := `
query TopPairs($n: Int!) {
  pairs(first: $n, orderBy: volumeUSD, orderDirection: desc) {
    id
    token0 { symbol }
    token1 { symbol }
    token0Price
    token1Price
    volumeUSD
  }
}`

	var resp struct {
		Pairs []struct {
			ID          string                  `json:"id"`
			Token0      struct{ Symbol string } `json:"token0"`
			Token1      struct{ Symbol string } `json:"token1"`
			Token0Price string                  `json:"token0Price"`
			Token1Price string                  `json:"token1Price"`
			VolumeUSD   string                  `json:"volumeUSD"`
		} `json:"pairs"`
	}

	if err := p.graph.Query(ctx, pancakeV2URL, query, map[string]interface{}{"n": p.limit}, &resp); err != nil {
		return nil, err
	}

	isStable := func(sym string) bool {
		s := strings.ToUpper(sym)
		return s == "USDC" || s == "USDT" || s == "BUSD" || s == "DAI"
	}

	now := time.Now()
	quotes := make([]model.Quote, 0, p.limit)

	for _, pair := range resp.Pairs {
		a := strings.ToUpper(pair.Token0.Symbol)
		b := strings.ToUpper(pair.Token1.Symbol)

		if isStable(b) {
			price := utils.ParseFloat(pair.Token0Price)
			if price > 0 {
				quotes = append(quotes, model.Quote{
					BaseAsset:  a,
					QuoteAsset: "USDT",
					Price:      price,
					Source:     p.Name(),
					Timestamp:  now,
				})
			}
		} else if isStable(a) {
			price := utils.ParseFloat(pair.Token1Price)
			if price > 0 {
				quotes = append(quotes, model.Quote{
					BaseAsset:  b,
					QuoteAsset: "USDT",
					Price:      price,
					Source:     p.Name(),
					Timestamp:  now,
				})
			}
		}
	}

	return quotes, nil
}
