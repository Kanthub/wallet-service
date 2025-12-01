package grpc

import (
	"context"

	"github.com/roothash-pay/wallet-services/proto/wallet"
)

func (prs *PhoenixRpcService) QueryFeeByChainId(ctx context.Context, in *wallet.QueryFeeRequest) (*wallet.QueryFeeResponse, error) {
	return &wallet.QueryFeeResponse{
		ReturnCode: 100,
		Message:    "query fee success",
	}, nil
}
