package resolver

import "github.com/roothash-pay/wallet-services/services/market/model"

type Resolver interface {
	Resolve(quotes []model.Quote) (map[string]model.Quote, error)
}
