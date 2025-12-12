package service

import (
	"context"
	"fmt"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type WalletAssetService interface {
	CreateWalletAsset(ctx context.Context, req CreateWalletAssetRequest) (*backend.WalletAsset, error)
	UpdateWalletAsset(ctx context.Context, req UpdateWalletAssetRequest) error
	GetWalletAsset(ctx context.Context, guid string) (*backend.WalletAsset, error)
	GetByTokenChain(ctx context.Context, tokenID, chainID string) (*backend.WalletAsset, error)
}

type CreateWalletAssetRequest struct {
	WalletUUID string `json:"wallet_uuid"`
	TokenID    string `json:"token_id"`
	ChainID    string `json:"chain_id"`
	Balance    string `json:"balance"`
	AssetUsdt  string `json:"asset_usdt"`
	AssetUsd   string `json:"asset_usd"`
}

type UpdateWalletAssetRequest struct {
	Guid    string                 `json:"guid"`
	Updates map[string]interface{} `json:"updates"`
}

type walletAssetService struct {
	db *database.DB
}

func NewWalletAssetService(db *database.DB) WalletAssetService {
	return &walletAssetService{db: db}
}

func (s *walletAssetService) CreateWalletAsset(
	ctx context.Context,
	req CreateWalletAssetRequest,
) (*backend.WalletAsset, error) {

	if req.WalletUUID == "" {
		return nil, fmt.Errorf("wallet_uuid required")
	}
	if req.Balance == "" {
		return nil, fmt.Errorf("balance required")
	}

	item := &backend.WalletAsset{
		WalletUUID: req.WalletUUID,
		TokenID:    req.TokenID,
		ChainID:    req.ChainID,
		Balance:    req.Balance,
		AssetUsdt:  req.AssetUsdt,
		AssetUsd:   req.AssetUsd,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := s.db.BackendWalletAsset.StoreWalletAsset(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *walletAssetService) UpdateWalletAsset(
	ctx context.Context,
	req UpdateWalletAssetRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid required")
	}
	if len(req.Updates) == 0 {
		return fmt.Errorf("updates empty")
	}

	return s.db.BackendWalletAsset.UpdateWalletAsset(req.Guid, req.Updates)
}

func (s *walletAssetService) GetWalletAsset(
	ctx context.Context,
	guid string,
) (*backend.WalletAsset, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	return s.db.BackendWalletAsset.GetByGuid(guid)
}

func (s *walletAssetService) GetByTokenChain(
	ctx context.Context,
	tokenID, chainID string,
) (*backend.WalletAsset, error) {

	if tokenID == "" || chainID == "" {
		return nil, fmt.Errorf("token_id and chain_id required")
	}
	return s.db.BackendWalletAsset.GetByTokenChain(tokenID, chainID)
}
