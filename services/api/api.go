package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/roothash-pay/wallet-services/common/redis"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/provider"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/provider/inch"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/provider/jupiter"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/provider/lifi"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/provider/zerox"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/store"
	"github.com/roothash-pay/wallet-services/services/api/aggregator/utils"
	"github.com/roothash-pay/wallet-services/services/api/validator"
	"github.com/roothash-pay/wallet-services/services/common/chaininfo"
	"github.com/roothash-pay/wallet-services/services/grpc_client/account"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/roothash-pay/wallet-services/common/httputil"
	"github.com/roothash-pay/wallet-services/config"
	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/services/api/routes"
	"github.com/roothash-pay/wallet-services/services/api/service"
	common2 "github.com/roothash-pay/wallet-services/services/common"
)

const (
	HealthPath = "/healthz"

	AdminLoginV1Path  = "/api/v1/admin/login"
	AdminLogoutV1Path = "/api/v1/admin/logout"
)

type APIConfig struct {
	HTTPServer    config.ServerConfig
	MetricsServer config.ServerConfig
}

type API struct {
	router    *chi.Mux
	apiServer *httputil.HTTPServer
	db        *database.DB
	stopped   atomic.Bool
}

func NewApi(ctx context.Context, cfg *config.Config) (*API, error) {
	out := &API{}
	if err := out.initFromConfig(ctx, cfg); err != nil {
		return nil, errors.Join(err, out.Stop(ctx))
	}
	return out, nil
}

func (a *API) initFromConfig(ctx context.Context, cfg *config.Config) error {
	if err := a.initDB(ctx, cfg); err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}
	a.initRouter(ctx, cfg)
	if err := a.startServer(cfg.HttpServer); err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	return nil
}

func (a *API) initRouter(ctx context.Context, cfg *config.Config) {
	allowedOrigins := []string{"http://localhost:8080", "http://127.0.0.1:8080"}
	allowAllOrigins := false
	if cfg.CORSAllowedOrigins != "" {
		if cfg.CORSAllowedOrigins == "*" {
			allowAllOrigins = true
		} else {
			allowedOrigins = parseCORSOrigins(cfg.CORSAllowedOrigins)
		}
	}

	v := new(validator.Validator)

	emailService, err := common2.NewEmailService(&cfg.EmailConfig)
	if err != nil {
		log.Error("failed to create email service", "err", err)
	}

	smsService, err := common2.NewSMSService(&cfg.SMSConfig)
	if err != nil {
		log.Error("failed to create sms service", "err", err)
	}

	authenticatorService := common2.NewAuthenticatorService("PHOENIX")

	var kodoService *common2.KodoService
	if cfg.KodoConfig.AccessKey != "" && cfg.KodoConfig.SecretKey != "" {
		kodoService, err = common2.NewKodoService(&cfg.KodoConfig)
		if err != nil {
			log.Error("failed to create kodo service", "err", err)
		} else {
			log.Info("kodo service initialized successfully")
		}
	}

	var s3Service *common2.S3Service
	if cfg.S3Config.AccessKey != "" && cfg.S3Config.SecretKey != "" {
		s3Service, err = common2.NewS3Service(&cfg.S3Config)
		if err != nil {
			log.Error("failed to create s3 service", "err", err)
		} else {
			log.Info("s3 service initialized successfully")
		}
	}

	svc := service.New(v, a.db, a.db.BackendAdmin, emailService, smsService, authenticatorService, kodoService, s3Service, cfg.JWTSecret, cfg.Domain)
	apiRouter := chi.NewRouter()
	h := routes.NewRoutes(apiRouter, svc)

	apiRouter.Use(middleware.Timeout(time.Second * 12))
	apiRouter.Use(middleware.Recoverer)

	// Add CORS middleware
	corsOptions := cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}

	if allowAllOrigins {
		corsOptions.AllowOriginFunc = func(r *http.Request, origin string) bool {
			return true
		}
	} else {
		corsOptions.AllowedOrigins = allowedOrigins
	}

	apiRouter.Use(cors.Handler(corsOptions))

	apiRouter.Use(middleware.Heartbeat(HealthPath))

	apiRouter.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Warn("NotFoundHandler hit", "path", r.URL.Path, "method", r.Method)
		http.Error(w, "route not found", http.StatusNotFound)
	})

	apiRouter.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		log.Warn("MethodNotAllowedHandler hit", "path", r.URL.Path, "method", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	/*
	 * ============== backend ===============
	 */
	apiRouter.Post(fmt.Sprintf(AdminLoginV1Path), h.AdminLoginHandler)
	apiRouter.Post(fmt.Sprintf(AdminLogoutV1Path), h.AdminLogoutHandler)

	/*
	 * ============== frontend ===============
	 */

	// Initialize Aggregator service and register routes
	aggregatorService, err := a.initAggregatorService(cfg)
	if err != nil {
		log.Error("failed to initialize Aggregator service", "err", err)
	} else if aggregatorService != nil {
		aggregatorRoutes := routes.NewAggregatorRoutes(aggregatorService)
		aggregatorRoutes.RegisterRoutes(apiRouter)
		log.Info("Aggregator routes registered successfully")
	}

	a.router = apiRouter
}

func (a *API) initDB(ctx context.Context, cfg *config.Config) error {
	var initDb *database.DB
	var err error
	if !cfg.SlaveDbEnable {
		initDb, err = database.NewDB(ctx, cfg.MasterDB)
		if err != nil {
			log.Error("failed to connect to master database", "err", err)
			return err
		}
	} else {
		initDb, err = database.NewDB(ctx, cfg.SlaveDB)
		if err != nil {
			log.Error("failed to connect to slave database", "err", err)
			return err
		}
	}
	a.db = initDb
	return nil
}

func (a *API) Start(ctx context.Context) error {
	return nil
}

func (a *API) Stop(ctx context.Context) error {
	var result error
	if a.apiServer != nil {
		if err := a.apiServer.Stop(ctx); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to stop API server: %w", err))
		}
	}
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("failed to close DB: %w", err))
		}
	}
	a.stopped.Store(true)
	log.Info("API service shutdown complete")
	return result
}

func (a *API) startServer(serverConfig config.ServerConfig) error {
	log.Debug("API server listening...", "port", serverConfig.Port)
	addr := net.JoinHostPort(serverConfig.Host, strconv.Itoa(serverConfig.Port))
	srv, err := httputil.StartHTTPServer(addr, a.router)
	if err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	log.Info("API server started", "addr", srv.Addr().String())
	a.apiServer = srv
	return nil
}

func (a *API) Stopped() bool {
	return a.stopped.Load()
}

// parseCORSOrigins parses comma-separated CORS origins string
func parseCORSOrigins(origins string) []string {
	var result []string
	for _, origin := range strings.Split(origins, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// initAggregatorService initializes the aggregator service with all dependencies
func (a *API) initAggregatorService(cfg *config.Config) (*service.AggregatorService, error) {
	// Skip initialization if wallet account address is not configured
	if cfg.AggregatorConfig.WalletAccountAddr == "" {
		log.Warn("Aggregator service not initialized: wallet_account_addr not configured")
		return nil, nil
	}

	// Create wallet account client
	accountClient, err := account.NewWalletAccountClient(cfg.AggregatorConfig.WalletAccountAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet account client: %w", err)
	}

	// Create providers
	var providers []provider.Provider

	// Initialize 0x provider if enabled
	if cfg.AggregatorConfig.EnableProviders["0x"] && cfg.AggregatorConfig.ZeroXAPIURL != "" {
		zeroXProvider := zerox.NewProvider(cfg.AggregatorConfig.ZeroXAPIURL, cfg.AggregatorConfig.ZeroXAPIKey)
		providers = append(providers, zeroXProvider)
		log.Info("0x provider initialized", "url", cfg.AggregatorConfig.ZeroXAPIURL)
	}

	// Initialize 1inch provider if enabled
	if cfg.AggregatorConfig.EnableProviders["1inch"] && cfg.AggregatorConfig.OneInchAPIURL != "" {
		oneInchProvider := inch.NewProvider(cfg.AggregatorConfig.OneInchAPIURL, cfg.AggregatorConfig.OneInchAPIKey)
		providers = append(providers, oneInchProvider)
		log.Info("1inch provider initialized", "url", cfg.AggregatorConfig.OneInchAPIURL)
	}

	// Initialize Jupiter provider if enabled
	if cfg.AggregatorConfig.EnableProviders["jupiter"] && cfg.AggregatorConfig.JupiterAPIURL != "" {
		jupiterProvider := jupiter.NewProvider(cfg.AggregatorConfig.JupiterAPIURL)
		providers = append(providers, jupiterProvider)
		log.Info("Jupiter provider initialized", "url", cfg.AggregatorConfig.JupiterAPIURL)
	}

	// Create Redis client
	var redisClient *redis.Client
	if cfg.RedisConfig.Addr != "" {
		var err error
		redisClient, err = redis.NewClient(&cfg.RedisConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create Redis client: %w", err)
		}
		log.Info("Redis client initialized", "addr", cfg.RedisConfig.Addr)
	} else {
		log.Warn("Redis not configured, using in-memory storage (not recommended for production)")
	}

	// Initialize chain metadata cache
	chainInfoManager := chaininfo.NewManager(
		a.db.BackendChain,
		redisClient,
		cfg.AggregatorConfig.WalletAccountConsumerToken,
		cfg.AggregatorConfig.ChainConsumerTokens,
	)
	if err := chainInfoManager.WarmUp(context.Background()); err != nil {
		log.Warn("Failed to warm up chain info cache", "err", err)
	}

	// Create EVM caller for contract interactions
	evmCaller := utils.NewEVMCaller(accountClient, chainInfoManager)

	// Initialize LiFi provider if enabled
	if cfg.AggregatorConfig.EnableProviders["lifi"] && cfg.AggregatorConfig.LiFiAPIURL != "" {
		lifiProvider := lifi.NewProvider(cfg.AggregatorConfig.LiFiAPIURL, cfg.AggregatorConfig.LiFiAPIKey, evmCaller)
		providers = append(providers, lifiProvider)
		log.Info("LiFi provider initialized", "url", cfg.AggregatorConfig.LiFiAPIURL)
	}

	if len(providers) == 0 {
		log.Warn("Aggregator service not initialized: no providers enabled")
		return nil, nil
	}

	// Create cache stores
	var quoteStore store.QuoteStore
	var swapStore store.SwapStore
	if redisClient != nil {
		quoteStore = store.NewRedisQuoteStore(redisClient.Client)
		swapStore = store.NewRedisSwapStore(redisClient.Client)
		log.Info("Using Redis-based storage")
	} else {
		quoteStore = store.NewInMemoryQuoteStore()
		swapStore = store.NewInMemorySwapStore()
		log.Warn("Cannot connect to Redis, using in-memory storage (data will be lost on restart)")
	}

	// Create validator
	validator := utils.NewValidator()

	// Create aggregator service
	aggregatorService := service.NewAggregatorService(
		providers,
		quoteStore,
		swapStore,
		validator,
		accountClient,
		chainInfoManager,
		a.db,
	)

	log.Info("Aggregator service initialized successfully", "providers", len(providers))
	return aggregatorService, nil
}
