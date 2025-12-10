package jupiter

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// Provider implements the Provider interface for Jupiter
type Provider struct {
	apiURL string
}

// NewProvider creates a new Jupiter provider
func NewProvider(apiURL string) *Provider {
	return &Provider{
		apiURL: apiURL,
	}
}

// GetQuote fetches a quote from Jupiter for Solana
func (p *Provider) GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error) {
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

// BuildSwap builds the swap transaction based on the quote for Solana
func (p *Provider) BuildSwap(ctx context.Context, quote *backend.Quote, userAddress string) (*backend.BuildSwapResponse, error) {
	// TODO: Implement actual HTTP call to Jupiter API to build transaction
	// Example endpoint: https://quote-api.jup.ag/v6/swap

	log.Info("Jupiter BuildSwap called", "chainID", quote.ChainID, "userAddress", userAddress)

	// Placeholder response
	return &backend.BuildSwapResponse{
		Actions: []*backend.Action{
			{
				ActionType: backend.ActionTypeSwap,
				ChainID:    quote.ChainID,
				SigningPayload: &backend.SigningPayload{
					SerializedTx: "TODO: Solana serialized transaction",
				},
				Description: fmt.Sprintf("Swap %s for %s on Solana", quote.FromToken, quote.ToToken),
			},
		},
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "Jupiter"
}

// SupportedChainType returns Solana
func (p *Provider) SupportedChainType() backend.ChainType {
	return backend.ChainTypeSolana
}
