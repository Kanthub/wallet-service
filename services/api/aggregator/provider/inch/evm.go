package inch

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// GetQuote fetches a quote from 1inch for EVM chains
func (p *Provider) GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error) {
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

// BuildSwap builds the swap transaction based on the quote for EVM chains
func (p *Provider) BuildSwap(ctx context.Context, quote *backend.Quote, userAddress string) (*backend.BuildSwapResponse, error) {
	// TODO: Implement actual HTTP call to 1inch API to build transaction
	// Example endpoint: https://api.1inch.dev/swap/v5.2/{chainId}/swap

	log.Info("1inch BuildSwap called", "chainID", quote.ChainID, "userAddress", userAddress)

	// Placeholder response
	return &backend.BuildSwapResponse{
		Actions: []*backend.Action{
			{
				ActionType: backend.ActionTypeApprove,
				ChainID:    quote.ChainID,
				SigningPayload: &backend.SigningPayload{
					To:      quote.FromToken,
					Data:    "0x", // TODO: Encode approve
					Value:   "0",
					Gas:     "50000",
					ChainID: quote.ChainID,
				},
				Description: fmt.Sprintf("Approve %s to spend %s", quote.Spender, quote.FromToken),
			},
			{
				ActionType: backend.ActionTypeSwap,
				ChainID:    quote.ChainID,
				SigningPayload: &backend.SigningPayload{
					To:      quote.Router,
					Data:    "0x", // TODO: Encode swap calldata from provider
					Value:   quote.FromAmount,
					Gas:     quote.GasEstimate,
					ChainID: quote.ChainID,
				},
				Description: fmt.Sprintf("Swap %s %s for %s %s", quote.FromAmount, quote.FromToken, quote.ToAmount, quote.ToToken),
			},
		},
	}, nil
}
