package main

import (
	"context"
	"os"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/common/opio"
)

// @title           Wallet Service API
// @version         1.0
// @description     Swagger UI for API testing
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelInfo, true)))
	app := NewCli()
	ctx := opio.WithInterruptBlocker(context.Background())
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Error("Application failed", "Err", err.Error())
		os.Exit(1)
	}
}
