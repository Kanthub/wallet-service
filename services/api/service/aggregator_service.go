package service

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/roothash-pay/wallet-services/common/redis"
	"github.com/roothash-pay/wallet-services/config"
	"github.com/roothash-pay/wallet-services/database"
	dbBackend "github.com/roothash-pay/wallet-services/database/backend"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/provider"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/provider/jupiter"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/provider/lifi"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/provider/oneinch"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/provider/zerox"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/store"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/utils"
	"github.com/roothash-pay/wallet-services/services/api/models/backend"

	"github.com/roothash-pay/wallet-services/services/common/chaininfo"
	"github.com/roothash-pay/wallet-services/services/grpc_client/account"
)

// AggregatorService handles swap aggregation operations
type AggregatorService struct {
	providers     []provider.Provider
	quoteStore    store.QuoteStore
	swapStore     store.SwapStore
	validator     *utils.Validator
	accountClient *account.WalletAccountClient
	db            *database.DB
	quoteTTL      time.Duration
	chainInfo     chaininfo.Provider
}

// initAggregatorService initializes the aggregator service with all dependencies
func InitAggregatorService(db *database.DB, cfg *config.Config) (*AggregatorService, error) {
	// Skip initialization if wallet account address is not configured
	if cfg.AggregatorConfig.WalletAccountAddr == "" {
		log.Warn("Aggregator service not initialized: wallet_account_addr not configured")
		return nil, nil
	}

	// Create wallet account client
	accountClient, err := account.NewWalletAccountClient(cfg.AggregatorConfig.WalletAccountAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet account client: %w", err)
	}

	// Create providers
	var providers []provider.Provider

	// Initialize 0x provider if enabled
	if cfg.AggregatorConfig.EnableProviders["0x"] && cfg.AggregatorConfig.ZeroXAPIURL != "" {
		zeroXProvider := zerox.NewProvider(cfg.AggregatorConfig.ZeroXAPIURL, cfg.AggregatorConfig.ZeroXAPIKey)
		providers = append(providers, zeroXProvider)
		log.Info("0x provider initialized", "url", cfg.AggregatorConfig.ZeroXAPIURL)
	}

	// Initialize 1inch provider if enabled
	if cfg.AggregatorConfig.EnableProviders["1inch"] && cfg.AggregatorConfig.OneInchAPIURL != "" {
		oneInchProvider := oneinch.NewProvider(cfg.AggregatorConfig.OneInchAPIURL, cfg.AggregatorConfig.OneInchAPIKey)
		providers = append(providers, oneInchProvider)
		log.Info("1inch provider initialized", "url", cfg.AggregatorConfig.OneInchAPIURL)
	}

	// Initialize Jupiter provider if enabled
	if cfg.AggregatorConfig.EnableProviders["jupiter"] && cfg.AggregatorConfig.JupiterAPIURL != "" {
		jupiterProvider := jupiter.NewProvider(cfg.AggregatorConfig.JupiterAPIURL)
		providers = append(providers, jupiterProvider)
		log.Info("Jupiter provider initialized", "url", cfg.AggregatorConfig.JupiterAPIURL)
	}

	// Create Redis client
	var redisClient *redis.Client
	if cfg.RedisConfig.Addr != "" {
		var err error
		redisClient, err = redis.NewClient(&cfg.RedisConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create Redis client: %w", err)
		}
		log.Info("Redis client initialized", "addr", cfg.RedisConfig.Addr)
	} else {
		log.Warn("Redis not configured, using in-memory storage (not recommended for production)")
	}

	// Initialize chain metadata cache
	chainInfoManager := chaininfo.NewManager(
		db.BackendChain,
		redisClient,
		cfg.AggregatorConfig.WalletAccountConsumerToken,
		cfg.AggregatorConfig.ChainConsumerTokens,
	)
	if err := chainInfoManager.WarmUp(context.Background()); err != nil {
		log.Warn("Failed to warm up chain info cache", "err", err)
	}

	// Create EVM caller for contract interactions
	evmCaller := utils.NewEVMCaller(accountClient, chainInfoManager)

	// Initialize LiFi provider if enabled
	if cfg.AggregatorConfig.EnableProviders["lifi"] && cfg.AggregatorConfig.LiFiAPIURL != "" {
		lifiProvider := lifi.NewProvider(cfg.AggregatorConfig.LiFiAPIURL, cfg.AggregatorConfig.LiFiAPIKey, evmCaller)
		providers = append(providers, lifiProvider)
		log.Info("LiFi provider initialized", "url", cfg.AggregatorConfig.LiFiAPIURL)
	}

	if len(providers) == 0 {
		log.Warn("Aggregator service not initialized: no providers enabled")
		return nil, nil
	}

	// Create cache stores
	var quoteStore store.QuoteStore
	var swapStore store.SwapStore
	if redisClient != nil {
		quoteStore = store.NewRedisQuoteStore(redisClient.Client)
		swapStore = store.NewRedisSwapStore(redisClient.Client)
		log.Info("Using Redis-based storage")
	} else {
		quoteStore = store.NewInMemoryQuoteStore()
		swapStore = store.NewInMemorySwapStore()
		log.Warn("Cannot connect to Redis, using in-memory storage (data will be lost on restart)")
	}

	// Create validator
	validator := utils.NewValidator()

	// Create aggregator service
	aggregatorService := NewAggregatorService(
		providers,
		quoteStore,
		swapStore,
		validator,
		accountClient,
		chainInfoManager,
		db,
	)

	log.Info("Aggregator service initialized successfully", "providers", len(providers))
	return aggregatorService, nil
}

// NewAggregatorService creates a new aggregator service
func NewAggregatorService(
	providers []provider.Provider,
	quoteStore store.QuoteStore,
	swapStore store.SwapStore,
	validator *utils.Validator,
	accountClient *account.WalletAccountClient,
	chainInfo chaininfo.Provider,
	db *database.DB,
) *AggregatorService {
	return &AggregatorService{
		providers:     providers,
		quoteStore:    quoteStore,
		swapStore:     swapStore,
		validator:     validator,
		accountClient: accountClient,
		db:            db,
		quoteTTL:      1 * time.Minute, // Default 1 minutes
		chainInfo:     chainInfo,
	}
}

func (s *AggregatorService) getChainInfo(ctx context.Context, chainID string) (*chaininfo.Info, error) {
	if s.chainInfo == nil {
		return nil, fmt.Errorf("chain info provider not configured")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return s.chainInfo.Get(ctx, chainID)
}

// GetQuotes aggregates quotes from multiple providers
func (s *AggregatorService) GetQuotes(ctx context.Context, req *backend.QuoteRequest) (*backend.QuoteResponse, error) {
	// TODO: 限流
	// Validate chain ID
	if err := s.validator.ValidateChainID(req.FromChainID); err != nil {
		return nil, err
	}

	// Fetch quotes from all providers concurrently
	quotes, err := s.aggregateQuotes(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("no quotes available")
	}

	// Sort quotes by toAmount (descending)
	// top quote is the best quote
	sort.Slice(quotes, func(i, j int) bool {
		amountI, _ := new(big.Int).SetString(quotes[i].ToAmount, 10)
		amountJ, _ := new(big.Int).SetString(quotes[j].ToAmount, 10)
		return amountI.Cmp(amountJ) > 0
	})

	// Prepare response
	quoteID := uuid.New().String()
	expiresAt := time.Now().Add(s.quoteTTL)

	storeData := &backend.QuoteStore{
		QuoteID:     quoteID,
		UserAddress: req.UserAddress,
		ExpiresAt:   expiresAt,
		WalletUUID:  req.WalletUUID,
		BestQuotes:  quotes,
	}

	//  Cache store quote snapshot
	if err = s.quoteStore.Save(ctx, quoteID, storeData, s.quoteTTL); err != nil {
		log.Error("Failed to cache quote", "err", err)
	}

	response := &backend.QuoteResponse{
		QuoteID:     quoteID,
		UserAddress: req.UserAddress,
		ExpiresAt:   expiresAt,
		WalletUUID:  req.WalletUUID,
		BestQuotes:  quotes,
	}

	for _, v := range response.BestQuotes {
		v.Raw = ""
	}

	return response, nil
}

// aggregateQuotes fetches quotes from all providers concurrently
func (s *AggregatorService) aggregateQuotes(ctx context.Context, req *backend.QuoteRequest) ([]*backend.Quote, error) {
	g, ctx := errgroup.WithContext(ctx)
	quoteChan := make(chan *backend.Quote, len(s.providers))

	// TODO: Filter providers based on chain type，比如evm 排除 solana的provider
	for _, provider := range s.providers {
		p := provider // Capture loop variable
		g.Go(func() error {
			quote, err := p.GetQuote(ctx, req)
			if err != nil {
				log.Warn("Provider failed", "provider", p.Name(), "err", err)
				return nil // Don't fail entire aggregation
			}
			quoteChan <- quote
			return nil
		})
	}

	// Wait for all providers
	_ = g.Wait()
	close(quoteChan)

	// Collect successful quotes
	var quotes []*backend.Quote
	for quote := range quoteChan {
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

// 前端获取报价后，用户接受该报价点击 swap，执行该方法
// 返回一个动作链路，让用户执行相关签名
// PrepareSwap generates a transaction plan for a swap
func (s *AggregatorService) PrepareSwap(ctx context.Context, quoteID string, bestQuotesIndex int) (*backend.PrepareSwapResponse, error) {
	// Validate quote
	cachedQuote, err := s.quoteStore.Get(ctx, quoteID)
	if err != nil {
		return nil, fmt.Errorf("quote not found or expired: %w", err)
	}
	if time.Now().After(cachedQuote.ExpiresAt) {
		return nil, fmt.Errorf("quote expired")
	}

	// 更新缓存
	cachedQuote.BestQuotesIndex = bestQuotesIndex
	err = s.quoteStore.Update(ctx, quoteID, cachedQuote, s.quoteTTL)
	if err != nil {
		return nil, fmt.Errorf("fail to update cache")
	}

	quote := cachedQuote.BestQuotes[bestQuotesIndex]

	// Generate swap ID
	swapID := uuid.New().String()

	// Find provider
	var selectedProvider provider.Provider
	for _, p := range s.providers {
		if p.Name() == quote.Provider {
			selectedProvider = p
			break
		}
	}
	if selectedProvider == nil {
		return nil, fmt.Errorf("provider not found: %s", quote.Provider)
	}

	// build swap tx
	buildResp, err := selectedProvider.BuildSwap(ctx, quote, cachedQuote.UserAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to build swap: %w", err)
	}

	actions := buildResp.Actions

	// Create swap record
	swap := &backend.Swap{
		SwapID:      swapID,
		QuoteID:     cachedQuote.QuoteID,
		UserAddress: cachedQuote.UserAddress,
		WalletUUID:  cachedQuote.WalletUUID,
		Status:      backend.TxStatusCreated, // 0 = CREATED
		Steps:       make([]*backend.Step, len(actions)),
	}

	// Initialize steps
	for i, action := range actions {
		swap.Steps[i] = &backend.Step{
			StepIndex:  i,
			ActionType: action.ActionType,
			Status:     backend.TxStatusCreated, // 0 = CREATED
		}

		if err = fillExpectedFromSigningPayload(swap.Steps[i], action.SigningPayload); err != nil {
			return nil, fmt.Errorf("failed to build expected tx snapshot (step %d): %w", i, err)
		}
	}

	if err = s.swapStore.CreateSwap(ctx, swap); err != nil {
		return nil, err
	}

	return &backend.PrepareSwapResponse{
		SwapID:  swapID,
		Actions: actions,
	}, nil
}

// SubmitSignedTx broadcasts a signed transaction
func (s *AggregatorService) SubmitSignedTx(ctx context.Context, req *backend.SubmitSignedTxRequest) (*backend.SubmitSignedTxResponse, error) {
	// Check idempotency
	if txHash, exists := s.swapStore.CheckIdempotency(ctx, req.SwapID, req.StepIndex, req.IdempotencyKey); exists {
		log.Info("Duplicate request detected", "swapID", req.SwapID, "stepIndex", req.StepIndex, "txHash", txHash)
		return &backend.SubmitSignedTxResponse{TxHash: txHash}, nil
	}

	// Get swap
	swap, err := s.swapStore.GetSwap(ctx, req.SwapID)
	if err != nil {
		return nil, err
	}

	// Validate step index
	if req.StepIndex < 0 || req.StepIndex >= len(swap.Steps) {
		return nil, fmt.Errorf("invalid step index: %d", req.StepIndex)
	}

	step := swap.Steps[req.StepIndex]

	// Get quote for chain info
	quoteResp, err := s.quoteStore.Get(ctx, swap.QuoteID)
	if err != nil {
		return nil, fmt.Errorf("quote not found: %w", err)
	}
	quote := quoteResp.BestQuotes[quoteResp.BestQuotesIndex]
	chainInfo, err := s.getChainInfo(ctx, quote.ChainID)
	if err != nil {
		return nil, err
	}

	// 防止滥用接口，只广播来自 prepare 的交易
	if err = validateSignedTxAgainstStepExpected(req.SignedTx, step); err != nil {
		return nil, fmt.Errorf("signed tx validation failed: %w", err)
	}

	// 1: Save to database with CREATED status (before broadcast)
	// This ensures we have a record even if broadcast fails
	recordGuid := s.saveStepTxStatusCreated(ctx, swap, quote, req.StepIndex)

	// 2: Broadcast transaction using SendTx
	result, err := s.accountClient.SendTx(ctx, account.SendTxParams{
		ConsumerToken: chainInfo.ConsumerToken,
		Chain:         chainInfo.WalletChain,
		Coin:          chainInfo.WalletCoin,
		Network:       chainInfo.WalletNetwork,
		RawTx:         req.SignedTx,
	})
	if err != nil {
		// Update step as failed
		step.Status = backend.TxStatusFailed // 2 = FAILED
		step.FailReasonCode = dbBackend.FailReasonBroadcastFailed
		step.FailMessage = err.Error()
		_ = s.swapStore.UpdateStep(ctx, req.SwapID, req.StepIndex, step)

		// Update database record to FAILED
		if recordGuid != "" {
			s.updateStepTxStatusFailed(ctx, recordGuid, dbBackend.FailReasonBroadcastFailed, err.Error())
		}

		return nil, err
	}

	txHash := result.TxHash

	// Update step
	now := time.Now()
	step.TxHash = txHash
	step.Status = backend.TxStatusPending // 1 = PENDING
	step.SubmittedAt = &now
	step.IdempotencyKey = req.IdempotencyKey

	if err := s.swapStore.UpdateStep(ctx, req.SwapID, req.StepIndex, step); err != nil {
		log.Error("Failed to update step", "err", err)
	}

	// Record idempotency
	_ = s.swapStore.RecordIdempotency(ctx, req.SwapID, req.StepIndex, req.IdempotencyKey, txHash)

	// Update swap status
	swap.Status = backend.TxStatusPending // 1 = PENDING
	_ = s.swapStore.UpdateSwap(ctx, swap)

	// 3: Update database record to PENDING (after successful broadcast)
	if recordGuid != "" {
		s.updateStepTxStatusPending(ctx, recordGuid, quote.ChainID, txHash)
	}

	return &backend.SubmitSignedTxResponse{TxHash: txHash}, nil
}

func normalizeValue(value string) (string, error) {
	if value == "" {
		return "0", nil
	}
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		v, err := hexutil.DecodeBig(value)
		if err != nil {
			return "", err
		}
		return v.String(), nil
	}
	// assume decimal string
	v := new(big.Int)
	if _, ok := v.SetString(value, 10); !ok {
		return "", fmt.Errorf("invalid value: %s", value)
	}
	return v.String(), nil
}

func fillExpectedFromSigningPayload(step *backend.Step, sp *backend.SigningPayload) error {
	if sp == nil {
		return nil
	}
	// 只对 EVM tx 做 expected（Solana 走 SerializedTx 另外处理）
	if sp.To == "" || sp.Data == "" || sp.ChainID == "" {
		return nil
	}

	dataBytes, err := hexutil.Decode(sp.Data)
	if err != nil {
		return fmt.Errorf("invalid signing payload data: %w", err)
	}

	value, err := normalizeValue(sp.Value)
	if err != nil {
		return fmt.Errorf("invalid signing payload value: %w", err)
	}

	step.ExpectedChainID = sp.ChainID
	step.ExpectedTo = strings.ToLower(sp.To)
	step.ExpectedValueWei = value
	step.ExpectedDataHash = crypto.Keccak256Hash(dataBytes).Hex()
	return nil
}

func validateSignedTxAgainstStepExpected(signedTxHex string, step *backend.Step) error {
	if step.ExpectedTo == "" || step.ExpectedDataHash == "" || step.ExpectedChainID == "" {
		return fmt.Errorf("missing expected tx snapshot in step")
	}

	rawBytes, err := hexutil.Decode(signedTxHex)
	if err != nil {
		return fmt.Errorf("invalid signedTx hex: %w", err)
	}

	var tx types.Transaction
	if err = tx.UnmarshalBinary(rawBytes); err != nil {
		return fmt.Errorf("failed to decode signed tx: %w", err)
	}

	// chainId
	expChainID, ok := new(big.Int).SetString(step.ExpectedChainID, 10)
	if !ok {
		return fmt.Errorf("invalid expected chainId: %s", step.ExpectedChainID)
	}
	if tx.ChainId() == nil || tx.ChainId().Cmp(expChainID) != 0 {
		return fmt.Errorf("chainId mismatch: got %v want %v", tx.ChainId(), expChainID)
	}

	// to
	if tx.To() == nil {
		return fmt.Errorf("contract creation not allowed")
	}
	if strings.ToLower(tx.To().Hex()) != strings.ToLower(step.ExpectedTo) {
		return fmt.Errorf("to mismatch: got %s want %s", tx.To().Hex(), step.ExpectedTo)
	}

	// value（你需要统一 step.ExpectedValueWei 的格式）
	expValue, ok := new(big.Int).SetString(step.ExpectedValueWei, 10)
	if !ok {
		return fmt.Errorf("invalid expected value: %s", step.ExpectedValueWei)
	}
	if tx.Value() == nil || tx.Value().Cmp(expValue) != 0 {
		return fmt.Errorf("value mismatch: got %s want %s", tx.Value(), expValue)
	}

	// data hash
	gotHash := crypto.Keccak256Hash(tx.Data()).Hex()
	if strings.ToLower(gotHash) != strings.ToLower(step.ExpectedDataHash) {
		return fmt.Errorf("data mismatch: got %s want %s", gotHash, step.ExpectedDataHash)
	}

	return nil
}

func (s *AggregatorService) SubmitTxHash(ctx context.Context, req *backend.SubmitTxHashRequest) (*backend.SubmitTxHashResponse, error) {
	// 1) 幂等
	if txHash, exists := s.swapStore.CheckIdempotency(ctx, req.SwapID, req.StepIndex, req.IdempotencyKey); exists {
		return &backend.SubmitTxHashResponse{TxHash: txHash}, nil
	}

	// 2) 获取 swap & step
	swap, err := s.swapStore.GetSwap(ctx, req.SwapID)
	if err != nil {
		return nil, err
	}
	if req.StepIndex < 0 || req.StepIndex >= len(swap.Steps) {
		return nil, fmt.Errorf("invalid step index: %d", req.StepIndex)
	}
	step := swap.Steps[req.StepIndex]

	// 3) 可选但强烈建议：链上回查交易内容，校验符合 step.Expected*
	//    （防止有人随便塞个 txHash 进来污染状态/对账）
	// 需要 chainInfo + accountClient.GetTxByHash 能返回 to/value/data 或至少 dataHash
	if err = s.validateTxHashAgainstStep(ctx, swap, step, req.TxHash); err != nil {
		return nil, fmt.Errorf("tx hash validation failed: %w", err)
	}

	// 4) 更新 step
	now := time.Now()
	step.TxHash = req.TxHash
	step.Status = backend.TxStatusPending
	step.SubmittedAt = &now
	step.IdempotencyKey = req.IdempotencyKey
	_ = s.swapStore.UpdateStep(ctx, req.SwapID, req.StepIndex, step)

	// 5) 更新 swap
	swap.Status = backend.TxStatusPending
	_ = s.swapStore.UpdateSwap(ctx, swap)

	// 6) 记录幂等
	_ = s.swapStore.RecordIdempotency(ctx, req.SwapID, req.StepIndex, req.IdempotencyKey, req.TxHash)

	// 7) 可选：写 wallet_tx_record（CREATED/PENDING）让 worker 接管后续状态
	// recordGuid := s.saveStepTxStatusCreated(...)
	// s.updateStepTxStatusPending(..., req.TxHash)

	return &backend.SubmitTxHashResponse{TxHash: req.TxHash}, nil
}

// validateTxHashAgainstStep: 用 txHash 回查链上交易并校验它属于该 step（在没有 input/data 的情况下是“半严格校验”）
func (s *AggregatorService) validateTxHashAgainstStep(
	ctx context.Context,
	swap *backend.Swap,
	step *backend.Step,
	txHash string,
) error {
	quoteResp, err := s.quoteStore.Get(ctx, swap.QuoteID)
	if err != nil {
		return fmt.Errorf("quote not found: %w", err)
	}
	quote := quoteResp.BestQuotes[quoteResp.BestQuotesIndex]

	if swap == nil || step == nil || quote == nil {
		return fmt.Errorf("nil swap/step/quote")
	}
	if txHash == "" {
		return fmt.Errorf("empty txHash")
	}

	// Step expected 必须存在（至少 to/value/chain）
	if step.ExpectedTo == "" || step.ExpectedChainID == "" || step.ExpectedValueWei == "" {
		return fmt.Errorf("missing expected snapshot in step")
	}

	// 1) 先检查链是否一致（我们是用 quote.ChainID 去查链的）
	//    这里的意义：防止 step.expectedChainId 写错 / quote 被串
	if quote.ChainID != step.ExpectedChainID {
		return fmt.Errorf("chainId mismatch: quote %s vs step expected %s", quote.ChainID, step.ExpectedChainID)
	}

	chainInfo, err := s.getChainInfo(ctx, quote.ChainID)
	if err != nil {
		return fmt.Errorf("chain info not found: %w", err)
	}

	txInfo, err := s.accountClient.GetTxByHash(
		ctx,
		chainInfo.ConsumerToken,
		chainInfo.WalletChain,
		chainInfo.WalletCoin,
		chainInfo.WalletNetwork,
		txHash,
	)
	if err != nil {
		return fmt.Errorf("failed to get tx by hash: %w", err)
	}
	if txInfo == nil {
		return fmt.Errorf("tx not found: %s", txHash)
	}

	// 2) from 必须是 swap.UserAddress（关键：防别人塞入任意 txHash 污染你的 swap）
	if swap.UserAddress != "" && !addrEq(txInfo.From, swap.UserAddress) {
		return fmt.Errorf("from mismatch: got %s want %s", txInfo.From, swap.UserAddress)
	}

	// 3) to 必须匹配 expected
	if !addrEq(txInfo.To, step.ExpectedTo) {
		return fmt.Errorf("to mismatch: got %s want %s", txInfo.To, step.ExpectedTo)
	}

	// 4) value 必须匹配 expected（统一成 wei 十进制）
	gotValueWei, err := normalizeWeiString(txInfo.Value) // 允许 "0x0" 或 "0"
	if err != nil {
		return fmt.Errorf("invalid tx value: %w", err)
	}
	expValueWei, ok := new(big.Int).SetString(step.ExpectedValueWei, 10)
	if !ok {
		return fmt.Errorf("invalid expected value in step: %s", step.ExpectedValueWei)
	}
	if gotValueWei.Cmp(expValueWei) != 0 {
		return fmt.Errorf("value mismatch: got %s want %s", gotValueWei.String(), expValueWei.String())
	}

	// 5) dataHash：当前 txInfo 没有 input/data，无法校验
	//    你可以在这里选择：
	//    - 直接放行（半严格）
	//    - 或者如果 step.ExpectedDataHash 不为空就拒绝（强制要求服务返回 input）
	if step.ExpectedDataHash != "" {
		return fmt.Errorf("cannot validate tx input/data hash: account GetTxByHash does not return calldata")
	}

	return nil
}

func addrEq(a, b string) bool {
	return strings.ToLower(strings.TrimSpace(a)) == strings.ToLower(strings.TrimSpace(b))
}

func normalizeWeiString(v string) (*big.Int, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return big.NewInt(0), nil
	}
	if strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X") {
		return hexutil.DecodeBig(v)
	}
	bi := new(big.Int)
	if _, ok := bi.SetString(v, 10); !ok {
		return nil, fmt.Errorf("invalid decimal: %s", v)
	}
	return bi, nil
}

// GetSwapStatus retrieves the status of a swap
// 主动更新状态
func (s *AggregatorService) GetSwapStatus(ctx context.Context, swapID string) (*backend.SwapStatusResponse, error) {
	swap, err := s.swapStore.GetSwap(ctx, swapID)
	if err != nil {
		return nil, err
	}

	var statusChainInfo *chaininfo.Info
	if quoteResp, err := s.quoteStore.Get(ctx, swap.QuoteID); err == nil && quoteResp != nil && quoteResp.BestQuotes[0] != nil {
		if info, err := s.getChainInfo(ctx, quoteResp.BestQuotes[0].ChainID); err == nil {
			statusChainInfo = info
		} else {
			log.Warn("Chain settings missing for swap, skip status refresh", "swapID", swapID, "chainID", quoteResp.BestQuotes[0].ChainID, "err", err)
		}
	} else if err != nil {
		log.Warn("Failed to load quote for swap status refresh", "swapID", swapID, "err", err)
	}

	// Query status for each submitted transaction
	for _, step := range swap.Steps {
		if step.TxHash != "" && step.Status == backend.TxStatusPending { // 1 = PENDING
			if statusChainInfo == nil {
				continue
			}

			txInfo, err := s.accountClient.GetTxByHash(
				ctx,
				statusChainInfo.ConsumerToken,
				statusChainInfo.WalletChain,
				statusChainInfo.WalletCoin,
				statusChainInfo.WalletNetwork,
				step.TxHash,
			)
			if err != nil {
				log.Warn("Failed to get tx status", "txHash", step.TxHash, "err", err)
				continue
			}

			// Update step based on status
			if txInfo.Status == 3 { // TxStatus_Success
				now := time.Now()
				step.Status = backend.TxStatusSuccess // 3 = SUCCESS
				step.ConfirmedAt = &now
				_ = s.swapStore.UpdateStep(ctx, swapID, step.StepIndex, step)
			} else if txInfo.Status == 2 { // TxStatus_Failed
				step.Status = backend.TxStatusFailed // 2 = FAILED
				step.FailReasonCode = "TX_FAILED"
				step.FailMessage = "Transaction failed on chain"
				_ = s.swapStore.UpdateStep(ctx, swapID, step.StepIndex, step)
			}
		}
	}

	// Determine overall swap status
	allSuccess := true
	anyFailed := false
	for _, step := range swap.Steps {
		if step.Status != backend.TxStatusSuccess { // 3 = SUCCESS
			allSuccess = false
		}
		if step.Status == backend.TxStatusFailed { // 2 = FAILED
			anyFailed = true
		}
	}

	previousStatus := swap.Status

	if anyFailed {
		swap.Status = backend.TxStatusFailed // 2 = FAILED
	} else if allSuccess {
		swap.Status = backend.TxStatusSuccess // 3 = SUCCESS
	}

	_ = s.swapStore.UpdateSwap(ctx, swap)

	// Update database status when swap status changes
	if previousStatus != swap.Status && (swap.Status == backend.TxStatusSuccess || swap.Status == backend.TxStatusFailed) {
		s.updateSwapTxStatus(ctx, swap)
	}

	return &backend.SwapStatusResponse{
		SwapID:         swap.SwapID,
		Status:         swap.Status,
		Steps:          swap.Steps,
		FailReasonCode: swap.FailReasonCode,
		FailMessage:    swap.FailMessage,
	}, nil
}

// saveStepTxStatusCreated saves a step to database with CREATED status (before broadcast)
// Returns the record GUID for later updates
func (s *AggregatorService) saveStepTxStatusCreated(ctx context.Context, swap *backend.Swap, quote *backend.Quote, stepIndex int) string {
	// Skip if database is not available
	if s.db == nil {
		log.Warn("Database not available, skip saving step history", "swapID", swap.SwapID, "stepIndex", stepIndex)
		return ""
	}

	// Validate wallet_uuid
	if swap.WalletUUID == "" {
		log.Warn("Wallet UUID not provided, skip saving step history", "swapID", swap.SwapID, "stepIndex", stepIndex)
		return ""
	}

	recordGuid := uuid.New().String()

	// Use amount string directly (no conversion needed, supports uint256)
	amount := quote.FromAmount
	if amount == "" {
		amount = "0"
	}

	step := swap.Steps[stepIndex]
	txType := strings.ToLower(string(step.ActionType))

	// Build memo based on action type
	var memo string
	memo = fmt.Sprintf("%s via %s (Step %d)",
		txType,
		quote.Provider,
		stepIndex,
	)

	// Create transaction record with CREATED status
	tokenID := s.resolveTokenID(ctx, quote.ChainID, quote.FromToken)

	record := &dbBackend.WalletTxRecord{
		Guid:        recordGuid,      // Unique UUID for this record
		OperationID: swap.SwapID,     // Associate with swap operation
		StepIndex:   stepIndex,       // Step index in the operation
		WalletUUID:  swap.WalletUUID, // Associate with wallet
		AddressUUID: "",              // Optional: can be filled if we have address_uuid
		TxTime:      time.Now().Format(time.RFC3339),
		ChainID:     quote.ChainID,
		TokenID:     tokenID,
		FromAddress: swap.UserAddress,
		ToAddress:   quote.Router,
		Amount:      amount, // Store as string (supports uint256)
		Memo:        memo,
		TxID:        "",                        // No hash yet
		BlockHeight: "",                        // Will be filled when confirmed
		TxType:      txType,                    // Transaction type: approve, swap, bridge, wrap, unwrap
		Status:      dbBackend.TxStatusCreated, // Status: CREATED (0)
	}

	// Save to database
	if err := s.db.BackendWalletTxRecord.StoreWalletTxRecord(record); err != nil {
		log.Error("Failed to save created step history", "swapID", swap.SwapID, "stepIndex", stepIndex, "actionType", step.ActionType, "walletUUID", swap.WalletUUID, "err", err)
		return ""
	}

	log.Info("Created step history saved", "swapID", swap.SwapID, "stepIndex", stepIndex, "actionType", step.ActionType, "recordGuid", recordGuid, "walletUUID", swap.WalletUUID, "status", "CREATED")
	return recordGuid
}

// updateStepTxStatusPending updates step history to PENDING status after successful broadcast
func (s *AggregatorService) updateStepTxStatusPending(_ context.Context, recordGuid string, _ string, txHash string) {
	if s.db == nil || recordGuid == "" {
		return
	}

	// Update record
	updates := map[string]interface{}{
		"hash":   txHash,
		"status": dbBackend.TxStatusPending, // Status: PENDING (1)
	}

	if err := s.db.BackendWalletTxRecord.UpdateWalletTxRecord(recordGuid, updates); err != nil {
		log.Error("Failed to update step history to pending", "recordGuid", recordGuid, "txHash", txHash, "err", err)
	} else {
		log.Info("Step history updated to pending", "recordGuid", recordGuid, "txHash", txHash, "status", "PENDING")
	}
}

// updateStepTxStatusFailed updates step history to FAILED status
func (s *AggregatorService) updateStepTxStatusFailed(ctx context.Context, recordGuid string, failReasonCode string, failReasonMsg string) {
	if s.db == nil || recordGuid == "" {
		return
	}

	// Update record
	updates := map[string]interface{}{
		"status":           dbBackend.TxStatusFailed, // Status: FAILED (2)
		"fail_reason_code": failReasonCode,
		"fail_reason_msg":  failReasonMsg,
	}

	if err := s.db.BackendWalletTxRecord.UpdateWalletTxRecord(recordGuid, updates); err != nil {
		log.Error("Failed to update step history to failed", "recordGuid", recordGuid, "err", err)
	} else {
		log.Info("Step history updated to failed", "recordGuid", recordGuid, "status", "FAILED", "reason", failReasonCode)
	}
}

// updateSwapTxStatus updates the swap history status when swap status changes
func (s *AggregatorService) updateSwapTxStatus(ctx context.Context, swap *backend.Swap) {
	// Skip if database is not available
	if s.db == nil {
		return
	}

	// Get quote information
	quoteResp, err := s.quoteStore.Get(ctx, swap.QuoteID)
	if err != nil {
		log.Warn("Failed to get quote for status update", "swapID", swap.SwapID, "err", err)
		return
	}

	quote := quoteResp.BestQuotes[0]
	var chainInfoForSwap *chaininfo.Info
	if info, err := s.getChainInfo(ctx, quote.ChainID); err == nil {
		chainInfoForSwap = info
	} else {
		log.Warn("Chain info not found while updating final block height", "swapID", swap.SwapID, "chainID", quote.ChainID, "err", err)
	}

	// Find the final swap transaction hash
	var finalTxHash string
	var finalBlockHeight string
	for i := len(swap.Steps) - 1; i >= 0; i-- {
		if swap.Steps[i].ActionType == backend.ActionTypeSwap && swap.Steps[i].TxHash != "" {
			finalTxHash = swap.Steps[i].TxHash

			// Try to get block height from chain
			if swap.Steps[i].Status == backend.TxStatusSuccess && chainInfoForSwap != nil { // 3 = SUCCESS
				txInfo, err := s.accountClient.GetTxByHash(
					ctx,
					chainInfoForSwap.ConsumerToken,
					chainInfoForSwap.WalletChain,
					chainInfoForSwap.WalletCoin,
					chainInfoForSwap.WalletNetwork,
					finalTxHash,
				)
				if err == nil && txInfo != nil {
					finalBlockHeight = txInfo.Height
				}
			}
			break
		}
	}

	if finalTxHash == "" {
		log.Warn("No swap transaction found in steps", "swapID", swap.SwapID)
		return
	}

	// Build updated memo and status based on final status
	var memo string
	var status int
	var failReasonCode string
	var failReasonMsg string

	if swap.Status == backend.TxStatusSuccess { // 3 = SUCCESS
		memo = fmt.Sprintf("Swap via %s: %s %s -> %s %s (Success)",
			quote.Provider,
			s.formatAmount(quote.FromAmount),
			s.getTokenSymbol(quote.FromToken),
			s.formatAmount(quote.ToAmount),
			s.getTokenSymbol(quote.ToToken),
		)
		status = dbBackend.TxStatusSuccess // Status: SUCCESS (3)
	} else if swap.Status == backend.TxStatusFailed { // 2 = FAILED
		memo = fmt.Sprintf("Swap via %s: %s %s -> %s %s (Failed: %s)",
			quote.Provider,
			s.formatAmount(quote.FromAmount),
			s.getTokenSymbol(quote.FromToken),
			s.formatAmount(quote.ToAmount),
			s.getTokenSymbol(quote.ToToken),
			swap.FailMessage,
		)
		status = dbBackend.TxStatusFailed // Status: FAILED (2)
		failReasonCode = swap.FailReasonCode
		if failReasonCode == "" {
			failReasonCode = dbBackend.FailReasonChainFailed
		}
		failReasonMsg = swap.FailMessage
	}

	// Update the record using SwapID as the record GUID
	updates := map[string]interface{}{
		"memo":         memo,
		"block_height": finalBlockHeight,
		"status":       status,
	}

	// Add failure info if failed
	if status == dbBackend.TxStatusFailed {
		updates["fail_reason_code"] = failReasonCode
		updates["fail_reason_msg"] = failReasonMsg
	}

	if err := s.db.BackendWalletTxRecord.UpdateWalletTxRecord(swap.SwapID, updates); err != nil {
		log.Error("Failed to update swap history status", "swapID", swap.SwapID, "walletUUID", swap.WalletUUID, "err", err)
	} else {
		log.Info("Swap history status updated", "swapID", swap.SwapID, "walletUUID", swap.WalletUUID, "status", swap.Status, "statusName", dbBackend.TxStatusNames[status])
	}
}

// parseAmount converts amount string to int64 (handles decimals by removing decimal point)
func (s *AggregatorService) parseAmount(amountStr string) (int64, error) {
	// Try to parse as big.Int first
	amount := new(big.Int)
	_, ok := amount.SetString(amountStr, 10)
	if !ok {
		return 0, fmt.Errorf("invalid amount format: %s", amountStr)
	}

	// If amount is too large for int64, use max int64
	if !amount.IsInt64() {
		log.Warn("Amount too large for int64, using max", "amount", amountStr)
		return 9223372036854775807, nil // max int64
	}

	return amount.Int64(), nil
}

// formatAmount formats amount for display (truncate if too long)
func (s *AggregatorService) formatAmount(amountStr string) string {
	// Convert to float for better display
	amount := new(big.Float)
	amount.SetString(amountStr)

	// Divide by 1e18 for typical ERC20 tokens (18 decimals)
	divisor := new(big.Float).SetFloat64(1e18)
	result := new(big.Float).Quo(amount, divisor)

	// Format with 6 decimal places
	return result.Text('f', 6)
}

// getTokenSymbol extracts token symbol from address (simplified)
func (s *AggregatorService) getTokenSymbol(tokenAddress string) string {
	// Common token addresses (simplified mapping)
	switch tokenAddress {
	case "0x0000000000000000000000000000000000000000", "ETH":
		return "ETH"
	case "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48":
		return "USDC"
	case "0xdAC17F958D2ee523a2206206994597C13D831ec7":
		return "USDT"
	case "0x6B175474E89094C44Da98b954EedeAC495271d0F":
		return "DAI"
	default:
		// Return shortened address if unknown
		if len(tokenAddress) > 10 {
			return tokenAddress[:6] + "..." + tokenAddress[len(tokenAddress)-4:]
		}
		return tokenAddress
	}
}

// resolveTokenID returns the token table GUID for (chainID, tokenAddress).
// If the token metadata is missing, we temporarily fall back to the raw address
func (s *AggregatorService) resolveTokenID(ctx context.Context, chainID, tokenAddress string) string {
	if tokenAddress == "" {
		return ""
	}
	if s.db == nil || s.db.BackendToken == nil {
		return tokenAddress
	}

	token, err := s.db.BackendToken.GetByContractAndChain(tokenAddress, chainID)
	if err != nil || token == nil {
		log.Warn("Token not found in token table, fallback to raw address", "chainID", chainID, "token", tokenAddress, "err", err)
		return tokenAddress
	}
	return token.Guid
}
