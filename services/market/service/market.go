package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/roothash-pay/wallet-services/services/market/cache"
	"github.com/roothash-pay/wallet-services/services/market/model"
)

type MarketService interface {
	// 读缓存里的最新价格（只读，不裁决）
	GetPrice(ctx context.Context, symbol string) (*model.Quote, error)
}

type marketService struct {
	cache cache.Cache
}

func NewMarketService(cache cache.Cache) MarketService {
	return &marketService{
		cache: cache,
	}
}

// 读缓存里的最新价格（只读，不裁决）
func (s *marketService) GetPrice(
	ctx context.Context,
	symbol string,
) (*model.Quote, error) {

	key := "price:" + strings.ToUpper(symbol)

	data, ok, err := s.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil // 或者自定义 ErrCacheMiss
	}

	var quote model.Quote
	if err := json.Unmarshal(data, &quote); err != nil {
		return nil, err
	}

	return &quote, nil
}
