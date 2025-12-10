package inch

import (
	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// Provider implements the Provider interface for 1inch
type Provider struct {
	apiURL string
	apiKey string
}

// NewProvider creates a new 1inch provider
func NewProvider(apiURL, apiKey string) *Provider {
	return &Provider{
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "1inch"
}

// SupportedChainType returns EVM
func (p *Provider) SupportedChainType() backend.ChainType {
	return backend.ChainTypeEVM
}
