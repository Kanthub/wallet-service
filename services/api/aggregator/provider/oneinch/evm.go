package oneinch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// SwapResponse represents the response from 1inch swap API
type SwapResponse struct {
	DstAmount string `json:"dstAmount"`
	Tx        struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Data     string `json:"data"`
		Value    string `json:"value"`
		Gas      int64  `json:"gas"`
		GasPrice string `json:"gasPrice"`
	} `json:"tx"`
}

// convertToQuoteEVM converts 1inch response to internal Quote format
func (p *Provider) convertToQuoteEVM(req *backend.QuoteRequest, quoteResp *OneInchQuoteResponse, raw string) *backend.Quote {
	return &backend.Quote{
		Provider:    p.Name(),
		ChainType:   backend.ChainTypeEVM,
		ChainID:     req.FromChainID,
		FromToken:   req.FromToken,
		ToToken:     req.ToToken,
		FromAmount:  req.Amount,
		ToAmount:    quoteResp.DstAmount,
		GasEstimate: fmt.Sprintf("%d", quoteResp.Gas),
		Fees:        "0", // 1inch API doesn't return fees directly in quote
		Spender:     "",  // 1inch router address usually returned in /approve/spender or handled in swap
		Router:      "",  // filled from swap response or config
		Raw:         raw,
	}
}

// buildSwapEVM builds the swap transaction for EVM chains
func (p *Provider) buildSwapEVM(ctx context.Context, quote *backend.Quote, userAddress string) (*backend.BuildSwapResponse, error) {
	log.Info("1inch BuildSwap called", "chainID", quote.ChainID, "userAddress", userAddress)

	// Call 1inch Swap API
	// /swap/v6.0/{chainID}/swap
	u, err := url.Parse(fmt.Sprintf("%s/%s/swap", p.apiURL, quote.ChainID))
	if err != nil {
		return nil, fmt.Errorf("invalid API URL: %w", err)
	}

	q := u.Query()
	q.Set("src", quote.FromToken)
	q.Set("dst", quote.ToToken)
	q.Set("amount", quote.FromAmount)
	q.Set("from", userAddress)
	q.Set("slippage", "1")           // Default 1% slippage for now, or parse from QuoteRequest ctx if passed
	q.Set("disableEstimate", "true") // Disable estimate to avoid failure if approval not yet set

	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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
		return nil, fmt.Errorf("1inch Swap API error (status %d): %s", resp.StatusCode, string(body))
	}

	var swapResp SwapResponse
	if err := json.Unmarshal(body, &swapResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 1inch swap API returns the transaction to execute (swap)
	// It assumes approval is done. If we receive a swap tx, we can assume it's the swap action.
	// NOTE: Ideally we should check allowance and add APPROVE action if needed.
	// Since we disabled estimate, we get the swap tx regardless of allowance.
	// The caller (client) or a separate check should handle approval.
	// For consistency with LiFi, we should probably check allowance here, but 1inch API is slightly different.
	// Current implementation: Return SWAP action only.

	action := &backend.Action{
		ActionType: backend.ActionTypeSwap,
		ChainID:    quote.ChainID,
		SigningPayload: &backend.SigningPayload{
			To:      swapResp.Tx.To,
			Data:    swapResp.Tx.Data,
			Value:   swapResp.Tx.Value,
			Gas:     fmt.Sprintf("%d", swapResp.Tx.Gas),
			ChainID: quote.ChainID,
		},
		Description: fmt.Sprintf("Swap %s to %s via 1inch",
			getTokenSymbol(quote.FromToken),
			getTokenSymbol(quote.ToToken)),
	}

	return &backend.BuildSwapResponse{
		Actions: []*backend.Action{action},
	}, nil
}
