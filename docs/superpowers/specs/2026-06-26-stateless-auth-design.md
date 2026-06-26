# wlt 无状态鉴权重构设计

日期:2026-06-26 ｜ 状态:已确认,待实现

## 背景

当前 `wlt` 的鉴权是「交互式 `wlt auth login` → token 存入 `~/.wlt/config.yaml` → 客户端自动刷新」。
作为 AI Agent 原生 CLI,这套交互式登录对 Agent 不友好:Agent 每个会话都要先跑 `auth status`/`auth login`,
且 token 持久化在磁盘上既不安全也无必要(Agent 自己持有新鲜 token)。

## 目标

改成**无状态鉴权**:每次调用由调用方(AI agent)在命令行直接传入 `token` 与 `tenant-id`。
CLI 不再持久化任何凭证,删除 login/logout/refresh 全部机制。

## 新契约

```bash
wlt customer page --token <accessToken> --tenant-id 22222
wlt customer page --token xxx --tenant-id 22222 --base-url https://erpapi.w-lian.com
```

| 来源 | 提供内容 |
|------|----------|
| `--profile`(默认 `sit`) | `base_url` + `api_prefix` + `enterprise_type`(连接信息) |
| `--token`(必填,API 命令) | bearer access token → `Authorization: Bearer <token>` |
| `--tenant-id`(必填,API 命令) | 租户 → `tenant-id` 头 |
| `--base-url`(可选) | 覆盖 profile 的 base_url |

## 校验边界

- 走 `EnsureClient()` 的命令(所有业务命令 + `wlt api`,含 `--dry-run`)必须提供 `--token` 与 `--tenant-id`,
  缺失 → 退出码 **4(validation)**。
- 服务端拒绝 token(body `code:401 账号未登录`)目前随现有错误映射表现为退出码 **5(api_error)**,
  message 中含 `code=401`;Agent 据此判断需重新获取并通过 `--token` 传入。
  (退出码 3 authentication 当前为保留值,未使用;统一错误分类为后续可选优化,不在本次范围。)
- `wlt version`/`completion`/`help`/`config show|init|set` 不需要凭证。

## 代码改动

### 删除
- `cmd/auth/` 整个包(login/logout/status)
- `internal/auth/` 整个包(Manager/Login/Logout/Refresh/GetValidToken,含 `auth_test.go`)
- `cmdutil.AuthMgr` 全局变量与 `GetAuthMgr()`

### 改动
- `internal/config/config.go`:`Profile` 删除 `AccessToken`/`RefreshToken`/`ExpiresTime`/`TenantID` 字段;
  `defaultConfig` 去掉 `TenantID:"1"`;`UpdateProfileField` 去掉 `tenant_id` 分支。
- `internal/client/client.go`:`NewClient` 改为接收 `RequestContext{BaseURL,APIPrefix,TenantID,EnterpriseType,Token}`;
  `buildRequest` 直接用 ctx.Token,不再调 `GetValidToken`;删 `authMgr` 字段与 `auth`/`config` 导入依赖。
- `internal/cmdutil/client.go`:新增 `AuthFlags`(token/tenant-id/base-url)全局与 setter;
  `EnsureClient` 从 profile + flag 覆盖组装 `RequestContext`,校验必填,构建 client。
- `cmd/root.go`:新增 3 个 persistent flag,`PersistentPreRunE` 解析后存入 `cmdutil`;删 `auth` import 与 `auth.Register`。
- `cmd/config/`:`config_init` 去掉租户表单项;`config_set` 去掉 `tenant_id`;`config_show` 去掉 `has_token`/`expires_time`/`tenant_id`。
- 测试:`internal/client/client_test.go`、`internal/config/config_test.go`、`internal/output/output_test.go` 更新引用。

## 文档改动
- `skills/SKILL.md`:删认证预检步骤,改为「每次必传 `--token`/`--tenant-id`」;更新 6 个工作流示例与退出码表。
- `skills/references/auth-config.md`:重写为 flag 鉴权。
- `skills/references/api-conventions.md`:注明 token 经 `--token` 传入。
- `CLAUDE.md`:更新认证段、项目结构(去 `cmd/auth`)、架构约定。

## 兼容性
破坏性变更。旧 config 中的 token 字段被 YAML 忽略(不报错)。版本号建议升 minor/major。
