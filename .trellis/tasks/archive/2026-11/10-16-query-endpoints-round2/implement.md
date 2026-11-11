# 执行计划 — Round 2 查询端点 + 财务写命令 + `-o` 方案

> 顺序：小命令热身 → screen → system 模块 → 财务写命令 → `-o` 方案 → 全量验收。
> 每批次的固定收尾：`go build ./... && go test ./...` + `./wlt <domain> --help` 命令树核对。
> 参数取法：先读 design.md §1 表中该模块的 Controller 锚点（方法签名 + VO 字段），再写命令；
> **不参考任何旧 erp.md 残留描述**。

## Batch 0 — 基线

- [x] `go build ./... && go test ./...` 通过；`git status --short` 干净（除 .trellis）

## Batch 1 — 小命令热身（5 个：position 2 + waybill 2 + invoice-item 1）

- [x] 新建 `cmd/position/position.go`（`Register`，模式照 `cmd/screen/screen.go` 包结构）
      + `latest-follow` / `history-follow`（参数：`ErpPositionController.java:30-77` 方法签名）
- [x] `cmd/waybill/`：`events`（`WaybillSourceController.java:98`）、
      `order-info`（`:104`）；模板 `waybill_get.go`
- [x] `cmd/invoice/invoice_item.go`：`item list`（`InvoiceItemController.java:31` + VO 字段；
      响应若是裸数组→`OutputJSON`，若是 `{list,total}`→`ParsePagedJSON`，以后端返回结构为准）
- [x] `cmd/root.go` init 注册 `position.Register(rootCmd)`
- [x] build/test + `--help` 树核对

## Batch 2 — screen 6 个

- [x] `cmd/screen/screen.go` 追加 6 命令（design §1 表锚点；模板 `:32-50`；
      核对各方法签名是否带 query 参数，带则加离散 flag）

## Batch 3 — system 模块（10 个）

- [x] `cmd/system/sys_account.go`：`account` 组 + `list`（分页）+ 8 个详情命令
      （锚点见 design §1；详情命令多为单 id/单 code 参数）
- [x] `cmd/system/sys_weight_account.go`：`weight-account get`
      （`ErpSystemWeightAccountController.java:24`）
- [x] build/test + 树核对

## Batch 4 — 财务三模块写命令 + 漂移清理

- [x] `finance_invoice.go`：加 `create/update/update-status/delete`
      （模板 `finance_receipt_payment.go:174-278`；锚点 `FinanceInvoiceController.java:33-52`）；
      **删 `page-count`**（`:108-131` + `financeInvoicePageCountFilters` 若独占）
- [x] `finance_payment_apply.go`：加 `create/update/update-status/delete/export`；
      **删 `page-count`**（锚点 `ErpFinancePaymentApplyController.java:37-90`）
- [x] `finance_prepayment_apply.go`：`list` 改名 `page`（design §4）；
      加 `create/update/update-status/delete/export`
      （锚点 `ErpFinancePrepaymentApplyController.java:37-90`）
- [x] 三模块 export 子命令直接用 Batch 5 的 `RunExport`（若 Batch 5 未完成，先留 TODO 本批内做完）
      ——实施顺序调整为：**先做 Batch 5 的 helper，再回来做本批 export**（见下）
- [x] build/test + 树核对

## Batch 5 — `-o` 方案（先 helper，供 Batch 4 的 export 用）

> 实际执行顺序：5a（flag+helper+单测）→ Batch 4 → 5b（改造存量 23 处）

- [x] 5a：`root.go` 加 `--output/-o` PersistentFlags；
      新建 `internal/cmdutil/export.go` 的 `RunExport(ctx, path, params)`（design §3）；
      `internal/cmdutil/export_test.go` 单测（临时目录写入 + 信封 JSON 断言 + flag 读取）
- [x] 5b：改造存量 23 处（清单见 `research/03-export-inventory.md`）：
      `Get()+fmt.Println` → `cmdutil.RunExport`；保留各自参数装配与 `--headers` 收集
- [x] build/test 全量

## Batch 6 — 全量验收 + 实测（需 token）+ 收尾

- [ ] 逐条走 prd Acceptance Criteria（更新后的范围版），每条给命令与输出
- [x] 命令树总核对：`wlt --help` 各域无缺失/错名
- [ ] **token 实测**（向用户索取有效 token）：
      - Round 1 遗留：`finance invoice list`、`report stock warehouse`、`profit-event list`、
        filter camelCase 各 1 条
      - 本轮新命令每域抽 1 条（position / screen / system account / weight-account /
        waybill events / invoice item / 财务写命令走 create 后 delete 或 --dry-run 类验证）
      - 已删 3 个 page-count 端点确认真机 404（佐证删除正确）
- [x] spec 更新完成：directory-structure（26 域/新文件行）、output-protocol §4 重写为 export 正式协议、flags 分页主名 11 处/finance(7)、backend-sync-playbook 增端点漂移规则
- [ ] 更新 `CLAUDE.md` 项目结构（cmd/position/、cmd/system/ sys_account 等）与 `skills/SKILL.md`
      对应表格（仅增改涉及行）
- [ ] 汇报 `git status` + 改动摘要，用户确认后提交（提交信息前缀 `feat:`，域分批或合一由用户定）

## 验证命令汇总

```bash
go build ./... && go test ./...
./wlt position --help && ./wlt screen --help && ./wlt system --help
./wlt waybill --help && ./wlt invoice --help && ./wlt finance --help
grep -rn "导出成功，返回数据" cmd/ | wc -l    # 期望 0（23 处全改造）
```

## 回滚点

| 批次 | 回滚方式 |
|---|---|
| Batch 1-4（纯新增/域内改动） | `git checkout -- cmd/<domain>/` 或按文件 revert |
| Batch 5a（-o flag + helper） | helper 独立文件，revert 无连带 |
| Batch 5b（23 处改造） | 集中一个提交，revert 即全撤 |
| Batch 6 文档 | 独立提交 |
