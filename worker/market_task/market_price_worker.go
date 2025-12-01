package market_task

import (
	"context"
	"fmt"

	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/common/tasks"
	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/services/websocket"
)

type MarketPriceWorkerConfig struct {
	LoopInterval time.Duration
}

type MarketPriceWorker struct {
	db             *database.DB
	wConf          *MarketPriceWorkerConfig
	wsHub          *websocket.Hub
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewMarketPriceWorker(db *database.DB, wConf *MarketPriceWorkerConfig, wsHub *websocket.Hub, shutdown context.CancelCauseFunc) (*MarketPriceWorker, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &MarketPriceWorker{
		db:             db,
		wConf:          wConf,
		wsHub:          wsHub,
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{
			HandleCrit: func(err error) {
				shutdown(fmt.Errorf("critical error in worker handle processor: %w", err))
			},
		},
	}, nil
}

func (mpw *MarketPriceWorker) Close() error {
	mpw.resourceCancel()
	return mpw.tasks.Wait()
}

func (mpw *MarketPriceWorker) Start() error {
	workerTicker := time.NewTicker(mpw.wConf.LoopInterval)
	mpw.tasks.Go(func() error {
		for range workerTicker.C {
			log.Info("==== star ======")
		}
		return nil
	})

	mpw.tasks.Go(func() error {
		for range workerTicker.C {
			log.Info("==== star ======")
		}
		return nil
	})
	return nil
}
