package service

import (
	"context"
	"fmt"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type TokenService interface {
	CreateToken(ctx context.Context, req CreateTokenRequest) (*backend.Token, error)
	UpdateToken(ctx context.Context, req UpdateTokenRequest) error
	GetToken(ctx context.Context, guid string) (*backend.Token, error)
	ListTokens(
		ctx context.Context,
		page, pageSize int,
		filters map[string]interface{},
	) ([]*backend.Token, int64, error)
}

type CreateTokenRequest struct {
	TokenName            string `json:"token_name"`
	TokenMark            string `json:"token_mark"`
	TokenLogo            string `json:"token_logo"`
	TokenActiveLogo      string `json:"token_active_logo"`
	TokenDecimal         string `json:"token_decimal"`
	TokenSymbol          string `json:"token_symbol"`
	TokenContractAddress string `json:"token_contract_address"`
	ChainID              string `json:"chain_id"`
	IsHot                string `json:"is_hot"`
}

type UpdateTokenRequest struct {
	Guid    string                 `json:"guid"`
	Updates map[string]interface{} `json:"updates"`
}

type tokenService struct {
	db *database.DB
}

func NewTokenService(db *database.DB) TokenService {
	return &tokenService{db: db}
}

func (s *tokenService) CreateToken(
	ctx context.Context,
	req CreateTokenRequest,
) (*backend.Token, error) {

	if req.TokenContractAddress == "" {
		return nil, fmt.Errorf("token_contract_address required")
	}

	token := &backend.Token{
		TokenName:            req.TokenName,
		TokenMark:            req.TokenMark,
		TokenLogo:            req.TokenLogo,
		TokenActiveLogo:      req.TokenActiveLogo,
		TokenDecimal:         req.TokenDecimal,
		TokenSymbol:          req.TokenSymbol,
		TokenContractAddress: req.TokenContractAddress,
		ChainID:              req.ChainID,
		IsHot:                req.IsHot,
		CreateTime:           time.Now(),
		UpdateTime:           time.Now(),
	}

	if err := s.db.BackendToken.StoreToken(token); err != nil {
		return nil, err
	}
	return token, nil
}

func (s *tokenService) UpdateToken(
	ctx context.Context,
	req UpdateTokenRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid required")
	}
	if len(req.Updates) == 0 {
		return fmt.Errorf("updates empty")
	}

	return s.db.BackendToken.UpdateToken(req.Guid, req.Updates)
}

func (s *tokenService) GetToken(
	ctx context.Context,
	guid string,
) (*backend.Token, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	return s.db.BackendToken.GetByGuid(guid)
}

func (s *tokenService) ListTokens(
	ctx context.Context,
	page, pageSize int,
	filters map[string]interface{},
) ([]*backend.Token, int64, error) {

	return s.db.BackendToken.GetTokenList(page, pageSize, filters)
}
