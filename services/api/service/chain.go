package service

import (
	"context"
	"fmt"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type ChainService interface {
	// 管理后台
	CreateChain(ctx context.Context, c *backend.Chain) (*backend.Chain, error)
	UpdateChain(ctx context.Context, guid string, updates map[string]interface{}) error

	// 查询
	GetChain(ctx context.Context, guid string) (*backend.Chain, error)
	GetChainByChainID(ctx context.Context, chainID string) (*backend.Chain, error)
	GetChainByName(ctx context.Context, name string) (*backend.Chain, error)

	ListChains(
		ctx context.Context,
		page, pageSize int,
		filters map[string]interface{},
	) ([]*backend.Chain, int64, error)

	// 给前端 / 聚合器
	ListAllChains(ctx context.Context) ([]*backend.Chain, error)
}

type chainService struct {
	db *database.DB
}

func NewChainService(db *database.DB) ChainService {
	return &chainService{db: db}
}

func (s *chainService) CreateChain(
	ctx context.Context,
	c *backend.Chain,
) (*backend.Chain, error) {

	if c.ChainID == "" {
		return nil, fmt.Errorf("chain_id required")
	}
	if c.ChainName == "" {
		return nil, fmt.Errorf("chain_name required")
	}
	if c.ChainMark == "" {
		return nil, fmt.Errorf("chain_mark required")
	}

	if err := s.db.BackendChain.StoreChain(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *chainService) UpdateChain(
	ctx context.Context,
	guid string,
	updates map[string]interface{},
) error {

	return s.db.BackendChain.UpdateChain(guid, updates)
}

func (s *chainService) GetChain(
	ctx context.Context,
	guid string,
) (*backend.Chain, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	return s.db.BackendChain.GetByGuid(guid)
}

func (s *chainService) GetChainByChainID(
	ctx context.Context,
	chainID string,
) (*backend.Chain, error) {

	if chainID == "" {
		return nil, fmt.Errorf("chain_id required")
	}
	return s.db.BackendChain.GetByChainID(chainID)
}

func (s *chainService) GetChainByName(
	ctx context.Context,
	name string,
) (*backend.Chain, error) {

	if name == "" {
		return nil, fmt.Errorf("chain_name required")
	}
	return s.db.BackendChain.GetByName(name)
}

func (s *chainService) ListChains(
	ctx context.Context,
	page, pageSize int,
	filters map[string]interface{},
) ([]*backend.Chain, int64, error) {

	return s.db.BackendChain.GetChainList(page, pageSize, filters)
}

func (s *chainService) ListAllChains(
	ctx context.Context,
) ([]*backend.Chain, error) {

	return s.db.BackendChain.ListAllChains()
}
