package service

import (
	"context"
	"fmt"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type WalletAddressNoteService interface {
	CreateWalletAddressNote(
		ctx context.Context,
		req CreateWalletAddressNoteRequest,
	) (*backend.WalletAddressNote, error)

	UpdateWalletAddressNote(
		ctx context.Context,
		req UpdateWalletAddressNoteRequest,
	) error

	GetWalletAddressNote(
		ctx context.Context,
		guid string,
	) (*backend.WalletAddressNote, error)

	GetByDeviceUUID(
		ctx context.Context,
		deviceUUID string,
	) ([]*backend.WalletAddressNote, error)
}

type CreateWalletAddressNoteRequest struct {
	DeviceUUID string `json:"device_uuid"`
	ChainID    string `json:"chain_id"`
	Address    string `json:"address"`
	Memo       string `json:"memo"`
}

type UpdateWalletAddressNoteRequest struct {
	Guid    string                 `json:"guid"`
	Updates map[string]interface{} `json:"updates"`
}

type walletAddressNoteService struct {
	db *database.DB
}

func NewWalletAddressNoteService(db *database.DB) WalletAddressNoteService {
	return &walletAddressNoteService{db: db}
}

func (s *walletAddressNoteService) CreateWalletAddressNote(
	ctx context.Context,
	req CreateWalletAddressNoteRequest,
) (*backend.WalletAddressNote, error) {

	if req.DeviceUUID == "" {
		return nil, fmt.Errorf("device_uuid required")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("address required")
	}
	if req.Memo == "" {
		return nil, fmt.Errorf("memo required")
	}

	item := &backend.WalletAddressNote{
		DeviceUUID: req.DeviceUUID,
		ChainID:    req.ChainID,
		Address:    req.Address,
		Memo:       req.Memo,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := s.db.BackendWalletAddressNote.StoreWalletAddressNote(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *walletAddressNoteService) UpdateWalletAddressNote(
	ctx context.Context,
	req UpdateWalletAddressNoteRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid required")
	}
	if len(req.Updates) == 0 {
		return fmt.Errorf("updates empty")
	}

	return s.db.BackendWalletAddressNote.UpdateWalletAddressNote(
		req.Guid,
		req.Updates,
	)
}

func (s *walletAddressNoteService) GetWalletAddressNote(
	ctx context.Context,
	guid string,
) (*backend.WalletAddressNote, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}

	return s.db.BackendWalletAddressNote.GetByGuid(guid)
}

func (s *walletAddressNoteService) GetByDeviceUUID(
	ctx context.Context,
	deviceUUID string,
) ([]*backend.WalletAddressNote, error) {

	if deviceUUID == "" {
		return nil, fmt.Errorf("device_uuid required")
	}

	return s.db.BackendWalletAddressNote.GetByDeviceUUID(deviceUUID)
}
