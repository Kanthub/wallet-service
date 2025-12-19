package market_task

import (
	"context"
	"fmt"

	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/common/tasks"
	"github.com/roothash-pay/wallet-services/config"
	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/services/market/cache"
	"github.com/roothash-pay/wallet-services/services/market/provider"
	"github.com/roothash-pay/wallet-services/services/market/resolver"
	"github.com/roothash-pay/wallet-services/services/market/service"
	"github.com/roothash-pay/wallet-services/services/websocket"
)

type MarketPriceWorker struct {
	db              *database.DB
	wConf           *config.MarketPriceWorkerConfig
	wsHub           *websocket.Hub
	marketCollector *service.MarketCollector
	resourceCtx     context.Context
	resourceCancel  context.CancelFunc
	tasks           tasks.Group
}

func NewMarketPriceWorker(db *database.DB, marketCache cache.Cache, wConf *config.MarketPriceWorkerConfig, wsHub *websocket.Hub, shutdown context.CancelCauseFunc) (*MarketPriceWorker, error) {
	resCtx, resCancel := context.WithCancel(context.Background())

	// 1. providers
	providers := []provider.Provider{
		provider.NewBinanceProvider(),
		provider.NewOkxProvider(),
		// CoinGecko（免费）
		provider.NewCoinGeckoProvider(nil),
		provider.NewDefiLlamaProvider(),
		// CMC（有 key 才启用）
		// provider.NewCMCProvider(cfg.CMC.APIKey, []string{"BTC","ETH","SOL"}),
		// DEX（The Graph）
		provider.NewUniswapV3GraphProvider(50),
		provider.NewPancakeSwapV2GraphProvider(50),

		// CoinAPI（有 key 才启用）
		// provider.NewCoinAPIProvider(""),
		// CryptoCompare（免费）
		provider.NewCryptoCompareProvider("087289638929828f202e6c05f93cb76441831720bd7778785105b94ab2ecf4b6"),
	}

	// 2. resolver
	resolver := resolver.NewInfoLevelResolver()

	// 3. cache (from parameter)

	// 4. Market Collector
	collector := service.NewMarketCollector(
		providers,
		resolver,
		marketCache,
	)
	return &MarketPriceWorker{
		db:              db,
		wConf:           wConf,
		wsHub:           wsHub,
		marketCollector: collector,
		resourceCtx:     resCtx,
		resourceCancel:  resCancel,
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
		for {
			select {
			case <-mpw.resourceCtx.Done():
				return nil
			case <-workerTicker.C:
				ctx, cancel := context.WithTimeout(
					mpw.resourceCtx,
					20*time.Second,
				)
				finalQuotes, err := mpw.marketCollector.Collect(ctx)
				if err != nil {
					log.Warn("market collect failed", "err", err)
					cancel()
					continue
				}

				if len(finalQuotes) > 0 {
					mpw.wsHub.Broadcast("market:price", finalQuotes)
				}

				cancel()
			}
		}
	})

	// mpw.tasks.Go(func() error {
	// 	for range workerTicker.C {
	// 		log.Info("==== star ======")
	// 	}
	// 	return nil
	// })
	return nil
}
