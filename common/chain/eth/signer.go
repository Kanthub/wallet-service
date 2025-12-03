package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

type TxParams struct {
	To       string
	Value    *big.Int
	GasLimit uint64
	Data     []byte
}

type Signer struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    common.Address
	chainID    *big.Int
}

func NewSigner(privateKeyHex string, chainID *big.Int) (*Signer, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Signer{
		privateKey: privateKey,
		publicKey:  publicKeyECDSA,
		address:    address,
		chainID:    chainID,
	}, nil
}

func (s *Signer) GetAddress() string {
	return s.address.Hex()
}

func (s *Signer) SignTransaction(ctx context.Context, rpcURL string, params TxParams) (string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(ctx, s.address)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	toAddress := common.HexToAddress(params.To)

	tx := types.NewTransaction(
		nonce,
		toAddress,
		params.Value,
		params.GasLimit,
		gasPrice,
		params.Data,
	)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(s.chainID), s.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("failed to marshal transaction: %w", err)
	}

	rawTxHex := hexutil.Encode(rawTxBytes)
	log.Info("Transaction signed", "hash", signedTx.Hash().Hex(), "rawTx", rawTxHex)

	return rawTxHex, nil
}
