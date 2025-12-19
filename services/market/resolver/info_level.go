package resolver

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/roothash-pay/wallet-services/services/market/model"
)

type InfoLevelResolver struct {
	MaxAge time.Duration
}

func NewInfoLevelResolver() *InfoLevelResolver {
	return &InfoLevelResolver{
		MaxAge: 2 * time.Minute,
	}
}

func (r *InfoLevelResolver) Resolve(quotes []model.Quote) (map[string]model.Quote, error) {
	if len(quotes) == 0 {
		return nil, errors.New("no quotes")
	}
	log.Info(
		"resolver start",
		"resolver", "info_level",
		"input_quotes", len(quotes),
	)

	result := make(map[string]model.Quote)
	now := time.Now()

	for _, q := range quotes {
		// 过期丢弃
		if now.Sub(q.Timestamp) > r.MaxAge {
			continue
		}

		// 价格非法
		if q.Price <= 0 {
			continue
		}

		// 极简规则：同一 BaseAsset，选 volume 最大的
		existing, ok := result[q.BaseAsset]
		if !ok || q.Volume24h > existing.Volume24h {
			log.Debug(
				"resolver select quote",
				"asset", q.BaseAsset,
				"price", q.Price,
				"volume24h", q.Volume24h,
				"source", q.Source,
			)
			result[q.BaseAsset] = q
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no valid quotes after resolve")
	}
	log.Info(
		"resolver finished",
		"resolver", "info_level",
		"output_assets", len(result),
	)

	return result, nil
}
