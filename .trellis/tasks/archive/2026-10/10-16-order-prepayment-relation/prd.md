# 订单预付关联命令（历史归档）

> **状态**：已完成（交付日期 2026-08-21，设计 `347a0d1`、计划 `2342e97`，实现 `af07f87`/`bb1a8df`/`2d8814b`）。
> 本文档为历史归档，不是待办。原文见同目录 `design.md`（设计）与 `implement.md`（实现计划）。

## Goal

在 order 域下新增 `wlt order prepayment-relation` 子命令组，覆盖后端 `/erp/order-prepayment-relation`
的 5 个接口（把订单与预付款来源：付款单 / 供应商期初，建立关联并维护关联金额）。
命令命名、flag 风格、输出协议与全仓既有约定一致。

## Requirements（原文摘要）

5 个接口与对应子命令：

| 子命令 | 端点 | 方法 |
|---|---|---|
| list | `/erp/order-prepayment-relation/page` | GET |
| create | `/erp/order-prepayment-relation/create` | POST |
| update-amount | `/erp/order-prepayment-relation/update-amount` | PUT |
| available-payments | `/erp/order-prepayment-relation/available-payments` | GET |
| available-initials | `/erp/order-prepayment-relation/available-initials` | GET |

- `relationType` 枚举：`PAYMENT`（付款单）/ `SUPPLIER`（供应商期初）。
- create / update-amount 用**离散 flag** 而非 `--data` JSON（字段少且固定：4 个 / 2 个），请求体以 `map[string]any` 构造。
- available-payments / available-initials 共享一个工厂函数（flag 与响应结构完全一致，仅端点与标签不同）。
- 新文件 `cmd/order/order_prepayment.go`（package order，经 `init()` 注册到既有 `orderCmd`）。

## Out of Scope（原文非目标）

- 不封装组合工作流（查询可用余额 → 创建关联 → 改金额）为一条命令，保持单接口单命令。
- 不改动 `DocumentCmds` / `AddCRUDToParent` 等通用工厂（本模块 5 个接口形状各异，不适合套用标准 CRUD 工厂）。

## Acceptance Criteria

- [x] `cmd/order/order_prepayment.go` 与 `cmd/order/order_prepayment_test.go` 存在
- [x] 5 个子命令齐全（list/create/update-amount/available-payments/available-initials）
- [x] `page` 接口响应为裸数组 → 用 `cmdutil.OutputJSON(resp.Data)` 而非 `ParsePagedJSON`
- [x] `go build ./...` + `go test ./...` 通过，`wlt order prepayment-relation --help` 命令树正确
- [x] `skills/references/order.md`、`skills/SKILL.md`、`CLAUDE.md` 同步更新

## Notes

- 后续在该命令组上叠加了 0923 基线增量 `--biz-type`（PURCHASE/LOGISTICS），见归档任务 `10-16-backend-sync-0923`。
- 「裸数组响应不可用 `ParsePagedJSON`」这一判据已提炼进 `.trellis/spec/cli/output-protocol.md`。