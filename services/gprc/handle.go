package grpc

import (
	"context"

	multimarket_grpc "github.com/roothash-pay/wallet-services/proto/multimarket"
)

func (prs *PhoenixRpcService) QueryFeeByChainId(ctx context.Context, in *multimarket_grpc.QueryFeeRequest) (*multimarket_grpc.QueryFeeResponse, error) {
	return &multimarket_grpc.QueryFeeResponse{
		ReturnCode: 100,
		Message:    "query fee success",
	}, nil
}
