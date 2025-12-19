package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
	"github.com/roothash-pay/wallet-services/services/market/cache"
	"github.com/roothash-pay/wallet-services/services/market/model"
)

type MarketPriceService interface {
	SetMarketPrice(ctx context.Context, req SetMarketPriceRequest) error
	GetByTokenID(ctx context.Context, tokenID string) (*backend.MarketPrice, error)
	GetByGuid(ctx context.Context, guid string) (*backend.MarketPrice, error)

	GetPrice(ctx context.Context, symbol string) (*model.Quote, error)
}

type SetMarketPriceRequest struct {
	ChainID     string `json:"chain_id"`
	TokenID     string `json:"token_id"`
	UsdtPrice   string `json:"usdt_price"`
	UsdPrice    string `json:"usd_price"`
	MarketCap   string `json:"market_cap"`
	Liquidity   string `json:"liquidity"`
	Volume24h   string `json:"24h_volume"`
	PriceChange string `json:"price_change"`
	Ranking     string `json:"ranking"`
}

type marketPriceService struct {
	db    *database.DB
	cache cache.Cache
}

func NewMarketPriceService(db *database.DB, cache cache.Cache) MarketPriceService {
	return &marketPriceService{db: db, cache: cache}
}

// 读缓存里的最新价格（只读，不裁决）
func (s *marketPriceService) GetPrice(
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

func (s *marketPriceService) SetMarketPrice(
	ctx context.Context,
	req SetMarketPriceRequest,
) error {

	if req.TokenID == "" {
		return fmt.Errorf("token_id required")
	}
	if req.UsdtPrice == "" || req.UsdPrice == "" {
		return fmt.Errorf("price required")
	}

	// 先查是否存在
	existing, err := s.db.BackendMarketPrice.GetByTokenID(req.TokenID)
	if err == nil && existing != nil {
		// update
		return s.db.BackendMarketPrice.UpdateMarketPrice(
			existing.Guid,
			map[string]interface{}{
				"chain_id":     req.ChainID,
				"usdt_price":   req.UsdtPrice,
				"usd_price":    req.UsdPrice,
				"market_cap":   req.MarketCap,
				"liquidity":    req.Liquidity,
				"24h_volume":   req.Volume24h,
				"price_change": req.PriceChange,
				"ranking":      req.Ranking,
				"updated_at":   time.Now(),
			},
		)
	}

	// create
	item := &backend.MarketPrice{
		ChainID:     req.ChainID,
		TokenID:     req.TokenID,
		UsdtPrice:   req.UsdtPrice,
		UsdPrice:    req.UsdPrice,
		MarketCap:   req.MarketCap,
		Liquidity:   req.Liquidity,
		Volume24h:   req.Volume24h,
		PriceChange: req.PriceChange,
		Ranking:     req.Ranking,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	return s.db.BackendMarketPrice.StoreMarketPrice(item)
}

func (s *marketPriceService) GetByTokenID(
	ctx context.Context,
	tokenID string,
) (*backend.MarketPrice, error) {

	if tokenID == "" {
		return nil, fmt.Errorf("token_id required")
	}
	return s.db.BackendMarketPrice.GetByTokenID(tokenID)
}

func (s *marketPriceService) GetByGuid(
	ctx context.Context,
	guid string,
) (*backend.MarketPrice, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	return s.db.BackendMarketPrice.GetByGuid(guid)
}
