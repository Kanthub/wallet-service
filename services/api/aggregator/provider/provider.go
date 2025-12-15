package provider

import (
	"context"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// Provider defines the interface for DEX quote providers
type Provider interface {
	// GetQuote fetches a quote for the given request
	GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error)

	// BuildSwap builds the swap transaction based on the quote
	BuildSwap(ctx context.Context, quoteID *backend.Quote, userAddress string) (*backend.BuildSwapResponse, error)

	// Name returns the provider name
	Name() string

	// SupportedChainType returns the chain type this provider supports
	SupportedChainType() backend.ChainType
}
