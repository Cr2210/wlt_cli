# 无状态鉴权重构（历史归档）

> **状态**：已完成（交付日期 2026-06-26，commit `c4b0c0e`）。本文档为历史归档，不是待办。
> 原文见同目录 `design.md`（设计）与 `implement.md`（实施清单）。

## Goal

把 `wlt` 的鉴权从「交互式 `wlt auth login` → token 存入 `~/.wlt/config.yaml` → 客户端自动刷新」
改为**无状态鉴权**：每次调用由调用方（AI Agent）在命令行直接传入 `--token` 与 `--tenant-id`。
CLI 不再持久化任何凭证，删除 login/logout/refresh 全部机制。

动机：作为 AI Agent 原生 CLI，交互式登录对 Agent 不友好（每个会话都要先跑 `auth status`/`auth login`），
且 token 持久化在磁盘上既不安全也无必要（Agent 自己持有新鲜 token）。

## Requirements（原文摘要）

- `--profile`（默认 `sit`）提供连接信息：`base_url` + `api_prefix` + `enterprise_type`。
- `--token`（API 命令必填）→ `Authorization: Bearer <token>` 请求头。
- `--tenant-id`（API 命令必填）→ `tenant-id` 请求头。
- `--base-url`（可选）覆盖 profile 的 `base_url`。
- 走 `EnsureClient()` 的命令（所有业务命令 + `wlt api`，含 `--dry-run`）缺失必填项 → 退出码 4（validation）。
- 服务端拒绝 token（body `code:401 账号未登录`）→ 退出码 5（api_error），message 含 `code=401`。
- 删除 `cmd/auth/`、`internal/auth/`、`cmdutil.AuthMgr`/`GetAuthMgr()`，以及 `Profile` 中的凭证字段。

## Out of Scope（原文）

- 统一错误分类（即把服务端 401 单独映射为退出码 3 authentication）——原文标注为「后续可选优化，不在本次范围」。
  代码现状：退出码 3 至今仍是保留值，无产出路径。

## Acceptance Criteria

- [x] `cmd/auth/`、`internal/auth/` 两包已删除（当前仓库路径不存在）
- [x] `internal/cmdutil/client.go` 含 `AuthFlags` 且 `EnsureClient()` 从 profile + flag 组装 `RequestContext`
- [x] `internal/client/client.go` 的 `NewClient` 接收 `RequestContext`，`buildHeaders` 自动加 `Bearer ` 前缀
- [x] `internal/config/config.go` 的 `Profile` 无凭证字段（仅 `BaseURL`/`APIPrefix`/`EnterpriseType`）
- [x] `cmd/root.go` 注册 3 个 persistent flag（`--token`/`--tenant-id`/`--base-url`），`PersistentPreRunE` 只暂存不校验
- [x] `skills/`、`CLAUDE.md` 已按新契约更新

## Notes

- 破坏性变更：旧 config 中的 token 字段被 YAML 忽略（不报错）。
- 现状核对与已知漂移见 `.trellis/spec/internal/auth-and-client.md` 与 `.trellis/spec/cli/output-protocol.md`。