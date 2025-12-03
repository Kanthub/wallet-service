# 系统架构文档

## 概述

Wallet Services 是一个企业级钱包服务平台，采用分层架构设计，提供多链资产管理、闪兑聚合、用户认证等核心功能。

---

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                        Client Apps                          │
│                    (Web / Mobile / SDK)                     │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/WebSocket
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Wallet Services API                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Auth API   │  │  Aggregator  │  │  Admin API   │     │
│  │              │  │     API      │  │              │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Service Layer                           │  │
│  │  - AggregatorService                                 │  │
│  │  - AuthService                                       │  │
│  │  - AdminService                                      │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Storage Layer                           │  │
│  │  - QuoteStore (In-Memory)                            │  │
│  │  - SwapStore (In-Memory)                             │  │
│  │  - Database (PostgreSQL)                             │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────┬────────────────────────────────────────────┘
                 │ gRPC
                 ▼
┌─────────────────────────────────────────────────────────────┐
│              wallet-chain-account Service                   │
│  - Transaction Broadcasting                                 │
│  - Transaction Query                                        │
│  - Multi-chain Support                                      │
└────────────────┬────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────┐
│                    Blockchain Networks                      │
│  Ethereum │ BSC │ Polygon │ Arbitrum │ Solana │ ...        │
└─────────────────────────────────────────────────────────────┘
```

---

## 闪兑聚合器架构

### 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                   Aggregator Service                        │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Provider Layer                          │  │
│  │                                                      │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐          │  │
│  │  │ 0x       │  │ 1inch    │  │ Jupiter  │          │  │
│  │  │ Provider │  │ Provider │  │ Provider │  ...     │  │
│  │  └──────────┘  └──────────┘  └──────────┘          │  │
│  │                                                      │  │
│  │  ┌──────────┐                                       │  │
│  │  │ LiFi     │  (支持跨链)                           │  │
│  │  │ Provider │                                       │  │
│  │  └──────────┘                                       │  │
│  └──────────────────────────────────────────────────────┘  │
│                         │                                   │
│                         ▼                                   │
│  ┌──────────────────────────────────────────────────────┐  │
│  │           Quote Aggregation Engine                   │  │
│  │  - Concurrent quote fetching                         │  │
│  │  - Best price selection                              │  │
│  │  - Quote caching (5 min TTL)                         │  │
│  └──────────────────────────────────────────────────────┘  │
│                         │                                   │
│                         ▼                                   │
│  ┌──────────────────────────────────────────────────────┐  │
│  │           Swap Execution Engine                      │  │
│  │  - Multi-step transaction generation                 │  │
│  │  - Approval + Swap flow                              │  │
│  │  - Idempotency handling                              │  │
│  │  - State tracking                                    │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 核心组件

#### 1. Provider Layer (提供商层)

**职责**：
- 与各个 DEX 聚合器 API 通信
- 统一接口抽象
- 错误处理和重试

**实现**：
```go
type Provider interface {
    GetQuote(ctx context.Context, req *backend.QuoteRequest) (*backend.Quote, error)
    Name() string
    SupportedChainType() backend.ChainType
}
```

**支持的 Providers**：
- **0x Protocol** (EVM) - 待实现
- **1inch** (EVM) - 待实现
- **Jupiter** (Solana) - 待实现
- **LiFi** (EVM + 跨链) - ✅ 已实现

#### 2. Quote Aggregation Engine (报价聚合引擎)

**职责**：
- 并发查询所有 providers
- 按价格排序
- 缓存报价结果

**实现细节**：
```go
func (s *AggregatorService) aggregateQuotes(ctx context.Context, req *backend.QuoteRequest) ([]*backend.Quote, error) {
    g, ctx := errgroup.WithContext(ctx)
    quoteChan := make(chan *backend.Quote, len(s.providers))
    
    // 并发查询所有 providers
    for _, p := range s.providers {
        p := p
        g.Go(func() error {
            quote, err := p.GetQuote(ctx, req)
            if err != nil {
                log.Warn("Provider failed", "provider", p.Name(), "err", err)
                return nil // 不中断整个聚合
            }
            quoteChan <- quote
            return nil
        })
    }
    
    g.Wait()
    close(quoteChan)
    
    // 收集成功的报价
    var quotes []*backend.Quote
    for quote := range quoteChan {
        quotes = append(quotes, quote)
    }
    
    return quotes, nil
}
```

**特性**：
- 使用 `errgroup` 并发查询
- 单个 provider 失败不影响其他
- 自动按 `to_amount` 降序排序
- 5 分钟 TTL 缓存

#### 3. Swap Execution Engine (交换执行引擎)

**职责**：
- 生成多步骤交易
- 管理交易状态
- 幂等性保证

**交易流程**：
```
Step 0: Approve (如需要)
  ↓
Step 1: Swap
  ↓
Step N: Additional steps (跨链场景)
```

**状态机**：
```
Swap States:
  pending → in_progress → completed
                       ↘ failed

Step States:
  pending → submitted → confirmed
                     ↘ failed
```

---

## 数据流

### 1. 获取报价流程

```
Client
  │
  │ POST /api/v1/aggregator/quotes
  ▼
AggregatorRoutes
  │
  │ ValidateRequest
  ▼
AggregatorService.GetQuotes()
  │
  │ Concurrent calls
  ├─→ 0x Provider
  ├─→ 1inch Provider
  ├─→ Jupiter Provider
  └─→ LiFi Provider
  │
  │ Aggregate & Sort
  ▼
QuoteStore.Save()
  │
  │ Return best + alternatives
  ▼
Client
```

### 2. 执行交换流程

```
Client
  │
  │ POST /api/v1/aggregator/swap/prepare
  ▼
AggregatorService.PrepareSwap()
  │
  │ Retrieve quote from cache
  ▼
Generate Actions
  │
  ├─→ Check allowance
  │   └─→ Generate Approve tx (if needed)
  │
  └─→ Generate Swap tx
  │
  ▼
SwapStore.CreateSwap()
  │
  │ Return unsigned transactions
  ▼
Client signs transactions
  │
  │ POST /api/v1/aggregator/tx/submitSigned
  ▼
AggregatorService.SubmitSignedTx()
  │
  │ Check idempotency
  ▼
WalletAccountClient.SendTx()
  │
  │ gRPC call
  ▼
wallet-chain-account
  │
  │ Broadcast to blockchain
  ▼
Blockchain Network
```

---

## 存储层设计

### QuoteStore (报价存储)

**实现**: In-Memory (内存)

**数据结构**:
```go
type quoteEntry struct {
    quote     *backend.QuoteResponse
    expiresAt time.Time
}

type InMemoryQuoteStore struct {
    mu     sync.RWMutex
    quotes map[string]*quoteEntry
}
```

**特性**:
- 5 分钟 TTL
- 自动清理过期数据
- 线程安全（RWMutex）

**未来优化**:
- 迁移到 Redis
- 支持分布式缓存
- 持久化选项

### SwapStore (交换存储)

**实现**: In-Memory (内存)

**数据结构**:
```go
type InMemorySwapStore struct {
    mu             sync.RWMutex
    swaps          map[string]*backend.Swap
    idempotencyMap map[string]string  // key: swapID+stepIndex+idempotencyKey → txHash
}
```

**特性**:
- 完整的交换状态管理
- 幂等性检查
- 步骤级别的状态追踪

**幂等性实现**:
```go
func (s *InMemorySwapStore) CheckIdempotency(ctx context.Context, swapID string, stepIndex int, idempotencyKey string) (string, bool) {
    key := fmt.Sprintf("%s:%d:%s", swapID, stepIndex, idempotencyKey)
    txHash, exists := s.idempotencyMap[key]
    return txHash, exists
}
```

---

## gRPC 集成

### WalletAccountClient

**职责**:
- 与 wallet-chain-account 服务通信
- 交易广播
- 交易查询

**接口**:
```go
type WalletAccountClient struct {
    conn   *grpc.ClientConn
    client pb.WalletAccountServiceClient
}

// 发送交易
func (c *WalletAccountClient) SendTx(ctx context.Context, params SendTxParams) (*SendTxResult, error)

// 查询交易
func (c *WalletAccountClient) GetTxByHash(ctx context.Context, consumerToken, chain, coin, network, txHash string) (*TxInfo, error)
```

**配置**:
```yaml
aggregator_config:
  wallet_account_addr: "127.0.0.1:8189"
```

**超时设置**:
- SendTx: 30 秒
- GetTxByHash: 10 秒

---

## 安全设计

### 1. 验证器 (Validator)

**职责**:
- 链 ID 验证
- Router/Spender 白名单检查
- 交易金额限制

**实现**:
```go
type Validator struct {
    whitelistedRouters  map[string]bool
    whitelistedSpenders map[string]bool
    maxValueWei         *big.Int
}
```

**验证规则**:
- 只允许白名单中的 router 地址
- 只允许白名单中的 spender 地址
- 单笔交易不超过最大金额限制（默认 100 ETH）

### 2. 幂等性保证

**目的**: 防止重复提交导致的双花

**实现**:
- 客户端生成唯一的 `idempotency_key`
- 服务端记录 `swapID + stepIndex + idempotencyKey → txHash` 映射
- 重复提交返回原 txHash

**示例**:
```go
// 检查幂等性
if txHash, exists := s.swapStore.CheckIdempotency(ctx, req.SwapID, req.StepIndex, req.IdempotencyKey); exists {
    return &backend.SubmitSignedTxResponse{TxHash: txHash}, nil
}
```

### 3. 报价过期机制

**目的**: 防止使用过期报价导致的滑点过大

**实现**:
- 报价生成时设置 `expires_at` (5 分钟后)
- 准备交换时检查报价是否过期
- 过期报价自动清理

---

## 性能优化

### 1. 并发查询

使用 `golang.org/x/sync/errgroup` 并发查询所有 providers：

```go
g, ctx := errgroup.WithContext(ctx)
for _, p := range s.providers {
    p := p
    g.Go(func() error {
        quote, err := p.GetQuote(ctx, req)
        // ...
    })
}
g.Wait()
```

**优势**:
- 查询时间 = max(provider 响应时间)，而非 sum
- 单个 provider 失败不影响其他
- 自动 context 取消传播

### 2. 内存缓存

**QuoteStore**:
- 避免重复查询相同报价
- 5 分钟 TTL 平衡新鲜度和性能

**优化空间**:
- 迁移到 Redis 支持分布式
- 添加 LRU 淘汰策略
- 预热常用交易对

### 3. HTTP 客户端复用

每个 Provider 使用单例 HTTP 客户端：

```go
type LiFiProvider struct {
    apiURL     string
    apiKey     string
    httpClient *http.Client  // 复用连接
}

func NewLiFiProvider(apiURL, apiKey string) *LiFiProvider {
    return &LiFiProvider{
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}
```

---

## 相关文档

- [API 文档](./api.md) - RESTful API 详细说明
- [部署文档](./deployment.md) - 部署和运维指南
- [返回 README](../README.md)

---

**最后更新**: 2024-12-03

