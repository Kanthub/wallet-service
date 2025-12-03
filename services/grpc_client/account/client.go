package account

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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
	Code   int32
	Msg    string
	TxHash string
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

func (c *WalletAccountClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
