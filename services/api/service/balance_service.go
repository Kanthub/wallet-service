package service

import (
	"context"
	"fmt"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type WalletBalanceService interface {
	// 查询钱包所有资产
	GetWalletBalances(ctx context.Context, walletUUID string) ([]*backend.WalletAsset, error)

	// 查询钱包指定 token + chain
	GetWalletBalanceByTokenChain(
		ctx context.Context,
		walletUUID string,
		tokenID string,
		chainID string,
	) (*backend.WalletAsset, error)

	// 查询钱包资产汇总
	GetWalletBalanceSummary(ctx context.Context, walletUUID string) (*WalletBalanceSummary, error)
}

type WalletBalanceSummary struct {
	WalletUUID string `json:"wallet_uuid"`
	TotalUSD   string `json:"total_usd"`
	TotalUSDT  string `json:"total_usdt"`
}

type walletBalanceService struct {
	db *database.DB
}

func NewWalletBalanceService(db *database.DB) WalletBalanceService {
	return &walletBalanceService{db: db}
}

func (s *walletBalanceService) GetWalletBalances(
	ctx context.Context,
	walletUUID string,
) ([]*backend.WalletAsset, error) {

	if walletUUID == "" {
		return nil, fmt.Errorf("wallet_uuid required")
	}

	return s.db.BackendWalletAsset.GetByWalletUUID(walletUUID)
}

func (s *walletBalanceService) GetWalletBalanceByTokenChain(
	ctx context.Context,
	walletUUID string,
	tokenID string,
	chainID string,
) (*backend.WalletAsset, error) {

	if walletUUID == "" || tokenID == "" || chainID == "" {
		return nil, fmt.Errorf("wallet_uuid, token_id and chain_id required")
	}

	return s.db.BackendWalletAsset.GetByWalletTokenChain(
		walletUUID,
		tokenID,
		chainID,
	)
}

func (s *walletBalanceService) GetWalletBalanceSummary(
	ctx context.Context,
	walletUUID string,
) (*WalletBalanceSummary, error) {

	if walletUUID == "" {
		return nil, fmt.Errorf("wallet_uuid required")
	}

	list, err := s.db.BackendWalletAsset.GetByWalletUUID(walletUUID)
	if err != nil {
		return nil, err
	}

	var totalUSD float64
	var totalUSDT float64

	for _, a := range list {
		// 简单版本：用 float 计算展示值（这是 OK 的）
		if a.AssetUsd != "" {
			var v float64
			fmt.Sscan(a.AssetUsd, &v)
			totalUSD += v
		}
		if a.AssetUsdt != "" {
			var v float64
			fmt.Sscan(a.AssetUsdt, &v)
			totalUSDT += v
		}
	}

	return &WalletBalanceSummary{
		WalletUUID: walletUUID,
		TotalUSD:   fmt.Sprintf("%.8f", totalUSD),
		TotalUSDT:  fmt.Sprintf("%.8f", totalUSDT),
	}, nil
}
