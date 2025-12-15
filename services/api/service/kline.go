package service

import (
	"context"
	"fmt"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type KlineService interface {
	SetKline(ctx context.Context, req SetKlineRequest) error
	GetKlines(
		ctx context.Context,
		tokenID string,
		interval string,
		start time.Time,
		end time.Time,
		limit int,
	) ([]*backend.Kline, error)
}

type SetKlineRequest struct {
	TokenID      string    `json:"token_id"`
	TimeInterval string    `json:"time_interval"`
	OpenTime     time.Time `json:"open_time"`

	OpenPrice  string `json:"open_price"`
	HighPrice  string `json:"high_price"`
	LowPrice   string `json:"low_price"`
	ClosePrice string `json:"close_price"`

	Volume      string `json:"volume"`
	QuoteVolume string `json:"quote_volume"`
	TradeCount  string `json:"trade_count"`
}

type klineService struct {
	db *database.DB
}

func NewKlineService(db *database.DB) KlineService {
	return &klineService{db: db}
}

// SetKline 幂等写入（token + interval + open_time 唯一）
func (s *klineService) SetKline(
	ctx context.Context,
	req SetKlineRequest,
) error {

	if req.TokenID == "" || req.TimeInterval == "" {
		return fmt.Errorf("token_id and time_interval required")
	}

	existing, err := s.db.BackendKline.GetByKey(
		req.TokenID,
		req.TimeInterval,
		req.OpenTime,
	)

	if err == nil && existing != nil {
		// update
		return s.db.BackendKline.UpdateKline(
			existing.Guid,
			map[string]interface{}{
				"open_price":   req.OpenPrice,
				"high_price":   req.HighPrice,
				"low_price":    req.LowPrice,
				"close_price":  req.ClosePrice,
				"volume":       req.Volume,
				"quote_volume": req.QuoteVolume,
				"trade_count":  req.TradeCount,
				"updated_at":   time.Now(),
			},
		)
	}

	// create
	item := &backend.Kline{
		TokenID:      req.TokenID,
		TimeInterval: req.TimeInterval,
		OpenTime:     req.OpenTime,
		OpenPrice:    req.OpenPrice,
		HighPrice:    req.HighPrice,
		LowPrice:     req.LowPrice,
		ClosePrice:   req.ClosePrice,
		Volume:       req.Volume,
		QuoteVolume:  req.QuoteVolume,
		TradeCount:   req.TradeCount,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}

	return s.db.BackendKline.StoreKline(item)
}

// GetKlines 给前端画图用
func (s *klineService) GetKlines(
	ctx context.Context,
	tokenID string,
	interval string,
	start time.Time,
	end time.Time,
	limit int,
) ([]*backend.Kline, error) {

	if tokenID == "" || interval == "" {
		return nil, fmt.Errorf("token_id and interval required")
	}

	if limit <= 0 {
		limit = 500
	}

	return s.db.BackendKline.GetKlines(
		tokenID,
		interval,
		start,
		end,
		limit,
	)
}
