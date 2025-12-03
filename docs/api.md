# API 文档

## 概述

Wallet Services 提供 RESTful API 用于闪兑聚合、用户管理、资产管理等功能。

**Base URL**: `http://localhost:8080`

**认证方式**: JWT Token (部分接口需要)

---

## 闪兑聚合器 API

### 1. 获取报价

从多个 DEX 聚合器获取最优报价。

**端点**: `POST /api/v1/aggregator/quotes`

**请求头**:
```
Content-Type: application/json
```

**请求体**:
```json
{
  "from_chain_id": "1",           // 源链 ID (1=Ethereum, 56=BSC, 137=Polygon)
  "to_chain_id": "1",             // 目标链 ID
  "from_token": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",  // USDC 地址
  "to_token": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",   // WETH 地址
  "amount": "1000000",            // 1 USDC (6 decimals)
  "slippage_bps": 50,             // 0.5% 滑点 (50 basis points)
  "user_address": "0x..."         // 可选：用户地址
}
```

**字段说明**:
- `from_chain_id`: 源链 ID，支持的链 ID 见[链 ID 列表](#链-id-列表)
- `to_chain_id`: 目标链 ID，跨链交换时可与源链不同（需要 LiFi）
- `from_token`: 源代币合约地址（EVM）或 mint 地址（Solana）
- `to_token`: 目标代币合约地址
- `amount`: 源代币数量（最小单位，需考虑 decimals）
- `slippage_bps`: 滑点容忍度，单位为基点（1 bps = 0.01%）
- `user_address`: 用户钱包地址，某些 provider 需要此字段

**响应**:
```json
{
  "quote_id": "550e8400-e29b-41d4-a716-446655440000",
  "expires_at": "2024-01-01T12:05:00Z",
  "best_quote": {
    "provider": "lifi",
    "chain_type": "evm",
    "chain_id": "1",
    "from_token": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    "to_token": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
    "from_amount": "1000000",
    "to_amount": "500000000000000000",  // 0.5 WETH (18 decimals)
    "gas_estimate": "150000",
    "fees": "5000000000000000",         // 0.005 ETH
    "spender": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
    "router": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
    "raw": "{...}"                      // 原始 provider 响应
  },
  "alternatives": [
    {
      "provider": "1inch",
      "to_amount": "495000000000000000",  // 稍差的价格
      "gas_estimate": "180000",
      ...
    },
    {
      "provider": "0x",
      "to_amount": "490000000000000000",
      ...
    }
  ]
}
```

**响应字段说明**:
- `quote_id`: 报价 ID，用于后续准备交换
- `expires_at`: 报价过期时间（默认 5 分钟）
- `best_quote`: 最优报价（按 to_amount 排序）
- `alternatives`: 备选报价列表
- `provider`: 提供报价的聚合器名称
- `to_amount`: 预期获得的目标代币数量
- `gas_estimate`: Gas 估算值
- `spender`: EVM 链上需要 approve 的地址
- `router`: 执行 swap 的合约地址

**错误响应**:
```json
{
  "error": "no quotes available",
  "code": "NO_QUOTES"
}
```

**状态码**:
- `200 OK`: 成功
- `400 Bad Request`: 请求参数错误
- `500 Internal Server Error`: 服务器错误

---

### 2. 准备交换

根据选定的报价生成待签名的交易。

**端点**: `POST /api/v1/aggregator/swap/prepare`

**请求体**:
```json
{
  "quote_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
}
```

**响应**:
```json
{
  "swap_id": "660e8400-e29b-41d4-a716-446655440001",
  "actions": [
    {
      "step_index": 0,
      "action_type": "approve",
      "description": "Approve USDC spending",
      "chain_id": "1",
      "unsigned_tx": {
        "from": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
        "to": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
        "data": "0x095ea7b3000000000000000000000000...",
        "value": "0",
        "gas_limit": "50000",
        "gas_price": "20000000000",
        "nonce": "42"
      }
    },
    {
      "step_index": 1,
      "action_type": "swap",
      "description": "Swap USDC to WETH via LiFi",
      "chain_id": "1",
      "unsigned_tx": {
        "from": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
        "to": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
        "data": "0x7c025200...",
        "value": "0",
        "gas_limit": "200000",
        "gas_price": "20000000000",
        "nonce": "43"
      }
    }
  ]
}
```

**响应字段说明**:
- `swap_id`: 交换 ID，用于后续提交和查询
- `actions`: 需要执行的交易步骤数组
- `step_index`: 步骤索引（从 0 开始）
- `action_type`: 操作类型（`approve` 或 `swap`）
- `unsigned_tx`: 待签名的交易数据

**注意事项**:
- 如果代币已经授权足够额度，可能不会包含 `approve` 步骤
- 用户需要按顺序签名并提交每个步骤
- 每个步骤的 `nonce` 已自动设置

---

### 3. 提交已签名交易

提交用户签名后的交易到区块链。

**端点**: `POST /api/v1/aggregator/tx/submitSigned`

**请求体**:
```json
{
  "swap_id": "660e8400-e29b-41d4-a716-446655440001",
  "step_index": 0,
  "signed_tx": "0xf86c808504a817c800825208...",
  "idempotency_key": "user-generated-unique-key-12345"
}
```

**字段说明**:
- `swap_id`: 交换 ID（从 prepare 接口获取）
- `step_index`: 步骤索引（从 0 开始）
- `signed_tx`: 用户签名后的完整交易（RLP 编码的十六进制字符串）
- `idempotency_key`: 幂等性密钥，防止重复提交（建议使用 UUID）

**响应**:
```json
{
  "tx_hash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
}
```

**幂等性说明**:
- 如果使用相同的 `idempotency_key` 重复提交，将返回之前的 `tx_hash`
- 建议客户端为每次提交生成唯一的 `idempotency_key`
- 幂等性密钥在同一 `swap_id` + `step_index` 下有效

**错误响应**:
```json
{
  "error": "swap not found",
  "code": "SWAP_NOT_FOUND"
}
```

**状态码**:
- `200 OK`: 成功
- `400 Bad Request`: 请求参数错误
- `404 Not Found`: Swap 不存在
- `500 Internal Server Error`: 广播失败

---

### 4. 查询交换状态

查询交换的当前状态和所有步骤的执行情况。

**端点**: `GET /api/v1/aggregator/swap/status`

**查询参数**:
- `swap_id`: 交换 ID（必需）

**示例**:
```
GET /api/v1/aggregator/swap/status?swap_id=660e8400-e29b-41d4-a716-446655440001
```

**响应**:
```json
{
  "swap_id": "660e8400-e29b-41d4-a716-446655440001",
  "quote_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
  "state": "in_progress",
  "steps": [
    {
      "step_index": 0,
      "action_type": "approve",
      "tx_hash": "0xabc123...",
      "state": "confirmed",
      "submitted_at": "2024-01-01T12:00:00Z",
      "confirmed_at": "2024-01-01T12:00:30Z"
    },
    {
      "step_index": 1,
      "action_type": "swap",
      "tx_hash": "0xdef456...",
      "state": "submitted",
      "submitted_at": "2024-01-01T12:01:00Z"
    }
  ],
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:01:00Z"
}
```

**状态说明**:

**Swap 状态** (`state`):
- `pending`: 等待用户签名和提交
- `in_progress`: 至少一个步骤已提交，但未全部完成
- `completed`: 所有步骤已确认
- `failed`: 某个步骤失败

**Step 状态** (`steps[].state`):
- `pending`: 等待签名
- `submitted`: 已提交到区块链，等待确认
- `confirmed`: 已在链上确认
- `failed`: 失败

**失败响应**:
```json
{
  "swap_id": "660e8400-e29b-41d4-a716-446655440001",
  "state": "failed",
  "steps": [
    {
      "step_index": 0,
      "state": "failed",
      "fail_reason_code": "INSUFFICIENT_GAS",
      "fail_message": "Transaction ran out of gas"
    }
  ],
  "fail_reason_code": "STEP_FAILED",
  "fail_message": "Step 0 failed: Transaction ran out of gas"
}
```

---

## 状态码和错误处理

### HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 OK | 请求成功 |
| 400 Bad Request | 请求参数错误 |
| 401 Unauthorized | 未授权（需要 JWT token） |
| 404 Not Found | 资源不存在 |
| 429 Too Many Requests | 请求过于频繁 |
| 500 Internal Server Error | 服务器内部错误 |
| 503 Service Unavailable | 服务暂时不可用 |

### 错误响应格式

```json
{
  "error": "错误描述信息",
  "code": "ERROR_CODE",
  "details": {
    "field": "具体错误字段",
    "reason": "详细原因"
  }
}
```

### 常见错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| `INVALID_CHAIN_ID` | 不支持的链 ID | 检查链 ID 是否正确 |
| `INVALID_TOKEN_ADDRESS` | 无效的代币地址 | 检查代币地址格式 |
| `INSUFFICIENT_LIQUIDITY` | 流动性不足 | 减少交易数量或稍后重试 |
| `QUOTE_EXPIRED` | 报价已过期 | 重新获取报价 |
| `SWAP_NOT_FOUND` | 交换不存在 | 检查 swap_id 是否正确 |
| `NO_QUOTES` | 没有可用报价 | 检查代币对是否支持 |
| `BROADCAST_FAILED` | 交易广播失败 | 检查交易签名和网络状态 |

---

## 链 ID 列表

### EVM 链

| 链名称 | Chain ID | 网络 |
|--------|----------|------|
| Ethereum | 1 | Mainnet |
| Ethereum Goerli | 5 | Testnet |
| Ethereum Sepolia | 11155111 | Testnet |
| BSC | 56 | Mainnet |
| BSC Testnet | 97 | Testnet |
| Polygon | 137 | Mainnet |
| Polygon Mumbai | 80001 | Testnet |
| Arbitrum One | 42161 | Mainnet |
| Arbitrum Goerli | 421613 | Testnet |
| Optimism | 10 | Mainnet |
| Optimism Goerli | 420 | Testnet |
| Avalanche C-Chain | 43114 | Mainnet |
| Avalanche Fuji | 43113 | Testnet |

### Solana

| 链名称 | Chain ID | 网络 |
|--------|----------|------|
| Solana | solana | Mainnet |
| Solana Devnet | solana-devnet | Devnet |
| Solana Testnet | solana-testnet | Testnet |

---

## 速率限制

为保证服务稳定性，API 实施以下速率限制：

| 端点 | 限制 |
|------|------|
| `/api/v1/aggregator/quotes` | 10 次/分钟 |
| `/api/v1/aggregator/swap/prepare` | 20 次/分钟 |
| `/api/v1/aggregator/tx/submitSigned` | 30 次/分钟 |
| `/api/v1/aggregator/swap/status` | 60 次/分钟 |

超过限制将返回 `429 Too Many Requests`。

---

## 最佳实践

### 1. 报价获取

```javascript
// 获取报价
const quoteResponse = await fetch('/api/v1/aggregator/quotes', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    from_chain_id: '1',
    to_chain_id: '1',
    from_token: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
    to_token: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2',
    amount: '1000000',
    slippage_bps: 50,
    user_address: userAddress
  })
});

const quote = await quoteResponse.json();
console.log('Best quote:', quote.best_quote);
```

### 2. 完整交换流程

```javascript
// 1. 获取报价
const quote = await getQuote(params);

// 2. 准备交换
const swap = await prepareSwap(quote.quote_id, userAddress);

// 3. 签名并提交每个步骤
for (const action of swap.actions) {
  // 用户签名
  const signedTx = await wallet.signTransaction(action.unsigned_tx);

  // 提交交易
  const result = await submitSignedTx({
    swap_id: swap.swap_id,
    step_index: action.step_index,
    signed_tx: signedTx,
    idempotency_key: generateUUID()
  });

  console.log('Tx hash:', result.tx_hash);

  // 等待确认
  await waitForConfirmation(swap.swap_id, action.step_index);
}

// 4. 检查最终状态
const status = await getSwapStatus(swap.swap_id);
console.log('Swap completed:', status.state === 'completed');
```

### 3. 错误处理

```javascript
try {
  const quote = await getQuote(params);
} catch (error) {
  if (error.code === 'NO_QUOTES') {
    // 没有可用报价，可能是代币对不支持
    console.error('No quotes available for this token pair');
  } else if (error.code === 'INSUFFICIENT_LIQUIDITY') {
    // 流动性不足，建议减少交易量
    console.error('Insufficient liquidity, try smaller amount');
  } else {
    // 其他错误
    console.error('Error:', error.message);
  }
}
```

### 4. 轮询状态

```javascript
async function waitForConfirmation(swapId, stepIndex, maxAttempts = 60) {
  for (let i = 0; i < maxAttempts; i++) {
    const status = await getSwapStatus(swapId);
    const step = status.steps[stepIndex];

    if (step.state === 'confirmed') {
      return true;
    } else if (step.state === 'failed') {
      throw new Error(`Step failed: ${step.fail_message}`);
    }

    // 等待 5 秒后重试
    await sleep(5000);
  }

  throw new Error('Confirmation timeout');
}
```

---

## WebSocket API (计划中)

未来将支持 WebSocket 实时推送交易状态更新，避免频繁轮询。

```javascript
// 计划中的 WebSocket API
const ws = new WebSocket('ws://localhost:8080/api/v1/aggregator/ws');

ws.on('open', () => {
  ws.send(JSON.stringify({
    action: 'subscribe',
    swap_id: '660e8400-e29b-41d4-a716-446655440001'
  }));
});

ws.on('message', (data) => {
  const update = JSON.parse(data);
  console.log('Status update:', update);
});
```

---

## 相关文档

- [架构文档](./architecture.md) - 系统架构和设计
- [部署文档](./deployment.md) - 部署和运维指南
- [开发指南](./development.md) - 开发者指南
- [返回 README](../README.md)

---

**最后更新**: 2024-12-03

