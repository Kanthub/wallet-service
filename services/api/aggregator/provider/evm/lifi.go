package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// LiFiProvider implements the Provider interface for LiFi
type LiFiProvider struct {
	apiURL     string
	apiKey     string
	httpClient *http.Client
}

// LiFiQuoteResponse represents the response from LiFi quote API
type LiFiQuoteResponse struct {
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

// NewLiFiProvider creates a new LiFi provider
func NewLiFiProvider(apiURL, apiKey string) *LiFiProvider {
	return &LiFiProvider{
		apiURL: apiURL,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetQuote fetches a quote from LiFi
func (p *LiFiProvider) GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error) {
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
	var lifiResp LiFiQuoteResponse
	if err := json.Unmarshal(body, &lifiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to internal Quote format
	quote := p.convertToQuote(req, &lifiResp)

	log.Info("LiFi GetQuote success",
		"fromToken", req.FromToken,
		"toToken", req.ToToken,
		"fromAmount", req.Amount,
		"toAmount", quote.ToAmount,
	)

	return quote, nil
}

// Name returns the provider name
func (p *LiFiProvider) Name() string {
	return "lifi"
}

// SupportedChainType returns EVM
func (p *LiFiProvider) SupportedChainType() backend.ChainType {
	return backend.ChainTypeEVM
}

// buildQuoteURL constructs the LiFi quote API URL with query parameters
func (p *LiFiProvider) buildQuoteURL(req *backend.QuoteRequest) (string, error) {
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

// convertToQuote converts LiFi response to internal Quote format
func (p *LiFiProvider) convertToQuote(req *backend.QuoteRequest, lifiResp *LiFiQuoteResponse) *backend.Quote {
	// Calculate total gas estimate
	gasEstimate := "0"
	if len(lifiResp.Estimate.GasCosts) > 0 {
		gasEstimate = lifiResp.Estimate.GasCosts[0].Estimate
	}

	// Get spender (approval address) and router (transaction to address)
	spender := lifiResp.Estimate.ApprovalAddress
	router := lifiResp.TransactionRequest.To

	// Marshal raw response for debugging
	rawJSON, _ := json.Marshal(lifiResp)

	return &backend.Quote{
		Provider:    p.Name(),
		ChainType:   backend.ChainTypeEVM,
		ChainID:     req.FromChainID,
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		FromAmount:  lifiResp.Estimate.FromAmount,
		ToAmount:    lifiResp.Estimate.ToAmount,
		GasEstimate: gasEstimate,
		Spender:     spender,
		Router:      router,
		Raw:         string(rawJSON),
	}
}
