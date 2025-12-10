package service

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

// TestSaveSwapHistory tests the swap history saving functionality
func TestSaveSwapHistory(t *testing.T) {
	// This is a unit test example showing how saveSwapHistory works
	// In a real test, you would mock the database and accountClient

	// Create a mock swap
	swap := &backend.Swap{
		SwapID:      uuid.New().String(),
		QuoteID:     uuid.New().String(),
		UserAddress: "0x1234567890123456789012345678901234567890",
		Status:      backend.TxStatusSuccess, // 3 = SUCCESS
		Steps: []*backend.Step{
			{
				StepIndex:  0,
				ActionType: backend.ActionTypeApprove,
				TxHash:     "0xapprove123...",
				Status:     backend.TxStatusSuccess, // 3 = SUCCESS
			},
			{
				StepIndex:  1,
				ActionType: backend.ActionTypeSwap,
				TxHash:     "0xswap456...",
				Status:     backend.TxStatusSuccess, // 3 = SUCCESS
			},
		},
		CreatedAt: time.Now().Add(-5 * time.Minute),
		UpdatedAt: time.Now(),
	}

	// Create a mock quote
	quote := &backend.Quote{
		Provider:    "lifi",
		ChainType:   backend.ChainTypeEVM,
		ChainID:     "1",
		FromToken:   "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", // USDC
		ToToken:     "0x0000000000000000000000000000000000000000", // ETH
		FromAmount:  "1000000000",                                 // 1000 USDC (6 decimals)
		ToAmount:    "500000000000000000",                         // 0.5 ETH (18 decimals)
		GasEstimate: "150000",
		Router:      "0xrouter123...",
	}

	t.Logf("Mock swap created: %s", swap.SwapID)
	t.Logf("Mock quote: %s %s -> %s %s", quote.FromAmount, quote.FromToken, quote.ToAmount, quote.ToToken)

	// In a real test, you would:
	// 1. Create a mock database
	// 2. Create a mock accountClient
	// 3. Call saveSwapHistory
	// 4. Verify the record was saved correctly

	t.Log("Test completed - saveSwapHistory would save this swap to wallet_tx_record table")
}

// TestParseAmount tests the amount parsing functionality
func TestParseAmount(t *testing.T) {
	s := &AggregatorService{}

	tests := []struct {
		name     string
		input    string
		expected int64
		hasError bool
	}{
		{
			name:     "Valid small amount",
			input:    "1000000",
			expected: 1000000,
			hasError: false,
		},
		{
			name:     "Valid large amount",
			input:    "1000000000000000000", // 1 ETH in wei
			expected: 1000000000000000000,
			hasError: false,
		},
		{
			name:     "Invalid format",
			input:    "invalid",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := s.parseAmount(tt.input)
			if tt.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestGetTokenSymbol tests the token symbol extraction
func TestGetTokenSymbol(t *testing.T) {
	s := &AggregatorService{}

	tests := []struct {
		name     string
		address  string
		expected string
	}{
		{
			name:     "ETH",
			address:  "0x0000000000000000000000000000000000000000",
			expected: "ETH",
		},
		{
			name:     "USDC",
			address:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			expected: "USDC",
		},
		{
			name:     "USDT",
			address:  "0xdAC17F958D2ee523a2206206994597C13D831ec7",
			expected: "USDT",
		},
		{
			name:     "Unknown token",
			address:  "0x1234567890123456789012345678901234567890",
			expected: "0x1234...7890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.getTokenSymbol(tt.address)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
