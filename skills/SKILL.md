---
name: wlt
description: 管理维链通 ERP 系统全部能力（库存/产品/客户/供应商/合同/销售/采购/财务/订单/生产/运单/质检/称重/统计/报表/结算/发票/系统管理等）。当用户需要查询库存、管理产品、操作客户供应商、管理合同、处理销售采购单据、管理财务收支核销、管理订单生产单、追踪运单物流、查看质检称重数据、查看统计分析报表、管理结算发票、操作系统用户角色权限字典时使用。
cli_version: ">=0.1.0"
---

# 维链通 ERP 全模块 Skill

通过 `wlt` 命令管理维链通 ERP 系统全部能力。

> ⚠️ **命令可用性取决于环境配置**。本文档基于 wlt v0.1.7 实测编写（已修复 customer/supplier/contract/quality-weight/stock 等模块的类型参数与解析 bug）。**鉴权为无状态模式**：CLI 不再保存登录态，每条业务命令都必须通过 `--token <accessToken>` 与 `--tenant-id <租户ID>` 传入鉴权信息（从用户/会话获取），缺失会以退出码 4 报错。使用 `--profile sit` 连接 SIT 环境，`--profile prod` 连接生产环境（profile 只提供 `base_url`/`api_prefix`）。验证命令结构用 `wlt <cmd> --help`；验证请求用 `wlt api GET <path> --dry-run`（注意：仅 `wlt api` 支持 `--dry-run`，业务子命令不支持）。各端点实测可用性见文末「实测可用性速查」。

## 严格禁止 (NEVER DO)

- 不要使用 wlt 命令以外的方式操作（禁止 curl、直接 HTTP 调用）
- 不要编造 ID、单号等标识符，必须从命令返回中提取
- 不要猜测字段名/参数值，操作前必须先查询确认
- 禁止编造命令路径、子命令或 flag；不确定时必须先运行对应层级的 `wlt --help` 查证
- 不要在未携带 `--token` 与 `--tenant-id` 的情况下执行业务命令（会直接报退出码 4）
- 不要对生产环境 (`--profile prod`) 执行删除操作而未获用户明确确认

## 严格要求 (MUST DO)

- 执行 `wlt` 命令前必须用当前 skill 资料确认命令；缺失或不确定时必须先用 `--help` 查证
- 所有命令默认输出 JSON 格式到 stdout，错误输出到 stderr
- 危险操作（删除、状态变更）必须先向用户确认后再执行
- 删除操作接受多个 ID（逗号分隔），确认时必须展示影响范围
- 状态更新需遵循业务流转规则（如：未审核→已审核→已拒绝）
- 分页查询默认 `--page-no=1`, `--page-size=20`（flag 用连字符短横线，不是下划线；后端字段为驼峰 `pageNo`/`pageSize`），大数据量需分页遍历
- **每条业务命令必须携带 `--token <accessToken>` 与 `--tenant-id <租户ID>`**（无状态鉴权，从用户/会话获取；token 即后端登录返回的 accessToken，CLI 自动加 `Bearer ` 前缀）
- 使用 `--profile sit|prod` 切换环境（提供 base_url/api_prefix），`--quiet` 静默模式，`--base-url` 可选覆盖后端地址

## 模块总览

> 若用户意图涉及多步操作或跨模块流程，**先匹配下方「常见工作流」**；仅当明确是单模块单步操作时，按本表路由。

| 模块 | 用途 | 参考文件 |
|------|------|----------|
| `config` | 配置管理（profile / base_url / api_prefix） | [auth-config.md](./references/auth-config.md) |
| `stock` | 库存管理：仓库 / 库存查询 / 入库 / 出库 / 调拨 / 盘点 / 库存明细 | [stock.md](./references/stock.md) |
| `product` | 产品管理：产品 CRUD / 单位 / 计量 / 分类 / 指标 | [product.md](./references/product.md) |
| `customer` / `supplier` | 客户供应商：CRUD / 发票 / 结算 / 信用额度 | [partner.md](./references/partner.md) |
| `contract` | 合同管理：采购合同&长协 / 销售合同&长协 / 运输合同&长协 / 服务合同&长协 | [contract.md](./references/contract.md) |
| `sale` | 销售管理：销售出库 / 销售退货 | [sale-purchase.md](./references/sale-purchase.md) |
| `purchase` | 采购管理：采购入库 / 采购退货 | [sale-purchase.md](./references/sale-purchase.md) |
| `finance` | 财务管理：账户 / 付款 / 收款 / 退款 / 收开票 / 付款申请 / 预付申请 / 结算 / 核销 / 开票申请 / 收付款 / 账户结算 | [finance.md](./references/finance.md) |
| `order` | 订单管理：主订单 / 计划（采购/销售 + CRUD） | [order.md](./references/order.md) |
| `produce` | 生产管理：生产单 / 生产计划 | [produce.md](./references/produce.md) |
| `waybill` | 运单管理：运单全生命周期 / 推送配置 | [waybill.md](./references/waybill.md) |
| `quality` | 质检管理：称重质检 / 质检单 | [quality-weight.md](./references/quality-weight.md) |
| `weight` | 称重管理：称重数据 / 关联运单 | [quality-weight.md](./references/quality-weight.md) |
| `stats` | 数据统计：总览 / 库存 / 财务 / 销售 / 采购 / 生产 | [stats-report.md](./references/stats-report.md) |
| `report` | 报表：库存报表 / 采购报表 / 销售报表 | [stats-report.md](./references/stats-report.md) |
| `homepage` / `screen` | 首页仪表盘 / 大屏数据 | [stats-report.md](./references/stats-report.md) |
| `settlement` | 结算管理：运单结算 | [settlement-invoice.md](./references/settlement-invoice.md) |
| `invoice` | 发票管理：发票 CRUD | [settlement-invoice.md](./references/settlement-invoice.md) |
| `system` | 系统管理：用户 / 部门 / 角色 / 菜单 / 字典 | [system.md](./references/system.md) |
| `operate-log` | 操作日志查询 | 本文档「辅助模块」章节 |
| `data-sync` | 数据同步消息查询 / 重发 | 本文档「辅助模块」章节 |
| `profit-event` | 利润事件查询 / 统计 / 类型 | 本文档「辅助模块」章节 |
| `profit-calculation` / `job-trigger` | 利润重算 / 定时任务触发（多为写操作） | 本文档「辅助模块」章节 |
| `api` | 通用 API 调用（兜底） | 本文档「通用 API 调用」章节 |

## 核心流程（每次请求必须执行）

1. **鉴权参数**：每条业务命令都必须携带 `--token <accessToken>` 与 `--tenant-id <租户ID>`（无状态鉴权，从用户/会话获取；`wlt version`/`config`/`completion`/`help` 除外）。若返回退出码 4 提示「缺少必填鉴权参数」，即未携带这两个 flag
2. **意图识别**：判断用户请求属于哪个模块（见「意图判断决策树」）
3. **参考文件加载**：按模块总览读取对应 `references/*.md`，按其中命令参考执行
4. **通用 API 兜底**：当快捷命令不覆盖所需操作时，使用 `wlt api` 兜底
5. **追问**：以上步骤都无法判断时，主动追问用户澄清

## 鉴权样例（无状态，每条业务命令必带 `--token`/`--tenant-id`）

> 以下为完整可执行样例。token 已脱敏，替换为真实 accessToken 即可；`--tenant-id 999` 为示例租户。所有业务命令只需固定带上 `--token fee383b0****fc0 --tenant-id 999` 这一段。

```bash
# 1. 分页列表查询
wlt customer list --token fee383b0****fc0 --tenant-id 999 --page-no 1 --page-size 10

# 2. 精简列表（下拉/选择控件用，返回 id+name）
wlt product simple-list --token fee383b0****fc0 --tenant-id 999
wlt stock warehouse simple-list --token fee383b0****fc0 --tenant-id 999

# 3. 详情查询（id 必须从上一步 list/simple-list 结果中提取，不要编造）
wlt customer get --token fee383b0****fc0 --tenant-id 999 --id <客户ID>

# 4. 新增（--data 为 JSON 对象，完整必填字段见 references/partner.md）
wlt customer create --token fee383b0****fc0 --tenant-id 999 --data '{"name":"测试客户","category":"ENTERPRISE",...}'

# 5. 切换环境 / 临时覆盖后端地址
wlt customer list --token fee383b0****fc0 --tenant-id 999 --profile prod
wlt customer list --token fee383b0****fc0 --tenant-id 999 --base-url https://erpapi.w-lian.com

# 6. 调试：只预览请求、不发送（仍需带 token/tenant-id，header 中 token 显示为 dry-run-token）
wlt api GET /erp/customer/page --token fee383b0****fc0 --tenant-id 999 --params '{"pageNo":1,"pageSize":10}' --dry-run
```

## 意图判断决策树

用户提到"库存/仓库/入库/出库/调拨/盘点/库存查询" → `stock`
用户提到"产品/商品/计量/单位/产品分类/产品指标" → `product`
用户提到"客户/供应商/合作伙伴/发票抬头/结算账户/信用额度" → `customer` / `supplier`
用户提到"合同/长协/采购长协/销售长协/运输长协/服务长协/销售合同/运输合同/服务合同" → `contract`
用户提到"销售/卖出/销售出库/销售退货" → `sale`
用户提到"采购/买入/采购入库/采购退货" → `purchase`
用户提到"财务/账户/付款/收款/退款/核销/开票/转账/调账" → `finance`
用户提到"订单/主订单/排产/关联运单/取消订单/完成订单/采购计划/销售计划/运输计划" → `order`
用户提到"生产/生产单/生产计划/质检数据" → `produce`
用户提到"运单/物流/发货/签收/装卸/推送配置" → `waybill`
用户提到"质检/检验/质检单/质检报告" → `quality`
用户提到"称重/地磅/磅单" → `weight`
用户提到"统计/数据总览/排名/趋势/数据分析" → `stats`
用户提到"报表/导出/明细报表/汇总" → `report`
用户提到"首页/仪表盘/大屏/看板" → `homepage` / `screen`
用户提到"结算/运单结算/未结算" → `settlement`
用户提到"发票/开票/发票管理" → `invoice`
用户提到"用户/部门/角色/权限/菜单/字典/系统设置" → `system`

关键区分:
- `stock`（库存数量查询） vs `report stock`（库存报表/统计）
- `finance settlement`（财务结算单据） vs `settlement`（运单结算）
- `customer credit`（客户信用额度） vs `finance account`（财务账户余额）
- `sale out`（销售出库单） vs `stock out`（其他出库单）
- `purchase in`（采购入库单） vs `stock in`（其他入库单）
- `order main`（业务订单） vs `produce main`（生产工单）
- `quality inspection`（质检单） vs `quality weight`（称重质检数据）
- `invoice`（业务发票） vs `partner invoice`（客户/供应商发票抬头）

## 通用 API 调用

对于未提供快捷命令的端点，使用通用 API 调用：

```bash
# GET 请求
wlt api GET /erp/warehouse/simple-list

# POST 请求
wlt api POST /erp/warehouse/create --data '{"name":"新仓库"}'

# 带查询参数
wlt api GET /erp/stock/page --params '{"pageNo":1,"pageSize":20}'

# 调试（不发送请求）
wlt api GET /erp/warehouse/page --dry-run
```

## 辅助模块

以下模块命令较少，未单独建参考文件，统一在此说明。**注意**：profit-event / data-sync / operate-log 等部分端点在 SIT 环境实测返回后端异常或权限不足，详见文末「实测可用性速查」。

### 操作日志 operate-log

```bash
wlt operate-log list --page-size 20          # 分页查询操作日志
# 可选过滤：--module <模块> --type <操作类型> --user-name <用户>
```

> 实测：该端点要求业务编号参数，SIT 纯列表查询返回后端异常。查特定单据的操作记录建议用 `wlt api GET /erp/operate-log/page --params '{...}'` 兜底。

### 数据同步 data-sync

```bash
wlt data-sync list --page-size 20            # 分页查询数据同步消息
wlt data-sync get --id <ID>                  # 获取同步消息详情
wlt data-sync resend --id <ID>               # 重新发送
```

> 实测：SIT 该端点返回 403（无权限）。

### 利润 profit-event / profit-calculation

```bash
wlt profit-event list --page-size 20         # 利润事件列表
wlt profit-event statistics                  # 利润事件统计
wlt profit-event types                       # 利润事件类型
wlt profit-event health                      # 健康检查
wlt profit-event retry --id <ID>             # 重试事件
wlt profit-calculation batch-recalculate-all # 批量重算（写操作，谨慎）
```

> 路径已修正为 `/erp/profit/event/*`（原 `/erp/profit-event/*` 返回 404）。待有效 token 复测。

### 定时任务 job-trigger

```bash
wlt job-trigger execute-product-cost         # 执行产品成本记录
wlt job-trigger execute-receivable-balance   # 执行应收余额计算
```

> 均为手动触发写操作，执行前必须获得用户确认。

## 常见工作流

> 以下示例为简洁起见**省略了 `--token` 与 `--tenant-id`**。实际执行时，**每条**业务命令都必须携带这两个 flag（例如 `wlt stock warehouse simple-list --token <accessToken> --tenant-id <租户ID>`），否则会以退出码 4 报「缺少必填鉴权参数」。

### 1. 入库全流程

```bash
wlt stock warehouse simple-list                    # 获取仓库列表
wlt product simple-list                            # 获取产品列表
wlt stock in create --data '{"warehouseId":1,...}' # 创建入库单
wlt stock in update-status --data '{"id":1,"status":2}' # 审核通过
```

### 2. 销售出库全流程

```bash
wlt sale out list --status 0                       # 查询待审核出库单
wlt sale out get --id <ID>                         # 查看详情
wlt sale out update-status --data '{"id":...,"status":2}' # 审核
```

### 3. 财务对账

```bash
wlt finance account list                           # 查看所有账户
wlt finance account settlement-page --account-id 1 # 查看账户流水
wlt finance receipt list --customer-id 1           # 查看客户收款
wlt finance payment list --supplier-id 1           # 查看供应商付款
```

### 4. 运单签收

```bash
wlt waybill source list --status 1                 # 查询在途运单
wlt waybill source sign --data '{"id":...}'        # 签收
wlt waybill source batch-sign --ids 1,2,3          # 批量签收
```

### 5. 生产质检

```bash
wlt produce main list --status 2                   # 查询已完工生产单
wlt produce main quality-page --produce-id 1       # 查看质检数据
wlt quality inspection create --data '...'         # 创建质检单
```

### 6. 统计分析

```bash
wlt stats overview --start-time 2024-01-01 --end-time 2024-12-31  # 总览
wlt stats finance data-overview --start-time 2024-01-01           # 财务统计
wlt stats sale customer-rankings --start-time 2024-01-01          # 客户排名
wlt report stock detail --warehouse-id 1 --start-time 2024-01-01  # 库存报表
```

## 危险操作确认

以下操作为不可逆或高影响操作，执行前**必须先向用户展示操作摘要并获得明确同意**：

| 模块 | 命令 | 说明 |
|------|------|------|
| `stock` | `warehouse delete` | 删除仓库 |
| `stock` | `in/out/move/check delete` | 删除出入库调拨盘点单 |
| `product` | `delete` | 删除产品 |
| `customer/supplier` | `delete` / `delete-list` | 删除客户/供应商 |
| `contract <子类>` | `delete` | 删除合同 |
| `finance` | `account delete` / 所有 delete | 删除财务单据 |
| `order` | `delete` / `cancel` | 删除/取消订单 |
| `produce` | `delete` | 删除生产单 |
| `waybill` | `delete` / `delete-list` | 删除运单 |
| `system` | `user/dept/role/menu delete` | 删除系统资源 |

### 确认流程

```
Step 1 → 展示操作摘要（操作类型 + 目标对象 + 影响范围）
Step 2 → 用户明确回复确认（如 "确认" / "好的"）
Step 3 → 执行命令
```

## 错误处理

所有错误以 JSON 格式输出到 stderr，包含 type、code、message、hint 字段。

| 退出码 | type | 含义 | 处理方式 |
|--------|------|------|---------|
| 0 | — | 成功 | — |
| 1 | `general` | 通用错误 | 检查错误信息 |
| 2 | `config` | 配置错误 | 运行 `wlt config init` |
| 3 | `authentication` | 鉴权错误（保留） | 当前服务端 token 拒绝统一表现为退出码 5（见下） |
| 4 | `validation` | 参数错误 | 检查命令参数；**缺少 `--token`/`--tenant-id` 也归此类**，补齐后重试 |
| 5 | `api_error` | API 错误 | 用 `wlt api GET <path> --dry-run` 调试；若 message 含 `code=401`/`账号未登录`，说明 token 失效，需重新获取并通过 `--token` 传入 |
| 6 | `network` | 网络错误 | 检查网络连接与 `--base-url`/profile 的 base_url |

### 错误处理流程

1. 遇到参数错误（退出码 4），先确认是否漏带 `--token`/`--tenant-id`，再使用 `--help` 查看命令参数说明
2. 遇到 API 错误（退出码 5），若 message 含 `code=401`/`账号未登录`，重新获取 token 并通过 `--token` 传入；否则用 `wlt api GET <path> --dry-run` 检查请求
3. 遇到网络错误（退出码 6），检查网络和 `--base-url`/profile 的 base_url
4. **严禁**连续重试超过 3 次相同命令；3 次仍失败必须停止并报告
5. 仍然失败，**立即停止**并报告完整错误信息

## 实测可用性速查（SIT 环境，v0.1.7）

> v0.1.7 在 SIT 实测结论。✅ 可用；⚠️ 需必填参数；❌ 当前后端/权限问题（非 wlt 缺陷，改用 `wlt api` 兜底或联系管理员）。
>
> **本轮（2026-06-27）变更**：修复 `profit-event` / `stock-report` / `homepage` / `job-trigger` 路径 bug；修复 filter flag 的 kebab→camelCase 参数名（此前 `--supplier-id` 等多词筛选不生效）；新增 finance invoice/payment-apply/prepayment-apply 查询、report direct、stock 与 stock-record 多个查询/统计端点、legacy 单据域 page-count。标 🛠 项为路径已对齐文档、待有效 token 实测。

### ✅ 完全可用
- **库存 stock**：warehouse / query(list·count·get) / record(list) / in / out / move / check
- **产品 product**：list / simple-list / category / unit / metrics / get-metrics / metrics item-list
- **客户/供应商**：list / simple-list / page-count / credit（仅客户）/ invoice / settlement
- **合同 contract**：7 个子类（purchase-long-cooperate / sale-contract / sale-long-cooperate / transport / transport-long / service-contract / service-long），每个子类含 list / page-count / get / update-status / create / update / delete
- **销售/采购**：sale out / purchase in
- **财务 finance**：account / account-settlement / invoice-apply / payment / receipt / receipt-payment / refund / settlement / write-off（list·summary·get）
- **订单/生产**：order main·plan / produce main(含 quality-page)·plan
- **运单/质检/称重**：waybill push-config get / quality inspection list·summary / weight waybill list
- **统计/报表/首页/大屏**：stats / report / homepage / screen 全部可用
- **系统/结算/发票**：system user·dept / settlement main / invoice main

### ⚠️ 需必填参数
- `product get-metrics --id <产品ID>`
- `product metrics item-list --metric-id <指标ID>`
- `stock query count --product-id <产品ID>`（后端强制）
- `quality weight order-page/order-summary --type SALE|PURCHASE`
- `quality inspection relate-list --business-type <类型> --business-id <ID>`

### ❌ 当前不可用（后端异常或权限）
- **后端 404/500**：`waybill source list/page-count`、`system role list`、`system dict type-list/data-list`、`stock record count`、`quality weight waybill-page/waybill-summary`、`operate-log list`
- **权限 403**：`sale return list`、`purchase return list`、`system menu list`、`data-sync list`
- **命令未注册**：`supplier credit *`（客户有 credit，供应商未实现；`supplier --help` 误列）

### 🛠 本轮新增/修复（路径已对齐文档，待有效 token 实测）
- **路径修复**：`profit-event *`（→`/erp/profit/event/*`）、`report stock warehouse/buy-send/finance/produce`（→对应 `*-page` 分页端点）、`job-trigger *`（→`/erp/job-trigger/*`）、`homepage dashboard2/inventory-backlog/product-ranking`（不再被强制到 dashboard6）
- **筛选修复**：多词 filter flag 现正确转 camelCase（`--supplier-id`→`supplierId`），此前此类筛选条件不生效
- **新增查询**：
  - 财务：`finance invoice / payment-apply / prepayment-apply`（list·get·summary）
  - 报表：`report direct detail / detail-count`、`report stock *-count`、`report purchase/sale *-count`
  - 库存：`stock query page-count / detail-count / stock-record-count`、`stock record page-count / record-page / total-cost`
  - 单据域：stock in/out/move/check、purchase in/return、sale out/return 现均含 `page-count`

## 详细参考（按需读取）

- [references/auth-config.md](./references/auth-config.md) — 认证与配置
- [references/stock.md](./references/stock.md) — 库存管理
- [references/product.md](./references/product.md) — 产品管理
- [references/partner.md](./references/partner.md) — 客户供应商
- [references/contract.md](./references/contract.md) — 合同管理
- [references/sale-purchase.md](./references/sale-purchase.md) — 销售采购
- [references/finance.md](./references/finance.md) — 财务管理
- [references/order.md](./references/order.md) — 订单管理
- [references/produce.md](./references/produce.md) — 生产管理
- [references/waybill.md](./references/waybill.md) — 运单管理
- [references/quality-weight.md](./references/quality-weight.md) — 质检称重
- [references/stats-report.md](./references/stats-report.md) — 统计报表
- [references/settlement-invoice.md](./references/settlement-invoice.md) — 结算发票
- [references/system.md](./references/system.md) — 系统管理
- [references/api-conventions.md](./references/api-conventions.md) — API 约定与通用模式
