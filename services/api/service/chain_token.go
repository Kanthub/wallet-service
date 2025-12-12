package service

import (
	"context"
	"fmt"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type ChainTokenService interface {
	CreateChainToken(ctx context.Context, req CreateChainTokenRequest) (*backend.ChainToken, error)
	UpdateChainToken(ctx context.Context, req UpdateChainTokenRequest) error
	GetChainToken(ctx context.Context, guid string) (*backend.ChainToken, error)
	ListByChainID(ctx context.Context, chainID string) ([]*backend.ChainToken, error)
	ListByTokenID(ctx context.Context, tokenID string) ([]*backend.ChainToken, error)
}

type CreateChainTokenRequest struct {
	ChainID string `json:"chain_id"`
	TokenID string `json:"token_id"`
}

type UpdateChainTokenRequest struct {
	Guid    string                 `json:"guid"`
	Updates map[string]interface{} `json:"updates"`
}

type chainTokenService struct {
	db *database.DB
}

func NewChainTokenService(db *database.DB) ChainTokenService {
	return &chainTokenService{db: db}
}

func (s *chainTokenService) CreateChainToken(
	ctx context.Context,
	req CreateChainTokenRequest,
) (*backend.ChainToken, error) {

	if req.ChainID == "" || req.TokenID == "" {
		return nil, fmt.Errorf("chain_id and token_id required")
	}

	item := &backend.ChainToken{
		ChainID:    req.ChainID,
		TokenID:    req.TokenID,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := s.db.BackendChainToken.StoreChainToken(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *chainTokenService) UpdateChainToken(
	ctx context.Context,
	req UpdateChainTokenRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid required")
	}
	if len(req.Updates) == 0 {
		return fmt.Errorf("updates empty")
	}

	return s.db.BackendChainToken.UpdateChainToken(req.Guid, req.Updates)
}

func (s *chainTokenService) GetChainToken(
	ctx context.Context,
	guid string,
) (*backend.ChainToken, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	return s.db.BackendChainToken.GetByGuid(guid)
}

func (s *chainTokenService) ListByChainID(
	ctx context.Context,
	chainID string,
) ([]*backend.ChainToken, error) {

	if chainID == "" {
		return nil, fmt.Errorf("chain_id required")
	}
	return s.db.BackendChainToken.GetByChainID(chainID)
}

func (s *chainTokenService) ListByTokenID(
	ctx context.Context,
	tokenID string,
) ([]*backend.ChainToken, error) {

	if tokenID == "" {
		return nil, fmt.Errorf("token_id required")
	}
	return s.db.BackendChainToken.GetByTokenID(tokenID)
}
