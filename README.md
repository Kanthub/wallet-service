# Wallet Services

企业级钱包服务平台，提供多链资产管理、闪兑聚合、用户认证等核心功能。

## 📋 目录

- [功能特性](#功能特性)
- [快速开始](#快速开始)
- [闪兑聚合器](#闪兑聚合器-swap-aggregator)
- [项目结构](#项目结构)
- [文档](#文档)
- [开发路线图](#开发路线图)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

---

## 🚀 功能特性

### 核心功能

- **🔄 闪兑聚合器 (Swap Aggregator)**
  - 聚合多个 DEX 提供商（0x, 1inch, Jupiter, LiFi）
  - 自动选择最优报价
  - 支持 EVM 链和 Solana
  - 跨链交换支持（通过 LiFi）
  - 交易状态追踪和幂等性保证

- **👤 用户管理**
  - JWT 认证
  - 角色权限管理
  - 多因素认证支持

- **💼 资产管理**
  - 多链钱包支持
  - 资产查询和转账
  - 交易历史记录

- **🔗 区块链集成**
  - 通过 gRPC 与 wallet-chain-account 服务通信
  - 支持交易广播和查询
  - 多链支持（Ethereum, BSC, Polygon, Arbitrum, Solana 等）

---

## 🚀 快速开始

### 前置要求

- Go 1.21+
- PostgreSQL 12+
- wallet-chain-account 服务（用于区块链交互）

### 安装和运行

```bash
# 克隆仓库
git clone https://github.com/roothash-pay/wallet-services.git
cd wallet-services

# 安装依赖
go mod download

# 创建配置文件（参考 config.example.yaml）
cp config.example.yaml config.yaml
# 编辑 config.yaml 填入你的配置

# 运行服务
go run cmd/wallet-services/main.go -config config.yaml
```

服务将在 `http://localhost:8080` 启动。

详细的部署说明请参考 [部署文档](./docs/deployment.md)。

---

## 🔄 闪兑聚合器 (Swap Aggregator)

闪兑聚合器是一个智能交易路由系统，它并发查询多个 DEX 聚合器，自动选择最优价格，生成交易并追踪状态直到完成。

### 支持的 Provider

| Provider | 链类型 | 特性 | 状态 |
|----------|--------|------|------|
| **0x Protocol** | EVM | 聚合多个 DEX | ⚠️ 待实现 |
| **1inch** | EVM | 最大的 DEX 聚合器 | ⚠️ 待实现 |
| **Jupiter** | Solana | Solana 最大聚合器 | ⚠️ 待实现 |
| **LiFi** | EVM + 跨链 | 支持 20+ 链和跨链 | ✅ 已实现 |

### API 示例

```bash
# 获取报价
curl -X POST http://localhost:8080/api/v1/aggregator/quotes \
  -H "Content-Type: application/json" \
  -d '{
    "from_chain_id": "1",
    "to_chain_id": "1",
    "from_token": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    "to_token": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
    "amount": "1000000",
    "slippage_bps": 50
  }'
```

详细的架构说明和 API 文档请参考：
- [架构文档](./docs/architecture.md) - 系统架构和设计
- [API 文档](./docs/api.md) - 完整的 API 参考



---

## 📂 项目结构

```
wallet-services/
├── cmd/wallet-services/            # 入口文件
├── config/                         # 配置结构
├── services/
│   ├── api/
│   │   ├── aggregator/             # 闪兑聚合器
│   │   │   ├── provider/           # DEX Provider 实现
│   │   │   │   ├── evm/            # EVM 链 providers
│   │   │   │   │   ├── 0x.go
│   │   │   │   │   ├── 1inch.go
│   │   │   │   │   └── lifi.go     # ✅ LiFi 实现
│   │   │   │   └── solana/         # Solana providers
│   │   │   ├── store/              # 存储层
│   │   │   └── utils/              # 工具函数
│   │   ├── models/                 # 数据模型
│   │   ├── routes/                 # HTTP 路由
│   │   └── service/                # 业务逻辑
│   └── grpc_client/                # gRPC 客户端
├── docs/                           # 📚 文档
│   ├── api.md                      # API 文档
│   ├── architecture.md             # 架构文档
│   └── deployment.md               # 部署文档
└── proto/                          # gRPC proto 定义
```

---

## 📚 文档

- **[API 文档](./docs/api.md)** - 完整的 RESTful API 参考
- **[架构文档](./docs/architecture.md)** - 系统架构和设计说明
- **[部署文档](./docs/deployment.md)** - 部署和运维指南
- **[闪兑历史持久化](./docs/swap_history.md)** - 交易历史保存机制

---

## 📝 开发路线图

### 已完成 ✅
- [x] 基础架构搭建
- [x] 闪兑聚合器框架
- [x] LiFi Provider 完整实现
- [x] gRPC 客户端集成
- [x] 交易状态追踪
- [x] 幂等性保证

### 进行中 🚧
- [ ] 用户认证系统
- [ ] 0x Protocol Provider 实现
- [ ] 1inch Provider 实现
- [ ] Jupiter Provider 实现
- [ ] 单元测试覆盖
- [ ] 集成测试

### 计划中 📋
- [ ] Redis 缓存支持
- [ ] WebSocket 实时推送
- [ ] 交易历史查询
- [ ] 费用估算优化
- [ ] 跨链桥接支持
- [ ] MEV 保护
- [ ] Gas 优化建议

---

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 代码规范

- 遵循 Go 官方代码风格
- 使用 `gofmt` 格式化代码
- 添加必要的注释和文档
- 编写单元测试

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

## 📞 联系方式

- 项目主页: https://github.com/roothash-pay/wallet-services
- 问题反馈: https://github.com/roothash-pay/wallet-services/issues

---

## 🙏 致谢

感谢以下开源项目：

- [0x Protocol](https://0x.org/) - DEX 聚合协议
- [1inch](https://1inch.io/) - DEX 聚合器
- [Jupiter](https://jup.ag/) - Solana 聚合器
- [LiFi](https://li.fi/) - 跨链聚合器
- [wallet-chain-account](https://github.com/dapplink-labs/wallet-chain-account) - 区块链账户服务

---

**Built with ❤️ by RootHash Pay Team**