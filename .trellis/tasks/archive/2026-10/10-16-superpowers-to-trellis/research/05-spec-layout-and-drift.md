<!-- 本文件是 `conventions-extraction` 的拆分件之一（按研究产物约定：一主题一文件）。
     初衷：单文件 57KB 超过 context_injection.max_file_bytes(32768) 会被截断。
     全量导航见同目录 `00-index.md`。文中形如「§C-1」「§4.5」的交叉引用按该导航定位。 -->

# A. 建议的 spec 文件布局（`.trellis/spec/cli/` + `.trellis/spec/internal/`）

约束（task prd 与工具链）：`.trellis/scripts/common/packages_context.py:30-40` 把 `spec/` 的**子目录**当 layer 扫描并**排除 `guides`**；`:198` 要求每个 layer 有 `index.md`。语言用中文（prd R7）。

```
.trellis/spec/
├── cli/
│   ├── index.md                      # cli 层索引 + Pre-Development Checklist + Quality Check
│   ├── auth-and-config.md            # 无状态鉴权契约 + profile/配置 + 退出码(4 的产生点)  [§1, §C-1]
│   ├── command-structure.md          # 分包/Register/init/新域步骤/文件命名/partner       [§2]
│   ├── flag-mapping.md               # kebab→camel、分页、时间数组、--headers、--data、离散 flag、Put/PutQuery [§4]
│   ├── output-protocol.md            # stdout/stderr 信封 + 完整退出码 0-6 + ParsePagedJSON vs OutputJSON [§5, §C-1]
│   └── testing-and-verification.md   # 测试清单/命名/风格 + Makefile 验证命令              [§7]
└── internal/
    ├── index.md                      # internal 层索引 + Checklist + Quality Check
    ├── cmdutil-factories.md          # 工厂全清单 + 选型判据 + legacy 已删                 [§3]
    ├── client-and-errors.md          # RequestContext/headers/重试/CommonResult + apierr   [§1.3-1.4, §1.6]
    └── backend-sync.md               # 薄透传原则 + 差异方法 + A/B/排除 + 后端仓库/基线    [§6, §D]
```

**逐文件职责与依据**
- `cli/index.md` / `internal/index.md`：**必需**（`packages_context.py:198`）；承载两段固定段落（prd R2）与索引表，不含规则正文。
- `cli/auth-and-config.md`：§1 规则量大且被 prd R3-1 点名；同时收纳 `internal/config` 的 profile 事实（不另建 internal/config 文件，避免碎片化）。
- `cli/command-structure.md`：§2 全部（prd R3-2）。
- `cli/flag-mapping.md`：§4 全部（最大一类，单一职责）；prd R3-4。
- `cli/output-protocol.md`：§5 + 完整退出码表；prd R3-5。与 auth-and-config 分工：后者只讲「缺失→4」与「401→5」，前者出全表。
- `cli/testing-and-verification.md`：§7；prd R3-7。
- `internal/cmdutil-factories.md`：§3；prd R3-3。
- `internal/client-and-errors.md`：`internal/client` + `internal/apierr` + `internal/output` 的实现级细节（与 cli/output-protocol 的「命令层怎么用」互补）。
- `internal/backend-sync.md`：§6；prd R3-6。

> 刻意**不建** `cli/domains.md`（逐域 flag 清单会迅速腐烂，属 `skills/references/*.md` 职责）与 `cli/exit-codes.md`（并入 output-protocol）。`internal/client` 仅 202 行，不必单列文件。

---

# B. 源映射（每个 spec 文件 ← 哪些 superpowers 段落 / 代码锚点）

| spec 文件 | superpowers 源 | 代码锚点（必带） |
|---|---|---|
| `cli/index.md` | 无（结构文件） | `.trellis/scripts/common/packages_context.py:30-40,198` |
| `cli/auth-and-config.md` | `auth-design` 全文（新契约 :18-37、校验边界 :30-36、代码改动 :39-61、兼容性 :63-64） | `cmd/root.go:45-61,88-93`；`internal/cmdutil/client.go:17-34,39-66`；`internal/client/client.go:29-35,192-201`；`internal/config/config.go:23-27,100-113,181-195` |
| `cli/command-structure.md` | `backlog:80`；`order-design:34`；`sync-design:46` | `cmd/root.go:87-123`；25 个 `Register`（样本 `cmd/product/product.go:12`）；`cmd/stock/stock.go:11-13`；`cmd/customer/customer.go:14-20`；`cmd/partner/*.go` |
| `cli/flag-mapping.md` | `backlog:17,74-84`；`stock-design:57-58,85,120-121`；`sync-design:46`；`sync-plan:16-19`；`order-design:51,53` | `internal/cmdutil/helpers.go:28-85`；`timeflags.go:8-21`；`crud.go:65-97,341-361,520-534,665-679,787-800,812-826,839-853,856-967`；`cmd/stock/stock_doc.go:36-37`；`internal/client/client.go:64,75-77,81-83` |
| `cli/output-protocol.md` | `order-design:55-59`（裸数组）；`backlog:59-62`（export 二进制冲突，见 §D） | `internal/output/output.go:13-19,26-33,60-98`；`internal/cmdutil/output.go:12-43`；`cmd/root.go:67-85`；`internal/apierr/apierr.go:43-68` |
| `cli/testing-and-verification.md` | `stock-design:124-129`；`stock-plan:18`；`sync-plan:20` | 14 个 `_test.go`（§7.1）；`Makefile:7-11` |
| `internal/cmdutil-factories.md` | `stock-design:36-122`（DocumentConfig 设计与决策）；`order-design:28`；`sync-design:30`；`backlog:74-79` | `internal/cmdutil/crud.go` 全表（§3.1）；8 个 `DocumentCmds` 调用点；6 个 CRUD 工厂调用点 |
| `internal/client-and-errors.md` | `auth-design:47-55`；`sync-design:18` | `internal/client/client.go:21-201`；`internal/apierr/apierr.go`；`internal/output/output.go` |
| `internal/backend-sync.md` | `sync-design` 全文（背景 :7-18、挂载 :33-42、非目标 :27-31、核实 :106-108、顺序 :114-118）；`sync-plan:7,22`；`backlog:7,83-84` | 无 Go 锚点；外部事实：后端仓库 `D:\Projects\weiliantong\weiliantong-api`、分支 `ge-dev-1.8`、基线 `6318e49`↔`2e52721b`、Swagger `D:\Download\erp.md` |

---

# C. 发现的矛盾 / 文档漂移（spec 必须记录现实，不可照抄旧文档）

| # | 漂移 | 文档说法 | 代码现实 | 锚点 | 处置建议 |
|---|---|---|---|---|---|
| C-1 | **退出码 6 `network` 从未产出** | `skills/SKILL.md:318` 与 `skills/references/api-conventions.md:105` 均列「6 = network 网络错误，检查网络连接」；`SKILL.md:324` 还给出处理流程 | 全仓 `grep -rho "NewExitError([0-9]"` **无 `(6`**；网络错误经 `client.do` 返回普通 error，被命令层包成 `NewExitError(5,...)`。`ExitNetwork` 仅定义 + `ExitType` 映射 | `internal/apierr/apierr.go:49,63`；`cmd/root.go:74-79`；`internal/client/client.go:138-141` | spec 写「0-6 常量已定义；**实际产出仅 1/2/4/5**；3 与 6 为保留值，当前无产出路径」 |
| C-2 | 退出码 3 保留 | 文档已正确标注「保留」（`SKILL.md:315`、`api-conventions.md:108`） | 一致；唯一使用处是测试 `internal/output/output_test.go:51` | `apierr.go:46` | 如实保留 |
| C-3 | **`api-conventions.md` 仍描述已删除的 legacy CRUD** | `skills/references/api-conventions.md:62-71`「### 2. Stock 子命令 CRUD（`newCRUDSubCmd`）」，并列为「三种 CRUD 模式」之二 | `internal/cmdutil/crud_legacy.go` 已删除，全仓 `NewCRUDSubCmd` 无命中；stock 单据已全走 `DocumentCmds` | `api-conventions.md:62-71`；`cmd/stock/stock_in.go:6`；`ls internal/cmdutil/` | spec 只列 `DocumentCmds`/`AddCRUDToParent`/`NewCRUDGroup`/`Crud*Cmd`/手写五类，**不写** legacy |
| C-4 | CLAUDE.md 曾称 stock 用 legacy（历史漂移，**已修**） | `stock-design:19`「文档失真：CLAUDE.md 称 stock 用 legacy 工厂，实际早全手写」；`stock-plan` Task 9 给改法 | 现 `CLAUDE.md:79-80` 已正确写 `DocumentCmds`+`AddCRUDToParent`；`CLAUDE.md:64` 已去掉 `crud_legacy.go` 行 | `CLAUDE.md:64,79-80` | 已收敛；spec 引用 CLAUDE.md 现状即可 |
| C-5 | **`DocumentCmds` vs `AddCRUDToParent` 无硬规则** | `CLAUDE.md:79-80` 表述为「通用单据 / 标准 CRUD…大多数域使用」，隐含清晰二分 | 并存且域间不一致：`product`/`stock_warehouse` 用 AddCRUDToParent，`product_unit`/`product_metrics`/`partner_credit`/`partner_invoice` 用 NewCRUDGroup | `cmd/product/product_main.go:16` vs `product_unit.go:6`；`cmd/stock/stock_warehouse.go:17` | 写「按语义选择」+ 调用点清单，并显式声明无强制规则 |
| C-6 | **「分页子命令统一叫 `list`」不成立** | `stock-design:15,87-89` 提出命名统一 `page`→`list`（作为该次目标/决策） | 仅 stock 单据 + 工厂族为 `list`；`finance`(6)、`produce`(2)、`quality`(1)、`waybill`(1) 仍为 `page` | `grep 'Use:   "page"'`（§4.5 清单） | 写「`list` 为工厂族/单据族主名；`page` 在 finance/produce/quality/waybill 仍是主名；`page` 仅 stock 单据作为 `list` 别名」 |
| C-7 | **`--headers` 只部分命令有** | `stock-design:92-93`「补齐 purchase in / sale out 的 `--headers`」；`sync-plan` Task 4/5 给 export 加 `--headers` | `settlement main export` **无** `--headers`；logistics/销售核算/运费申请 export 有；`DocumentCmds` 的 page-count 刻意不带 | `cmd/settlement/settlement_main.go:331-347`；`cmd/settlement/logistics.go:297,307`；`internal/cmdutil/crud.go:944` | 写清「list 收集、page-count 不收集；export 是否支持因命令而异，以各命令 flag 为准」 |
| C-8 | **时间范围两套 helper 并存** | `stock-design:120` 决策 2「flag 统一为 `--start-time/--end-time`」；`sync-plan:19`「结算日期复用 `--start-date/--end-date`」 | `crud.go` 内 6 组域内 helper（flag 名各异）+ `timeflags.go` 通用 helper；同域 `stock query/record` 用单值 | `internal/cmdutil/crud.go:341-853`；`internal/cmdutil/timeflags.go:8-21`；`cmd/stock/stock_doc.go:36-37` | 写「两条路径：通用 `AddTimeRangeFlags/CollectTimeRangeFlags`（新代码优先）+ 域内专用 helper；stock query/record 为单值反例」 |
| C-9 | **export 输出非 JSON 信封** | `backlog:59-62` 指出「返回二进制 Excel，与现有 JSON stdout 协议冲突」，待 `-o` 方案 | 21 处 `fmt.Println("导出成功，返回数据：", string(resp.Data))`，输出文本行，**不符** `{data:...}` 信封 | `cmd/finance/finance_account.go:365` 等 21 处；`cmd/settlement/settlement_main.go:342` | spec 记「export 类为已知例外（非 JSON 信封），待 §D 的 `-o` 方案统一」 |
| C-10 | `finance receipt-payment page` 与测试耦合（无冲突，作交叉验证） | `sync-plan` Task 1 测试片段引用 `newFinanceReceiptPaymentPageCmd()` | 一致：`cmd/finance/finance_receipt_payment.go:86` 为 `Use: "page"`，测试 `finance_receipt_payment_test.go:6` 调用之 | 同上 | 无冲突 |

---

