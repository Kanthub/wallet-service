package chaininfo

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/common/redis"
	dbBackend "github.com/roothash-pay/wallet-services/database/backend"
)

const (
	cacheKeyPrefix   = "chain_info:"
	cacheTTL         = 24 * time.Hour
	defaultChainType = "EVM"
)

// Info represents cached chain metadata required across services.
type Info struct {
	ChainID       string `json:"chain_id"`
	ChainType     string `json:"chain_type"`
	ChainName     string `json:"chain_name"`
	Network       string `json:"network"`
	NativeSymbol  string `json:"native_symbol"`
	ExplorerURL   string `json:"explorer_url"`
	WalletChain   string `json:"wallet_chain"`
	WalletNetwork string `json:"wallet_network"`
	WalletCoin    string `json:"wallet_coin"`
	ConsumerToken string `json:"consumer_token"`
	RPCURL        string `json:"rpc_url"`
	IsEnabled     bool   `json:"is_enabled"`
}

// Provider exposes the functionality required by API/worker layers.
type Provider interface {
	WarmUp(ctx context.Context) error
	Get(ctx context.Context, chainID string) (*Info, error)
	Refresh(ctx context.Context, chainID string) (*Info, error)
}

// Manager implements Provider with DB + Redis + in-memory caches.
type Manager struct {
	db                   dbBackend.ChainDB
	redis                *redis.Client
	cache                map[string]*Info
	cacheM               sync.RWMutex
	defaultConsumerToken string
	chainConsumerTokens  map[string]string
}

// NewManager creates a new chain info manager.
func NewManager(db dbBackend.ChainDB, redisClient *redis.Client, defaultToken string, chainTokens map[string]string) *Manager {
	return &Manager{
		db:                   db,
		redis:                redisClient,
		cache:                make(map[string]*Info),
		defaultConsumerToken: defaultToken,
		chainConsumerTokens:  chainTokens,
	}
}

// WarmUp loads all enabled chain info into memory (and Redis if configured).
func (m *Manager) WarmUp(ctx context.Context) error {
	if m == nil || m.db == nil {
		return fmt.Errorf("chain info manager not initialized")
	}
	list, err := m.db.ListAllChains()
	if err != nil {
		return fmt.Errorf("failed to list chain info: %w", err)
	}
	for _, item := range list {
		info := m.convertInfo(item)
		m.setLocal(info)
		if err := m.saveRedis(ctx, info); err != nil {
			log.Warn("Failed to cache chain info into redis", "chainID", info.ChainID, "err", err)
		}
	}
	log.Info("Chain info warmup completed", "count", len(list))
	return nil
}

// Get returns chain info with caching fallback order: memory -> redis -> DB.
func (m *Manager) Get(ctx context.Context, chainID string) (*Info, error) {
	if m == nil {
		return nil, fmt.Errorf("chain info manager not initialized")
	}
	if chainID == "" {
		return nil, fmt.Errorf("chainID is required")
	}

	if info := m.getLocal(chainID); info != nil {
		return info, nil
	}

	if info := m.getFromRedis(ctx, chainID); info != nil {
		m.setLocal(info)
		return info, nil
	}

	return m.Refresh(ctx, chainID)
}

// Refresh forces a DB lookup for the latest data and updates caches.
func (m *Manager) Refresh(ctx context.Context, chainID string) (*Info, error) {
	if m == nil || m.db == nil {
		return nil, fmt.Errorf("chain info manager not initialized")
	}
	item, err := m.db.GetByChainID(chainID)
	if err != nil {
		return nil, err
	}
	info := m.convertInfo(item)
	m.setLocal(info)
	if err := m.saveRedis(ctx, info); err != nil {
		log.Warn("Failed to refresh redis cache for chain info", "chainID", chainID, "err", err)
	}
	return info, nil
}

func (m *Manager) cacheKey(chainID string) string {
	return cacheKeyPrefix + chainID
}

func (m *Manager) setLocal(info *Info) {
	if info == nil {
		return
	}
	m.cacheM.Lock()
	defer m.cacheM.Unlock()
	m.cache[info.ChainID] = info
}

func (m *Manager) getLocal(chainID string) *Info {
	m.cacheM.RLock()
	defer m.cacheM.RUnlock()
	return m.cache[chainID]
}

func (m *Manager) getFromRedis(ctx context.Context, chainID string) *Info {
	if m.redis == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	data, err := m.redis.Get(ctx, m.cacheKey(chainID)).Result()
	if err != nil {
		return nil
	}
	var info Info
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		log.Warn("Failed to unmarshal chain info from redis", "chainID", chainID, "err", err)
		return nil
	}
	info.ConsumerToken = m.resolveConsumerToken(chainID)
	return &info
}

func (m *Manager) saveRedis(ctx context.Context, info *Info) error {
	if m.redis == nil || info == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return m.redis.Set(ctx, m.cacheKey(info.ChainID), payload, cacheTTL).Err()
}

func (m *Manager) resolveConsumerToken(chainID string) string {
	if m.chainConsumerTokens != nil {
		if token, ok := m.chainConsumerTokens[chainID]; ok {
			return token
		}
	}
	return m.defaultConsumerToken
}

func (m *Manager) convertInfo(src *dbBackend.Chain) *Info {
	if src == nil {
		return nil
	}
	chainType := src.ChainType
	if chainType == "" {
		chainType = defaultChainType
	}
	info := &Info{
		ChainID:       src.ChainID,
		ChainType:     chainType,
		ChainName:     src.ChainName,
		Network:       src.Network,
		NativeSymbol:  src.NativeSymbol,
		ExplorerURL:   src.ExplorerURL,
		WalletChain:   src.WalletChain,
		WalletNetwork: src.WalletNetwork,
		WalletCoin:    src.WalletCoin,
		RPCURL:        src.RpcURL,
		IsEnabled:     src.IsEnabled,
	}
	info.ConsumerToken = m.resolveConsumerToken(src.ChainID)
	return info
}
