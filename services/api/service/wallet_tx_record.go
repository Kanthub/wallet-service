package service

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type WalletTxRecordService interface {
	CreateWalletTx(ctx context.Context, req CreateWalletTxRequest) (*backend.WalletTxRecord, error)
	UpdateWalletTx(ctx context.Context, req UpdateWalletTxRequest) error
	GetWalletTx(ctx context.Context, guid string) (*backend.WalletTxRecord, error)
	GetByOperationID(ctx context.Context, operationID string) ([]*backend.WalletTxRecord, error)
}

type CreateWalletTxRequest struct {
	OperationID string `json:"operation_id"`
	StepIndex   int    `json:"step_index"`
	WalletUUID  string `json:"wallet_uuid"`
	AddressUUID string `json:"address_uuid"`
	TxTime      string `json:"tx_time"`
	ChainID     string `json:"chain_id"`
	TokenID     string `json:"token_id"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      string `json:"amount"`
	Memo        string `json:"memo"`
	TxType      string `json:"tx_type"`
}

type UpdateWalletTxRequest struct {
	Guid    string                 `json:"guid"`
	Updates map[string]interface{} `json:"updates"`
}

type walletTxRecordService struct {
	db *database.DB
}

func NewWalletTxRecordService(db *database.DB) WalletTxRecordService {
	return &walletTxRecordService{db: db}
}

func (s *walletTxRecordService) CreateWalletTx(
	ctx context.Context,
	req CreateWalletTxRequest,
) (*backend.WalletTxRecord, error) {

	if req.WalletUUID == "" {
		return nil, fmt.Errorf("wallet_uuid required")
	}
	if req.FromAddress == "" || req.ToAddress == "" {
		return nil, fmt.Errorf("from_address and to_address required")
	}
	if req.Amount == "" {
		return nil, fmt.Errorf("amount required")
	}

	item := &backend.WalletTxRecord{
		OperationID: req.OperationID,
		StepIndex:   req.StepIndex,
		WalletUUID:  req.WalletUUID,
		AddressUUID: req.AddressUUID,
		TxTime:      req.TxTime,
		ChainID:     req.ChainID,
		TokenID:     req.TokenID,
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
		Amount:      req.Amount,
		Memo:        req.Memo,
		TxType:      req.TxType,
		Status:      backend.TxStatusCreated,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	if err := s.db.BackendWalletTxRecord.StoreWalletTxRecord(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *walletTxRecordService) UpdateWalletTx(
	ctx context.Context,
	req UpdateWalletTxRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid required")
	}
	if len(req.Updates) == 0 {
		return fmt.Errorf("updates empty")
	}

	// 1. 先查旧交易
	oldTx, err := s.db.BackendWalletTxRecord.GetByGuid(req.Guid)
	if err != nil {
		return err
	}
	oldStatus := oldTx.Status

	// 2. 更新交易本身
	req.Updates["updated_at"] = time.Now()
	if err := s.db.BackendWalletTxRecord.UpdateWalletTxRecord(
		req.Guid,
		req.Updates,
	); err != nil {
		return err
	}

	// 3. 判断：是否是「第一次成功」
	newStatus, ok := req.Updates["status"].(int)
	if !ok {
		return nil // 本次更新没涉及 status，直接结束
	}

	if oldStatus == backend.TxStatusSuccess ||
		newStatus != backend.TxStatusSuccess {
		return nil // 不是第一次成功，不动余额
	}

	// 4️. 交易成功 → 更新钱包资产
	return s.applyTxToWalletAsset(ctx, oldTx)
}

func (s *walletTxRecordService) GetWalletTx(
	ctx context.Context,
	guid string,
) (*backend.WalletTxRecord, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	return s.db.BackendWalletTxRecord.GetByGuid(guid)
}

func (s *walletTxRecordService) GetByOperationID(
	ctx context.Context,
	operationID string,
) ([]*backend.WalletTxRecord, error) {

	if operationID == "" {
		return nil, fmt.Errorf("operation_id required")
	}
	return s.db.BackendWalletTxRecord.GetByOperationID(operationID)
}

func (s *walletTxRecordService) applyTxToWalletAsset(
	ctx context.Context,
	tx *backend.WalletTxRecord,
) error {

	// 1. 查钱包资产（token + chain）
	asset, err := s.db.BackendWalletAsset.GetByTokenChain(
		tx.TokenID,
		tx.ChainID,
	)
	if err != nil {
		return err
	}

	// 2. 用 big.Int 算余额（因为 uint256）
	balance := new(big.Int)
	balance.SetString(asset.Balance, 10)

	amount := new(big.Int)
	amount.SetString(tx.Amount, 10)

	// 3. 判断方向：转出 or 转入
	if tx.FromAddress != "" {
		// 从这个钱包转出
		balance.Sub(balance, amount)
	} else {
		// 转入
		balance.Add(balance, amount)
	}

	// 4. 写回 wallet_asset
	return s.db.BackendWalletAsset.UpdateWalletAsset(
		asset.Guid,
		map[string]interface{}{
			"balance": balance.String(),
		},
	)
}
