package utils

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/roothash-pay/wallet-services/services/common/chaininfo"
	"github.com/roothash-pay/wallet-services/services/grpc_client/account"
)

// EVMCaller provides utilities for calling EVM contract methods
type EVMCaller struct {
	accountClient *account.WalletAccountClient
	chainInfo     chaininfo.Provider
}

// NewEVMCaller creates a new EVM caller
func NewEVMCaller(accountClient *account.WalletAccountClient, chainInfo chaininfo.Provider) *EVMCaller {
	return &EVMCaller{
		accountClient: accountClient,
		chainInfo:     chainInfo,
	}
}

func (c *EVMCaller) getChainInfo(ctx context.Context, chainID string) (*chaininfo.Info, error) {
	if c.chainInfo == nil {
		return nil, fmt.Errorf("chain info provider not configured")
	}
	return c.chainInfo.Get(ctx, chainID)
}

// GetERC20Allowance checks the ERC20 allowance for a spender
// Returns the current allowance amount
func (c *EVMCaller) GetERC20Allowance(
	ctx context.Context,
	chainID string,
	tokenAddress string,
	ownerAddress string,
	spenderAddress string,
) (*big.Int, error) {
	callData, err := encodeERC20Allowance(ownerAddress, spenderAddress)
	if err != nil {
		return nil, err
	}

	resultHex, err := c.CallContract(ctx, chainID, tokenAddress, callData)
	if err != nil {
		return nil, err
	}

	allowance, err := DecodeUint256(resultHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode allowance: %w", err)
	}

	return allowance, nil
}

// encodeERC20Allowance encodes the ERC20 allowance function call
func encodeERC20Allowance(owner, spender string) (string, error) {
	// ERC20 allowance function selector
	selector := "dd62ed3e"

	// Remove 0x prefix if present
	owner = strings.TrimPrefix(owner, "0x")
	spender = strings.TrimPrefix(spender, "0x")

	// Pad addresses to 32 bytes (64 hex chars)
	ownerPadded := fmt.Sprintf("%064s", owner)
	spenderPadded := fmt.Sprintf("%064s", spender)

	return "0x" + selector + ownerPadded + spenderPadded, nil
}

// DecodeUint256 decodes a uint256 value from hex string
func DecodeUint256(hexData string) (*big.Int, error) {
	// Remove 0x prefix
	hexData = strings.TrimPrefix(hexData, "0x")

	// Decode hex to bytes
	data, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex: %w", err)
	}

	// Convert to big.Int
	value := new(big.Int).SetBytes(data)
	return value, nil
}

// CallContract is a generic method to call any contract function
// This can be extended to support other EVM contract calls
func (c *EVMCaller) CallContract(
	ctx context.Context,
	chainID string,
	contractAddress string,
	callData string,
) (string, error) {
	if c.accountClient == nil {
		return "", fmt.Errorf("wallet account client not configured")
	}

	info, err := c.getChainInfo(ctx, chainID)
	if err != nil {
		return "", err
	}

	result, err := c.accountClient.CallContract(ctx, account.CallContractParams{
		ConsumerToken:   info.ConsumerToken,
		Chain:           info.WalletChain,
		Network:         info.WalletNetwork,
		ContractAddress: contractAddress,
		Data:            callData,
	})
	if err != nil {
		return "", err
	}

	return result.Result, nil
}

// Helper function to create ERC20 ABI for common operations
func getERC20ABI() (abi.ABI, error) {
	// Minimal ERC20 ABI for allowance and balanceOf
	const erc20ABI = `[
		{
			"constant": true,
			"inputs": [
				{"name": "owner", "type": "address"},
				{"name": "spender", "type": "address"}
			],
			"name": "allowance",
			"outputs": [{"name": "", "type": "uint256"}],
			"type": "function"
		},
		{
			"constant": true,
			"inputs": [{"name": "account", "type": "address"}],
			"name": "balanceOf",
			"outputs": [{"name": "", "type": "uint256"}],
			"type": "function"
		}
	]`

	return abi.JSON(strings.NewReader(erc20ABI))
}

// GetERC20Balance gets the token balance for an address
func (c *EVMCaller) GetERC20Balance(
	ctx context.Context,
	chainID string,
	tokenAddress string,
	accountAddress string,
) (*big.Int, error) {
	// Encode balanceOf call
	erc20ABI, err := getERC20ABI()
	if err != nil {
		return nil, fmt.Errorf("failed to parse ERC20 ABI: %w", err)
	}

	account := common.HexToAddress(accountAddress)
	callData, err := erc20ABI.Pack("balanceOf", account)
	if err != nil {
		return nil, fmt.Errorf("failed to encode balanceOf call: %w", err)
	}

	callDataHex := "0x" + hex.EncodeToString(callData)

	resultHex, err := c.CallContract(ctx, chainID, tokenAddress, callDataHex)
	if err != nil {
		return nil, err
	}

	balance, err := DecodeUint256(resultHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode balance: %w", err)
	}

	return balance, nil
}
