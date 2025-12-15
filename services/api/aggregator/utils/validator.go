package utils

import (
	"fmt"
	"math/big"
	"strings"
)

// Validator provides validation utilities for DEX operations
type Validator struct {
	whitelistedRouters  map[string]bool
	whitelistedSpenders map[string]bool
	maxValueWei         *big.Int
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	// TODO: Load from configuration
	return &Validator{
		whitelistedRouters: map[string]bool{
			"0x0000000000000000000000000000000000000000": true, // Placeholder
		},
		whitelistedSpenders: map[string]bool{
			"0x0000000000000000000000000000000000000000": true, // Placeholder
		},
		maxValueWei: new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)), // 100 ETH
	}
}

// ValidateChainID validates that the chain ID is supported
func (v *Validator) ValidateChainID(chainID string) error {
	// TODO: Implement chain ID validation
	supportedChains := map[string]bool{
		// Ethereum & L2s
		"1":        true, // Ethereum Mainnet
		"5":        true, // Goerli
		"10":       true, // Optimism
		"420":      true, // Optimism Goerli
		"42161":    true, // Arbitrum One
		"42170":    true, // Arbitrum Nova
		"421613":   true, // Arbitrum Goerli
		"11155111": true, // Sepolia
		// BNB Chain
		"56": true, // BNB Chain
		"97": true, // BNB Testnet
		// Polygon
		"137":   true, // Polygon
		"80001": true, // Polygon Mumbai
		// Avalanche
		"43114": true, // Avalanche C-Chain
		"43113": true, // Avalanche Fuji
		// Fantom
		"250":  true, // Fantom
		"4002": true, // Fantom Testnet
		// zkSync
		"324": true, // zkSync Era
		"280": true, // zkSync Testnet
		// Base
		"8453":  true, // Base
		"84531": true, // Base Goerli
		// Linea
		"59144": true, // Linea Mainnet
		"59140": true, // Linea Sepolia
		// Scroll
		"534352": true, // Scroll
		"534351": true, // Scroll Sepolia
		// Misc
		"solana-mainnet": true,
		"solana-devnet":  true,
	}

	if !supportedChains[chainID] {
		return fmt.Errorf("unsupported chain ID: %s", chainID)
	}

	return nil
}

// ValidateRouter validates that the router address is whitelisted
func (v *Validator) ValidateRouter(router string) error {
	router = strings.ToLower(router)
	if !v.whitelistedRouters[router] {
		return fmt.Errorf("router not whitelisted: %s", router)
	}
	return nil
}

// ValidateSpender validates that the spender address is whitelisted
func (v *Validator) ValidateSpender(spender string) error {
	spender = strings.ToLower(spender)
	if !v.whitelistedSpenders[spender] {
		return fmt.Errorf("spender not whitelisted: %s", spender)
	}
	return nil
}

// ValidateValue validates that the transaction value is within limits
func (v *Validator) ValidateValue(valueWei *big.Int) error {
	if valueWei.Cmp(v.maxValueWei) > 0 {
		return fmt.Errorf("value exceeds maximum: %s > %s", valueWei.String(), v.maxValueWei.String())
	}
	return nil
}
