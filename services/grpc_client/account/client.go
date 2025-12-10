package account

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
	pb "github.com/roothash-pay/wallet-services/proto/account"
)

type WalletAccountClient struct {
	conn   *grpc.ClientConn
	client pb.WalletAccountServiceClient
	addr   string
}

type SendTxParams struct {
	ConsumerToken string
	Chain         string
	Coin          string
	Network       string
	RawTx         string
}

type SendTxResult struct {
	Code   common.ReturnCode
	Msg    string
	TxHash string
}

type CallContractParams struct {
	ConsumerToken   string
	Chain           string
	Network         string
	ContractAddress string
	Data            string
}

type CallContractResult struct {
	Code   common.ReturnCode
	Msg    string
	Result string
}

func NewWalletAccountClient(addr string) (*WalletAccountClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to wallet-chain-account service: %w", err)
	}

	client := pb.NewWalletAccountServiceClient(conn)

	return &WalletAccountClient{
		conn:   conn,
		client: client,
		addr:   addr,
	}, nil
}

func (c *WalletAccountClient) SendTx(ctx context.Context, params SendTxParams) (*SendTxResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req := &pb.SendTxRequest{
		ConsumerToken: params.ConsumerToken,
		Chain:         params.Chain,
		Coin:          params.Coin,
		Network:       params.Network,
		RawTx:         params.RawTx,
	}

	resp, err := c.client.SendTx(ctx, req)
	if err != nil {
		log.Error("SendTx RPC failed", "err", err)
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return &SendTxResult{
		Code:   resp.Code,
		Msg:    resp.Msg,
		TxHash: resp.TxHash,
	}, nil
}

func (c *WalletAccountClient) CallContract(ctx context.Context, params CallContractParams) (*CallContractResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req := &pb.CallContractRequest{
		ConsumerToken:   params.ConsumerToken,
		Chain:           params.Chain,
		Network:         params.Network,
		ContractAddress: params.ContractAddress,
		Data:            params.Data,
	}

	resp, err := c.client.CallContract(ctx, req)
	if err != nil {
		log.Error("CallContract RPC failed", "err", err)
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	if resp.Code != common.ReturnCode_SUCCESS {
		return nil, fmt.Errorf("call contract failed: %s", resp.Msg)
	}

	return &CallContractResult{
		Code:   resp.Code,
		Msg:    resp.Msg,
		Result: resp.Result,
	}, nil
}

// TxInfo represents transaction information
type TxInfo struct {
	Hash            string
	Status          pb.TxStatus
	Height          string
	From            string
	To              string
	Value           string
	Fee             string
	ContractAddress string
	Datetime        string
}

// GetTxByHash queries transaction details by hash
func (c *WalletAccountClient) GetTxByHash(ctx context.Context, consumerToken, chain, coin, network, txHash string) (*TxInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := &pb.TxHashRequest{
		ConsumerToken: consumerToken,
		Chain:         chain,
		Coin:          coin,
		Network:       network,
		Hash:          txHash,
	}

	resp, err := c.client.GetTxByHash(ctx, req)
	if err != nil {
		log.Error("GetTxByHash RPC failed", "err", err)
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	if resp.Code != common.ReturnCode_SUCCESS {
		return nil, fmt.Errorf("get transaction failed: %s", resp.Msg)
	}

	if resp.Tx == nil {
		return nil, fmt.Errorf("transaction not found")
	}

	return &TxInfo{
		Hash:            resp.Tx.Hash,
		Status:          resp.Tx.Status,
		Height:          resp.Tx.Height,
		From:            resp.Tx.From,
		To:              resp.Tx.To,
		Value:           resp.Tx.Value,
		Fee:             resp.Tx.Fee,
		ContractAddress: resp.Tx.ContractAddress,
		Datetime:        resp.Tx.Datetime,
	}, nil
}

func (c *WalletAccountClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
