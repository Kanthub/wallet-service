package service

import (
	"context"
	"fmt"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type AssetAmountStatService interface {
	CreateStat(ctx context.Context, req CreateAssetAmountStatRequest) (*backend.AssetAmountStat, error)
	UpdateStat(ctx context.Context, guid string, updates map[string]interface{}) error
	GetByAssetAndDate(ctx context.Context, assetUUID, date string) (*backend.AssetAmountStat, error)
}

type CreateAssetAmountStatRequest struct {
	AssetUUID string `json:"asset_uuid"`
	TimeDate  string `json:"time_date"`
	Amount    string `json:"amount"`
}

type assetAmountStatService struct {
	db *database.DB
}

func NewAssetAmountStatService(db *database.DB) AssetAmountStatService {
	return &assetAmountStatService{db: db}
}

func (s *assetAmountStatService) CreateStat(
	ctx context.Context,
	req CreateAssetAmountStatRequest,
) (*backend.AssetAmountStat, error) {

	if req.AssetUUID == "" {
		return nil, fmt.Errorf("asset_uuid required")
	}
	if req.TimeDate == "" {
		return nil, fmt.Errorf("time_date required")
	}
	if req.Amount == "" {
		return nil, fmt.Errorf("amount required")
	}

	item := &backend.AssetAmountStat{
		AssetUUID:  req.AssetUUID,
		TimeDate:   req.TimeDate,
		Amount:     req.Amount,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := s.db.BackendAssetAmountStat.StoreAssetAmountStat(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *assetAmountStatService) UpdateStat(
	ctx context.Context,
	guid string,
	updates map[string]interface{},
) error {

	if guid == "" {
		return fmt.Errorf("guid required")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates empty")
	}

	return s.db.BackendAssetAmountStat.UpdateAssetAmountStat(guid, updates)
}

func (s *assetAmountStatService) GetByAssetAndDate(
	ctx context.Context,
	assetUUID, date string,
) (*backend.AssetAmountStat, error) {

	if assetUUID == "" || date == "" {
		return nil, fmt.Errorf("asset_uuid and date required")
	}
	return s.db.BackendAssetAmountStat.GetByAssetAndDate(assetUUID, date)
}
