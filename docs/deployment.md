# 部署文档

## 概述

部署 Wallet Services

---

## 前置要求

### 软件依赖

- **Go**: 1.21+
- **PostgreSQL**: 12+
- **wallet-chain-account**: gRPC 服务（用于区块链交互）

### 可选依赖

- **Redis**: 用于分布式缓存（计划中）
- **Nginx**: 用于反向代理和负载均衡
- **Docker**: 容器化部署
- **Docker Compose**: 多容器编排

---

## 本地开发环境

### 1. 克隆代码

```bash
git clone https://github.com/roothash-pay/wallet-services.git
cd wallet-services
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置数据库

```bash
# 创建数据库
createdb wallet_services

# 运行迁移（如果有）
# go run cmd/migrate/main.go
```

### 4. 创建配置文件

创建 `config.yaml`:

```yaml
# HTTP 服务器
http_server:
  host: "0.0.0.0"
  port: 8080

# 数据库
master_db:
  host: "localhost"
  port: 5432
  name: "wallet_services"
  user: "postgres"
  password: "your_password"

# JWT
jwt_secret: "your-dev-secret-key"

# 闪兑聚合器
aggregator_config:
  wallet_account_addr: "127.0.0.1:8189"
  
  # 0x Protocol
  zerox_api_url: "https://api.0x.org"
  zerox_api_key: ""
  
  # 1inch
  oneinch_api_url: "https://api.1inch.dev"
  oneinch_api_key: ""
  
  # Jupiter
  jupiter_api_url: "https://quote-api.jup.ag"
  
  # LiFi
  lifi_api_url: "https://li.quest/v1"
  lifi_api_key: ""
  
  enable_providers:
    0x: false
    1inch: false
    jupiter: false
    lifi: true
```

### 5. 启动服务

```bash
# 方式 1: 直接运行
go run cmd/wallet-services/main.go -config config.yaml

# 方式 2: 编译后运行
go build -o wallet-services ./cmd/wallet-services
./wallet-services -config config.yaml
```

### 6. 验证服务

```bash
# 健康检查
curl http://localhost:8080/health

# 测试 API
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

---

## Docker 部署

### 1. 创建 Dockerfile

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wallet-services ./cmd/wallet-services

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 复制编译好的二进制文件
COPY --from=builder /app/wallet-services .

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./wallet-services", "-config", "/config/config.yaml"]
```

### 2. 构建镜像

```bash
docker build -t wallet-services:latest .
```

### 3. 运行容器

```bash
docker run -d \
  --name wallet-services \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/config/config.yaml:ro \
  -e LOG_LEVEL=info \
  wallet-services:latest
```

### 4. 查看日志

```bash
docker logs -f wallet-services
```

---

## Docker Compose 部署

### 1. 创建 docker-compose.yml

```yaml
version: '3.8'

services:
  wallet-services:
    build: .
    ports:
      - "8080:8080"
    environment:
      - CONFIG_PATH=/config/config.yaml
    volumes:
      - ./config.yaml:/config/config.yaml:ro
    depends_on:
      - postgres
    restart: unless-stopped

  postgres:
    image: postgres:14-alpine
    environment:
      POSTGRES_DB: wallet_services
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### 2. 启动服务

```bash
docker-compose up -d
```

---

## 生产环境部署

### Systemd 服务

创建 `/etc/systemd/system/wallet-services.service`:

```ini
[Unit]
Description=Wallet Services API
After=network.target

[Service]
Type=simple
User=wallet
WorkingDirectory=/opt/wallet-services
ExecStart=/opt/wallet-services/wallet-services -config /opt/wallet-services/config.yaml
Restart=always

[Install]
WantedBy=multi-user.target
```

启动:

```bash
sudo systemctl start wallet-services
sudo systemctl enable wallet-services
```

---

## 监控和日志

### 日志配置

使用 journald 查看日志:

```bash
sudo journalctl -u wallet-services -f
```

### 健康检查

```bash
curl http://localhost:8080/health
```

---

## 相关文档

- [API 文档](./api.md)
- [架构文档](./architecture.md)
- [返回 README](../README.md)

---

**最后更新**: 2024-12-03

