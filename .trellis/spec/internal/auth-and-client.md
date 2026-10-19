# 无状态鉴权契约与 HTTP 客户端

> 核心：**CLI 不保存登录态、无 `auth` 命令**。每次调用由调用方（AI Agent / 脚本）
> 通过 `--token` / `--tenant-id` 传入凭证，`--profile` 提供连接信息。
> 设计来源：2026-06-26 无状态鉴权重构（归档任务 `10-16-stateless-auth`）。

---

## 1. 三段链路：flag → 暂存 → 校验 → 成头

```
cmd/root.go PersistentPreRunE（只暂存，不校验）
   └─► cmdutil.SetAuthFlags(token, tenantID, baseURL)
          └─► 全局 cmdutil.AuthFlags
                 └─► 命令 RunE 首行 cmdutil.EnsureClient()（懒加载 + 校验 + 组装）
                        └─► client.NewClient(client.RequestContext{...})
                               └─► buildHeaders() 把 token/tenant-id 变成请求头
```

### 1.1 `--profile` 提供什么

| 来源 | 提供内容 | 锚点 |
|------|----------|------|
| `--profile <name>`（persistent，默认 `sit`） | profile 名 → `base_url` / `api_prefix` / `enterprise_type` | `cmd/root.go:88`、`cmd/root.go:45-61` |
| `Profile` 结构体 | `Profile{BaseURL, APIPrefix, EnterpriseType}` —— **只含连接信息，无凭证字段** | `internal/config/config.go:23-27` |
| 默认值 | `sit` = `https://erpsit.api.w-lian.com` + `/admin-api`；`prod` = `https://erpapi.w-lian.com` + `/weiliantong/admin-api` | `internal/config/config.go:181-195` |
| 激活解析 | `ActiveProfile()`：`profileOverride`（flag）非空优先，否则 `cfg.Active`；找不到 → error | `internal/config/config.go:100-113` |
| 配置文件 | `~/.wlt/config.yaml`：目录权限 `0700`、文件权限 `0600` | `internal/config/config.go:86-95` |

详见 [config.md](./config.md)。

### 1.2 必填 flag、校验路径、缺失行为

- `--token` / `--tenant-id` / `--base-url` 定义于 `cmd/root.go:91-93`。
- `PersistentPreRunE` **只暂存不校验**：`cmd/root.go:45-61` → `cmdutil.SetAuthFlags`
  （`internal/cmdutil/client.go:30-34`），暂存于全局 `cmdutil.AuthFlags`（`internal/cmdutil/client.go:17-22`）。
- **校验点在懒加载的 `EnsureClient()`**（`internal/cmdutil/client.go:39-66`），顺序：

  1. `APIClient != nil` 直接返回（进程内单例，`:40-42`）；
  2. `CfgMgr == nil` → 普通 error（`:43-45`）；
  3. `ActiveProfile()` 失败 → `NewExitError(2, "当前环境配置不存在: ...", "运行 wlt config init 初始化")`（`:46-49`）；
  4. `AuthFlags.Token == "" || AuthFlags.TenantID == ""` → **`output.NewExitError(4, "缺少必填鉴权参数",
     "请通过 --token 与 --tenant-id 传入鉴权信息(对应 Authorization 与 tenant-id 请求头)")`**（`:50-53`）
     ← **这是「缺失 → 退出码 4」的唯一产生点**；
  5. `--base-url` flag 非空则覆盖 `profile.BaseURL`（`:54-57`）；
  6. `client.NewClient(client.RequestContext{...})` 组装并缓存（`:58-64`）。

> 懒加载的含义：不调用 `EnsureClient()` 的命令（`version`/`config`/`completion`/`help`）
> 不需要 `--token`/`--tenant-id`。清单见 [../cli/directory-structure.md](../cli/directory-structure.md) §6。

### 1.3 鉴权如何变成请求头

`internal/client/client.go:192-201` `buildHeaders`：

```go
h.Set("Authorization", "Bearer "+token)   // CLI 自动加 Bearer 前缀
h.Set("tenant-id", c.ctx.TenantID)
h.Set("Content-Type", "application/json")
if c.ctx.EnterpriseType != "" {
	h.Set("enterprise-type", c.ctx.EnterpriseType)  // 可选，空则不发
}
```

- `RequestContext` 定义：`internal/client/client.go:29-35`（`BaseURL`/`APIPrefix`/`TenantID`/`EnterpriseType`/`Token`）。
- `NewClient` 设 60s HTTP 超时：`internal/client/client.go:44-51`。
- URL 拼接 `BaseURL + APIPrefix + path`，query 由 `url.Values` 编码：`internal/client/client.go:163-171`。
- dry-run 会用占位 token（`"dry-run-token"`）构建头，避免把真实 token 回显到 stdout：
  `internal/client/client.go:118-120`。

---

## 2. 服务端 401 如何浮现（重点）

- 响应体 `code != 0` → `client.do` 调 `apierr.FromCommonResult(resp.StatusCode, raw)` 并返回该 error
  （`internal/client/client.go:150-155`）；`code == 0` 才取 `data`（`:157-158`）。
- `apierr.APIError.Error()` = `fmt.Sprintf("API error (http=%d, code=%d): %s", HTTPStatus, Code, Msg)`
  （`internal/apierr/apierr.go:17-19`）。
- 各命令统一包装：`output.NewExitError(5, fmt.Sprintf("查询%s失败: %s", label, err), "")`
  （例 `internal/cmdutil/crud.go:86`）。

### ⇒ 结论

**token 失效（后端 body 为 `{"code":401,"msg":"账号未登录"}`，HTTP 200）表现为退出码 5**，
stderr JSON 的 `message` 内含 `code=401` 与后端 `msg`，形如：

```
API error (http=200, code=401): 账号未登录
```

- `ExitAuthentication = 3` 是**保留常量，代码中从未产出**（`internal/apierr/apierr.go:46`；
  唯一使用处是测试 `internal/output/output_test.go:51`）。
- 测试佐证：`internal/apierr/apierr_test.go:15-30` 用 HTTP 200 + body `code:401` 断言 `err.Code == 401`、`err.HTTPStatus == 200`。
- Agent 判断「需重新取 token」的依据是 **stderr message 含 `code=401`**，不是退出码 3。

---

## 3. 多租户头

| 头 | 来源 | 必填？ |
|----|------|--------|
| `Authorization: Bearer <token>` | `--token` | 是（走 API 的命令） |
| `tenant-id` | `--tenant-id` | 是（走 API 的命令） |
| `enterprise-type` | profile 的 `enterprise_type` | 否（空则不发） |
| `Content-Type: application/json` | 固定 | — |

---

## 4. 重试语义

`client.do`（`internal/client/client.go:123-158`）：

- **只在网络错误上重试一次**：`for attempt := 0; attempt < 2; attempt++`（`:126`），
  第二次前等待 500ms（`:131`）；ctx 取消即返回（`:129-130`）。
- **HTTP 4xx/5xx 或 body `code != 0` 不重试**（`:150-155` 直接返回 error）。
- 请求体注意：`resp.Body.Close()` **内联在循环里**（`:143-146`），
  注释明确禁止 `defer`-in-loop（会泄漏到函数返回）。

---

## 5. `CommonResult<T>` 解析

- 统一响应：`{ "code": 0, "data": T, "msg": "" }`。
- `client.do` 用 `gjson` 读 `code`（`:151`）与 `data`（`:157-158`），
  返回 `&Response{Data: json.RawMessage(dataBytes)}` —— **命令层拿到的是原始 JSON 字节**，
  由 `cmdutil.OutputJSON` / `ParsePagedJSON` 决定如何信封化。
- `code != 0` 时经 `apierr.FromCommonResult` 转成 `*APIError`。
  `FromCommonResult` 的细节：
  - 非法 JSON body → `code = -1`（`internal/apierr/apierr.go:34-36`，测试 `apierr_test.go:32-41`）；
  - `code == 0` → 返回 `nil`。

---

## 6. 客户端 API 一览

| 方法 | 行 | 用途 |
|------|----|------|
| `Get(ctx, path, params)` | `internal/client/client.go` | GET + query |
| `Post(ctx, path, body)` | 同上 | POST + JSON body |
| `Put(ctx, path, body)` | `:64` | PUT + JSON body（`@RequestBody` 型） |
| `PutQuery(ctx, path, params)` | `:75` | PUT + query，无 body（`@RequestParam` 型，如 `update-status`） |
| `DeleteWithBody(ctx, path, body)` | `:81` | DELETE + JSON body |
| `BuildDryRun(...)` | `:118` 附近 | `wlt api --dry-run` 用 |

PUT 形态判据见 [../cli/flags-and-params.md](../cli/flags-and-params.md) §8。

---

## 7. 相关

- 配置与 profile：[config.md](./config.md)
- 退出码完整表与 `output`/`apierr` 实现：[output-and-errors.md](./output-and-errors.md)
- 历史归档（设计原文）：`.trellis/tasks/archive/2026-10/10-16-stateless-auth/`