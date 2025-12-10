package lifi

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// GetQuote fetches a quote from LiFi for EVM chains
func (p *Provider) GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error) {
	lifiResp, err := p.fetchQuote(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert to internal Quote format
	quote := p.convertToQuote(req, lifiResp)

	log.Info("LiFi GetQuote success",
		"fromToken", req.FromToken,
		"toToken", req.ToToken,
		"fromAmount", req.Amount,
		"toAmount", quote.ToAmount,
	)

	return quote, nil
}

// BuildSwap builds the swap transaction based on the quote for EVM chains
func (p *Provider) BuildSwap(ctx context.Context, quote *backend.Quote, userAddress string) (*backend.BuildSwapResponse, error) {
	log.Info("LiFi BuildSwap called", "chainID", quote.ChainID, "userAddress", userAddress)

	// Parse the raw LiFi response from quote
	if quote.Raw == "" {
		return nil, fmt.Errorf("quote.Raw is empty, cannot build swap")
	}

	var lifiResp QuoteResponse
	if err := json.Unmarshal([]byte(quote.Raw), &lifiResp); err != nil {
		return nil, fmt.Errorf("failed to parse quote.Raw: %w", err)
	}

	// Build actions list
	actions := []*backend.Action{}

	// Check if approval is needed
	if !isNativeToken(quote.FromToken) && quote.Spender != "" {
		// Query allowance using EVMCaller
		var allowance *big.Int
		var err error

		if p.evmCaller != nil {
			allowance, err = p.evmCaller.GetERC20Allowance(ctx, quote.ChainID, quote.FromToken, userAddress, quote.Spender)
			if err != nil {
				// If allowance check fails, log warning and assume zero allowance (will add approve action)
				log.Warn("Failed to check allowance, assuming zero", "err", err)
				allowance = big.NewInt(0)
			}
		} else {
			// If no EVMCaller, assume zero allowance (will add approve action)
			log.Warn("EVMCaller not available, assuming zero allowance")
			allowance = big.NewInt(0)
		}

		// Parse needed amount
		need := new(big.Int)
		need, ok := need.SetString(quote.FromAmount, 10)
		if !ok {
			return nil, fmt.Errorf("invalid fromAmount: %s", quote.FromAmount)
		}

		// Add approve action if allowance is insufficient
		if allowance.Cmp(need) < 0 {
			approveAmt := need.String()
			approveData := buildApproveData(quote.Spender, approveAmt)

			approveAction := &backend.Action{
				ActionType: backend.ActionTypeApprove,
				ChainID:    quote.ChainID,
				SigningPayload: &backend.SigningPayload{
					To:      quote.FromToken,
					Data:    approveData,
					Value:   "0",
					Gas:     "",
					ChainID: quote.ChainID,
				},
				Description: fmt.Sprintf("Approve %s to spend %s", quote.Spender, getTokenSymbol(quote.FromToken)),
			}

			actions = append(actions, approveAction)
		}
	}

	// Add SWAP action
	swapAction := &backend.Action{
		ActionType: backend.ActionTypeSwap,
		ChainID:    quote.ChainID,
		SigningPayload: &backend.SigningPayload{
			To:      lifiResp.TransactionRequest.To,
			Data:    lifiResp.TransactionRequest.Data,
			Value:   lifiResp.TransactionRequest.Value,
			Gas:     lifiResp.TransactionRequest.GasLimit,
			ChainID: quote.ChainID,
		},
		Description: fmt.Sprintf("Swap %s to %s via LiFi",
			getTokenSymbol(quote.FromToken),
			getTokenSymbol(quote.ToToken)),
	}
	actions = append(actions, swapAction)

	log.Info("LiFi BuildSwap success",
		"actionCount", len(actions),
		"hasApprove", len(actions) > 1,
		"userAddress", userAddress)

	return &backend.BuildSwapResponse{
		Actions: actions,
	}, nil
}

// convertToQuote converts LiFi response to internal Quote format
func (p *Provider) convertToQuote(req *backend.QuoteRequest, lifiResp *QuoteResponse) *backend.Quote {
	// Calculate total gas estimate
	gasEstimate := "0"
	if len(lifiResp.Estimate.GasCosts) > 0 {
		gasEstimate = lifiResp.Estimate.GasCosts[0].Estimate
	}

	// Get spender (approval address) and router (transaction to address)
	spender := lifiResp.Estimate.ApprovalAddress
	router := lifiResp.TransactionRequest.To

	// Marshal raw response for debugging
	rawJSON, _ := json.Marshal(lifiResp)

	return &backend.Quote{
		Provider:    p.Name(),
		ChainType:   backend.ChainTypeEVM,
		ChainID:     req.FromChainID,
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		FromAmount:  lifiResp.Estimate.FromAmount,
		ToAmount:    lifiResp.Estimate.ToAmount,
		GasEstimate: gasEstimate,
		Spender:     spender,
		Router:      router,
		Raw:         string(rawJSON),
	}
}
