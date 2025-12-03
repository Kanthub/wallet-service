package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	relayer_node "github.com/roothash-pay/wallet-services"
	"github.com/roothash-pay/wallet-services/common/chain/eth"
	"github.com/roothash-pay/wallet-services/common/cliapp"
	"github.com/roothash-pay/wallet-services/common/opio"
	"github.com/roothash-pay/wallet-services/config"
	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/services/api"
	grpc "github.com/roothash-pay/wallet-services/services/gprc"
	"github.com/roothash-pay/wallet-services/services/grpc_client/account"
)

var (
	ConfigFlag = &cli.StringFlag{
		Name:    "config",
		Value:   "./wallet-services-config.local.yaml",
		Aliases: []string{"c"},
		Usage:   "path to config file",
		EnvVars: []string{"WALLET_SERVICES_CONFIG"},
	}
	MigrationsFlag = &cli.StringFlag{
		Name:    "migrations-dir",
		Value:   "./migrations",
		Usage:   "path to migrations folder",
		EnvVars: []string{"WALLET_SERVICES_MIGRATIONS_DIR"},
	}
)

func runMigrations(ctx *cli.Context) error {
	ctx.Context = opio.CancelOnInterrupt(ctx.Context)
	log.Info("running migrations...")
	cfg, err := config.New(ctx.String(ConfigFlag.Name))
	if err != nil {
		log.Error("failed to load config", "err", err)
		return err
	}
	db, err := database.NewDB(ctx.Context, cfg.MasterDB)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		return err
	}
	defer func(db *database.DB) {
		err := db.Close()
		if err != nil {
			log.Error("fail to close database", "err", err)
		}
	}(db)
	return db.ExecuteSQLMigration(cfg.Migrations)
}

func runWalletServices(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	log.Info("running phoenix node...")
	cfg, err := config.New(ctx.String(ConfigFlag.Name))
	if err != nil {
		log.Error("failed to load config", "err", err)
		return nil, err
	}
	return relayer_node.NewWalletServices(ctx.Context, cfg, shutdown)
}

func runApi(ctx *cli.Context, _ context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	log.Info("running api...")
	cfg, err := config.New(ctx.String(ConfigFlag.Name))
	if err != nil {
		log.Error("failed to load config", "err", err)
		return nil, err
	}
	return api.NewApi(ctx.Context, cfg)
}

func runRpc(ctx *cli.Context, _ context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	fmt.Println("running grpc services...")
	cfg, err := config.New(ctx.String(ConfigFlag.Name))
	if err != nil {
		log.Error("config error", "err", err)
		return nil, err
	}

	grpcServerCfg := &grpc.PhoenixRpcConfig{
		Host: cfg.RpcServer.Host,
		Port: cfg.RpcServer.Port,
	}

	db, err := database.NewDB(ctx.Context, cfg.MasterDB)
	if err != nil {
		log.Error("new database fail", "err", err)
		return nil, err
	}
	return grpc.NewPhoenixRpcService(grpcServerCfg, db)
}

func runSendTx(ctx *cli.Context, _ context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	log.Info("Testing transaction signing and broadcasting...")

	// Example parameters - in production these would come from config or flags
	privateKeyHex := "your_private_key_here"                      // TODO: Replace with actual private key
	rpcURL := "https://eth-sepolia.g.alchemy.com/v2/your-api-key" // TODO: Replace with actual RPC URL
	chainID := big.NewInt(11155111)                               // Sepolia testnet
	walletAccountAddr := "127.0.0.1:8189"

	// Create signer
	signer, err := eth.NewSigner(privateKeyHex, chainID)
	if err != nil {
		log.Error("Failed to create signer", "err", err)
		return nil, err
	}

	log.Info("Signer created", "address", signer.GetAddress())

	// Sign transaction
	txParams := eth.TxParams{
		To:       "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		Value:    big.NewInt(1000000000000000), // 0.001 ETH
		GasLimit: 21000,
		Data:     []byte{},
	}

	rawTx, err := signer.SignTransaction(ctx.Context, rpcURL, txParams)
	if err != nil {
		log.Error("Failed to sign transaction", "err", err)
		return nil, err
	}

	log.Info("Transaction signed", "rawTx", rawTx)

	// Create wallet account client
	client, err := account.NewWalletAccountClient(walletAccountAddr)
	if err != nil {
		log.Error("Failed to create wallet account client", "err", err)
		return nil, err
	}
	defer client.Close()

	// Send transaction
	result, err := client.SendTx(ctx.Context, account.SendTxParams{
		Chain:   "Ethereum",
		Coin:    "ETH",
		Network: "testnet",
		RawTx:   rawTx,
	})

	if err != nil {
		log.Error("Failed to send transaction", "err", err)
		return nil, err
	}

	log.Info("Transaction sent", "code", result.Code, "msg", result.Msg, "txHash", result.TxHash)

	return nil, nil
}

func NewCli() *cli.App {
	flags := []cli.Flag{ConfigFlag}
	migrationFlags := []cli.Flag{MigrationsFlag, ConfigFlag}
	return &cli.App{
		Version:              "0.0.1",
		Description:          "An Services For RootHash Wallet",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "api",
				Flags:       flags,
				Description: "Run event http api service",
				Action:      cliapp.LifecycleCmd(runApi),
			},
			{
				Name:        "Run event node task",
				Flags:       flags,
				Description: "Run ",
				Action:      cliapp.LifecycleCmd(runWalletServices),
			},
			{
				Name:        "rpc",
				Flags:       flags,
				Description: "Run event grpc service",
				Action:      cliapp.LifecycleCmd(runRpc),
			},
			{
				Name:        "event migrate",
				Flags:       migrationFlags,
				Description: "Run event database migrations",
				Action:      runMigrations,
			},
			{
				Name:        "send-tx",
				Flags:       flags,
				Description: "Test sending a signed transaction via wallet-chain-account service",
				Action:      cliapp.LifecycleCmd(runSendTx),
			},
			{
				Name:        "event version",
				Description: "Show event services project version",
				Action: func(ctx *cli.Context) error {
					cli.ShowVersion(ctx)
					return nil
				},
			},
		},
	}
}
