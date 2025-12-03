package evm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// OneInchProvider implements the Provider interface for 1inch
type OneInchProvider struct {
	apiURL string
	apiKey string
}

// NewOneInchProvider creates a new 1inch provider
func NewOneInchProvider(apiURL, apiKey string) *OneInchProvider {
	return &OneInchProvider{
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

// GetQuote fetches a quote from 1inch
func (p *OneInchProvider) GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error) {
	// TODO: Implement actual HTTP call to 1inch API
	// Example endpoint: https://api.1inch.dev/swap/v5.2/{chainId}/quote
	// Request params: src, dst, amount, protocols

	log.Info("1inch GetQuote called", "fromToken", req.FromToken, "toToken", req.ToToken)

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
	}, fmt.Errorf("TODO: implement 1inch API integration")
}

// Name returns the provider name
func (p *OneInchProvider) Name() string {
	return "1inch"
}

// SupportedChainType returns EVM
func (p *OneInchProvider) SupportedChainType() backend.ChainType {
	return backend.ChainTypeEVM
}
