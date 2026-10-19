<!-- 历史归档 · 非待办事项
     原始路径: docs/superpowers/specs/2026-06-26-stateless-auth-design.md（「代码改动」+「文档改动」两段）
     交付日期: 2026-06-26
     对应 commit: c4b0c0e
     文档类型: 实施清单（由设计文档的实施段生成，原文逐字摘录）
     说明: 该交付已完成并从 superpowers 工作流迁移至 Trellis。原文逐字保留（未改写），
           仅按 Trellis 三件套归位。事实以仓库当前代码为准，本文档仅供溯源。 -->

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


---

## 兼容性（原文）

破坏性变更。旧 config 中的 token 字段被 YAML 忽略(不报错)。版本号建议升 minor/major。
