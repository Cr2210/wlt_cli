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

## 付款单 (`wlt finance payment`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance payment list` | 分页查询付款 | `--payment-no`, `--supplier-id`, `--status`, `--account-id`, `--reviewer-id`, `--page-no`, `--page-size` |
| `wlt finance payment get --id <N>` | 获取付款详情 | `--id`（必填） |
| `wlt finance payment create --data '<json>'` | 创建付款 | `--data`（必填） |
| `wlt finance payment update --data '<json>'` | 更新付款 | `--data`（必填） |
| `wlt finance payment delete --id <N>` | 删除付款 | `--id`（必填） |
| `wlt finance payment update-status --data '<json>'` | 更新付款状态 | `--data`（必填） |
| `wlt finance payment summary` | 付款汇总 | 无 |
| `wlt finance payment export` | 导出 Excel | `--payment-no`, `--supplier-id`, `--status`, `--account-id`, `--reviewer-id` |

## 收款单 (`wlt finance receipt`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance receipt list` | 分页查询收款 | `--receipt-no`, `--customer-id`, `--status`, `--account-id`, `--reviewer-id`, `--page-no`, `--page-size` |
| `wlt finance receipt get --id <N>` | 获取收款详情 | `--id`（必填） |
| `wlt finance receipt create --data '<json>'` | 创建收款 | `--data`（必填） |
| `wlt finance receipt update --data '<json>'` | 更新收款 | `--data`（必填） |
| `wlt finance receipt delete --id <N>` | 删除收款 | `--id`（必填） |
| `wlt finance receipt update-status --data '<json>'` | 更新收款状态 | `--data`（必填） |
| `wlt finance receipt summary` | 收款汇总 | 无 |
| `wlt finance receipt export` | 导出 Excel | `--receipt-no`, `--customer-id`, `--status`, `--account-id`, `--reviewer-id` |

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

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance write-off list` | 分页查询 | `--write-off-no`, `--customer-id`, `--supplier-id`, `--status`, `--account-id`, `--reviewer-id`, `--page-no`, `--page-size` |
| `wlt finance write-off get --id <N>` | 获取详情 | `--id`（必填） |
| `wlt finance write-off create --data '<json>'` | 创建 | `--data`（必填） |
| `wlt finance write-off update --data '<json>'` | 更新 | `--data`（必填） |
| `wlt finance write-off delete --id <N>` | 删除 | `--id`（必填） |
| `wlt finance write-off update-status --data '<json>'` | 更新状态 | `--data`（必填） |
| `wlt finance write-off summary` | 汇总 | 无 |
| `wlt finance write-off get-item-relate --id <N>` | 获取关联项 | `--id`（必填） |
| `wlt finance write-off export` | 导出 Excel | `--write-off-no`, `--customer-id`, `--supplier-id`, `--status`, `--account-id`, `--reviewer-id` |

## 开票申请 (`wlt finance invoice-apply`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance invoice-apply list` | 分页查询 | `--apply-no`, `--customer-id`, `--supplier-id`, `--status`, `--account-id`, `--reviewer-id`, `--page-no`, `--page-size` |
| `wlt finance invoice-apply get --id <N>` | 获取详情 | `--id`（必填） |
| `wlt finance invoice-apply create --data '<json>'` | 创建 | `--data`（必填） |
| `wlt finance invoice-apply update --data '<json>'` | 更新 | `--data`（必填） |
| `wlt finance invoice-apply delete --id <N>` | 删除 | `--id`（必填） |
| `wlt finance invoice-apply update-status --data '<json>'` | 更新状态 | `--data`（必填） |
| `wlt finance invoice-apply summary` | 汇总 | 无 |
| `wlt finance invoice-apply export` | 导出 Excel | `--apply-no`, `--customer-id`, `--supplier-id`, `--status`, `--account-id`, `--reviewer-id` |

## 收付款记录 (`wlt finance receipt-payment`)

> `type=RECEIPT` 收款 / `type=PAYMENT` 付款

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt finance receipt-payment page` | 分页查询收付款记录 | `--page-no`, `--page-size`, `--no`, `--type`, `--pay-date`, `--account-id`, `--account-name`, `--account-no`, `--partner-id`, `--partner-name`, `--service-user-id`, `--service-user-name`, `--status`, `--approve-status`, `--remark`, `--creator-name`, `--updater-name`, `--create-time`, `--update-time`, `--keyword`, `--custom-order`, `--headers` |
| `wlt finance receipt-payment page-count` | 按筛选统计收付款记录数量 | 同 `receipt-payment page`（去 `--headers`） |
| `wlt finance receipt-payment get --id <N>` | 获取收付款详情 | `--id`（必填） |
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
