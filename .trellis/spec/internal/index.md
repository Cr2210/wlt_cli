# 共享内部包约定（`.trellis/spec/internal/`）

> 本层对应仓库顶层源码目录 `internal/` —— 不对外暴露的共享包：
> `apierr/`、`build/`、`client/`、`cmdutil/`、`config/`、`output/`。
> 这一层是 cli 层（`cmd/`）的依赖底座：改这里的契约会波及所有命令。
>
> **语言：中文**（与 `CLAUDE.md`、`skills/`、`docs/`、commit message 一致）。
>
> 所有约定以**代码现实**为准。

---

## Pre-Development Checklist

在 `internal/` 下新增或修改代码前，逐条确认：

- [ ] 改动落在**正确的包**：
      - 鉴权/请求组装/多租户头 → `internal/client/` 与 `internal/cmdutil/client.go`
      - 参数收集/时间折叠/JSON 解析 → `internal/cmdutil/helpers.go`、`timeflags.go`、`output.go`
      - CRUD/单据工厂 → `internal/cmdutil/crud.go`
      - 配置/profile → `internal/config/`
      - 输出信封/退出码类型 → `internal/output/`；退出码常量与类型映射 → `internal/apierr/`
- [ ] 改了 `internal/client` 的请求契约后，确认 `internal/cmdutil` 与 `cmd/**` 的调用方仍编译（`go build ./...`）。
- [ ] 新增退出码产出点时，确认 code 对应的 `ExitType` 已存在（`internal/apierr/apierr.go:53-68`），
      且**没有**引入 code 3 或 6 的新产出（这两个是保留值，见 [output-and-errors.md](./output-and-errors.md)）。
- [ ] 若改了 `internal/config` 的 `Profile` 字段，同步检查 `UpdateProfileField` 的 switch 与 `defaultConfig()`。
- [ ] 有 HTTP 层行为改动时补 `httptest` 测试（`internal/client/client_test.go` 是样例）。

---

## Quality Check

- [ ] `go build ./... && go vet ./... && go test ./...` 通过。
- [ ] 工厂/参数逻辑的改动已补表驱动单测（`helpers_test.go`、`timeflags_test.go`、`crud_test.go`）。
- [ ] 没有把凭证写入磁盘：`internal/config` 的 `Profile` 只有连接信息（`BaseURL`/`APIPrefix`/`EnterpriseType`），
      **无** `AccessToken`/`RefreshToken`/`TenantID` 等字段。
- [ ] 没有把 `tenant-id` 当成配置项 —— 它由必填 flag `--tenant-id` 每次传入。
- [ ] 错误信息的语言与仓库一致（中文面向用户的 hint；`apierr.APIError.Error()` 保持英文格式串）。
- [ ] 改了契约后，`cli/` 层对应 spec 已同步。

---

## Guidelines Index

| Guide | 内容 | Status |
|-------|------|--------|
| [auth-and-client.md](./auth-and-client.md) | 无状态鉴权契约、`RequestContext`、请求头、重试、`CommonResult` 解析、`apierr` | 已填充 |
| [config.md](./config.md) | `Profile` 结构、`~/.wlt/config.yaml`、`ActiveProfile()`、`UpdateProfileField`、默认 profile | 已填充 |
| [output-and-errors.md](./output-and-errors.md) | `internal/output` 信封与 `ExitError`、`internal/apierr` 退出码映射来源 | 已填充 |
| [quality-guidelines.md](./quality-guidelines.md) | 内部包测试约定与验证命令 | 已填充 |

相关层：
- `cli/` 层（[../cli/index.md](../cli/index.md)）：命令分包、工厂选型、flag 映射、输出协议（命令层视角）。
- `guides/`（[../guides/index.md](../guides/index.md)）：跨层思维指南。

---

## 层间职责边界（避免重复表述）

| 主题 | 写在 `internal/` | 写在 `cli/` |
|------|------------------|-------------|
| 鉴权 | 契约本身、`RequestContext`、请求头、403/401 如何浮现（**实现级**） | "缺参数 → 退出码 4" 的使用视角（见 cli 的 Pre-Development Checklist） |
| 输出 | `output` 包函数与 `ExitError` 结构（**实现级**） | 退出码 0-6 的**完整表**与 `ParsePagedJSON` vs `OutputJSON` 的**选择判据** |
| 配置 | `Profile` / 配置文件路径与权限 / `ActiveProfile()` | 只引用，不复述 |

> 同一约定不要在两层重复展开——重复表述是新的漂移源（见 `guides/code-reuse-thinking-guide.md`）。