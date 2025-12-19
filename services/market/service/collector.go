package service

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/log"
	"github.com/roothash-pay/wallet-services/services/market/cache"
	"github.com/roothash-pay/wallet-services/services/market/model"
	"github.com/roothash-pay/wallet-services/services/market/provider"
	"github.com/roothash-pay/wallet-services/services/market/resolver"
)

type MarketCollector struct {
	providers []provider.Provider
	resolver  resolver.Resolver
	cache     cache.Cache
}

func NewMarketCollector(
	providers []provider.Provider,
	resolver resolver.Resolver,
	cache cache.Cache,
) *MarketCollector {
	return &MarketCollector{
		providers: providers,
		resolver:  resolver,
		cache:     cache,
	}
}

func (mc *MarketCollector) Collect(ctx context.Context) (map[string]model.Quote, error) {
	allQuotes := make([]model.Quote, 0, 256)

	for _, p := range mc.providers {
		quotes, err := p.FetchQuotes(ctx)
		if err != nil {
			log.Warn("fetch quotes failed", "provider", p.Name(), "err", err)
			continue
		}
		log.Info("quotes fetched", "provider", p.Name(), "count", len(quotes))
		allQuotes = append(allQuotes, quotes...)
	}
	if len(allQuotes) == 0 {
		log.Warn("no quotes collected")
		return nil, nil
	}

	finalQuotes, err := mc.resolver.Resolve(allQuotes)
	if err != nil {
		return nil, err
	}

	// 比如 resolver 返回 map[string]Quote
	for symbol, q := range finalQuotes {
		key := "price:" + symbol

		buf, err := json.Marshal(q)
		if err != nil {
			log.Error("marshal quote failed", "key", key, "err", err)
			continue
		}

		// log.Info("cache set try", "key", key, "price", q.Price)

		if err := mc.cache.Set(ctx, key, buf); err != nil {
			log.Error("cache set failed", "key", key, "err", err)
		} else {
			log.Info("cache set success", "key", key)
		}

		// // debug 读取确认
		// val, ok, err := mc.cache.Get(ctx, key)
		// if err != nil {
		// 	log.Error("cache get after set failed", "key", key, "err", err)
		// } else {
		// 	log.Info("cache get after set ok", "key", key, "type", fmt.Sprintf("%T", val), "ok", ok)
		// }

	}

	log.Info("market quotes collected", "symbols", len(finalQuotes))
	return finalQuotes, nil
}
