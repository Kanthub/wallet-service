package service

import (
	"context"
	"fmt"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type FiatCurrencyRateService interface {
	SetRate(ctx context.Context, key, value string) error
	GetRate(ctx context.Context, key string) (*backend.FiatCurrencyRate, error)
}

type fiatCurrencyRateService struct {
	db *database.DB
}

func NewFiatCurrencyRateService(db *database.DB) FiatCurrencyRateService {
	return &fiatCurrencyRateService{db: db}
}

// SetRate 设置 / 更新某个法币汇率（幂等）
// key: 例如 USD_CNY
// value: 例如 7.25
func (s *fiatCurrencyRateService) SetRate(
	ctx context.Context,
	key string,
	value string,
) error {

	if key == "" {
		return fmt.Errorf("key_name required")
	}
	if value == "" {
		return fmt.Errorf("value_data required")
	}

	// 先查是否存在
	existing, err := s.db.BackendFiatCurrencyRate.GetByKeyName(key)
	if err == nil && existing != nil {
		// 已存在，更新
		return s.db.BackendFiatCurrencyRate.UpdateFiatCurrencyRate(
			existing.Guid,
			map[string]interface{}{
				"value_data": value,
				"updated_at": time.Now(),
			},
		)
	}

	// 不存在，创建
	item := &backend.FiatCurrencyRate{
		KeyName:    key,
		ValueData:  value,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	return s.db.BackendFiatCurrencyRate.StoreFiatCurrencyRate(item)
}

func (s *fiatCurrencyRateService) GetRate(
	ctx context.Context,
	key string,
) (*backend.FiatCurrencyRate, error) {

	if key == "" {
		return nil, fmt.Errorf("key_name required")
	}

	return s.db.BackendFiatCurrencyRate.GetByKeyName(key)
}
