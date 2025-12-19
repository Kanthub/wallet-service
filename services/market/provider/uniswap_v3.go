package provider

import (
	"context"

	"strings"
	"time"

	"github.com/roothash-pay/wallet-services/services/market/model"
	"github.com/roothash-pay/wallet-services/services/market/provider/utils"
)

type UniswapV3GraphProvider struct {
	graph *utils.GraphClient
	limit int
}

func NewUniswapV3GraphProvider(limit int) *UniswapV3GraphProvider {
	if limit <= 0 {
		limit = 50
	}
	return &UniswapV3GraphProvider{
		graph: utils.NewGraphClient(),
		limit: limit,
	}
}

func (p *UniswapV3GraphProvider) Name() string { return "uniswap_v3_graph" }

// const uniswapV3URL = "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3"
const uniswapV3URL = "https://gateway.thegraph.com/api/fc75acc62ddc7e0abbb1276c0dfdd386/subgraphs/id/5zvR82QoaXYFyDEKLZ9t6v9adgnptxYpKpSbxtgVENFV"

func (p *UniswapV3GraphProvider) FetchQuotes(ctx context.Context) ([]model.Quote, error) {
	// 取 top pools
	// 	query := `
	// query TopPools($n: Int!) {
	//   pools(first: $n, orderBy: volumeUSD, orderDirection: desc) {
	//     id
	//     token0 { symbol decimals }
	//     token1 { symbol decimals }
	//     token0Price
	//     token1Price
	//     volumeUSD
	//     totalValueLockedUSD
	//   }
	// }
	//   `

	query := `
query Pools($first: Int = 10) {
  pools(first: $first, orderBy: volumeUSD, orderDirection: desc) {
    id
    feeTier
    token0 {
      symbol
      decimals
    }
    token1 {
      symbol
      decimals
    }
    totalValueLockedUSD
    }
}
  `
	var resp struct {
		Pools []struct {
			ID          string                  `json:"id"`
			Token0      struct{ Symbol string } `json:"token0"`
			Token1      struct{ Symbol string } `json:"token1"`
			Token0Price string                  `json:"token0Price"` // token0 in token1
			Token1Price string                  `json:"token1Price"` // token1 in token0
			VolumeUSD   string                  `json:"volumeUSD"`
			TVLUSD      string                  `json:"totalValueLockedUSD"`
		} `json:"pools"`
	}

	if err := p.graph.Query(ctx, uniswapV3URL, query, map[string]interface{}{"n": p.limit}, &resp); err != nil {
		return nil, err
	}

	isStable := func(sym string) bool {
		s := strings.ToUpper(sym)
		return s == "USDC" || s == "USDT" || s == "DAI"
	}

	quotes := make([]model.Quote, 0, p.limit)
	now := time.Now()

	for _, pool := range resp.Pools {
		a := strings.ToUpper(pool.Token0.Symbol)
		b := strings.ToUpper(pool.Token1.Symbol)
		volumeUSD := utils.ParseFloat(pool.VolumeUSD)
		if volumeUSD < 1000 {
			continue
		}

		// 只输出 stable 计价的那一边（信息级）
		// token0Price 表示 token0 = token0Price * token1
		// token1Price 表示 token1 = token1Price * token0
		if isStable(b) {
			// token0 priced in stable token1
			price := utils.ParseFloat(pool.Token0Price)
			if price > 0 {
				quotes = append(quotes, model.Quote{
					BaseAsset:  a,
					QuoteAsset: "USDT",
					Price:      price,
					Volume24h:  volumeUSD,
					Source:     p.Name(),
					Timestamp:  now,
				})
			}
		} else if isStable(a) {
			// token1 priced in stable token0
			price := utils.ParseFloat(pool.Token1Price)
			if price > 0 {
				quotes = append(quotes, model.Quote{
					BaseAsset:  b,
					QuoteAsset: "USDT",
					Price:      price,
					Volume24h:  volumeUSD,
					Source:     p.Name(),
					Timestamp:  now,
				})
			}
		}
	}

	return quotes, nil
}
