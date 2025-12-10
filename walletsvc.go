package relayer_node

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/common/httputil"
	"github.com/roothash-pay/wallet-services/common/redis"
	"github.com/roothash-pay/wallet-services/config"
	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/metrics"
	"github.com/roothash-pay/wallet-services/services/common/chaininfo"
	"github.com/roothash-pay/wallet-services/services/grpc_client/account"
	"github.com/roothash-pay/wallet-services/services/websocket"
	"github.com/roothash-pay/wallet-services/worker/aggregator_task"
	"github.com/roothash-pay/wallet-services/worker/market_task"
)

type WalletServices struct {
	DB                 *database.DB
	metricsServer      *httputil.HTTPServer
	metricsRegistry    *prometheus.Registry
	phoenixMetrics     *metrics.PhoenixMetrics
	marketPriceWorker  *market_task.MarketPriceWorker
	fiatCurrencyWorker *market_task.FiatCurrencyWorker
	txRecordWorker     *aggregator_task.WalletTxRecordWorker
	wsHub              *websocket.Hub
	wsServer           *httputil.HTTPServer
	shutdown           context.CancelCauseFunc
	stopped            atomic.Bool
	chainIdList        []uint64
}

type RpcServerConfig struct {
	GrpcHostname string
	GrpcPort     int
}

func NewWalletServices(ctx context.Context, cfg *config.Config, shutdown context.CancelCauseFunc) (*WalletServices, error) {
	log.Info("New wallet services start️ 🕖")

	metricsRegistry := metrics.NewRegistry()

	PhoenixMetrics := metrics.NewPhoenixMetrics(metricsRegistry, "phoenix")

	out := &WalletServices{
		metricsRegistry: metricsRegistry,
		phoenixMetrics:  PhoenixMetrics,
		shutdown:        shutdown,
	}
	if err := out.initFromConfig(ctx, cfg); err != nil {
		return nil, errors.Join(err, out.Stop(ctx))
	}
	log.Info("New wallet services success🏅️")
	return out, nil
}

func (as *WalletServices) Start(ctx context.Context) error {
	errMpWorker := as.marketPriceWorker.Start()
	if errMpWorker != nil {
		log.Error("start market price worker fail", "err", errMpWorker)
		return errMpWorker
	}

	errFcwWorker := as.fiatCurrencyWorker.Start()
	if errFcwWorker != nil {
		log.Error("start fiat currency worker fail", "err", errFcwWorker)
		return errFcwWorker
	}

	// Start tx record worker if initialized
	if as.txRecordWorker != nil {
		as.txRecordWorker.Start()
		log.Info("Wallet tx record worker started")
	}

	return nil
}

func (as *WalletServices) Stop(ctx context.Context) error {
	var result error

	// Stop tx record worker first
	if as.txRecordWorker != nil {
		log.Info("Stopping wallet tx record worker...")
		as.txRecordWorker.Stop()
	}

	if as.DB != nil {
		if err := as.DB.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to close DB: %w", err))
		}
	}

	if as.metricsServer != nil {
		if err := as.metricsServer.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to close metrics server: %w", err))
		}
	}

	if as.wsHub != nil {
		as.wsHub.CloseAllClients()
	}

	if as.wsServer != nil {
		if err := as.wsServer.Stop(ctx); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to stop WebSocket server: %w", err))
		}
	}

	as.stopped.Store(true)

	log.Info("phoenix services stopped")

	return result
}

func (as *WalletServices) Stopped() bool {
	return as.stopped.Load()
}

func (as *WalletServices) initFromConfig(ctx context.Context, cfg *config.Config) error {
	if err := as.initDB(ctx, cfg.MasterDB); err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}

	as.wsHub = websocket.NewHub()
	go as.wsHub.Run()

	if err := as.startWebSocketServer(cfg.WebsocketServer); err != nil {
		return fmt.Errorf("failed to start web socket server: %w", err)
	}

	if err := as.initWorker(cfg); err != nil {
		return fmt.Errorf("failed to init worker processor: %w", err)
	}

	err := as.startMetricsServer(cfg.MetricsServer)
	if err != nil {
		log.Error("start metrics server fail", "err", err)
		return err
	}
	return nil
}

func (as *WalletServices) startWebSocketServer(serverConfig config.ServerConfig) error {
	addr := net.JoinHostPort(serverConfig.Host, strconv.Itoa(serverConfig.Port))

	wsRouter := chi.NewRouter()
	wsRouter.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWebSocket(as.wsHub, w, r)
	})

	srv, err := httputil.StartHTTPServer(addr, wsRouter)
	if err != nil {
		return fmt.Errorf("failed to start WebSocket server: %w", err)
	}
	log.Info("WebSocket server started", "addr", srv.Addr().String())
	as.wsServer = srv
	return nil
}

func (as *WalletServices) initDB(ctx context.Context, cfg config.DBConfig) error {
	db, err := database.NewDB(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	as.DB = db
	log.Info("Init database success")
	return nil
}

func (as *WalletServices) initWorker(config *config.Config) error {
	mwConfig := &market_task.MarketPriceWorkerConfig{
		LoopInterval: time.Second * 5,
	}
	marketPriceWorker, err := market_task.NewMarketPriceWorker(as.DB, mwConfig, as.wsHub, as.shutdown)
	if err != nil {
		log.Error("new market price worker fail", "err", err)
		return err
	}
	as.marketPriceWorker = marketPriceWorker

	fiatCurrencyWorkerConfig := &market_task.FiatCurrencyWorkerConfig{
		LoopInterval: time.Second * 5,
	}
	fiatCurrencyWorker, err := market_task.NewFiatCurrencyWorker(as.DB, fiatCurrencyWorkerConfig, as.wsHub, as.shutdown)
	if err != nil {
		log.Error("new fiat currency worker fail", "err", err)
		return err
	}
	as.fiatCurrencyWorker = fiatCurrencyWorker

	// Initialize wallet tx record worker if aggregator is enabled
	if config.AggregatorConfig.WalletAccountAddr != "" {
		accountClient, err := account.NewWalletAccountClient(config.AggregatorConfig.WalletAccountAddr)
		if err != nil {
			log.Warn("failed to create wallet account client for tx worker", "err", err)
		} else {
			var chainInfoRedis *redis.Client
			if config.RedisConfig.Addr != "" {
				chainInfoRedis, err = redis.NewClient(&config.RedisConfig)
				if err != nil {
					log.Warn("failed to init redis for chain info cache", "err", err)
				}
			}
			chainInfoManager := chaininfo.NewManager(
				as.DB.BackendChain,
				chainInfoRedis,
				config.AggregatorConfig.WalletAccountConsumerToken,
				config.AggregatorConfig.ChainConsumerTokens,
			)
			if warmErr := chainInfoManager.WarmUp(context.Background()); warmErr != nil {
				log.Warn("failed to warm up chain info cache for worker", "err", warmErr)
			}

			txWorkerConfig := aggregator_task.WalletTxRecordWorkerConfig{
				ScanInterval:         10,   // 10 seconds
				LastCheckedThreshold: 5,    // 5 seconds
				BatchSize:            100,  // 100 records per batch
				Concurrency:          10,   // 10 concurrent workers
				TimeoutThreshold:     3600, // 1 hour timeout
			}
			txRecordWorker := aggregator_task.NewWalletTxRecordWorker(
				as.DB.BackendWalletTxRecord,
				accountClient,
				chainInfoManager,
				txWorkerConfig,
			)
			as.txRecordWorker = txRecordWorker
			log.Info("Wallet tx record worker initialized",
				"scanInterval", txWorkerConfig.ScanInterval,
				"concurrency", txWorkerConfig.Concurrency,
				"timeout", txWorkerConfig.TimeoutThreshold)
		}
	} else {
		log.Info("Wallet tx record worker not initialized: wallet_account_addr not configured")
	}

	return nil
}

func (as *WalletServices) startMetricsServer(cfg config.ServerConfig) error {
	srv, err := metrics.StartServer(as.metricsRegistry, cfg.Host, cfg.Port)
	if err != nil {
		return fmt.Errorf("metrics server failed to start: %w", err)
	}
	as.metricsServer = srv
	log.Info("metrics server started", "port", cfg.Port, "addr", srv.Addr())
	return nil
}
