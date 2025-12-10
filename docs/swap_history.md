# 闪兑历史持久化

## 概述

闪兑聚合器采用**主动持久化 + 定时扫链**策略：
- ✅ **提交时立即保存** - 用户提交交易时立即保存到数据库（状态：CREATED → PENDING）
- ✅ **定时任务自动更新** - 后台定时扫描链上状态并更新数据库（PENDING → SUCCESS/FAILED）
- ✅ **不依赖用户查询** - 即使用户不查询状态，交易记录也会被自动跟踪

## 状态设计

### 4 个状态（0/1/2/3）

| 状态值 | 状态名 | 说明 | 触发时机 |
|-------|--------|------|---------|
| **0** | CREATED | 后端收到 signedTx 请求并写入记录，但尚未广播 | `SubmitSignedTx` 开始时 |
| **1** | PENDING | 广播成功并拿到 txHash | `SendTx` 成功后 |
| **2** | FAILED | 广播失败或链上执行失败或超时 | 广播失败 / 链上失败 / 超时 |
| **3** | SUCCESS | 链上确认成功 | 链上确认成功 |

### 状态流转

```
CREATED (0) → PENDING (1) → SUCCESS (3)
     ↓             ↓
   FAILED (2) ← FAILED (2)
```

**流转规则：**
1. **CREATED → PENDING**: 广播成功，获得 txHash
2. **CREATED → FAILED**: 广播失败
3. **PENDING → SUCCESS**: 链上确认成功（EVM receipt status=1；Solana 确认且无 error）
4. **PENDING → FAILED**: 链上执行失败或超时（>1小时）

## 设计原则

### 只存储最终结果，不存储中间步骤

**存储内容：**
- ✅ 完成的交换记录（成功或失败）
- ✅ 最终的 swap 交易哈希
- ✅ 交易金额、代币信息
- ✅ 区块链浏览器链接

**不存储内容：**
- ❌ Quote（报价）- 5分钟过期，临时数据
- ❌ Step（步骤）- 技术实现细节，用户不关心
- ❌ Approve 交易 - 中间步骤

### 为什么不存储步骤？

1. **用户视角** - 用户只关心"我用 1000 USDC 换了 0.5 ETH"，不关心中间有几个步骤
2. **链上可查** - 所有交易都有 tx_hash，可以在区块链浏览器查看完整信息
3. **生命周期短** - 步骤数据只在交易进行中有用（< 10 分钟），完成后就没用了
4. **降低复杂度** - 不需要额外的表和 JOIN 查询

## 数据流

### 新设计（主动持久化 + 定时扫链）✅

```
┌─────────────────────────────────────────┐
│         用户发起交换                      │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│    内存中创建 Swap + Steps               │
│    (用于实时状态追踪)                     │
└─────────────────────────────────────────┘
                  ↓
         用户签名并提交交易
                  ↓
┌─────────────────────────────────────────┐
│  Step 1: 保存 CREATED 状态到数据库        │
│  (广播前立即保存) ✅                      │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│  Step 2: 调用 SendTx 广播交易             │
└─────────────────────────────────────────┘
         ↓                    ↓
    广播成功              广播失败
         ↓                    ↓
┌──────────────────┐  ┌──────────────────┐
│ Step 3a: 更新为   │  │ Step 3b: 更新为   │
│ PENDING + txHash │  │ FAILED + 失败原因 │
└──────────────────┘  └──────────────────┘
         ↓
┌─────────────────────────────────────────┐
│    定时任务并发扫链 (每 10 秒)            │
│    - 查询 PENDING 交易的链上状态          │
│    - 并发度: 10 workers                  │
│    - 超时阈值: 1 小时                     │
└─────────────────────────────────────────┘
         ↓
    ┌────┴────┐
    ↓         ↓
链上成功   链上失败/超时
    ↓         ↓
┌─────┐   ┌─────┐
│SUCCESS│   │FAILED│
└─────┘   └─────┘
    ↓         ↓
┌─────────────────────────────────────────┐
│    用户可以在交易历史中查看               │
│    (不依赖用户查询，自动更新) ✅          │
└─────────────────────────────────────────┘
```

### 优势

| 特性 | 旧设计（被动） | 新设计（主动 + 扫链） |
|------|---------------|---------------------|
| **保存时机** | ❌ 用户查询状态时 | ✅ 提交交易时（广播前） |
| **依赖查询** | ❌ 必须查询才保存 | ✅ 完全不依赖查询 |
| **记录完整性** | ❌ 可能丢失记录 | ✅ 所有交易都记录 |
| **状态实时性** | ❌ 取决于查询频率 | ✅ 定时任务自动更新（10秒） |
| **失败追踪** | ❌ 广播失败无记录 | ✅ 完整记录失败原因 |

## 实现细节

### 1. 提交时保存（SubmitSignedTx）- 三步流程

#### Step 1: 保存 CREATED 状态（广播前）

```go
// Step 1: Save to database with CREATED status (before broadcast)
var recordGuid string
if step.ActionType == backend.ActionTypeSwap {
    recordGuid = s.saveSwapHistoryCreated(ctx, swap, quote)
}
```

**保存内容：**
- Status: `0` (CREATED)
- Hash: 空（尚未广播）
- Memo: `"Swap via lifi: 1000.000000 USDC -> 0.500000 ETH (Created)"`

#### Step 2: 广播交易

```go
// Step 2: Broadcast transaction using SendTx
result, err := s.accountClient.SendTx(ctx, account.SendTxParams{
    ConsumerToken: "",
    Chain:         chain,
    Coin:          coin,
    Network:       network,
    RawTx:         step.SignedTx,
})
```

#### Step 3a: 广播成功 → PENDING

```go
if err == nil {
    // Step 3a: Update database record to PENDING (after successful broadcast)
    if recordGuid != "" {
        s.updateSwapHistoryPending(ctx, recordGuid, txHash, quote)
    }
}
```

**更新内容：**
- Status: `1` (PENDING)
- Hash: txHash（广播返回的交易哈希）
- Memo: `"Swap via lifi: 1000.000000 USDC -> 0.500000 ETH (Pending)"`
- ExplorerURL: 区块链浏览器链接

#### Step 3b: 广播失败 → FAILED

```go
if err != nil {
    // Step 3b: Update database record to FAILED
    if recordGuid != "" {
        s.updateSwapHistoryFailed(ctx, recordGuid, FailReasonBroadcastFailed, err.Error())
    }
}
```

**更新内容：**
- Status: `2` (FAILED)
- FailReasonCode: `"BROADCAST_FAILED"`
- FailReasonMsg: 错误详情
- Memo: `"Swap via lifi: 1000.000000 USDC -> 0.500000 ETH (Failed: broadcast failed)"`

### 2. 定时任务扫链（WalletTxRecordWorker）

#### 自动启动

Worker 已集成到 WalletServices 中，会在 worker 服务启动时自动启动：

```go
// walletsvc.go
func (as *WalletServices) initWorker(config *config.Config) error {
    // ... 初始化其他 workers

    // Initialize wallet tx record worker if aggregator is enabled
    if config.AggregatorConfig.WalletAccountAddr != "" {
        accountClient, err := account.NewWalletAccountClient(config.AggregatorConfig.WalletAccountAddr)
        if err != nil {
            log.Warn("failed to create wallet account client for tx worker", "err", err)
        } else {
            txRecordWorker := aggregator_task.NewWalletTxRecordWorker(
                as.DB.BackendWalletTxRecord,
                accountClient,
                txWorkerConfig,
            )
            as.txRecordWorker = txRecordWorker
        }
    }
    return nil
}

func (as *WalletServices) Start(ctx context.Context) error {
    // ... 启动其他 workers

    // Start tx record worker if initialized
    if as.txRecordWorker != nil {
        as.txRecordWorker.Start()
        log.Info("Wallet tx record worker started")
    }
    return nil
}
```

**启动条件：**
- ✅ `aggregator_config.wallet_account_addr` 已配置
- ✅ WalletServices 启动时自动初始化和启动
- ✅ WalletServices 停止时自动停止

#### 配置参数

```go
type WalletTxRecordWorkerConfig struct {
    ScanInterval         int // 扫描间隔（秒），默认: 10
    LastCheckedThreshold int // 上次检查阈值（秒），默认: 5
    BatchSize            int // 每批处理数量，默认: 100
    Concurrency          int // 并发度，默认: 10
    TimeoutThreshold     int // 超时阈值（秒），默认: 3600 (1小时)
}
```

**当前使用默认值（硬编码）：**
- 扫描间隔: 10 秒
- 上次检查阈值: 5 秒
- 批处理大小: 100 条
- 并发度: 10 workers
- 超时阈值: 3600 秒（1 小时）

#### 扫描逻辑

```go
// 1. 查询待检查的 PENDING 交易
records, err := w.db.GetPendingTxsForCheck(lastCheckedBefore, w.config.BatchSize)

// 2. 并发处理（worker pool）
for _, record := range records {
    jobs <- record
}

// 3. 查询链上状态
txInfo, err := w.accountClient.GetTxByHash(ctx, "", chain, coin, network, record.Hash)

// 4. 根据状态更新数据库
if txInfo.Status == 3 { // TxStatus_Success
    w.markAsSuccess(record, txInfo.Height, txInfo.Datetime)
} else if txInfo.Status == 2 || txInfo.Status == 4 { // Failed or ContractExecuteFailed
    w.markAsFailed(record, FailReasonChainFailed, "Transaction failed on chain")
}

// 5. 超时处理
if time.Since(record.CreateTime) > timeout {
    w.markAsFailed(record, FailReasonNotFoundTimeout, "Transaction timeout")
}
```

#### 更新内容

**成功时：**
- Status: `3` (SUCCESS)
- BlockHeight: 从链上查询
- TxTime: 交易时间
- Memo: `"Swap via lifi: 1000.000000 USDC -> 0.500000 ETH (Success)"`

**失败时：**
- Status: `2` (FAILED)
- FailReasonCode: `"CHAIN_FAILED"` 或 `"NOT_FOUND_TIMEOUT"`
- FailReasonMsg: 失败详情
- Memo: `"Swap via lifi: 1000.000000 USDC -> 0.500000 ETH (Failed: ...)"`

### 3. 数据库字段设计

**重要改进：** 完善 `WalletTxRecord` 表结构，支持多维度查询和状态追踪

#### 新增字段

```go
type WalletTxRecord struct {
    Guid           string     // 主键：使用 SwapID（每条交易记录的唯一 ID）
    WalletUUID     string     // 新增：关联到 Wallet 表（必填）
    AddressUUID    string     // 新增：关联到 WalletAddress 表（可选）
    TxType         string     // 新增：交易类型（transfer/swap/approve）
    Status         int        // 新增：交易状态（0/1/2/3）
    FailReasonCode string     // 新增：失败原因代码（BROADCAST_FAILED/CHAIN_FAILED/NOT_FOUND_TIMEOUT）
    FailReasonMsg  string     // 新增：失败原因详情
    LastCheckedAt  *time.Time // 新增：上次检查时间（用于定时任务调度）
    // ... 其他字段
}
```

#### 状态常量定义

```go
// TxStatus 状态常量
const (
    TxStatusCreated = 0 // CREATED: 后端收到 signedTx 请求并写入记录，但尚未广播
    TxStatusPending = 1 // PENDING: 广播成功并拿到 txHash
    TxStatusFailed  = 2 // FAILED: 广播失败或链上执行失败或超时
    TxStatusSuccess = 3 // SUCCESS: 链上确认成功
)

// 失败原因代码常量
const (
    FailReasonBroadcastFailed = "BROADCAST_FAILED"  // 广播失败
    FailReasonChainFailed     = "CHAIN_FAILED"      // 链上执行失败
    FailReasonNotFoundTimeout = "NOT_FOUND_TIMEOUT" // 查不到且超时
    FailReasonUnknown         = "UNKNOWN"           // 未知错误
)
```

#### 请求参数更新

```go
type PrepareSwapRequest struct {
    QuoteID     string `json:"quote_id" validate:"required"`
    UserAddress string `json:"user_address" validate:"required"`
    WalletUUID  string `json:"wallet_uuid,omitempty"` // 新增：钱包 UUID
}
```

#### Swap 模型更新

```go
type Swap struct {
    SwapID      string `json:"swap_id"`
    QuoteID     string `json:"quote_id"`
    UserAddress string `json:"user_address"`
    WalletUUID  string `json:"wallet_uuid,omitempty"` // 新增：钱包 UUID
    State       SwapState `json:"state"`
    // ...
}
```

#### 数据库记录

```go
record := &WalletTxRecord{
    Guid:        swap.SwapID,     // 交易记录的唯一 ID（使用 SwapID）
    WalletUUID:  swap.WalletUUID, // 关联到钱包（钱包的 UUID）
    TxType:      "swap",          // 交易类型
    Status:      0,               // 交易状态（CREATED）
    // ...
}
```

**优势：**
- ✅ 关联到具体的用户钱包（通过 wallet_uuid）
- ✅ 支持按钱包查询所有交易
- ✅ 支持按交易类型过滤（transfer/swap/approve）
- ✅ 支持按交易状态过滤（pending/confirmed/failed）
- ✅ 支持多钱包场景

### 保存的数据

保存到 `wallet_tx_record` 表的字段：

| 字段 | 说明 | 示例 |
|------|------|------|
| `guid` | **交易记录 ID** (使用 SwapID) | "swap-abc123..." |
| `wallet_uuid` | **钱包 UUID** (关联 wallet 表) | "wallet-123..." |
| `address_uuid` | 地址 UUID (关联 wallet_address 表) | "" (可选) |
| `tx_time` | 交易时间 | "2024-01-01T12:00:00Z" |
| `chain_id` | 链 ID | "1" (Ethereum) |
| `token_id` | Token 主键（找不到则回落到合约地址） | "usdc-mainnet-guid" |
| `from_address` | 用户地址 | "0x1234...5678" |
| `to_address` | Router 地址 | "0xrouter..." |
| `amount` | 交易金额 | 1000000000 |
| `memo` | 交易备注 | "Swap via lifi: 1000 USDC -> 0.5 ETH (Pending)" |
| `hash` | 交易哈希 | "0xabc123..." |
| `block_height` | 区块高度 | "" (pending) / "12345678" (confirmed) |
| `tx_type` | **交易类型** | "swap" |
| `tx_status` | **交易状态** | "pending" / "confirmed" / "failed" |

---

## 使用示例

### 查询交易历史

```go
// 1. 按钱包查询所有交易（推荐）
records, total, err := db.BackendWalletTxRecord.GetTxList(1, 20, map[string]interface{}{
    "wallet_uuid": "wallet-123",
})

// 2. 按钱包查询所有 swap 交易
records, total, err := db.BackendWalletTxRecord.GetTxList(1, 20, map[string]interface{}{
    "wallet_uuid": "wallet-123",
    "tx_type": "swap",
})

// 3. 按钱包查询所有 pending 交易
records, total, err := db.BackendWalletTxRecord.GetTxList(1, 20, map[string]interface{}{
    "wallet_uuid": "wallet-123",
    "tx_status": "pending",
})

// 4. 按链查询交易
records, total, err := db.BackendWalletTxRecord.GetTxList(1, 20, map[string]interface{}{
    "chain_id": "1", // Ethereum
})

// 5. 按用户地址查询（兼容旧方式）
records, total, err := db.BackendWalletTxRecord.GetTxList(1, 20, map[string]interface{}{
    "from_address": "0x1234...5678",
})

// 6. 通过 swap_id 查询
record, err := db.BackendWalletTxRecord.GetByGuid("swap-id-123")

// 7. 通过 tx_hash 查询
record, err := db.BackendWalletTxRecord.GetByHash("0xabc123...")
```

---

### 辅助方法

#### 1. `parseAmount` - 金额解析
将字符串金额转换为 int64：
```go
amount, err := s.parseAmount("1000000000")
// Returns: 1000000000, nil
```

#### 2. `formatAmount` - 金额格式化
将 wei 格式化为可读格式（假设 18 位小数）：
```go
formatted := s.formatAmount("1000000000000000000")
// Returns: "1.000000"
```

#### 3. `getTokenSymbol` - 代币符号
根据地址获取代币符号：
```go
symbol := s.getTokenSymbol("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
// Returns: "USDC"
```

支持的代币：
- ETH: `0x0000000000000000000000000000000000000000`
- USDC: `0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48`
- USDT: `0xdAC17F958D2ee523a2206206994597C13D831ec7`
- DAI: `0x6B175474E89094C44Da98b954EedeAC495271d0F`

#### 4. `buildExplorerURL` - 浏览器链接
根据链 ID 构建区块链浏览器链接：
```go
url := s.buildExplorerURL("1", "0xabc123...")
// Returns: "https://etherscan.io/tx/0xabc123..."
```

支持的链：
- Ethereum (1): etherscan.io
- BSC (56): bscscan.com
- Polygon (137): polygonscan.com
- Arbitrum (42161): arbiscan.io
- Optimism (10): optimistic.etherscan.io
- Avalanche (43114): snowtrace.io

## 错误处理

`saveSwapHistory` 方法采用**非阻塞**设计：

- ✅ 如果保存失败，只记录错误日志，不影响 swap 状态查询
- ✅ 如果数据库不可用，跳过保存
- ✅ 如果 quote 已过期，记录警告但继续尝试保存

```go
if s.db == nil {
    log.Warn("Database not available, skip saving swap history")
    return
}

if err := s.db.BackendWalletTxRecord.StoreWalletTxRecord(record); err != nil {
    log.Error("Failed to save swap history", "err", err)
    // 不返回错误，继续执行
}
```

## 查询交易历史

用户可以通过现有的 `wallet_tx_record` API 查询交易历史：

```go
// 按用户地址查询
records, total, err := db.BackendWalletTxRecord.GetTxList(1, 20, map[string]interface{}{
    "from_address": "0x1234...5678",
})

// 按交易哈希查询
record, err := db.BackendWalletTxRecord.GetByHash("0xabc123...")
```

## 测试

运行单元测试：

```bash
go test -v ./services/api/service -run "TestParseAmount|TestGetTokenSymbol|TestBuildExplorerURL"
```

## 未来优化

### 可选优化（按需实现）

1. **更丰富的代币映射**
   - 从链上查询代币符号和小数位
   - 支持更多常用代币

2. **更准确的金额格式化**
   - 根据代币的实际小数位格式化
   - 支持不同小数位的代币（6, 8, 18）

3. **区块高度自动获取**
   - 在保存时自动查询区块高度
   - 缓存链上查询结果

4. **失败原因记录**
   - 记录失败的 swap 的详细原因
   - 用于用户查询和问题排查

## 🚀 使用方式

### 1. 启动 WalletServices（包含 Worker）

```bash
# 使用配置文件启动 wallet services（包含所有 workers）
./wallet-services "Run event node task" --config ./wallet-services-config.local.yaml
```

### 2. 查看启动日志

```
INFO [12-04|10:00:00] New wallet services start️ 🕖
INFO [12-04|10:00:00] Init database success
INFO [12-04|10:00:00] Wallet tx record worker initialized      scanInterval=10 concurrency=10 timeout=3600
INFO [12-04|10:00:00] New wallet services success🏅️
INFO [12-04|10:00:00] Wallet tx record worker started
```

### 3. 同时启动 API 服务（可选）

如果需要 HTTP API，可以单独启动 API 服务：

```bash
# 启动 API 服务
./wallet-services api --config ./wallet-services-config.local.yaml
```

---

## 总结

这个设计遵循**简单实用**的原则：

- ✅ 只存储用户关心的最终结果
- ✅ 复用现有的 `wallet_tx_record` 表
- ✅ 不增加额外的表和复杂度
- ✅ 非阻塞设计，不影响主流程
- ✅ 完整的错误处理和日志记录
- ✅ 自动后台扫链更新交易状态
- ✅ Worker 随 WalletServices 自动启动和停止

用户可以在交易历史中看到所有交易（包括普通转账和闪兑），体验统一且简洁。
