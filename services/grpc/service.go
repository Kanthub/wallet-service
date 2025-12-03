package grpc

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/proto/wallet"
)

const MaxRecvMessageSize = 1024 * 1024 * 30000

type PhoenixRpcConfig struct {
	Host string
	Port int
}

type PhoenixRpcService struct {
	*PhoenixRpcConfig
	db *database.DB
	wallet.UnimplementedPhoenixServicesServer
	stopped atomic.Bool
}

func NewPhoenixRpcService(conf *PhoenixRpcConfig, db *database.DB) (*PhoenixRpcService, error) {
	return &PhoenixRpcService{
		PhoenixRpcConfig: conf,
		db:               db,
	}, nil
}

func (prs *PhoenixRpcService) Start(ctx context.Context) error {
	go func(prs *PhoenixRpcService) {
		rpcAddr := fmt.Sprintf("%s:%d", prs.PhoenixRpcConfig.Host, prs.PhoenixRpcConfig.Port)
		listener, err := net.Listen("tcp", rpcAddr)
		if err != nil {
			log.Error("Could not start tcp listener. ")
		}

		opt := grpc.MaxRecvMsgSize(MaxRecvMessageSize)

		gs := grpc.NewServer(
			opt,
			grpc.ChainUnaryInterceptor(
				nil,
			),
		)

		reflection.Register(gs)
		wallet.RegisterPhoenixServicesServer(gs, prs)

		log.Info("grpc info", "addr", listener.Addr())

		if err := gs.Serve(listener); err != nil {
			log.Error("start rpc server fail", "err", err)
		}
	}(prs)
	return nil
}

func (prs *PhoenixRpcService) Stop(ctx context.Context) error {
	prs.stopped.Store(true)
	return nil
}

func (prs *PhoenixRpcService) Stopped() bool {
	return prs.stopped.Load()
}
