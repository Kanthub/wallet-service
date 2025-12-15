package service

import (
	"context"
	"fmt"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type WalletAddressService interface {
	CreateWalletAddress(ctx context.Context, req CreateWalletAddressRequest) (*backend.WalletAddress, error)
	UpdateWalletAddress(ctx context.Context, req UpdateWalletAddressRequest) error
	GetWalletAddress(ctx context.Context, guid string) (*backend.WalletAddress, error)
	GetByAddress(ctx context.Context, address string) (*backend.WalletAddress, error)
	ListByWalletUUID(ctx context.Context, walletUUID string) ([]*backend.WalletAddress, error)
}

type CreateWalletAddressRequest struct {
	AddressIndex int64  `json:"address_index"`
	Address      string `json:"address"`
	WalletUUID   string `json:"wallet_uuid"`
	ChainID      string `json:"chain_id"`
}

type UpdateWalletAddressRequest struct {
	Guid    string                 `json:"guid"`
	Updates map[string]interface{} `json:"updates"`
}

type walletAddressService struct {
	db *database.DB
}

func NewWalletAddressService(db *database.DB) WalletAddressService {
	return &walletAddressService{db: db}
}

func (s *walletAddressService) CreateWalletAddress(
	ctx context.Context,
	req CreateWalletAddressRequest,
) (*backend.WalletAddress, error) {

	if req.Address == "" || req.WalletUUID == "" {
		return nil, fmt.Errorf("address and wallet_uuid required")
	}
	if req.AddressIndex <= 0 {
		return nil, fmt.Errorf("address_index must be > 0")
	}

	item := &backend.WalletAddress{
		AddressIndex: req.AddressIndex,
		Address:      req.Address,
		WalletUUID:   req.WalletUUID,
		ChainID:      req.ChainID,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}

	if err := s.db.BackendWalletAddress.StoreWalletAddress(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *walletAddressService) UpdateWalletAddress(
	ctx context.Context,
	req UpdateWalletAddressRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid required")
	}
	if len(req.Updates) == 0 {
		return fmt.Errorf("updates empty")
	}

	return s.db.BackendWalletAddress.UpdateWalletAddress(req.Guid, req.Updates)
}

func (s *walletAddressService) GetWalletAddress(
	ctx context.Context,
	guid string,
) (*backend.WalletAddress, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	return s.db.BackendWalletAddress.GetByGuid(guid)
}

func (s *walletAddressService) GetByAddress(
	ctx context.Context,
	address string,
) (*backend.WalletAddress, error) {

	if address == "" {
		return nil, fmt.Errorf("address required")
	}
	return s.db.BackendWalletAddress.GetByAddress(address)
}

func (s *walletAddressService) ListByWalletUUID(
	ctx context.Context,
	walletUUID string,
) ([]*backend.WalletAddress, error) {

	if walletUUID == "" {
		return nil, fmt.Errorf("wallet_uuid required")
	}
	return s.db.BackendWalletAddress.GetByWalletUUID(walletUUID)
}
