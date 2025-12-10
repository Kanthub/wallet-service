package lifi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/roothash-pay/wallet-services/services/api/aggregator/utils"
	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// Provider implements the Provider interface for LiFi
type Provider struct {
	apiURL     string
	apiKey     string
	httpClient *http.Client
	evmCaller  *utils.EVMCaller
}

// QuoteResponse represents the response from LiFi quote API
type QuoteResponse struct {
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
		FromAmount string `json:"fromAmount"`
		ToAmount   string `json:"toAmount"`
		Slippage   string `json:"slippage"`
	} `json:"action"`
	Estimate struct {
		FromAmount        string `json:"fromAmount"`
		ToAmount          string `json:"toAmount"`
		ToAmountMin       string `json:"toAmountMin"`
		ApprovalAddress   string `json:"approvalAddress"`
		ExecutionDuration int    `json:"executionDuration"`
		GasCosts          []struct {
			Type     string `json:"type"`
			Estimate string `json:"estimate"`
			Limit    string `json:"limit"`
			Amount   string `json:"amount"`
			Token    struct {
				Address  string `json:"address"`
				Symbol   string `json:"symbol"`
				Decimals int    `json:"decimals"`
			} `json:"token"`
		} `json:"gasCosts"`
	} `json:"estimate"`
	TransactionRequest struct {
		From     string `json:"from"`
		To       string `json:"to"`
		ChainID  int    `json:"chainId"`
		Data     string `json:"data"`
		Value    string `json:"value"`
		GasLimit string `json:"gasLimit"`
		GasPrice string `json:"gasPrice"`
	} `json:"transactionRequest"`
}

// NewProvider creates a new LiFi provider
func NewProvider(apiURL, apiKey string, evmCaller *utils.EVMCaller) *Provider {
	return &Provider{
		apiURL: apiURL,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		evmCaller: evmCaller,
	}
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "lifi"
}

// SupportedChainType returns EVM (LiFi supports both EVM and Solana, but we start with EVM)
func (p *Provider) SupportedChainType() backend.ChainType {
	return backend.ChainTypeEVM
}

// buildQuoteURL constructs the LiFi quote API URL with query parameters
func (p *Provider) buildQuoteURL(req *backend.QuoteRequest) (string, error) {
	baseURL := p.apiURL
	if baseURL == "" {
		baseURL = "https://li.quest/v1"
	}

	// Parse base URL
	u, err := url.Parse(baseURL + "/quote")
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Build query parameters
	q := u.Query()
	q.Set("fromChain", req.FromChainID)
	q.Set("toChain", req.ToChainID)
	q.Set("fromToken", req.FromToken)
	q.Set("toToken", req.ToToken)
	q.Set("fromAmount", req.Amount)

	// Convert slippage from bps to decimal (e.g., 50 bps = 0.005)
	slippageDecimal := float64(req.SlippageBps) / 10000.0
	q.Set("slippage", strconv.FormatFloat(slippageDecimal, 'f', 4, 64))

	// Add user address if provided
	if req.UserAddress != "" {
		q.Set("fromAddress", req.UserAddress)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// fetchQuote makes HTTP request to LiFi API
func (p *Provider) fetchQuote(ctx context.Context, req *backend.QuoteRequest) (*QuoteResponse, error) {
	// Build request URL
	quoteURL, err := p.buildQuoteURL(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build quote URL: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "GET", quoteURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	if p.apiKey != "" {
		httpReq.Header.Set("x-lifi-api-key", p.apiKey)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LiFi API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var lifiResp QuoteResponse
	if err := json.Unmarshal(body, &lifiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &lifiResp, nil
}

// isNativeToken checks if the token is a native token
func isNativeToken(token string) bool {
	return token == "0x0000000000000000000000000000000000000000" ||
		token == "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
}

// buildApproveData builds the ERC20 approve function call data
// Function signature: approve(address spender, uint256 amount)
// Selector: 0x095ea7b3
func buildApproveData(spender string, amount string) string {
	// ERC20 approve function selector
	selector := "0x095ea7b3"

	// Remove "0x" prefix from spender if present
	spenderAddr := spender
	if len(spenderAddr) > 2 && spenderAddr[:2] == "0x" {
		spenderAddr = spenderAddr[2:]
	}

	// Pad spender address to 32 bytes (64 hex chars)
	spenderPadded := fmt.Sprintf("%064s", spenderAddr)

	// For amount, we use max uint256 for unlimited approval
	// This is a common practice to avoid multiple approval transactions
	amountPadded := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"

	return selector + spenderPadded + amountPadded
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
