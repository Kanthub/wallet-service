package solana

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// JupiterProvider implements the Provider interface for Jupiter
type JupiterProvider struct {
	apiURL string
}

// NewJupiterProvider creates a new Jupiter provider
func NewJupiterProvider(apiURL string) *JupiterProvider {
	return &JupiterProvider{
		apiURL: apiURL,
	}
}

// GetQuote fetches a quote from Jupiter
func (p *JupiterProvider) GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error) {
	// TODO: Implement actual HTTP call to Jupiter API
	// Example endpoint: https://quote-api.jup.ag/v6/quote
	// Request params: inputMint, outputMint, amount, slippageBps

	log.Info("Jupiter GetQuote called", "fromToken", req.FromToken, "toToken", req.ToToken)

	// Placeholder response
	return &backend.Quote{
		Provider:    p.Name(),
		ChainType:   backend.ChainTypeSolana,
		ChainID:     req.FromChainID,
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		FromAmount:  req.Amount,
		ToAmount:    "0",    // TODO: Parse from API response
		GasEstimate: "5000", // Solana compute units
	}, fmt.Errorf("TODO: implement Jupiter API integration")
}

// Name returns the provider name
func (p *JupiterProvider) Name() string {
	return "Jupiter"
}

// SupportedChainType returns Solana
func (p *JupiterProvider) SupportedChainType() backend.ChainType {
	return backend.ChainTypeSolana
}
