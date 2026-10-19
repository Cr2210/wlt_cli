# 库存命令优化（历史归档）

> **状态**：已完成（交付日期 2026-08-07，设计 `eb2d679`、计划 `95fb7dd`，实现 `e4fedd4`..`9e586d9`）。
> 本文档为历史归档，不是待办。原文见同目录 `design.md`（设计）与 `implement.md`（实现计划）。

## Goal

把全仓单据命令（采购入库/销售出库/其他入库/其他出库/调拨/盘点/退货）归一到**一套通用工厂 `DocumentCmds`**，
消除「工厂 / legacy / 手写」三套并存风格；统一库存域命名、筛选 flag 与输出行为；削减约 1000 行重复样板，
删除 legacy 工厂 `crud_legacy.go`；修正失真文档。

## 背景（原文摘要）

库存相关命令分布在 `cmd/stock/`、`cmd/purchase/`、`cmd/sale/`，实现风格三套并存，衍生真实问题：

1. 命名割裂：工厂生成的分页子命令叫 `list`，stock 手写的叫 `page`（`wlt purchase in list` 但 `wlt stock in page`）。
2. `--headers` 自定义导出缺失：只有 stock 手写命令支持，采购入库/销售出库（工厂生成）没有。
3. 分页解析不健壮：stock 部分命令用手写 `int64 total` 解析（遇字符串型 `total` 报错）。
4. flag 列表重复：stock in/out/move/check 的 page 与 page-count 各抄一份 20+ flag。
5. 文档失真：`CLAUDE.md` 称 stock 用 legacy 工厂，实际早全手写。

## Requirements（原文摘要）

- 引入通用单据工厂 `DocumentConfig`（`Name`/`APIPath`/`Label`/`Filters`/`TimeKey`/`Headers`/`PageAlias`）与 `DocumentCmds`。
- stock in/out/move/check 与采购/销售退货迁移到 `DocumentCmds`；`stock warehouse` 迁移到 `AddCRUDToParent` + `CrudSimpleListCmd`。
- stock 单据的 `page` 改名 `list`，但保留 `page` 作为 cobra `Aliases` 兼容别名。
- `--headers` 导出补齐（`purchase in`/`sale out`）。
- stock query/record 不套工厂（非 CRUD 端点有真实业务逻辑），改在 stock 包内用共享 `FlagSpec` + 本地辅助函数消除 flag 重复。
- 手写 `int64 total` 分页解析全部替换为 `cmdutil.ParsePagedJSON`；删除 `internal/cmdutil/crud_legacy.go`（224 行，全仓再无调用方）。
- 修正 `CLAUDE.md` 关于 stock 用 legacy 工厂的失真描述。

## Out of Scope（原文非目标）

- 不改变后端 API 契约（仅调整 CLI 层 flag 与传参映射）。
- 不重构 `cmd/report/`、`cmd/stats/` 下的库存报表/统计命令。
- 不改变鉴权、输出协议（stdout JSON / stderr 错误 / 退出码 0-6）。

## Acceptance Criteria

- [x] `internal/cmdutil/crud.go` 含 `DocumentConfig` / `DocumentCmds`（含 `TimeKey`/`Headers`/`PageAlias`）
- [x] `internal/cmdutil/crud_legacy.go` 已删除（全仓无 `NewCRUDSubCmd` 命中）
- [x] stock in/out/move/check 四个单据的 `list` 子命令含 `page` 别名
- [x] `make test` / `go build ./...` 通过（迁移类任务靠编译 + `make test` + `--help` 结构核对验证等价性）
- [x] `CLAUDE.md` 的 CRUD 模式描述已修正为 `DocumentCmds` + `AddCRUDToParent`

## Notes

- 迁移类任务无法纯单测（`RunE` 依赖全局 `APIClient`，见 `internal/cmdutil/client.go`），
  工厂层逻辑（时间转换、命令结构）由 `internal/cmdutil/crud_test.go` 守护。
- 现状与已知不一致见 `.trellis/spec/cli/command-factories.md`。