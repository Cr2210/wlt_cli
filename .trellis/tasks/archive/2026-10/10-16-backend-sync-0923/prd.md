# 后端 0923 基线同步（历史归档）

> **状态**：已完成（交付日期 2026-09-28，设计 `0f05b25`、计划 `5eee8a5`，
> 实现 `40de0c9`/`4888596`/`13e7f7e`/`67aa36a`/`cbd785b`/`398d6ee`）。
> 本文档为历史归档，不是待办。原文见同目录 `design.md`（设计）与 `implement.md`（实现计划）。

## Goal

为 wlt CLI 补全后端 2026-09-23 基线（后端 commit `2e52721b`，分支 dev-1.7/ge-dev-1.8）的新功能：
物流结算（13 端点）、销售核算（11 端点）、运费申请（11 端点）三个新命令组，
并处理结算忽略剩余、预付关联 `bizType`、收付款 `isPrepaid` 三处增量。

CLI 侧上次同步基线：2026-08-21（commit `6318e49`）。差异分析结论：**3 个全新模块 + 3 处已有域增量**，
其余为后端内部口径调整或响应字段透传（**CLI 是薄透传层，无需改动**）。

## Requirements（原文摘要）

A / B / 排除 三类分类与处置：

| 类别 | 模块 | 后端前缀 | 端点数 | 处置 |
|---|---|---|---|---|
| A1 | 物流结算 | `/erp/logistics-settlement` | 13 | `wlt settlement logistics <sub>`（`cmd/settlement/logistics.go`） |
| A2 | 销售核算 | `/erp/sale-accounting` | 11 | `wlt sale accounting <sub>`（`cmd/sale/accounting.go`） |
| A3 | 运费申请 | `/erp/finance-freight-apply` | 11 | `wlt finance freight-apply <sub>`（`cmd/finance/freight_apply.go`） |
| B1 | 结算单忽略剩余 | `/erp/settlement/ignore-remainder` + `invoiceFlag` | 1 | `wlt settlement main ignore-remainder` |
| B2 | 预付关联双业务类型 | `bizType`（PURCHASE/LOGISTICS） | flag 增量 | prepayment-relation list/create 加 `--biz-type` |
| B3 | 收付款预付标识 | `isPrepaid` 查询筛选 | flag 增量 | 收付款 list 加 `--is-prepaid` |

排除项（后端内部逻辑或响应透传，CLI 无改动）：发票核销/开票申请卡口过滤、出入库详情磅单回填、
驾驶舱/资金占用/物流预付统计口径、WebSocket 鉴权、通知公告评论回复（CLI 未接 notice）、扫码接单/费用下发。

## Out of Scope（原文非目标）

- 不封装组合工作流（如「试算→创建核算单」一键命令），保持单接口单命令。
- 不改动 `DocumentCmds` / `AddCRUDToParent` / `CrudXxxCmd` 通用工厂签名；不足之处用域内自定义命令补齐。
- 不补 notice（通知公告）域——不在本次后端差异主线上，如需要另立任务。

## Acceptance Criteria

- [x] `cmd/settlement/logistics.go`（13 子命令）、`cmd/sale/accounting.go`（11）、`cmd/finance/freight_apply.go`（11）存在
- [x] `wlt settlement main ignore-remainder` 已注册（`cmd/settlement/settlement_main.go`）
- [x] prepayment-relation 的 list/create 支持 `--biz-type`；收付款 list 支持 `--is-prepaid`
- [x] `go build ./...` 通过；`--help` 命令树核对完成
- [x] `CLAUDE.md`、`skills/` 文档同步；`query-endpoints backlog` 增量更新

## Notes

- 踩坑记录：后端一个域可能有多个 Controller（如 `/erp/settlement` 与 `/erp/logistics-settlement`），
  统计端点数时须分开 grep，否则会数错（初次分析曾把物流结算误判为 14 端点，纠正为 13）。
- 该次同步遗留的待办（第 2 轮端点 + export `-o` 方案）已转为开放任务 `10-16-query-endpoints-round2`。
- 方法论已提炼进 `.trellis/spec/cli/backend-sync-playbook.md`。