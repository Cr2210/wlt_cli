# 订单管理 (order)

## 总览

订单模块覆盖：

| 子命令 | 说明 |
|---|---|
| `wlt order main` | 主订单：CRUD / 状态管理 / 运单关联 / 导出 |
| `wlt order plan` | 订单计划：采购计划 + 销售计划（分页） / CRUD / 业务动作 |
| `wlt order prepayment-relation` | 订单预付关联：关联记录查询 / 创建 / 改金额 / 可用付款单与期初查询 |

## 主订单 (`wlt order main`)

后端统一通过 `/erp/order` 端点 + `type` 字段区分子类型：

| 子命令 | 后端 type |
|---|---|
| `wlt order main purchase` | `PURCHASE`（`/erp/order`） |
| `wlt order main sale` | `SALE`（`/erp/order`） |

每个子类型自动拥有 `list` + `page-count`。

### 共用的筛选与分页标志

所有 order main purchase/sale 的 list / page-count 命令均支持：

| 标志 | 后端字段 | 说明 |
|---|---|---|
| `--no` | `no` | 订单号 |
| `--enterprise-id` | `enterpriseId` | 企业 ID |
| `--product-id` | `productId` | 产品 ID |
| `--order-start` | `orderTime[0]` | 下单时间起始 |
| `--order-end` | `orderTime[1]` | 下单时间结束 |
| `--page-no` | `pageNo` | 页码（默认 1） |
| `--page-size` | `pageSize` | 每页数量（默认 20） |

### 采购订单

```bash
wlt order main purchase list \
  --no CGDD20260402000002 \
  --enterprise-id 2001552070697033730 \
  --product-id 1927670802287808513 \
  --order-start "2026-07-17 00:00:00" --order-end "2026-08-06 23:59:59" \
  --page-no 1 --page-size 10

wlt order main purchase page-count --no CGDD20260402000002
```

### 销售订单

```bash
wlt order main sale list \
  --no XSDD20260402000001 \
  --enterprise-id 2001494968037298178 \
  --product-id 1927612799140261889 \
  --order-start "2026-07-14 00:00:00" --order-end "2026-08-11 23:59:59"

wlt order main sale page-count --enterprise-id 2001494968037298178
```

### 订单通用子命令（不区分 type）

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt order main get --id <N>` | 获取订单详情 | `--id`（必填） |
| `wlt order main get-linkorder-by-orderId` | 获取关联运单 | `--order-id`（必填） |
| `wlt order main create --data '<json>'` | 创建订单 | `--data`（必填） |
| `wlt order main update --data '<json>'` | 更新订单 | `--data`（必填） |
| `wlt order main delete --id <N>` | 删除订单 | `--id`（必填） |
| `wlt order main update-status --data '<json>'` | 更新订单状态 | `--data`（必填） |
| `wlt order main cancel --data '<json>'` | 取消订单 | `--data`（必填） |
| `wlt order main reopen --data '<json>'` | 重新打开订单 | `--data`（必填） |
| `wlt order main complete --data '<json>'` | 完成订单 | `--data`（必填） |
| `wlt order main link-waybill --data '<json>'` | 关联运单 | `--data`（必填） |
| `wlt order main unlink-waybill --data '<json>'` | 取消关联运单 | `--data`（必填） |
| `wlt order main export` | 导出 Excel | `--no`, `--enterprise-id`, `--product-id`, `--order-start`, `--order-end` |

## 订单计划 (`wlt order plan`)

后端统一通过 `/erp/order-plan` 端点 + `type` 字段区分子类型：

| 子命令 | 后端 type |
|---|---|
| `wlt order plan purchase` | `PURCHASE_TRANSPORT_PLAN`（`/erp/order-plan`）|
| `wlt order plan sale` | `SALE_TRANSPORT_PLAN`（`/erp/order-plan`）|

每个子类型自动拥有 `list` + `page-count`。

### 共用的筛选与分页标志

所有 plan purchase/sale 的 list / page-count 命令均支持：

| 标志 | 后端字段 | 说明 |
|---|---|---|
| `--no` | `no` | 计划单号 |
| `--product-id` | `productId` | 产品 ID |
| `--supplier-id` | `supplierId` | 供应商 ID（采购） |
| `--customer-id` | `customerId` | 客户 ID（销售） |
| `--start` | `startDate[0]` | 计划开始日期起始 |
| `--end` | `startDate[1]` | 计划开始日期结束 |
| `--page-no` | `pageNo` | 页码（默认 1） |
| `--page-size` | `pageSize` | 每页数量（默认 20） |

### 采购计划

```bash
wlt order plan purchase list \
  --no CGJH20260528000001 \
  --supplier-id 2001552070697033730 \
  --product-id 1927285600675729409 \
  --start "2026-07-21 00:00:00" --end "2026-08-19 23:59:59" \
  --page-no 1 --page-size 10

wlt order plan purchase page-count --no CGJH20260528000001
```

### 销售计划

```bash
wlt order plan sale list \
  --no XSJH20260401000001 \
  --customer-id 2001494968037298178 \
  --product-id 1927670802287808513 \
  --start "2026-07-20 00:00:00" --end "2026-08-20 23:59:59"

wlt order plan sale page-count --customer-id 2001494968037298178
```

### 计划通用子命令（不区分 type）

```bash
wlt order plan get --id <计划ID>
wlt order plan create --data '<json>'
wlt order plan update --data '<json>'
wlt order plan delete --id <计划ID>
wlt order plan update-status --data '{"id":1,"status":2}'
wlt order plan cancel --data '...'
wlt order plan reopen --data '...'
wlt order plan complete --data '...'
wlt order plan export [--no ...] [--product-id ...] [--supplier-id ...] [--customer-id ...] [--start ...] [--end ...]
```

## 订单预付关联 (`wlt order prepayment-relation`)

把订单与预付款来源（付款单 / 供应商期初）建立关联并维护关联金额。后端端点 `/erp/order-prepayment-relation`，`relationType` 枚举：`PAYMENT`（付款单）/ `SUPPLIER`（供应商期初）。

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `list` | 分页查询订单的关联记录 | `--order-id`（必填）、`--relation-type`、`--headers`、`--page-no`、`--page-size` |
| `create` | 创建关联，输出新关联 ID | `--order-id`、`--relation-type`、`--relation-id`、`--relation-amount`（均必填） |
| `update-amount` | 更新关联金额 | `--id`、`--relation-amount`（均必填） |
| `available-payments` | 查询供应商可用的付款单（含可用余额） | `--supplier-id`（必填）、`--order-id`、`--biz-date-start`、`--biz-date-end`、`--no`、`--headers`、`--page-no`、`--page-size` |
| `available-initials` | 查询供应商可用的期初 | 同 `available-payments` |

`list` 响应为裸数组（无 total）；`available-*` 为标准分页结果。金额单位均为元。

```bash
# 查某订单已关联的预付记录
wlt order prepayment-relation list --order-id 123

# 查供应商可用付款单余额 → 挑选后创建关联
wlt order prepayment-relation available-payments --supplier-id 456 --order-id 123
wlt order prepayment-relation create --order-id 123 --relation-type PAYMENT --relation-id 789 --relation-amount 1000

# 调整关联金额
wlt order prepayment-relation update-amount --id 1 --relation-amount 800
```

## 常见工作流

### 创建订单并关联运单

```bash
wlt order main create --data '{"orderNo":"ORD-001","customerId":1,...}'
wlt order main link-waybill --data '{"orderId":1,"waybillId":1}'
```

### 订单生命周期

```bash
wlt order main list --status 0               # 查询待处理
wlt order main update-status --data '...'    # 更新状态
wlt order main complete --data '...'         # 完成订单
```

### 采购/销售计划 → 订单

```bash
# 查询采购计划明细
wlt order plan purchase list --no CGJH20260528000001

# 并据此创建订单（从计划数据中提取字段）
wlt order main create --data '{...}'
```
