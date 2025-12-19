package provider

import (
	"context"

	"github.com/roothash-pay/wallet-services/services/market/model"
)

type Provider interface {
	Name() string
	FetchQuotes(ctx context.Context) ([]model.Quote, error)
}
