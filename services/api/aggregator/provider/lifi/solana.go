package lifi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// SolanaQuoteResponse represents the Solana-specific fields in LiFi quote response
type SolanaQuoteResponse struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Tool   string `json:"tool"`
	Action struct {
		FromChainID int `json:"fromChainId"`
		ToChainID   int `json:"toChainId"`
		FromToken   struct {
			Address  string `json:"address"`
			Symbol   string `json:"symbol"`
			Decimals int    `json:"decimals"`
			ChainID  int    `json:"chainId"`
		} `json:"fromToken"`
		ToToken struct {
			Address  string `json:"address"`
			Symbol   string `json:"symbol"`
			Decimals int    `json:"decimals"`
			ChainID  int    `json:"chainId"`
		} `json:"toToken"`
		FromAmount string  `json:"fromAmount"`
		ToAmount   string  `json:"toAmount"`
		Slippage   float64 `json:"slippage"`
	} `json:"action"`
	Estimate struct {
		FromAmount        string `json:"fromAmount"`
		ToAmount          string `json:"toAmount"`
		ToAmountMin       string `json:"toAmountMin"`
		ExecutionDuration int    `json:"executionDuration"`
		FeeCosts          []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Token       struct {
				Address  string `json:"address"`
				Symbol   string `json:"symbol"`
				Decimals int    `json:"decimals"`
			} `json:"token"`
			Amount     string `json:"amount"`
			AmountUSD  string `json:"amountUSD"`
			Percentage string `json:"percentage"`
			Included   bool   `json:"included"`
		} `json:"feeCosts"`
	} `json:"estimate"`
	TransactionRequest struct {
		// Solana uses a different transaction format
		Data string `json:"data"` // Base64-encoded serialized transaction
	} `json:"transactionRequest"`
}

// convertToQuoteSolana converts LiFi Solana response to internal Quote format
func (p *Provider) convertToQuoteSolana(req *backend.QuoteRequest, lifiResp *LifiQuoteResponse) *backend.Quote {
	// For Solana, gas is typically measured in compute units
	// LiFi might provide this in feeCosts
	gasEstimate := "5000" // Default Solana compute units

	// Calculate total fees if available
	fees := "0"
	if len(lifiResp.Estimate.GasCosts) > 0 {
		// Try to extract gas/fee information
		gasEstimate = lifiResp.Estimate.GasCosts[0].Estimate
		fees = lifiResp.Estimate.GasCosts[0].Amount
	}

	// Marshal raw response for debugging
	rawJSON, _ := json.Marshal(lifiResp)

	return &backend.Quote{
		Provider:    p.Name(),
		ChainType:   backend.ChainTypeSolana,
		ChainID:     req.FromChainID,
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		FromAmount:  req.Amount,
		ToAmount:    lifiResp.Estimate.ToAmount,
		GasEstimate: gasEstimate,
		Fees:        fees,
		Spender:     "", // Solana doesn't use spender concept like EVM
		Router:      "", // Solana doesn't use router concept like EVM
		Raw:         string(rawJSON),
	}
}

// buildSwapSolana builds the swap transaction for Solana chains
func (p *Provider) buildSwapSolana(ctx context.Context, quote *backend.Quote, userAddress string) (*backend.BuildSwapResponse, error) {
	log.Info("LiFi BuildSwap called for Solana", "chainID", quote.ChainID, "userAddress", userAddress)

	// Parse the raw LiFi response from quote
	if quote.Raw == "" {
		return nil, fmt.Errorf("quote.Raw is empty, cannot build swap")
	}

	var lifiResp LifiQuoteResponse
	if err := json.Unmarshal([]byte(quote.Raw), &lifiResp); err != nil {
		return nil, fmt.Errorf("failed to parse quote.Raw: %w", err)
	}

	// For Solana, there's no approval step needed
	// Just create a single SWAP action with the serialized transaction

	// Extract the serialized transaction data
	serializedTx := lifiResp.TransactionRequest.Data
	if serializedTx == "" {
		return nil, fmt.Errorf("no transaction data found in LiFi response")
	}

	// Create the swap action
	swapAction := &backend.Action{
		ActionType: backend.ActionTypeSwap,
		ChainID:    quote.ChainID,
		SigningPayload: &backend.SigningPayload{
			SerializedTx: serializedTx,
			ChainID:      quote.ChainID,
		},
		Description: fmt.Sprintf("Swap %s to %s via LiFi on Solana",
			getTokenSymbol(quote.FromToken),
			getTokenSymbol(quote.ToToken)),
	}

	actions := []*backend.Action{swapAction}

	log.Info("LiFi BuildSwap success for Solana",
		"actionCount", len(actions),
		"userAddress", userAddress)

	return &backend.BuildSwapResponse{
		Actions: actions,
	}, nil
}
