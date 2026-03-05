# User-Center

基于 `go-zero` 重组后的用户中心示例，当前将用户能力拆成了两个独立服务：

- `apps/user-api`: HTTP 接入层，负责参数解析、统一响应、会话管理、RPC 转发
- `apps/user-rpc`: 用户领域 RPC 服务，负责登录、注册、密码查询等核心逻辑

## 目录

```text
user-center/
├── common/         # 通用错误、响应、session、context、拦截器
├── apps/
│   ├── user-api/   # go-zero api 服务
│   └── user-rpc/   # goctl 生成的 rpc 服务
│       └── internal/
│           ├── dao/        # 原 dao 已移动到 rpc 内部
│           └── httpclient/ # 原 httpclient 已移动到 rpc 内部
├── deploy/         # docker / k8s / scripts 占位目录
└── sql/            # 用户表、session 表初始化脚本
```

## Run

项目根目录使用一份统一配置：[config.yaml](/Users/mgg/Project/Golang/User-Center/config.yaml)。

先启动 RPC 服务：

```bash
go run ./apps/user-rpc
```

再启动 API 服务：

```bash
go run ./apps/user-api
```

## Verify

```bash
go test ./common/... ./apps/user-api/... ./apps/user-rpc/...
go build ./...
```

更新 DAO 生成代码：

```bash
go run ./cmd/gen
```

## Notes

- `apps/user-rpc/user.go`、`apps/user-rpc/internal/server`、`apps/user-rpc/usercenterservice` 和 `apps/user-rpc/pb` 均由 `goctl rpc protoc` 生成。
- `apps/user-api` 通过 goctl 生成的 HTTP 骨架承接路由和逻辑，并直接使用 `apps/user-rpc/usercenterservice` 作为 RPC client。
- `UserAPI` 和 `UserRPC` 配置已经合并到根目录 `config.yaml`，各服务只读取自己的配置段。
- `UserRPC.Mysql` 已展开为独立字段，运行时会自动拼成 DSN。
- 会话实现默认使用签名 Cookie；如需切换服务端 Session，可参考 [sql/session.sql](/Users/mgg/Project/Golang/User-Center/sql/session.sql)。
