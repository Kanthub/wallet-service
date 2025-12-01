package driver

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/bindings"
	"github.com/roothash-pay/wallet-services/metrics"
	"github.com/roothash-pay/wallet-services/relayer/txmgr"
)

var (
	errMaxPriorityFeePerGasNotFound = errors.New("Method eth_maxPriorityFeePerGas not found")
	FallbackGasTipCap               = big.NewInt(1500000000)
)

type DriverEngineConfig struct {
	ChainClient               *ethclient.Client
	ChainId                   *big.Int
	RrmManagerAddress         common.Address
	CallerAddress             common.Address
	PrivateKey                *ecdsa.PrivateKey
	NumConfirmations          uint64
	SafeAbortNonceTooLowCount uint64
}

type DriverEngine struct {
	Ctx                    context.Context
	Cfg                    *DriverEngineConfig
	RrmManagerContract     *bindings.ReferralRewardManager
	RawPoolManagerContract *bind.BoundContract
	RrmManagerContractAbi  *abi.ABI
	TxMgr                  txmgr.TxManager
	PhoenixMetrics         *metrics.PhoenixMetrics
	cancel                 func()
}

func NewDriverEngine(ctx context.Context, phoenixMetrics *metrics.PhoenixMetrics, cfg *DriverEngineConfig) (*DriverEngine, error) {
	_, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	rrmManagerContract, err := bindings.NewReferralRewardManager(cfg.RrmManagerAddress, cfg.ChainClient)
	if err != nil {
		log.Error("new rrm manager fail", "err", err)
		return nil, err
	}

	parsed, err := abi.JSON(strings.NewReader(bindings.ReferralRewardManagerMetaData.ABI))
	if err != nil {
		log.Error("parsed abi fail", "err", err)
		return nil, err
	}

	rrmManagerContractAbi, err := bindings.ReferralRewardManagerMetaData.GetAbi()
	if err != nil {
		log.Error("get rrm manager meta data fail", "err", err)
		return nil, err
	}

	rawPoolManagerContract := bind.NewBoundContract(cfg.RrmManagerAddress, parsed, cfg.ChainClient, cfg.ChainClient, cfg.ChainClient)

	txManagerConfig := txmgr.Config{
		ResubmissionTimeout:       time.Second * 5,
		ReceiptQueryInterval:      time.Second,
		NumConfirmations:          cfg.NumConfirmations,
		SafeAbortNonceTooLowCount: cfg.SafeAbortNonceTooLowCount,
	}

	txManager := txmgr.NewSimpleTxManager(txManagerConfig, cfg.ChainClient)

	return &DriverEngine{
		Ctx:                    ctx,
		Cfg:                    cfg,
		RrmManagerContract:     rrmManagerContract,
		RawPoolManagerContract: rawPoolManagerContract,
		RrmManagerContractAbi:  rrmManagerContractAbi,
		TxMgr:                  txManager,
		PhoenixMetrics:         phoenixMetrics,
		cancel:                 cancel,
	}, nil
}

func (de *DriverEngine) UpdateGasPrice(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	var opts *bind.TransactOpts
	var err error
	opts, err = bind.NewKeyedTransactorWithChainID(de.Cfg.PrivateKey, de.Cfg.ChainId)
	if err != nil {
		log.Error("new keyed transactor with chain id fail", "err", err)
		return nil, err
	}
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(tx.Nonce())
	opts.NoSend = true

	finalTx, err := de.RawPoolManagerContract.RawTransact(opts, tx.Data())
	switch {
	case err == nil:
		return finalTx, nil

	case de.isMaxPriorityFeePerGasNotFoundError(err):
		log.Info("Don't support priority fee")
		opts.GasTipCap = FallbackGasTipCap
		return de.RawPoolManagerContract.RawTransact(opts, tx.Data())

	default:
		return nil, err
	}
}

func (de *DriverEngine) isMaxPriorityFeePerGasNotFoundError(err error) bool {
	return strings.Contains(err.Error(), errMaxPriorityFeePerGasNotFound.Error())
}

func (de *DriverEngine) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return de.Cfg.ChainClient.SendTransaction(ctx, tx)
}

func (de *DriverEngine) batchAddRewards(ctx context.Context, userAddresses []common.Address, amount []*big.Int) (*types.Transaction, error) {
	nonce, err := de.Cfg.ChainClient.NonceAt(ctx, de.Cfg.CallerAddress, nil)
	if err != nil {
		log.Error("get nonce error", "err", err)
		return nil, err
	}
	de.PhoenixMetrics.RecordChainAddressNonce(de.Cfg.ChainId.String(), big.NewInt(int64(nonce)))
	log.Info("bridge finalized eth tx nonce", "CallerAddress", de.Cfg.CallerAddress, "nonce", nonce)

	balance, err := de.Cfg.ChainClient.BalanceAt(ctx, de.Cfg.CallerAddress, nil)
	if err != nil {
		log.Error("get nonce error", "err", err)
		return nil, err
	}
	de.PhoenixMetrics.RecordNativeTokenBalance(de.Cfg.ChainId.String(), balance)
	log.Info("bridge eth balance", "CallerAddress", de.Cfg.CallerAddress, "balance", balance.String())

	opts, err := bind.NewKeyedTransactorWithChainID(de.Cfg.PrivateKey, de.Cfg.ChainId)
	if err != nil {
		log.Error("new keyed transactor with chain id fail", "err", err)
		return nil, err
	}
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(nonce)
	opts.NoSend = true

	tx, err := de.RrmManagerContract.BatchAddRewards(opts, userAddresses, amount)
	switch {
	case err == nil:
		return tx, nil

	case de.isMaxPriorityFeePerGasNotFoundError(err):
		opts.GasTipCap = FallbackGasTipCap
		return de.RrmManagerContract.BatchAddRewards(opts, userAddresses, amount)

	default:
		return nil, err
	}
}

func (de *DriverEngine) BatchAddRewards(userAddresses []common.Address, amount []*big.Int) (*types.Receipt, error) {
	tx, err := de.batchAddRewards(de.Ctx, userAddresses, amount)
	if err != nil {
		log.Error("build bridge finalized tx fail", "err", err)
		return nil, err
	}
	updateGasPrice := func(ctx context.Context) (*types.Transaction, error) {
		return de.UpdateGasPrice(ctx, tx)
	}

	log.Info("build bridge finalized eth tx", "txHash", tx.Hash().String(), "data", tx.Data())

	receipt, err := de.TxMgr.Send(de.Ctx, updateGasPrice, de.SendTransaction)
	if err != nil {
		log.Error("send tx fail", "err", err)
		return nil, err
	}
	return receipt, nil
}
