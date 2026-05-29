# 收银台配置管理平台 (Cashier Configuration Management Platform)

收银台终端配置的统一管理平台，支持多门店、多设备的配置远程下发、模块执行监控与同步报告查询。

## 技术栈

| 用途 | 技术 |
|------|------|
| 后端框架 | Go + Gin |
| 前端框架 | React 18 + TypeScript + Ant Design 5 |
| 构建工具 | Vite |
| 数据库 | PostgreSQL 16 |
| ORM | GORM |
| 认证 | JWT + API Key + LDAP/OIDC (可选) |
| API 文档 | swaggo/swag + gin-swagger (OpenAPI 2.0) |
| 状态管理 | Zustand |
| CSV 处理 | Go encoding/csv + papaparse (浏览器端) |

## 项目结构

```
收银台配置管理平台/
├── backend/                    # Go 后端
│   ├── cmd/server/main.go      # 入口
│   ├── docs/                   # Swagger 自动生成的 API 文档
│   ├── config.yaml             # 配置文件（本地开发用）
│   ├── internal/
│   │   ├── config/             # 配置加载
│   │   ├── database/           # 数据库连接 + 种子数据
│   │   ├── model/              # 数据模型
│   │   ├── service/            # 业务逻辑层
│   │   ├── handler/            # HTTP 请求处理层 + 路由注册
│   │   ├── middleware/         # JWT / API Key 认证中间件
│   │   └── router/             # 路由组装（含 /swagger/*any）
│   └── generated/              # CSV 生成文件（构建产物，不追踪）
├── frontend/                   # React 前端
│   ├── src/
│   │   ├── api/                # API 调用封装
│   │   ├── store/              # Zustand 状态管理
│   │   ├── pages/              # 页面组件
│   │   ├── components/         # 通用组件
│   │   └── App.tsx             # 路由 + 布局
│   └── vite.config.ts          # Vite 配置（含 API 代理）
├── csv/                        # 源数据 CSV 文件
├── data.yaml                   # CSV 数据源
└── Readme.md                   # 本文件
```

## 快速开始

### 前置要求

- Go 1.22+
- Node.js 18+
- PostgreSQL 16
- 创建数据库 `cashier_config`

### 启动后端

```bash
cd backend
# 修改 config.yaml 中的数据库配置
go run cmd/server/main.go
# 监听 :8080，首次启动自动建表 + 种子数据（admin / admin123）
# Swagger UI: http://localhost:8080/swagger/index.html
```

### 启动前端

```bash
cd frontend
npm install
npm run dev
# 监听 :3000，/api 自动代理到 :8080
```

## 认证方式

支持两种认证方式：

| 方式 | Header | 适用场景 |
|------|--------|----------|
| JWT Bearer Token | `Authorization: Bearer <token>` | 浏览器登录 |
| API Key | `X-API-Key: <key>` | 脚本/自动化（如收银台 checkin） |

> 默认 API Key: `my-secret-api-key-2026`（生产环境请修改）

## 生产部署

### 后端

```bash
cd backend
go build -o server cmd/server/main.go
# 将二进制文件和 config.yaml（修改为生产配置）部署到服务器
./server
```

### 前端

```bash
cd frontend
npm run build
# 输出到 dist/，部署到 Nginx，反向代理 /api 到后端
```

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;

    root /path/to/frontend/dist;
    index index.html;

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

## Swagger API 文档

集成 swaggo/swag，启动后端后访问：

```
http://localhost:8080/swagger/index.html
```

支持交互式调试，可在 Swagger UI 中直接发送请求测试 API。支持两种认证方式：

1. **API Key** — 点击 `Authorize`，输入 `X-API-Key` 的值
2. **Bearer Token** — 先调用 `/api/auth/login`，然后将 token 填入 `Authorize`

### 更新 API 文档

修改 handler 注解后，在 `backend/` 目录下重新生成：

```bash
cd backend
swag init -g cmd/server/main.go -o docs
```

## API 概览

主要路由分组：

| 分组 | 路径 | 说明 |
|------|------|------|
| 认证 | `/api/auth/*` | 登录、登出、刷新令牌 |
| 门店 | `/api/stores/*` | 门店 CRUD |
| 设备 | `/api/devices/*` | 设备管理 |
| 收银台 | `/api/till-lists/*` | 收银台设备 + CheckIn + 同步报告 |
| CSV 数据 | `/api/csv/{modules,rules,stores,tills,vars}/*` | CSV 数据 CRUD |
| 生成 | `/api/csv/generate*`, `/api/csv/download*` | CSV 文件生成与下载 |
| 用户 | `/api/system/users*` | 用户管理 |
| 角色 | `/api/system/roles*` | 角色管理 |

所有 API 统一返回格式: `{ code: 0, message: "ok", data: ... }`

## 高可用部署

### 架构

```
Nginx LB → Go 实例 xN → PgBouncer → PostgreSQL (Patroni 集群)
```

- Go 后端无状态，可水平扩展
- 通过 Nginx 健康检查自动摘除故障节点
- PostgreSQL 使用 Patroni + etcd 实现自动主从切换
- PgBouncer 减少 PG 连接数压力

### 并发承载（估算）

| 实例数 | CheckIn 并发 | 管理页面并发 |
|--------|-------------|-------------|
| 1      | ~3,000      | ~500        |
| 2      | ~6,000      | ~1,000      |
| 4      | ~12,000     | ~2,000      |

> 瓶颈在 PostgreSQL 写入能力。详情见 `开发与运维指南.md → 四、高可用部署方案`。

## License

Internal Use Only
