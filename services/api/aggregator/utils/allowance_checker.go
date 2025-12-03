package utils

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

// AllowanceChecker checks ERC20 token allowances
type AllowanceChecker struct {
	client *ethclient.Client
}

// NewAllowanceChecker creates a new allowance checker
func NewAllowanceChecker(rpcURL string) (*AllowanceChecker, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	return &AllowanceChecker{
		client: client,
	}, nil
}

// CheckAllowance checks if the user has sufficient allowance for the spender
func (c *AllowanceChecker) CheckAllowance(ctx context.Context, tokenAddress, ownerAddress, spenderAddress string, requiredAmount *big.Int) (*big.Int, error) {
	// TODO: Implement actual ERC20 allowance check
	// This would involve:
	// 1. Creating an ERC20 contract instance
	// 2. Calling the allowance(owner, spender) function
	// 3. Comparing with requiredAmount

	// Placeholder: return zero allowance
	return big.NewInt(0), nil
}

// Close closes the Ethereum client connection
func (c *AllowanceChecker) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
