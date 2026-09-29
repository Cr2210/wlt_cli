# 业务主线与数据关联 (business flows)

> 本文回答「一个业务对象从头到尾经过哪些模块、用什么命令衔接」。各环节的完整命令/参数明细见对应的 `references/*.md`；本文只讲主线、衔接点与传递的关键 ID 字段。
>
> **通用约定**：`enterprise-id`（企业 ID）即客户/供应商档案 ID（`customer`/`supplier` 的 ID），在各模块筛选中通用。单号前缀：CGDD=采购订单、XSDD=销售订单、CGJH=采购计划、XSJH=销售计划、CGRK=采购入库、XSCK=销售出库、JSD=结算单、YD=运单。

## 三大主线速览

```text
采购（付钱）：供应商 → 采购合同 → 采购计划 → 采购订单 → 采购入库 → 质检/称重 → 库存
                → 采购结算单 → 收票/付款 → 核销
销售（收钱）：客户 → 销售合同 → 销售计划 → 销售订单 → 销售出库 → 运单签收
                → 销售结算单 + 销售核算 → 开票申请 → 发票 → 收款 → 核销
物流（付运费）：承运商 → 运输合同 → 运单（装→卸→签） → 物流结算单 → 运费申请 → 付款 → 核销
```

## 主线一：采购（供应商侧，从建档到核销）

| # | 环节 | 命令入口 | 衔接点（传给下一步什么） |
|---|------|----------|--------------------------|
| 1 | 供应商建档 | `wlt supplier create` | 供应商 ID（后续作为 `--supplier-id` / `--enterprise-id`） |
| 2 | 采购合同/长协 | `wlt contract purchase-long-cooperate`（详见 [contract.md](./contract.md)） | 合同约束价格；订单/入库可引用 |
| 3 | 采购计划 | `wlt order plan purchase list`（CGJH…，`--supplier-id`） | 计划 → 生成采购订单（`order main create` 从计划数据提取字段） |
| 4 | 采购订单 | `wlt order main purchase list`（CGDD…） | 订单 ID → `--order-id`；可 `link-waybill` 关联运单 |
| 5 | 采购入库 | `wlt purchase in list`（CGRK…，`--supplier-id`） | 入库审核后形成库存；`stock record --biz-no` 可反查 |
| 6 | 质检/称重 | `wlt quality inspection relate-list --business-type --business-id`；`wlt quality weight order-page --type PURCHASE` | 质检指标（如含水）影响结算计价 |
| 7 | 库存 | `wlt stock query page --supplier-id --plan-no`；明细 `wlt stock record page --biz-type --biz-no` | 库存可追溯到来源计划/供应商 |
| 8 | 采购结算 | `wlt settlement main list --type PURCHASE_SETTLEMENT`（JSD…）；创建前 `unsettle-waybill --supplier-id` 查可结算运单 | 结算单 ID → 付款/收票/核销的对象 |
| 9 | 收票 | `wlt finance invoice page --type INVOICE_PAYMENT` | 发票核销 `WRITE_OFF_IN` |
| 10 | 付款 | `wlt finance receipt-payment page --type PAYMENT`（预付见下方支线） | 资金核销 `WRITE_OFF_PURCHASE` |
| 11 | 核销 | `wlt finance write-off page --type WRITE_OFF_PURCHASE`（`writeOffType=AMOUNT` 资金 / `WRITE_OFF_IN` 发票） | 主线终点 |

## 主线二：销售（客户侧，从建档到核销）

| # | 环节 | 命令入口 | 衔接点 |
|---|------|----------|--------|
| 1 | 客户建档 | `wlt customer create`（信用额度 `wlt customer credit`） | 客户 ID（`--customer-id` / `--enterprise-id`） |
| 2 | 销售合同/长协 | `wlt contract sale-contract`（HT…）/ `sale-long-cooperate`（XY…） | 同采购链 |
| 3 | 销售计划 | `wlt order plan sale list`（XSJH…，`--customer-id`） | → 销售订单 |
| 4 | 销售订单 | `wlt order main sale list`（XSDD…） | `link-waybill` 关联运单；`get-linkorder-by-orderId` 反查运单 |
| 5 | 销售出库 | `wlt sale out list`（XSCK…，`--customer-id --batch-no`） | 出库扣减库存；运单 `--order-type SALE_OUT` |
| 6 | 运单 | `wlt waybill page --order-type SALE_OUT`；`load`（UN_LOAD）→ `unload`（ON_LOAD）→ `sign-batch` | 签收后的运单进入结算；`waybill get` 详情含关联订单/合同/**结算标志** |
| 7 | 销售结算 | `wlt settlement main list --type SALE_SETTLEMENT`；`unsettle-waybill --customer-id` | 结算单 → 开票/收款/核销 |
| 8 | 销售核算 | `wlt sale accounting list`（按客户 + 核算周期）；候选运单 `waybill-list`（`settlementId` 非空 = 已被结算单锁定） | 核算价可 `calculate` 试算（取运单最新已审核质检单指标） |
| 9 | 开票申请 | `wlt finance invoice-apply page`（走审批流 `update-status`） | 申请 → 开票 |
| 10 | 发票 | `wlt finance invoice page --type INVOICE_RECEIPT`；业务发票 `wlt invoice main` | 发票核销 `WRITE_OFF_OUT` |
| 11 | 收款 | `wlt finance receipt-payment page --type RECEIPT` | 资金核销 `WRITE_OFF_SALE` |
| 12 | 核销 | `wlt finance write-off page --type WRITE_OFF_SALE` | 主线终点 |

## 主线三：物流（承运商侧，运费从发生到支付）

| # | 环节 | 命令入口 | 衔接点 |
|---|------|----------|--------|
| 1 | 承运商/运输合同 | `wlt contract transport`（HT…）；运单侧承运商即 `--capacity-name` / `--carrier-name` | 承运商 ID（`--carrier-id` / `--carrier-enterprise-id`） |
| 2 | 运单生命周期 | `wlt waybill page --capacity-name <承运商>`；`load` → `unload` → `sign-batch` | 已签收且未结算的运单是下游的结算素材 |
| 3 | 物流结算 | `wlt settlement logistics list`（承运商维度）；`available-waybill --carrier-id --settlement-id` 选可选运单；`mismatch-list` 查金额差异运单 | 物流结算单 ID → 运费申请 |
| 4 | 运费申请 | `wlt finance freight-apply list`；`available-waybill` / `available-settlement`（`availableCount=0` 置灰）带出明细；`update-status` 审批 | 审批通过后进入付款环节 |
| 5 | 付款 | `wlt finance receipt-payment page --type PAYMENT` | 核销 |
| 6 | 核销 | `wlt finance write-off page` | 主线终点 |

> 忽略剩余金额（尾差）：`settlement main|logistics ignore-remainder --type SETTLE|INVOICE`（`INVOICE` 要求 `invoiceFlag=无需开票` 且已审核）。

## 支线一：生产（produce）

```text
生产方案(produce plan, 含 --customer-id) → 生产单(produce main, 可按 --order-id/--plan-no 关联订单与方案)
  → 生产质检(produce main quality-page --produce-id / quality inspection)
  → 产成品入库形成库存（stock record --biz-type 可追溯） → 进入销售主线
```

## 支线二：预付款（prepayment）

```text
付款单标记预付(finance receipt-payment --is-prepaid)
  → 关联订单(order prepayment-relation create --relation-type PAYMENT|SUPPLIER)
     · PAYMENT=付款单 / SUPPLIER=供应商期初
     · 先用 available-payments / available-initials 查供应商可用余额
     · --biz-type PURCHASE(采购订单, 默认) | LOGISTICS(物流订单, 仅 PAYMENT)
  → 调整金额 update-amount；随订单后续核销
```

## 支线三：称重 → 运单 → 质检计价（weight / quality）

```text
磅单(weight waybill list --weighing-no) → 匹配运单(weight waybill match-link-waybill-source)
  → 运单装卸重量(waybill load --load-weight / unload --unload-weight)
  → 称重质检(quality weight order-page|waybill-page, 按 send/receive inspection no 区分发货/收货质检)
  → 计价(sale accounting calculate 按规则取运单最新已审核质检单指标)
```

## 三处「选运单」命令对比（易混淆）

| 场景 | 命令 | 筛选维度 | 说明 |
|------|------|----------|------|
| 销售核算选候选运单 | `wlt sale accounting waybill-list` | `--customer-id`（客户） | 返回游离 + 已绑定运单；`settlementId` 非空 = locked（已被结算单占用） |
| 物流结算选运单 | `wlt settlement logistics available-waybill` | `--carrier-id`（承运商）、`--settlement-id` | 创建物流结算单时选已签收未结算运单；结算单详情回显 |
| 运费申请选运单/结算单 | `wlt finance freight-apply available-waybill` / `available-settlement` | `--carrier-id`、`--settlement-id`、`--apply-id` | 从运单或物流结算单带出明细；`available-settlement` 的 `availableCount=0` 表示不可选 |

## 关联字段/命令速查

| 从 | 到 | 用法 |
|----|----|------|
| 订单 | 运单 | `order main link-waybill --data '{"orderId":..,"waybillId":..}'`；反查 `get-linkorder-by-orderId --order-id` |
| 运单 | 订单/合同/结算标志 | `waybill get --id`（详情 VO 含 metricsList、关联订单、合同、图片 URL、结算标志） |
| 生产单 | 订单/方案 | `produce main page --order-id / --order-no / --plan-no` |
| 库存 | 计划/供应商 | `stock query page --plan-no / --supplier-id` |
| 库存明细 | 任意业务单据 | `stock record page --biz-type / --biz-no` |
| 质检单 | 业务单据 | `quality inspection relate-list --business-type --business-id` |
| 结算单 | 未结算运单 | `settlement main unsettle-waybill --customer-id / --supplier-id / --warehouse-id` |
| 预付付款单 | 订单 | `order prepayment-relation list/create --biz-type PURCHASE\|LOGISTICS` |
| 账户流水 | 业务单据 | `finance account-settlement list --business-no / --business-id` |
