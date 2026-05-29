# 收银台配置管理平台 — AI Code Agent 项目理解指南

> 本文档旨在让新的 AI Code Agent 或开发者快速理解项目全貌、架构约定和关键实现细节，避免重复探索。

---

## 一、项目定位

收银台配置管理平台是一个**内部后台管理系统**，用于管理收银台（Till）设备的配置、规则（Rule）、模块（Module）、门店（Store）、变量（Var）和分组（Group）。收银台设备通过 PowerShell 脚本在启动时上报状态（CheckIn），管理员通过 Web 前端查看、配置和导出 CSV。

**项目目录:** `C:\Users\Jingjing\收银台配置管理平台`

---

## 二、技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 后端框架 | **Gin** (Go) | HTTP 路由，`backend/cmd/server/main.go` 入口 |
| ORM | **GORM** | PostgreSQL 16，AutoMigrate 管理表结构 |
| API 文档 | **swaggo/swag + gin-swagger** | OpenAPI 2.0，注解驱动，`/swagger/*any` |
| 认证 | **golang-jwt + bcrypt** | JWT Bearer Token |
| 扩展认证 | **go-ldap + go-oidc** | LDAP（AD/OpenLDAP）和 SSO（Azure AD/Okta/Keycloak） |
| 日志 | **zap** | JSON 格式 stdout |
| 前端框架 | **React 18 + Ant Design 5** | `frontend/` |
| 状态管理 | **Zustand** + persist（localStorage） |
| HTTP | **Axios** | 前端 API 调用 |
| 构建 | **Vite** | 开发代理 `/api` → `:8080` |

---

## 三、目录结构

```
收银台配置管理平台/
├── backend/
│   ├── cmd/
│   │   ├── server/main.go        # 入口：初始化 DB、AutoMigrate、注册 Service/Handler、启动（含 swagger 注解）
│   │   ├── sync_groups/main.go   # 工具脚本：从 Module 数据同步 Group
│   │   └── initdb/main.go        # 工具脚本：创建 cashier_config 数据库
│   ├── docs/                     # swag init 自动生成的 OpenAPI 文档（需要提交）
│   ├── config.yaml               # 配置文件（DB、JWT、API Key、LDAP、SSO）
│   ├── internal/
│   │   ├── config/config.go      # 配置加载
│   │   ├── database/             # 数据库连接 + 种子数据 (admin/admin123)
│   │   ├── model/                # GORM 数据模型
│   │   ├── service/              # 业务逻辑层
│   │   ├── handler/              # HTTP 请求处理层（含 swagger 注解）
│   │   ├── middleware/           # JWT 认证中间件
│   │   └── router/router.go      # 路由注册（含 /swagger/*any）
│   └── generated/                # CSV 生成文件存储（按时间戳子目录）
├── frontend/
│   ├── src/
│   │   ├── api/                  # Axios API 调用（每个资源一个文件）
│   │   ├── store/                # Zustand 状态管理
│   │   ├── pages/                # 页面组件
│   │   ├── App.tsx               # 路由 + 侧边栏菜单
│   │   └── main.tsx              # React 入口
│   └── vite.config.ts            # 开发代理配置
├── win-till-modules/             # 收银台本地执行脚本（PowerShell）
│   ├── install.ps1               # 主安装脚本，构建 JSON → POST checkin
│   └── Assistive Tool/           # 辅助工具
├── 开发与运维指南.md              # 开发文档
├── 需求文档.md                    # 需求文档
└── 开发进度与待办.md               # 任务跟踪
```

---

## 四、核心数据模型

所有模型在 `backend/internal/model/` 中，GORM AutoMigrate 自动建表（`main.go:34`）。

| 模型 | 表名 | 用途 |
|------|------|------|
| User | users | 用户，支持 local/ldap/sso 三种 auth_source |
| Role | roles | 角色，Permissions 为 JSON 数组字符串 |
| Module | modules | 配置模块定义（名称、步骤列表等） |
| Rule | rules | 配置规则（条件 + 目标） |
| Store | stores | 门店 |
| TillList | till_lists | 收银台设备 registry |
| **SyncReport** | sync_reports | **每次 checkin 的同步报告（2026-05-28 新增）** |
| Var | vars | 配置变量（键值对） |
| Group | groups | 分组 |
| CsvGenerationLog | csv_generation_logs | CSV 生成历史 |
| OperationLog | operation_logs | 操作审计日志 |
| SSOState | sso_states | SSO 认证状态（防 CSRF） |

### TillList + SyncReport 关系（关键）

- **TillList** — 设备注册表，按 `host_name` 唯一索引，`request_body` 存最新一次上报的 JSON
- **SyncReport** — 每次 `POST /api/till-lists/checkin` 创建一条记录，`till_list_id` 关联到 TillList，`module_execution` 存储解析后的 ModuleExecution JSON 数组

```go
type SyncReport struct {
    ID              int64     // 主键
    TillListID      int64     // FK → till_lists.id
    RequestBody     string    // 完整原始 JSON
    ModuleExecution string    // 解析后的 ModuleExecution 数组 JSON
    Status          int       // 0=成功 1=失败
    Duration        int       // 总耗时（秒）
    SyncTime        string    // 同步时间（取 Execution.StartTime）
    CreatedAt       time.Time // 记录创建时间
}
```

---

## 五、架构模式

### 分层结构（后端）

```
handler/*.go     ← 解析请求参数、调用 Service、返回 JSON + 注册路由 (Register)
       ↓
service/*.go     ← 业务逻辑、数据库查询
       ↓
model/*.go       ← 数据模型定义
       ↓
router/router.go ← 只负责创建 Engine + 认证中间件，遍历 Handler 列表调用 Register
```

**依赖注入方式:** `main.go` 中手动创建 `*Service` → 传入 `*Handler` → 组装 `[]handler.Handler` → 传入 `Setup()`。

```go
// main.go 典型模式
tillListSvc := service.NewTillListService(db)
tillListH := handler.NewTillListHandler(tillListSvc)

r := router.Setup([]handler.Handler{
    authH, scriptH, dashH, csvGenH, moduleH, ruleH, storeH,
    tillListH, varH, groupH, userH, roleH, logH,
}, authSvc, cfg.APIKey)
```

**路由注册机制（`handler.Handler` 接口）:**

每个 handler 实现 `Register(api gin.IRouter, authed gin.IRouter)` 方法，在自己的文件中完成路由注册。`router.go` 不再维护路由表，只需遍历 `[]Handler` 调用 Register。

`router.go` 额外注册了 Swagger UI 路由（仅在主进程生效）:
```go
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

```go
// handler/registry.go
type Handler interface {
    Register(api gin.IRouter, authed gin.IRouter)
}

// handler/till_list.go
func (h *TillListHandler) Register(api gin.IRouter, authed gin.IRouter) {
    tills := authed.Group("/till-lists")
    tills.GET("", h.List)
    tills.POST("/checkin", h.CheckIn)
    // ...
}
```

### 统一返回格式

所有 API 返回 `{ code: 0, message: "ok", data: ... }`，错误时 `code` 为非零。

### Swagger 注解规范

每个 handler 方法上方添加 swaggo/swag 注解，格式：

```go
// @Summary      List till lists
// @Description  List till lists with optional filters
// @Tags         TillList
// @Security     ApiKeyAuth
// @Security     BearerAuth
// @Produce      json
// @Param        host_name  query  string  false  "Filter by host name"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /till-lists [get]
func (h *TillListHandler) List(c *gin.Context) { ...
```

- 公开路由不加 `@Security`（如 `/auth/login`）
- 需认证路由加两条 `@Security`：`ApiKeyAuth` 和 `BearerAuth`
- `Parameters()` 注释在 main.go 顶部定义 API 元信息：`@title`、`@version`、`@BasePath`、`@securityDefinitions` 等
- 路由路径以 `/api` 为根，如 `/till-lists` 对应实际路径 `/api/till-lists`

### 重新生成 API 文档

修改注解后需重新生成：

```bash
cd backend
swag init -g cmd/server/main.go -o docs
```

`backend/docs/` 中的文件需要提交到 git，编译时会嵌入二进制。

### 认证流程

支持两种认证方式：

**方式一：JWT Bearer Token**
1. `POST /api/auth/login` → 返回 JWT token
2. 前端 localStorage 存储 token
3. Axios 拦截器自动附加 `Authorization: Bearer <token>`
4. Gin 中间件 `middleware.JWTAuth()` 解析 token → 注入 `user_id`、`username`、`role_code` 到 context

**方式二：API Key**
1. 在 `config.yaml` 中配置 `api_key`（默认 `my-secret-api-key-2026`）
2. 请求时附加 `X-API-Key` 请求头
3. 中间件优先检查 `X-API-Key`，匹配则放行（`user_id=0`, `username="api-key"`, `role_code="api"`）
4. 不匹配或未提供时回退到 JWT 验证
5. 适用于收银台 checkin 等自动化场景

### 路由分组

```
公开（无需认证）:          注册在 api 分组
  POST /api/auth/login
  POST /api/auth/login/ldap
  GET  /api/auth/sso/login
  GET  /api/auth/sso/callback
  POST /api/auth/logout
  GET  /api/script-files

需认证（JWT 或 X-API-Key）: 注册在 authed 分组
  ├── /api/auth/refresh, /api/auth/permissions
  ├── /api/dashboard/stats
  ├── /api/system/users, roles, logs
  ├── /api/csv/generate, download, history, diff
  ├── /api/groups
  ├── /api/modules
  ├── /api/rules
  ├── /api/stores
  ├── /api/till-lists（含 /checkin、/reports?host_name=&mac_address=）
  └── /api/vars
```

---

## 六、关键 API 端点

### till-lists

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/till-lists` | 列表，支持 `host_name`、`location`、`env` 过滤 |
| GET | `/api/till-lists/:id` | 详情 |
| POST | `/api/till-lists/checkin` | 收银台上报，**同时创建 SyncReport** |
| GET | `/api/till-lists/:id/reports` | 获取该设备所有同步报告（按 id DESC） |
| GET | `/api/till-lists/:id/reports/:reportId` | 获取单条报告 |
| GET | `/api/till-lists/reports`?`host_name=&mac_address=` | **按主机名/MAC 查询所有设备的报告历史（按设备分组）** |
| POST | `/api/till-lists/import` | CSV 导入 |
| GET | `/api/till-lists/export/csv` | CSV 导出 |

### auth

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/login` | 登录，body: `{"username","password"}` |
| POST | `/api/auth/login/ldap` | LDAP 登录 |
| GET | `/api/auth/sso/login` | SSO OIDC 登录跳转 |
| GET | `/api/auth/sso/callback` | SSO 回调 |
| POST | `/api/auth/refresh` | 刷新 token |
| GET | `/api/auth/permissions` | 获取权限 |

### CSV

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/csv/generate` | 生成所有类型 CSV |
| POST | `/api/csv/generate/:type` | 生成指定类型 |
| GET | `/api/csv/download/:type` | 下载最新 |
| GET | `/api/csv/download/all` | 下载全部（zip） |
| GET | `/api/csv/history` | 生成历史 |
| GET | `/api/csv/diff/:type` | 版本对比 |

---

## 七、CheckIn 完整流程

这是项目最核心的业务流程，涉及收银台客户端和后端两个部分：

```
┌─ 收银台 (PowerShell install.ps1) ─────────────────────────┐
│ 1. 执行各模块安装脚本                                       │
│ 2. 收集 Execution 结果（ModuleExecution[].Steps[].Output）  │
│ 3. 构建 JSON（含 Ip, Name, MACAddress, Execution, Fact 等） │
│ 4. POST /api/till-lists/checkin                            │
└────────────────────────────────────────────────────────────┘
                            ↓
┌─ 后端 (TillListService.CheckIn) ───────────────────────────┐
│ 1. json.Unmarshal 解析请求体                                │
│ 2. 提取字段 (Name→host_name, MACAddress→mac_address 等)    │
│ 3. 按 host_name 查找已有 TillList                           │
│    ├─ 不存在 → 创建新 TillList                              │
│    └─ 存在 → 更新 TillList 字段                             │
│ 4. 解析 ModuleExecution 为 JSON 字符串                      │
│ 5. 创建 SyncReport（保留完整历史）                          │
└────────────────────────────────────────────────────────────┘
```

**关键细节:**
- `host_name` 是唯一标识（对应 JSON 的 `Name` 字段）
- `TillList.request_body` 始终存储最新的原始 JSON（覆盖更新）
- `SyncReport` 每次 checkin 追加一条（不覆盖），形成执行历史
- JSON 中 `Group` 字段须为数字（`float64` → `int64`），不能是字符串
- `mac_address` 优先取 JSON 顶层 `MACAddress`，若不存在则回退到 `Fact` 数组中 `Key === "MACAddress"` 的值

---

## 八、前端页面结构

| 路由 | 页面 | 说明 |
|------|------|------|
| `/dashboard` | Dashboard.tsx | 仪表盘统计 |
| `/modules` | ModuleList.tsx | 模块列表 |
| `/modules/:id` | ModuleDetail.tsx | 模块详情 |
| `/rules` | RuleList.tsx | 规则列表 |
| `/rules/:id` | RuleDetail.tsx | 规则详情 |
| `/stores` | StoreList.tsx | 门店列表 |
| `/stores/:id` | StoreDetail.tsx | 门店详情 |
| `/till-lists` | TillListPage.tsx | 设备列表（搜索/导出/导入） |
| `/till-lists/:id` | TillListDetail.tsx | 设备详情（基本信息 + 同步报告时间线） |
| `/vars` | VarList.tsx | 变量列表 |
| `/groups` | GroupList.tsx | 分组列表 |
| `/csv-generate` | CsvGenerate.tsx | CSV 生成 |
| `/system/users` | UserList.tsx | 用户管理 |
| `/system/roles` | RoleList.tsx | 角色管理 |
| `/system/logs` | OperationLog.tsx | 操作日志 |
| `/auth/sso/callback` | SSOCallback.tsx | SSO 回调处理 |

前端状态管理全部使用 **Zustand**，每个资源一个 store 文件（`frontend/src/store/*.ts`），API 调用在 `frontend/src/api/*.ts`。

---

## 九、常见的开发任务指引

### 新增一个模块（含 CRUD）

按以下顺序创建/修改 7 个文件（router.go 无需改动）：

1. `model/*.go` — 数据模型 + 添加 `TableName()` 方法
2. `cmd/server/main.go` — AutoMigrate 注册
3. `service/*.go` — 业务逻辑
4. `handler/*.go` — HTTP handler + 实现 `Register()` 方法注册路由
5. `frontend/src/api/*.ts` — API 类型和调用函数
6. `frontend/src/pages/*.tsx` — 列表页和详情页
7. `frontend/src/App.tsx` — 注册路由 + 侧边栏菜单
8. `cmd/server/main.go` — 往 `[]handler.Handler` 切片中追加一项

### 修改数据库字段

编辑 Model 结构体 → 重启后端 → GORM AutoMigrate 自动新增列（**不会删除/修改已有列**）。

### 修改 CSV 导出/导入

编辑对应 Service 中的 `ExportCSV` / `ImportCSV` 方法。

### 调试 checkin 数据

1. 查看 `till_lists` 表的 `request_body` 字段（最新原始 JSON）
2. 查看 `sync_reports` 表的 `module_execution`（解析后的步骤数据）
3. 前端访问 `/till-lists/:id` → 「同步报告」标签页 → 时间线选择报告 → 查看步骤详情

---

## 十、已知的坑和注意事项

1. **JSON 中的 CRLF 问题** — request_body.json 文件可能有 Windows CRLF 行尾（`\r\n`），Go 的 `json.Unmarshal` 会拒绝字符串中的 `\r`。如需手动导入需先转义。
2. **中文路径问题** — 项目目录名含中文 `收银台配置管理平台`，PowerShell 的 `-Command` 参数中引用该路径可能失败（FileNotFoundException），需复制到无中文路径处理。
3. **Go 后端认证** — 所有 `authed` 组路由都需要认证。收银台 checkin 端点使用 `X-API-Key` 请求头（配置在 `config.yaml`），无需 JWT 登录。如果收银台脚本未携带 API Key，请求会被 401 拒绝。
4. **`golang-migrate` 已引入但未启用** — 当前使用 AutoMigrate 管理表结构，`golang-migrate` 依赖已添加可后续启用版本化迁移。
5. **windows 日志与 JSON** — 后端使用 zap 输出 JSON 格式日志到 stdout，GORM 调试日志会直接打印 SQL。

---

## 十一、启动方式

```powershell
# 后端
cd backend
go run cmd/server/main.go    # 监听 :8080，Swagger UI: http://localhost:8080/swagger/index.html

# 前端
cd frontend
npm run dev                   # 监听 :3000，代理 /api → :8080

# 默认登录
用户名: admin  密码: admin123
```

---

## 十二、测试账号与数据

- 种子数据自动创建（首次启动）：admin (super_admin) / admin123
- 测试 checkin 数据：`request_body.json`（需修复 CRLF 后通过 checkin API 导入）
- 已验证的导入方式：Python 脚本修复 → POST `/api/till-lists/checkin`
