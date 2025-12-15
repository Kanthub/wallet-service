package lifi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/log"

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

// LifiQuoteResponse represents the response from LiFi quote API
type LifiQuoteResponse struct {
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
		FromAmount string  `json:"fromAmount"`
		ToAmount   string  `json:"toAmount"`
		Slippage   float64 `json:"slippage"`
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

// 获取报价
func (p *Provider) GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error) {
	lifiResp, err := p.fetchQuote(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert to internal Quote format
	quote := p.convertToQuote(req, lifiResp)

	log.Info("LiFi GetQuote success",
		"fromToken", req.FromToken,
		"toToken", req.ToToken,
		"fromAmount", req.Amount,
		"toAmount", quote.ToAmount,
	)

	return quote, nil
}

// 把报价转为标准形式
// 把报价转为标准形式
func (p *Provider) convertToQuote(req *backend.QuoteRequest, resp *LifiQuoteResponse) *backend.Quote {
	// check chain type
	if req.FromChainID == "solana" || req.FromChainID == "solana-mainnet" || req.FromChainID == "solana-devnet" {
		return p.convertToQuoteSolana(req, resp)
	}

	quote := p.convertToQuoteEVM(req, resp)
	return quote
}

// 构建 swap actions
func (p *Provider) BuildSwap(ctx context.Context, quote *backend.Quote, userAddress string) (*backend.BuildSwapResponse, error) {
	// check chain type
	if quote.ChainID == "solana" || quote.ChainID == "solana-mainnet" || quote.ChainID == "solana-devnet" {
		return p.buildSwapSolana(ctx, quote, userAddress)
	}

	swapActions, err := p.buildSwapEVM(ctx, quote, userAddress)
	if err != nil {
		return nil, err
	}
	return swapActions, nil
}

// buildQuoteRequest 构建请求URL和请求体
// 跨链：返回 /advanced/routes URL + JSON请求体
// 同链：返回 /quote URL + nil请求体（GET参数）
func (p *Provider) buildQuoteRequest(req *backend.QuoteRequest) (reqURL string, reqBody []byte, err error) {
	baseURL := p.apiURL
	if baseURL == "" {
		baseURL = "https://li.quest/v1"
	}

	// 若跨链：使用 advanced/routes（POST + JSON体）
	if req.FromChainID != req.ToChainID {
		reqURL = baseURL + "/advanced/routes"

		// 滑点转换（bps转小数，保留4位）
		slippageDecimal := float64(req.SlippageBps) / 10000.0
		//slippageStr := strconv.FormatFloat(slippageDecimal, 'f', 4, 64)

		// 构建跨链请求体
		routesReq := backend.RoutesRequest{
			FromChainId:      req.FromChainID,
			FromAmount:       req.Amount,
			FromTokenAddress: req.FromToken,
			ToChainId:        req.ToChainID,
			ToTokenAddress:   req.ToToken,
			FromAddress:      req.UserAddress,
			Slippage:         slippageDecimal,
			Options: backend.RoutesRequestOptions{
				Bridges: backend.RoutesAllowList{
					Allow: []string{"all"},
				},
				Exchanges: backend.RoutesAllowList{
					Allow: []string{"all"},
				},
				AllowSwitchChain:     false,
				AllowDestinationCall: true,
			},
		}

		// 序列化JSON
		reqBody, err = json.Marshal(routesReq)
		if err != nil {
			return "", nil, fmt.Errorf("marshal routes request failed: %w", err)
		}
		return reqURL, reqBody, nil
	}

	// 若同链：使用 quote（GET + URL参数）
	u, err := url.Parse(baseURL + "/quote")
	if err != nil {
		return "", nil, fmt.Errorf("invalid quote URL: %w", err)
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
	return u.String(), nil, nil
}

// fetchQuote makes HTTP request to LiFi API
func (p *Provider) fetchQuote(ctx context.Context, req *backend.QuoteRequest) (*LifiQuoteResponse, error) {
	// 构建请求URL和请求体
	reqURL, reqBody, err := p.buildQuoteRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build quote request: %w", err)
	}

	// 创建HTTP请求（区分跨链POST/同链GET）
	var httpReq *http.Request
	if req.FromChainID != req.ToChainID {
		// 跨链：POST请求
		httpReq, err = http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(reqBody))
	} else {
		// 同链：GET请求
		httpReq, err = http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	}
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
	var lifiResp LifiQuoteResponse
	if err := json.Unmarshal(body, &lifiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
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
