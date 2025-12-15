package service

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type UpsertAddressAssetRequest struct {
	TokenID     string `json:"token_id"`     // 合约 token_id 或 token 合约地址/标识
	WalletUUID  string `json:"wallet_uuid"`  // 所属钱包
	AddressUUID string `json:"address_uuid"` // 地址记录 ID

	AssetUsdt string `json:"asset_usdt"` // 以 USDT 计价
	AssetUsd  string `json:"asset_usd"`  // 以 USD 计价
	Balance   string `json:"balance"`    // 原始 token 数量（wei 单位）
}

type AddressAssetService interface {
	UpsertAddressAsset(ctx context.Context, req UpsertAddressAssetRequest) (*backend.AddressAsset, error)
	ListAddressAssetsByAddress(ctx context.Context, addressUUID string) ([]*backend.AddressAsset, error)
	GetAddressAsset(ctx context.Context, addressUUID, tokenID string) (*backend.AddressAsset, error)
}

type addressAssetService struct {
	Db *database.DB
}

func NewAddressAssetService(db *database.DB) AddressAssetService {
	return &addressAssetService{Db: db}
}

func (s *addressAssetService) UpsertAddressAsset(
	ctx context.Context,
	req UpsertAddressAssetRequest,
) (*backend.AddressAsset, error) {
	// 1. 基础校验
	if req.AddressUUID == "" {
		return nil, fmt.Errorf("address_uuid is required")
	}
	if req.TokenID == "" {
		return nil, fmt.Errorf("token_id is required")
	}

	assetUsdt, err := decimal.NewFromString(req.AssetUsdt)
	if err != nil {
		return nil, fmt.Errorf("invalid asset_usdt: %w", err)
	}
	assetUsd, err := decimal.NewFromString(req.AssetUsd)
	if err != nil {
		return nil, fmt.Errorf("invalid asset_usd: %w", err)
	}
	balance, err := decimal.NewFromString(req.Balance)
	if err != nil {
		return nil, fmt.Errorf("invalid balance: %w", err)
	}
	if balance.IsNegative() {
		return nil, fmt.Errorf("balance cannot be negative")
	}

	a := &backend.AddressAsset{
		TokenID:     req.TokenID,
		WalletUUID:  req.WalletUUID,
		AddressUUID: req.AddressUUID,
		AssetUsdt:   assetUsdt,
		AssetUsd:    assetUsd,
		Balance:     balance,
	}

	// 2. Upsert 到 DB
	if err := s.Db.BackendAddressAsset.UpsertAddressAsset(a); err != nil {
		return nil, err
	}

	// 3. 返回更新后的记录（可以再读一次，视情况而定）
	// 这里简单返回 a 即可
	return a, nil
}

func (s *addressAssetService) ListAddressAssetsByAddress(
	ctx context.Context,
	addressUUID string,
) ([]*backend.AddressAsset, error) {
	if addressUUID == "" {
		return nil, fmt.Errorf("address_uuid is required")
	}
	return s.Db.BackendAddressAsset.GetByAddressUUID(addressUUID)
}

func (s *addressAssetService) GetAddressAsset(
	ctx context.Context,
	addressUUID, tokenID string,
) (*backend.AddressAsset, error) {
	if addressUUID == "" || tokenID == "" {
		return nil, fmt.Errorf("address_uuid and token_id are required")
	}
	return s.Db.BackendAddressAsset.GetByAddressAndToken(addressUUID, tokenID)
}
