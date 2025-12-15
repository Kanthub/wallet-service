package oneinch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

type Provider struct {
	apiURL     string
	apiKey     string
	httpClient *http.Client
}

// OneInchQuoteResponse represents the response from 1inch quote API
type OneInchQuoteResponse struct {
	DstAmount string `json:"dstAmount"`
	Gas       int64  `json:"gas"`
}

// NewProvider creates a new 1inch provider
func NewProvider(apiURL, apiKey string) *Provider {
	if apiURL == "" {
		apiURL = "https://api.1inch.dev/swap/v6.0"
	}
	return &Provider{
		apiURL: apiURL,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *Provider) Name() string {
	return "1inch"
}

func (p *Provider) SupportedChainType() backend.ChainType {
	return backend.ChainTypeEVM
}

func (p *Provider) GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error) {
	resp, err := p.fetchQuote(ctx, req)
	if err != nil {
		return nil, err
	}
	// Store raw response
	rawBytes, _ := json.Marshal(resp)

	return p.convertToQuoteEVM(req, resp, string(rawBytes)), nil
}

func (p *Provider) BuildSwap(ctx context.Context, quote *backend.Quote, userAddress string) (*backend.BuildSwapResponse, error) {
	return p.buildSwapEVM(ctx, quote, userAddress)
}

// buildQuoteRequest builds the quote request URL
func (p *Provider) buildQuoteRequest(req *backend.QuoteRequest) (string, error) {
	// API: /swap/v6.0/{chainID}/quote
	u, err := url.Parse(fmt.Sprintf("%s/%s/quote", p.apiURL, req.FromChainID))
	if err != nil {
		return "", fmt.Errorf("invalid API URL: %w", err)
	}

	q := u.Query()
	q.Set("src", req.FromToken)
	q.Set("dst", req.ToToken)
	q.Set("amount", req.Amount)
	//q.Set("slippage", fmt.Sprintf("%f", req.SlippageBps/100.0)) // 1inch uses percentage (e.g. 1 for 1%)

	// Optional: add fee/referral params if needed

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// fetchQuote fetches quote from 1inch API
func (p *Provider) fetchQuote(ctx context.Context, req *backend.QuoteRequest) (*OneInchQuoteResponse, error) {
	reqURL, err := p.buildQuoteRequest(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if p.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
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
		return nil, fmt.Errorf("1inch API error (status %d): %s", resp.StatusCode, string(body))
	}

	var quoteResp OneInchQuoteResponse
	if err := json.Unmarshal(body, &quoteResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &quoteResp, nil
}

// getTokenSymbol extracts a simple token symbol from address
func getTokenSymbol(tokenAddress string) string {
	// Common token addresses (lowercase for comparison)
	tokens := map[string]string{
		"0x0000000000000000000000000000000000000000": "ETH",
		"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee": "ETH",
		"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48": "USDC",
		"0xdac17f958d2ee523a2206206994597c13d831ec7": "USDT",
		"0x6b175474e89094c44da98b954eedeac495271d0f": "DAI",
		"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599": "WBTC",
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": "WETH",
	}

	// Normalize to lowercase for lookup
	normalized := ""
	if len(tokenAddress) > 0 {
		normalized = tokenAddress
		if len(normalized) > 2 && normalized[:2] == "0x" {
			normalized = "0x" + normalized[2:]
		}
	}

	// Try to find in map (case-insensitive)
	for addr, symbol := range tokens {
		if len(normalized) == len(addr) && normalized == addr {
			return symbol
		}
	}

	// Return shortened address if not found
	if len(tokenAddress) > 10 {
		return tokenAddress[:6] + "..." + tokenAddress[len(tokenAddress)-4:]
	}
	return tokenAddress
}
