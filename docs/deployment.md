# 部署指南

## 部署方式

支持两种方式：

1. **容器部署**（推荐）— Docker Compose 或 Podman Compose 一键部署
2. **手动部署** — 分别运行后端二进制和前端 nginx

## 容器部署

### 前置条件

- Docker 20.10+ / Podman 4.0+ 及对应 Compose 插件
- 至少 2GB RAM

### 步骤

```bash
# 1. 克隆仓库
git clone https://github.com/<owner>/ai-interview.git
cd ai-interview

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env，设置以下变量：
#   DB_DSN        - MySQL 连接串 (容器内使用 mysql 容器名)
#   JWT_SECRET    - JWT 签名密钥 (随机字符串)
#   ENCRYPTION_KEY - AES-256 密钥 (64个十六进制字符)

# 3. 启动全部服务
docker compose up -d
# 或
podman compose up -d

# 4. 检查状态
docker compose ps
docker compose logs -f backend
```

访问 http://localhost（前端）

### 服务端口

| 服务 | 容器端口 | 主机端口 |
|------|---------|---------|
| frontend (nginx) | 80 | 80 |
| backend (Go) | 8000 | 8000 |
| MySQL | 3306 | 3306 |
| Redis | 6379 | 6379 |

### 数据持久化

- `mysql_data` — MySQL 数据目录
- `redis_data` — Redis 数据目录

### 构建参数

Backend Dockerfile 支持 `VERSION` build-arg:

```bash
docker compose build --build-arg VERSION=v1.0.0
```

## 开发模式

开发时只需运行基础设施（MySQL + Redis）：

```bash
# 启动 MySQL + Redis，跳过 backend/frontend 容器
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

然后本地运行：

```bash
# 后端
make dev

# 前端
cd frontend && npm run dev
```

## 手动部署

### 后端

```bash
# 编译
cd backend
CGO_ENABLED=0 go build -ldflags "-X main.Version=v1.0.0" -o bin/server ./cmd/server/

# 部署文件
#   bin/server         — 二进制
#   configs/config.yaml — 配置文件
#   .env               — 环境变量

# 运行
./bin/server -conf configs/
```

### 前端

```bash
cd frontend
npm ci
npm run build-only
# dist/ 目录用 nginx 托管，参考 frontend/nginx.conf
```

Nginx 关键配置：
- `try_files $uri /index.html` — Vue Router history mode
- `/api/` 反向代理到后端
- WebSocket 升级支持

## 使用 GHCR 镜像

CI 自动构建的镜像：

```bash
docker pull ghcr.io/<owner>/ai-interview/backend:latest
docker pull ghcr.io/<owner>/ai-interview/frontend:latest
```

标签格式：`v1.0.0`、`1.0`、`sha-abc1234`

## 数据库迁移

SQL 迁移文件在 `backend/sql/migrations/`：

```bash
# 上行迁移
make migrate-up

# 下行迁移
make migrate-down
```

容器部署时 MySQL 初次启动会自动创建 `ai_interview` 数据库，但表需要手动迁移。

## 环境变量说明

| 变量 | 必须 | 说明 | 示例 |
|------|------|------|------|
| DB_DSN | ✅ | MySQL 连接串 | `root:password@tcp(mysql:3306)/ai_interview?parseTime=true&charset=utf8mb4` |
| REDIS_ADDR | | Redis 地址 (容器部署时设置) | `redis:6379` |
| REDIS_PASSWORD | | Redis 密码 | |
| JWT_SECRET | ✅ | JWT 签名密钥 | 随机 32+ 字符 |
| ENCRYPTION_KEY | ✅ | AES-GCM 密钥 | 64 hex chars (= 32 bytes) |

## 健康检查

```bash
curl http://localhost:8000/health
# {"status":"ok"}
```
