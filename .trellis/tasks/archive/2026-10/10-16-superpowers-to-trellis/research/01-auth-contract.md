<!-- 本文件是 `conventions-extraction` 的拆分件之一（按研究产物约定：一主题一文件）。
     初衷：单文件 57KB 超过 context_injection.max_file_bytes(32768) 会被截断。
     全量导航见同目录 `00-index.md`。文中形如「§C-1」「§4.5」的交叉引用按该导航定位。 -->

## 1. 无状态鉴权契约

### 1.1 `--profile` 提供什么

| 来源 | 提供内容 | 锚点 |
|---|---|---|
| `--profile <name>`（persistent，默认 `sit`） | profile 名 → `base_url`/`api_prefix`/`enterprise_type` | `cmd/root.go:88`、`cmd/root.go:46-48` |
| profile 结构体 | `Profile{BaseURL, APIPrefix, EnterpriseType}` — **只含连接信息，无凭证字段** | `internal/config/config.go:23-27` |
| 默认值 | `sit`=`https://erpsit.api.w-lian.com`+`/admin-api`；`prod`=`https://erpapi.w-lian.com`+`/weiliantong/admin-api` | `internal/config/config.go:181-195` |
| 激活解析 | `ActiveProfile()`：`profileOverride`（flag）非空优先，否则 `cfg.Active`；找不到 → error | `internal/config/config.go:100-113` |
| 配置文件 | `~/.wlt/config.yaml`，目录 `0700` / 文件 `0600` | `internal/config/config.go:39-45`、`:86-96` |

### 1.2 必填 flag、校验路径、缺失行为

- `--token`/`--tenant-id`/`--base-url` 定义于 `cmd/root.go:91-93`。
- `PersistentPreRunE` **只暂存不校验**：`cmd/root.go:56-59` → `cmdutil.SetAuthFlags`；暂存于全局 `cmdutil.AuthFlags`（`internal/cmdutil/client.go:17-22`，setter `:30-34`）。
- **校验点在懒加载的 `EnsureClient()`**（`internal/cmdutil/client.go:39-66`）：
  1. `APIClient != nil` 直接返回（进程内单例，`:40-42`）；
  2. `CfgMgr == nil` → 普通 error（`:43-45`）；
  3. `ActiveProfile()` 失败 → `NewExitError(2, "当前环境配置不存在: ...", "运行 wlt config init 初始化")`（`:46-49`）；
  4. `AuthFlags.Token == "" || AuthFlags.TenantID == ""` → **`output.NewExitError(4, "缺少必填鉴权参数", "请通过 --token 与 --tenant-id 传入鉴权信息(对应 Authorization 与 tenant-id 请求头)")`**（`:50-53`）← **这是「缺失 → 退出码 4」的唯一产生点**；
  5. `base-url` flag 非空则覆盖 `profile.BaseURL`（`:54-57`）；
  6. `client.NewClient(client.RequestContext{...})` 组装并缓存（`:58-64`）。
- 退出码 4 最终由 `cmd/root.go:67-85` 的 `Execute()` 写出：类型断言 `*output.ExitError`（`:72-73`），非 ExitError 兜底 `Code: 1, Type: "general"`（`:74-80`），`output.WriteError(os.Stderr, exitErr)`（`:81`），**`os.Exit(exitErr.Code)`（`:82`）**。
- `main.go:7-10` 的 `os.Exit(1)` 只是兜底（`os.Exit` 已在上一步发生，实际不可达）。
- 实测计数（`grep -rho "NewExitError([0-9]" | uniq -c`）：`(1)×1, (2)×6, (3)×1, (4)×102, (5)×319`。

### 1.3 鉴权如何变成请求头

`internal/client/client.go:192-201` `buildHeaders`：

```go
h.Set("Authorization", "Bearer "+token)   // CLI 自动加 Bearer 前缀
h.Set("tenant-id", c.ctx.TenantID)
h.Set("Content-Type", "application/json")
if c.ctx.EnterpriseType != "" { h.Set("enterprise-type", c.ctx.EnterpriseType) }  // 可选，空则不发
```

`RequestContext` 定义见 `internal/client/client.go:29-35`；`NewClient` 60s HTTP 超时 `:44-51`；`buildRequest` 拼接 `BaseURL+APIPrefix+path` 与 query（`:163-171`、`:188`）。

### 1.4 服务端 401 如何浮现

- 响应体 `code != 0` → `client.do` 调 `apierr.FromCommonResult(resp.StatusCode, raw)` 并返回该 error（`internal/client/client.go:150-155`）；`code == 0` 才取 `data`（`:157-158`）。
- `apierr.APIError.Error()` = `fmt.Sprintf("API error (http=%d, code=%d): %s", HTTPStatus, Code, Msg)`（`internal/apierr/apierr.go:17-19`）。
- 各命令统一包装：`output.NewExitError(5, fmt.Sprintf("查询%s失败: %s", label, err), "")`（例：`internal/cmdutil/crud.go:86`、`:112`、`:139`）。
- ⇒ **token 失效（body `code:401 账号未登录`）表现为退出码 5**，stderr JSON 的 `message` 内含 `code=401` 与后端 `msg`。与 `auth-design:32-36` 一致；`ExitAuthentication = 3` 是保留常量，**代码中从未产出**（`internal/apierr/apierr.go:46`，唯一使用处是测试 `internal/output/output_test.go:51`）。

### 1.5 无需凭证的命令

判定规则：**只要该命令不调用 `EnsureClient()` 就不需要 token/tenant-id**（`grep -rl "EnsureClient()" --include=*.go .` 命中 60 个文件，全部位于 `cmd/<domain>/` 与 `cmd/api.go`）。具体：

| 命令 | 为何不需要 | 锚点 |
|---|---|---|
| `wlt version` | 覆盖 `PersistentPreRun` 为空函数，绕过配置加载与全局 flag | `cmd/version.go:12-16` |
| `wlt config *` | `cmd/config/` 下三个文件均无 `EnsureClient` 调用 | `cmd/config/config.go:11-13`（仅 `parent.AddCommand(configCmd)`） |
| `wlt completion` | 手写命令，不调 API | `cmd/completion.go:11+` |
| `wlt help` / `--help` | `rootCmd.SetHelpCommand(...)`，`RunE` 缺省走 `cmd.Help()` | `cmd/help.go:11+`、`cmd/root.go:62-64` |
| `wlt api --dry-run` | **仍需要**（先 `EnsureClient()` 才能取 client 调 `BuildDryRun`，`cmd/api.go:43-45`、`:57-66`） | `cmd/api.go` |

### 1.6 退出码表 0–6 及定义位置

| code | type | 定义处 | 产出处 |
|---|---|---|---|
| 0 | —（无 type） | 成功路径直接 return nil | `cmd/root.go:84` |
| 1 | `general` | `internal/apierr/apierr.go:44` | 仅 `Execute()` 兜底 + `ExitType` default 分支（`cmd/root.go:75-79`、`apierr.go:65-66`） |
| 2 | `config` | `apierr.go:45` | 配置加载失败 `cmd/root.go:50`；profile 不存在 `cmdutil/client.go:48`（共 6 处） |
| 3 | `authentication` | `apierr.go:46` | **UNVERIFIED：无任何产出点**（保留常量） |
| 4 | `validation` | `apierr.go:47` | 缺 `--token`/`--tenant-id`（`cmdutil/client.go:51`）；`--data`/`--params` JSON 解析失败（`crud.go:135,162,232`；`cmd/api.go:60,65`）；不支持的 HTTP 方法（`cmd/api.go:36`）；ParsePagedJSON 失败见下（共 102 处） |
| 5 | `api_error` | `apierr.go:48` | 所有 API 调用失败（共 319 处，例 `crud.go:86`）；**含服务端 401** |
| 6 | `network` | `apierr.go:49` | **UNVERIFIED：无任何产出点**（网络错误实际落到退出码 5，见 §C-1） |

`ExitType(code)` 映射表：`internal/apierr/apierr.go:53-68`（default → `"general"`；测试覆盖 `apierr_test.go:52-70`，含 `99 → "general"`）。
`ExitError` 结构：`internal/output/output.go:13-19`；构造函数 `:26-33`（`Type` 自动由 `apierr.ExitType(code)` 派生）与 `:36-44`（带 `Detail`）。

**重试语义**：`client.do` 只在**网络错误**上重试一次（`for attempt := 0; attempt < 2; attempt++`，`:126-141`，间隔 500ms `:131`；ctx 取消即返回 `:129-130`）。HTTP 4xx/5xx 或 `code != 0` **不重试**（`:150-155`）。请求体注意：`resp.Body.Close()` 内联在循环里（`:143-146`，注释明确禁止 defer-in-loop）。

---

