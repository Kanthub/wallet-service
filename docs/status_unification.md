# 状态系统统一设计

## 📋 背景

之前系统中存在两套状态系统：
1. **SwapState / StepState**（字符串）- 用于内存中的 Swap 对象
   - `PENDING` / `SUBMITTED` / `CONFIRMED` / `FAILED`
2. **TxStatus**（整数 0/1/2/3）- 用于数据库持久化
   - `0=CREATED` / `1=PENDING` / `2=FAILED` / `3=SUCCESS`

这导致了状态不一致、重复维护、映射复杂等问题。

---

## ✅ 解决方案：统一使用 TxStatus

### 核心改动

**废弃 SwapState 和 StepState，全部使用 TxStatus (0/1/2/3)**

---

## 📊 统一状态定义

### TxStatus 常量

```go
// services/api/models/backend/aggregator.go
const (
    TxStatusCreated = 0 // CREATED: 后端收到 signedTx 请求并写入记录，但尚未广播
    TxStatusPending = 1 // PENDING: 广播成功并拿到 txHash
    TxStatusFailed  = 2 // FAILED: 广播失败或链上执行失败或超时
    TxStatusSuccess = 3 // SUCCESS: 链上确认成功
)

// TxStatusNames provides human-readable names for status codes
var TxStatusNames = map[int]string{
    TxStatusCreated: "CREATED",
    TxStatusPending: "PENDING",
    TxStatusFailed:  "FAILED",
    TxStatusSuccess: "SUCCESS",
}
```

### 状态流转

```
CREATED (0) → PENDING (1) → SUCCESS (3)
     ↓             ↓
   FAILED (2) ← FAILED (2)
```

---

## 🔧 数据模型改动

### Step 结构体

```go
// Before
type Step struct {
    StepIndex  int
    ActionType ActionType
    TxHash     string
    State      StepState  // ❌ 字符串类型
    // ...
}

// After
type Step struct {
    StepIndex  int
    ActionType ActionType
    TxHash     string
    Status     int        // ✅ 整数类型: 0=CREATED, 1=PENDING, 2=FAILED, 3=SUCCESS
    // ...
}
```

### Swap 结构体

```go
// Before
type Swap struct {
    SwapID      string
    QuoteID     string
    UserAddress string
    WalletUUID  string
    State       SwapState  // ❌ 字符串类型
    Steps       []*Step
    // ...
}

// After
type Swap struct {
    SwapID      string
    QuoteID     string
    UserAddress string
    WalletUUID  string
    Status      int        // ✅ 整数类型: 0=CREATED, 1=PENDING, 2=FAILED, 3=SUCCESS
    Steps       []*Step
    // ...
}
```

### SwapStatusResponse

```go
// Before
type SwapStatusResponse struct {
    SwapID         string
    State          SwapState  // ❌ 字符串类型
    Steps          []*Step
    FailReasonCode string
    FailMessage    string
}

// After
type SwapStatusResponse struct {
    SwapID         string
    Status         int        // ✅ 整数类型: 0=CREATED, 1=PENDING, 2=FAILED, 3=SUCCESS
    Steps          []*Step
    FailReasonCode string
    FailMessage    string
}
```

---

## 📡 API 响应示例

### 旧格式（字符串状态）

```json
{
  "swap_id": "550e8400-e29b-41d4-a716-446655440000",
  "state": "SUBMITTED",
  "steps": [
    {
      "step_index": 0,
      "action_type": "APPROVE",
      "state": "CONFIRMED",
      "tx_hash": "0xabc123..."
    },
    {
      "step_index": 1,
      "action_type": "SWAP",
      "state": "PENDING",
      "tx_hash": "0xdef456..."
    }
  ]
}
```

### 新格式（整数状态）

```json
{
  "swap_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": 1,
  "steps": [
    {
      "step_index": 0,
      "action_type": "APPROVE",
      "status": 3,
      "tx_hash": "0xabc123..."
    },
    {
      "step_index": 1,
      "action_type": "SWAP",
      "status": 1,
      "tx_hash": "0xdef456..."
    }
  ]
}
```

---

## 💡 前端处理建议

前端可以定义状态名称映射：

```typescript
const STATUS_NAMES = {
  0: 'Created',
  1: 'Pending',
  2: 'Failed',
  3: 'Success'
}

const STATUS_COLORS = {
  0: 'gray',
  1: 'blue',
  2: 'red',
  3: 'green'
}

// 使用示例
function renderStatus(status: number) {
  return (
    <Badge color={STATUS_COLORS[status]}>
      {STATUS_NAMES[status]}
    </Badge>
  )
}
```

---

## ✅ 优势

1. **单一数据源** - API 和数据库使用相同状态
2. **状态定义清晰** - 整数状态更简洁
3. **无需状态映射** - 减少转换逻辑
4. **简化代码逻辑** - 减少状态同步错误
5. **性能更好** - 整数比较比字符串比较更快
6. **数据库友好** - 整数索引效率更高

---

## 📝 修改清单

### 代码文件
- ✅ `services/api/models/backend/aggregator.go` - 删除 SwapState/StepState，添加 TxStatus 常量
- ✅ `services/api/service/aggregator_service.go` - 所有状态操作改为使用 TxStatus
- ✅ `services/api/service/aggregator_service_test.go` - 测试用例更新
- ✅ `services/api/aggregator/store/swap_store.go` - 无需修改（状态无关）

### 文档文件
- ✅ `docs/status_unification.md` - 本文档
- 🔄 `docs/api.md` - 需要更新 API 响应示例
- 🔄 `docs/swap_history.md` - 已经使用 TxStatus，无需修改
- 🔄 `docs/operation_multi_step_design.md` - 需要更新状态说明

