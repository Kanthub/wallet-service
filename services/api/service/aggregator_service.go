package service

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/roothash-pay/wallet-services/services/api/aggregator/provider"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/store"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/utils"
	"github.com/roothash-pay/wallet-services/services/api/models/backend"
	"github.com/roothash-pay/wallet-services/services/grpc_client/account"
)

// AggregatorService handles swap aggregation operations
type AggregatorService struct {
	providers     []provider.Provider
	quoteStore    store.QuoteStore
	swapStore     store.SwapStore
	validator     *utils.Validator
	accountClient *account.WalletAccountClient
	quoteTTL      time.Duration
}

// NewAggregatorService creates a new aggregator service
func NewAggregatorService(
	providers []provider.Provider,
	quoteStore store.QuoteStore,
	swapStore store.SwapStore,
	validator *utils.Validator,
	accountClient *account.WalletAccountClient,
) *AggregatorService {
	return &AggregatorService{
		providers:     providers,
		quoteStore:    quoteStore,
		swapStore:     swapStore,
		validator:     validator,
		accountClient: accountClient,
		quoteTTL:      5 * time.Minute, // Default 5 minutes
	}
}

// GetQuotes aggregates quotes from multiple providers
func (s *AggregatorService) GetQuotes(ctx context.Context, req *backend.QuoteRequest) (*backend.QuoteResponse, error) {
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
	sort.Slice(quotes, func(i, j int) bool {
		amountI, _ := new(big.Int).SetString(quotes[i].ToAmount, 10)
		amountJ, _ := new(big.Int).SetString(quotes[j].ToAmount, 10)
		return amountI.Cmp(amountJ) > 0
	})

	// Prepare response
	quoteID := uuid.New().String()
	expiresAt := time.Now().Add(s.quoteTTL)

	bestQuote := quotes[0]
	var alternatives []*backend.Quote
	if len(quotes) > 1 {
		alternatives = quotes[1:]
	}

	response := &backend.QuoteResponse{
		QuoteID:      quoteID,
		ExpiresAt:    expiresAt,
		BestQuote:    bestQuote,
		Alternatives: alternatives,
	}

	// Store quote snapshot
	if err := s.quoteStore.Save(ctx, quoteID, response, s.quoteTTL); err != nil {
		log.Error("Failed to save quote", "err", err)
		// Non-fatal, continue
	}

	return response, nil
}

// aggregateQuotes fetches quotes from all providers concurrently
func (s *AggregatorService) aggregateQuotes(ctx context.Context, req *backend.QuoteRequest) ([]*backend.Quote, error) {
	g, ctx := errgroup.WithContext(ctx)
	quoteChan := make(chan *backend.Quote, len(s.providers))

	for _, p := range s.providers {
		p := p // Capture loop variable
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

// PrepareSwap generates a transaction plan for a swap
func (s *AggregatorService) PrepareSwap(ctx context.Context, req *backend.PrepareSwapRequest) (*backend.PrepareSwapResponse, error) {
	// Retrieve quote from store
	quoteResp, err := s.quoteStore.Get(ctx, req.QuoteID)
	if err != nil {
		return nil, fmt.Errorf("quote not found or expired: %w", err)
	}

	// Validate expiration
	if time.Now().After(quoteResp.ExpiresAt) {
		return nil, fmt.Errorf("quote expired")
	}

	bestQuote := quoteResp.BestQuote

	// Generate swap ID
	swapID := uuid.New().String()

	// Generate actions based on chain type
	var actions []*backend.Action

	switch bestQuote.ChainType {
	case backend.ChainTypeEVM:
		actions, err = s.generateEVMActions(ctx, bestQuote, req.UserAddress)
		if err != nil {
			return nil, err
		}
	case backend.ChainTypeSolana:
		actions, err = s.generateSolanaActions(ctx, bestQuote, req.UserAddress)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported chain type: %s", bestQuote.ChainType)
	}

	// Create swap record
	swap := &backend.Swap{
		SwapID:      swapID,
		QuoteID:     req.QuoteID,
		UserAddress: req.UserAddress,
		State:       backend.SwapStatePending,
		Steps:       make([]*backend.Step, len(actions)),
	}

	// Initialize steps
	for i, action := range actions {
		swap.Steps[i] = &backend.Step{
			StepIndex:  i,
			ActionType: action.ActionType,
			State:      backend.StepStatePending,
		}
	}

	if err := s.swapStore.CreateSwap(ctx, swap); err != nil {
		return nil, err
	}

	return &backend.PrepareSwapResponse{
		SwapID:  swapID,
		Actions: actions,
	}, nil
}

// generateEVMActions generates actions for EVM chains
func (s *AggregatorService) generateEVMActions(ctx context.Context, quote *backend.Quote, userAddress string) ([]*backend.Action, error) {
	var actions []*backend.Action

	// Check if approval is needed (for non-native tokens)
	isNativeToken := quote.FromToken == "0x0000000000000000000000000000000000000000" || quote.FromToken == "ETH"

	if !isNativeToken && quote.Spender != "" {
		// TODO: Check actual allowance
		// For now, always add APPROVE action
		actions = append(actions, &backend.Action{
			ActionType: backend.ActionTypeApprove,
			ChainID:    quote.ChainID,
			SigningPayload: &backend.SigningPayload{
				To:      quote.FromToken,
				Data:    fmt.Sprintf("approve(%s,%s)", quote.Spender, quote.FromAmount), // TODO: Encode properly
				Value:   "0",
				Gas:     "50000",
				ChainID: quote.ChainID,
			},
			Description: fmt.Sprintf("Approve %s to spend %s", quote.Spender, quote.FromToken),
		})
	}

	// Add SWAP action
	value := "0"
	if isNativeToken {
		value = quote.FromAmount
	}

	actions = append(actions, &backend.Action{
		ActionType: backend.ActionTypeSwap,
		ChainID:    quote.ChainID,
		SigningPayload: &backend.SigningPayload{
			To:      quote.Router,
			Data:    "0x", // TODO: Encode swap calldata from provider
			Value:   value,
			Gas:     quote.GasEstimate,
			ChainID: quote.ChainID,
		},
		Description: fmt.Sprintf("Swap %s %s for %s %s", quote.FromAmount, quote.FromToken, quote.ToAmount, quote.ToToken),
	})

	return actions, nil
}

// generateSolanaActions generates actions for Solana
func (s *AggregatorService) generateSolanaActions(ctx context.Context, quote *backend.Quote, userAddress string) ([]*backend.Action, error) {
	// TODO: Implement Solana transaction generation
	return []*backend.Action{
		{
			ActionType: backend.ActionTypeSwap,
			ChainID:    quote.ChainID,
			SigningPayload: &backend.SigningPayload{
				SerializedTx: "TODO: Solana serialized transaction",
			},
			Description: fmt.Sprintf("Swap %s for %s on Solana", quote.FromToken, quote.ToToken),
		},
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

	// TODO: Validate signed transaction against quote snapshot
	// - Check chainId
	// - Check router/spender whitelist
	// - Check value limits

	// Get quote for chain info
	_, err = s.quoteStore.Get(ctx, swap.QuoteID)
	if err != nil {
		return nil, fmt.Errorf("quote not found: %w", err)
	}

	chain := "Ethereum"  // TODO: Map from chainID
	network := "mainnet" // TODO: Get from config

	// Broadcast transaction using SendTx
	result, err := s.accountClient.SendTx(ctx, account.SendTxParams{
		ConsumerToken: "", // TODO: Get from config
		Chain:         chain,
		Coin:          "ETH", // TODO: Map from chainID
		Network:       network,
		RawTx:         req.SignedTx,
	})
	if err != nil {
		// Update step as failed
		step.State = backend.StepStateFailed
		step.FailReasonCode = "BROADCAST_FAILED"
		step.FailMessage = err.Error()
		_ = s.swapStore.UpdateStep(ctx, req.SwapID, req.StepIndex, step)

		return nil, err
	}

	txHash := result.TxHash

	// Update step
	now := time.Now()
	step.TxHash = txHash
	step.State = backend.StepStateSubmitted
	step.SubmittedAt = &now
	step.IdempotencyKey = req.IdempotencyKey

	if err := s.swapStore.UpdateStep(ctx, req.SwapID, req.StepIndex, step); err != nil {
		log.Error("Failed to update step", "err", err)
	}

	// Record idempotency
	if swapStore, ok := s.swapStore.(*store.InMemorySwapStore); ok {
		_ = swapStore.RecordIdempotency(ctx, req.SwapID, req.StepIndex, req.IdempotencyKey, txHash)
	}

	// Update swap state
	swap.State = backend.SwapStateSubmitted
	_ = s.swapStore.UpdateSwap(ctx, swap)

	return &backend.SubmitSignedTxResponse{TxHash: txHash}, nil
}

// GetSwapStatus retrieves the status of a swap
func (s *AggregatorService) GetSwapStatus(ctx context.Context, swapID string) (*backend.SwapStatusResponse, error) {
	swap, err := s.swapStore.GetSwap(ctx, swapID)
	if err != nil {
		return nil, err
	}

	// Query status for each submitted transaction
	for _, step := range swap.Steps {
		if step.TxHash != "" && step.State == backend.StepStateSubmitted {
			// Get quote for chain info
			quoteResp, _ := s.quoteStore.Get(ctx, swap.QuoteID)
			if quoteResp != nil {
				chain := "Ethereum" // TODO: Map from chainID
				network := "mainnet"

				txInfo, err := s.accountClient.GetTxByHash(ctx, "", chain, "ETH", network, step.TxHash)
				if err != nil {
					log.Warn("Failed to get tx status", "txHash", step.TxHash, "err", err)
					continue
				}

				// Update step based on status
				if txInfo.Status == 3 { // TxStatus_Success
					now := time.Now()
					step.State = backend.StepStateConfirmed
					step.ConfirmedAt = &now
					_ = s.swapStore.UpdateStep(ctx, swapID, step.StepIndex, step)
				} else if txInfo.Status == 2 { // TxStatus_Failed
					step.State = backend.StepStateFailed
					step.FailReasonCode = "TX_FAILED"
					step.FailMessage = "Transaction failed on chain"
					_ = s.swapStore.UpdateStep(ctx, swapID, step.StepIndex, step)
				}
			}
		}
	}

	// Determine overall swap state
	allConfirmed := true
	anyFailed := false
	for _, step := range swap.Steps {
		if step.State != backend.StepStateConfirmed {
			allConfirmed = false
		}
		if step.State == backend.StepStateFailed {
			anyFailed = true
		}
	}

	if anyFailed {
		swap.State = backend.SwapStateFailed
	} else if allConfirmed {
		swap.State = backend.SwapStateConfirmed
	}

	_ = s.swapStore.UpdateSwap(ctx, swap)

	return &backend.SwapStatusResponse{
		SwapID:         swap.SwapID,
		State:          swap.State,
		Steps:          swap.Steps,
		FailReasonCode: swap.FailReasonCode,
		FailMessage:    swap.FailMessage,
	}, nil
}
