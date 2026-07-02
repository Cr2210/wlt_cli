# 财务管理 (finance)

## 财务账户 (`wlt finance account`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance account list` | 分页查询账户 | `--account-id`, `--account-no`, `--name`, `--status`, `--remark`, `--page-no`, `--page-size` |
| `wlt finance account get --id <N>` | 获取账户详情 | `--id`（必填） |
| `wlt finance account create --data '<json>'` | 创建账户 | `--data`（必填） |
| `wlt finance account update --data '<json>'` | 更新账户 | `--data`（必填） |
| `wlt finance account delete --id <N>` | 删除账户 | `--id`（必填） |
| `wlt finance account update-status --data '<json>'` | 更新账户状态 | `--data`（必填） |
| `wlt finance account update-default-status` | 设置默认账户 | `--id`（必填）, `--default` |
| `wlt finance account simple-list` | 账户精简列表 | 无 |
| `wlt finance account adjust --data '<json>'` | 账户调账 | `--data`（必填） |
| `wlt finance account transfer --data '<json>'` | 账户转账 | `--data`（必填） |
| `wlt finance account settlement-page` | 账户流水 | `--account-id`, `--page-no`, `--page-size` |
| `wlt finance account export` | 导出 Excel | `--account-id`, `--account-no`, `--name`, `--status`, `--remark` |

## 收开票 (`wlt finance invoice`)

> `type=INVOICE_PAYMENT` 收票 / `type=INVOICE_RECEIPT` 开票

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance invoice page` | 分页查询收开票 | `--page-no`, `--page-size`, `--no`, `--type`, `--invoice-type`, `--invoice-no`, `--partner-id`, `--partner-name`, `--service-user-id`, `--service-user-name`, `--status`, `--approve-status`, `--invoice-date`, `--remark`, `--creator-name`, `--updater-name`, `--create-time`, `--update-time`, `--keyword`, `--custom-order`, `--headers` |
| `wlt finance invoice page-count` | 按筛选统计收开票数量 | 同 `invoice page`（去 `--headers`） |
| `wlt finance invoice get --id <N>` | 获取收开票详情 | `--id`（必填） |
| `wlt finance invoice summary` | 获取收开票合计（返回 totalCount/unSettledCount/partSettledCount/settledCount/totalAmount/checkedAmount/unCheckedAmount，支持同 `page-count` 的筛选） | 同 `invoice page-count` |
| `wlt finance invoice export` | 导出收开票 Excel | 同 `invoice page` |

## ⚠️ 付款单 / 收款单（已废弃）

`/erp/finance-payment/*` 与 `/erp/finance-receipt/*` 已废弃，统一使用 `/erp/finance-receipt-payment/*`（`wlt finance receipt-payment`，`type=RECEIPT` 收款 / `type=PAYMENT` 付款）。

## 退款单 (`wlt finance refund`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance refund list` | 分页查询退款 | `--refund-no`, `--customer-id`, `--supplier-id`, `--status`, `--account-id`, `--reviewer-id`, `--page-no`, `--page-size` |
| `wlt finance refund get --id <N>` | 获取退款详情 | `--id`（必填） |
| `wlt finance refund create --data '<json>'` | 创建退款 | `--data`（必填） |
| `wlt finance refund update --data '<json>'` | 更新退款 | `--data`（必填） |
| `wlt finance refund delete --id <N>` | 删除退款 | `--id`（必填） |
| `wlt finance refund update-status --data '<json>'` | 更新退款状态 | `--data`（必填） |
| `wlt finance refund summary` | 退款汇总 | 无 |
| `wlt finance refund export` | 导出 Excel | `--refund-no`, `--customer-id`, `--supplier-id`, `--status`, `--account-id`, `--reviewer-id` |

## 财务结算 (`wlt finance settlement`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance settlement list` | 分页查询 | `--settlement-no`, `--customer-id`, `--supplier-id`, `--status`, `--account-id`, `--reviewer-id`, `--page-no`, `--page-size` |
| `wlt finance settlement get --id <N>` | 获取详情 | `--id`（必填） |
| `wlt finance settlement create --data '<json>'` | 创建 | `--data`（必填） |
| `wlt finance settlement update --data '<json>'` | 更新 | `--data`（必填） |
| `wlt finance settlement delete --id <N>` | 删除 | `--id`（必填） |
| `wlt finance settlement update-status --data '<json>'` | 更新状态 | `--data`（必填） |
| `wlt finance settlement summary` | 汇总 | 无 |
| `wlt finance settlement finance-record --id <N>` | 财务记录 | `--id`（必填） |
| `wlt finance settlement export` | 导出 Excel | `--settlement-no`, `--customer-id`, `--supplier-id`, `--status`, `--account-id`, `--reviewer-id` |

## 核销单 (`wlt finance write-off`)

> 资金核销: `type=WRITE_OFF_PURCHASE/WRITE_OFF_SALE`, `writeOffType=AMOUNT`
> 发票核销: `type=WRITE_OFF_IN/WRITE_OFF_OUT`, `writeOffType=INVOICE`

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance write-off page` | 分页查询核销单 | `--page-no`, `--page-size`, `--no`, `--type`, `--write-off-type`, `--write-off-date`, `--approve-status`, `--partner-id`, `--partner-name`, `--service-user-id`, `--service-user-name`, `--source-no`, `--invoice-no`, `--target-no`, `--remark`, `--creator-name`, `--updater-name`, `--create-time`, `--update-time`, `--keyword`, `--custom-order`, `--headers` |
| `wlt finance write-off page-count` | 按筛选统计核销单数量 | 同 `write-off page`（去 `--headers`） |
| `wlt finance write-off get --id <N>` | 获取详情 | `--id`（必填） |
| `wlt finance write-off create --data '<json>'` | 创建 | `--data`（必填） |
| `wlt finance write-off update --data '<json>'` | 更新 | `--data`（必填） |
| `wlt finance write-off delete --id <N>` | 删除 | `--id`（必填） |
| `wlt finance write-off update-status --data '<json>'` | 更新状态 | `--data`（必填） |
| `wlt finance write-off summary` | 获取核销单合计（返回 inCount/outCount/totalAmount，支持同 `page-count` 的筛选） | 同 `write-off page-count` |
| `wlt finance write-off get-item-relate --id <N>` | 获取关联项 | `--id`（必填） |
| `wlt finance write-off export` | 导出 Excel | 同 `write-off page` |

## 开票申请 (`wlt finance invoice-apply`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance invoice-apply page` | 分页查询开票申请 | `--page-no`, `--page-size`, `--no`, `--partner-id`, `--partner-name`, `--invoice-title`, `--tax-no`, `--service-user-id`, `--service-user-name`, `--approve-status`, `--invoice-date`, `--remark`, `--creator-name`, `--updater-name`, `--create-time`, `--update-time`, `--custom-order`, `--keyword`, `--headers` |
| `wlt finance invoice-apply get --id <N>` | 获取详情 | `--id`（必填） |
| `wlt finance invoice-apply create --data '<json>'` | 创建 | `--data`（必填） |
| `wlt finance invoice-apply update --data '<json>'` | 更新 | `--data`（必填） |
| `wlt finance invoice-apply delete --id <N>` | 删除 | `--id`（必填） |
| `wlt finance invoice-apply update-status --data '<json>'` | 更新状态 | `--data`（必填） |
| `wlt finance invoice-apply summary` | 开票申请合计（返回 totalCount/checkCount/processCount/approveCount/rejectCount/totalAmount，支持同 page-count 的筛选） | 同 `invoice-apply page-count` |
| `wlt finance invoice-apply export` | 导出 Excel | 同 `invoice-apply page` |

## 收付款记录 (`wlt finance receipt-payment`)

> `type=RECEIPT` 收款 / `type=PAYMENT` 付款

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance receipt-payment page` | 分页查询收付款记录 | `--page-no`, `--page-size`, `--no`, `--type`, `--pay-date`, `--account-id`, `--account-name`, `--account-no`, `--partner-id`, `--partner-name`, `--service-user-id`, `--service-user-name`, `--status`, `--approve-status`, `--remark`, `--creator-name`, `--updater-name`, `--create-time`, `--update-time`, `--keyword`, `--custom-order`, `--headers` |
| `wlt finance receipt-payment page-count` | 按筛选统计收付款记录数量 | 同 `receipt-payment page`（去 `--headers`） |
| `wlt finance receipt-payment get` | 获取收付款详情 | `--id` 或 `--no`（任选其一） |
| `wlt finance receipt-payment create --data '<json>'` | 创建收付款 | `--data`（必填） |
| `wlt finance receipt-payment update --data '<json>'` | 更新收付款 | `--data`（必填） |
| `wlt finance receipt-payment delete --id <N>` | 删除收付款 | `--id`（必填） |
| `wlt finance receipt-payment update-status --data '<json>'` | 更新收付款状态 | `--data`（必填） |
| `wlt finance receipt-payment summary` | 获取收付款汇总数据（支持同 `page-count` 的筛选） | 同 `receipt-payment page-count` |
| `wlt finance receipt-payment export` | 导出收付款记录 Excel | 同 `receipt-payment page` |

## 账户结算 (`wlt finance account-settlement`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance account-settlement list` | 分页查询 | `--account-id`, `--settlement-no`, `--business-no`, `--business-id`, `--type`, `--page-no`, `--page-size` |
| `wlt finance account-settlement get --id <N>` | 获取详情 | `--id`（必填） |
| `wlt finance account-settlement create --data '<json>'` | 创建 | `--data`（必填） |
| `wlt finance account-settlement update --data '<json>'` | 更新 | `--data`（必填） |
| `wlt finance account-settlement delete --id <N>` | 删除 | `--id`（必填） |
| `wlt finance account-settlement export` | 导出 Excel | `--account-id`, `--settlement-no`, `--business-no`, `--business-id`, `--type` |
