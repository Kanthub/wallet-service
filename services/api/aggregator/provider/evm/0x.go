package evm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// ZeroXProvider implements the Provider interface for 0x Protocol
type ZeroXProvider struct {
	apiURL string
	apiKey string
}

// NewZeroXProvider creates a new 0x provider
func NewZeroXProvider(apiURL, apiKey string) *ZeroXProvider {
	return &ZeroXProvider{
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

// GetQuote fetches a quote from 0x Protocol
func (p *ZeroXProvider) GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error) {
	// TODO: Implement actual HTTP call to 0x API
	// Example endpoint: https://api.0x.org/swap/v1/quote
	// Request params: sellToken, buyToken, sellAmount, slippagePercentage

	log.Info("ZeroX GetQuote called", "fromToken", req.FromToken, "toToken", req.ToToken)

	// Placeholder response
	return &backend.Quote{
		Provider:    p.Name(),
		ChainType:   backend.ChainTypeEVM,
		ChainID:     req.FromChainID,
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		FromAmount:  req.Amount,
		ToAmount:    "0", // TODO: Parse from API response
		GasEstimate: "21000",
		Spender:     "0x0000000000000000000000000000000000000000", // TODO: Get from API
		Router:      "0x0000000000000000000000000000000000000000", // TODO: Get from API
	}, fmt.Errorf("TODO: implement 0x API integration")
}

// Name returns the provider name
func (p *ZeroXProvider) Name() string {
	return "0x"
}

// SupportedChainType returns EVM
func (p *ZeroXProvider) SupportedChainType() backend.ChainType {
	return backend.ChainTypeEVM
}
