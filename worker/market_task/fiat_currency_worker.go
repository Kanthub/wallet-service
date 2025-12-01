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

type FiatCurrencyWorkerConfig struct {
	LoopInterval time.Duration
}

type FiatCurrencyWorker struct {
	db             *database.DB
	wConf          *FiatCurrencyWorkerConfig
	wsHub          *websocket.Hub
	resourceCtx    context.Context
	resourceCancel context.CancelFunc
	tasks          tasks.Group
}

func NewFiatCurrencyWorker(db *database.DB, wConf *FiatCurrencyWorkerConfig, wsHub *websocket.Hub, shutdown context.CancelCauseFunc) (*FiatCurrencyWorker, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &FiatCurrencyWorker{
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

func (mpw *FiatCurrencyWorker) Close() error {
	mpw.resourceCancel()
	return mpw.tasks.Wait()
}

func (mpw *FiatCurrencyWorker) Start() error {
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
