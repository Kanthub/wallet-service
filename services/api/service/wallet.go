package service

import (
	"context"
	"fmt"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type WalletService interface {
	CreateWallet(ctx context.Context, req CreateWalletRequest) (*backend.Wallet, error)
	UpdateWallet(ctx context.Context, req UpdateWalletRequest) error
	GetWallet(ctx context.Context, guid string) (*backend.Wallet, error)
	GetByWalletUUID(ctx context.Context, walletUUID string) (*backend.Wallet, error)
	ListWallets(
		ctx context.Context,
		page, pageSize int,
		filters map[string]interface{},
	) ([]*backend.Wallet, int64, error)
}

type CreateWalletRequest struct {
	DeviceUUID string `json:"device_uuid"`
	WalletUUID string `json:"wallet_uuid"`
	ChainID    string `json:"chain_id"`
	WalletName string `json:"wallet_name"`
	AssetUsdt  string `json:"asset_usdt"`
	AssetUsd   string `json:"asset_usd"`
}

type UpdateWalletRequest struct {
	Guid    string                 `json:"guid"`
	Updates map[string]interface{} `json:"updates"`
}

type walletService struct {
	db *database.DB
}

func NewWalletService(db *database.DB) WalletService {
	return &walletService{db: db}
}

func (s *walletService) CreateWallet(
	ctx context.Context,
	req CreateWalletRequest,
) (*backend.Wallet, error) {

	if req.DeviceUUID == "" || req.WalletUUID == "" {
		return nil, fmt.Errorf("device_uuid and wallet_uuid required")
	}

	wallet := &backend.Wallet{
		DeviceUUID: req.DeviceUUID,
		WalletUUID: req.WalletUUID,
		ChainID:    req.ChainID,
		WalletName: req.WalletName,
		AssetUsdt:  req.AssetUsdt,
		AssetUsd:   req.AssetUsd,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := s.db.BackendWallet.StoreWallet(wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) UpdateWallet(
	ctx context.Context,
	req UpdateWalletRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid required")
	}
	if len(req.Updates) == 0 {
		return fmt.Errorf("updates empty")
	}

	return s.db.BackendWallet.UpdateWallet(req.Guid, req.Updates)
}

func (s *walletService) GetWallet(
	ctx context.Context,
	guid string,
) (*backend.Wallet, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	return s.db.BackendWallet.GetByGuid(guid)
}

func (s *walletService) GetByWalletUUID(
	ctx context.Context,
	walletUUID string,
) (*backend.Wallet, error) {

	if walletUUID == "" {
		return nil, fmt.Errorf("wallet_uuid required")
	}
	return s.db.BackendWallet.GetByWalletUUID(walletUUID)
}

func (s *walletService) ListWallets(
	ctx context.Context,
	page, pageSize int,
	filters map[string]interface{},
) ([]*backend.Wallet, int64, error) {

	return s.db.BackendWallet.GetWalletList(page, pageSize, filters)
}
