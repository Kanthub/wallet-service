package jupiter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// Provider implements the Provider interface for Jupiter
type Provider struct {
	apiURL     string
	httpClient *http.Client
}

// NewProvider creates a new Jupiter provider
func NewProvider(apiURL string) *Provider {
	if apiURL == "" {
		apiURL = "https://quote-api.jup.ag/v6"
	}
	return &Provider{
		apiURL: apiURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// JupiterQuoteResponse represents the response from Jupiter Quote API
type JupiterQuoteResponse struct {
	InputMint            string      `json:"inputMint"`
	InAmount             string      `json:"inAmount"`
	OutputMint           string      `json:"outputMint"`
	OutAmount            string      `json:"outAmount"`
	OtherAmountThreshold string      `json:"otherAmountThreshold"`
	SwapMode             string      `json:"swapMode"`
	SlippageBps          int         `json:"slippageBps"`
	PlatformFee          interface{} `json:"platformFee"`
	PriceImpactPct       string      `json:"priceImpactPct"`
	RoutePlan            []struct {
		SwapInfo struct {
			AmmKey     string `json:"ammKey"`
			Label      string `json:"label"`
			InputMint  string `json:"inputMint"`
			OutputMint string `json:"outputMint"`
			InAmount   string `json:"inAmount"`
			OutAmount  string `json:"outAmount"`
			FeeAmount  string `json:"feeAmount"`
			FeeMint    string `json:"feeMint"`
		} `json:"swapInfo"`
		Percent int `json:"percent"`
	} `json:"routePlan"`
	ContextSlot int64 `json:"contextSlot"`
	TimeTaken   int64 `json:"timeTaken"`
}

// JupiterSwapRequest represents the request to Jupiter Swap API
type JupiterSwapRequest struct {
	QuoteResponse    json.RawMessage `json:"quoteResponse"`
	UserPublicKey    string          `json:"userPublicKey"`
	WrapAndUnwrapSol bool            `json:"wrapAndUnwrapSol"`
	// UseSharedAccounts bool   `json:"useSharedAccounts,omitempty"`
	// FeeAccount        string `json:"feeAccount,omitempty"`
	// ComputeUnitPriceMicroLamports int `json:"computeUnitPriceMicroLamports,omitempty"`
}

// JupiterSwapResponse represents the response from Jupiter Swap API
type JupiterSwapResponse struct {
	SwapTransaction      string  `json:"swapTransaction"` // Base64 encoded transaction
	LastValidBlockHeight float64 `json:"lastValidBlockHeight"`
}

// GetQuote fetches a quote from Jupiter for Solana
func (p *Provider) GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error) {
	// Build request URL
	u, err := url.Parse(p.apiURL + "/quote")
	if err != nil {
		return nil, fmt.Errorf("invalid API URL: %w", err)
	}

	q := u.Query()
	q.Set("inputMint", req.FromToken)
	q.Set("outputMint", req.ToToken)
	q.Set("amount", req.Amount)
	q.Set("slippageBps", fmt.Sprintf("%d", int(req.SlippageBps)))
	// Optional: q.Set("onlyDirectRoutes", "true")
	// Optional: q.Set("asLegacyTransaction", "true")
	u.RawQuery = q.Encode()

	// Execute request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Jupiter API error (status %d): %s", resp.StatusCode, string(body))
	}

	var jupResp JupiterQuoteResponse
	if err := json.Unmarshal(body, &jupResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Calculate estimated gas (Solana compute units cost is minimal, but we can try to estimate)
	// For now, use a default placeholder or 0 if not provided by API directly in a simple way
	gasEstimate := "5000"

	return &backend.Quote{
		Provider:    p.Name(),
		ChainType:   backend.ChainTypeSolana,
		ChainID:     req.FromChainID,
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		FromAmount:  req.Amount,
		ToAmount:    jupResp.OutAmount,
		GasEstimate: gasEstimate,
		Fees:        "0",          // Complex to calculate exactly from route plan without token prices
		Raw:         string(body), // Store raw JSON to pass back to swap endpoint
	}, nil
}

// BuildSwap builds the swap transaction based on the quote for Solana
func (p *Provider) BuildSwap(ctx context.Context, quote *backend.Quote, userAddress string) (*backend.BuildSwapResponse, error) {
	log.Info("Jupiter BuildSwap called", "chainID", quote.ChainID, "userAddress", userAddress)

	if quote.Raw == "" {
		return nil, fmt.Errorf("quote raw data is empty, cannot build swap")
	}

	// Prepare Swap Request
	// Note: We need to pass the exact JSON object received from quote API
	swapReq := JupiterSwapRequest{
		QuoteResponse:    json.RawMessage(quote.Raw),
		UserPublicKey:    userAddress,
		WrapAndUnwrapSol: true,
	}

	reqBody, err := json.Marshal(swapReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal swap request: %w", err)
	}

	// Execute Request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiURL+"/swap", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create swap request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute swap request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read swap response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Jupiter Swap API error (status %d): %s", resp.StatusCode, string(body))
	}

	var jupSwapResp JupiterSwapResponse
	if err := json.Unmarshal(body, &jupSwapResp); err != nil {
		return nil, fmt.Errorf("failed to parse swap response: %w", err)
	}

	// Create Action
	action := &backend.Action{
		ActionType: backend.ActionTypeSwap,
		ChainID:    quote.ChainID,
		SigningPayload: &backend.SigningPayload{
			SerializedTx: jupSwapResp.SwapTransaction,
			ChainID:      quote.ChainID,
		},
		Description: fmt.Sprintf("Swap %s to %s via Jupiter",
			getTokenSymbol(quote.FromToken),
			getTokenSymbol(quote.ToToken)),
	}

	return &backend.BuildSwapResponse{
		Actions: []*backend.Action{action},
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

// getTokenSymbol extracts a simple token symbol from address (simulated)
func getTokenSymbol(tokenAddress string) string {
	if len(tokenAddress) > 10 {
		return tokenAddress[:4] + "..." + tokenAddress[len(tokenAddress)-4:]
	}
	return tokenAddress
}
