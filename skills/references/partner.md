# 客户供应商 (customer / supplier)

客户 (`customer`, partnerType=1) 和供应商 (`supplier`, partnerType=2) 共享相同的命令结构。

## 客户管理 (`wlt customer`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt customer list` | 分页查询客户 | `--page-no`, `--page-size`, `--name`, `--status`, `--category-id` |
| `wlt customer get --id <N>` | 获取客户详情 | `--id`（必填） |
| `wlt customer create --data '<json>'` | 创建客户 | `--data`（必填） |
| `wlt customer update --data '<json>'` | 更新客户 | `--data`（必填） |
| `wlt customer delete --id <N>` | 删除客户 | `--id`（必填） |
| `wlt customer update-status` | 更新客户状态 | `--id`, `--status` |
| `wlt customer simple-list` | 客户精简列表 | 无 |
| `wlt customer page-count` | 统计客户数量 | `--name`, `--status` |
| `wlt customer update-audit-status` | 更新审核状态 | `--id`（必填）, `--status`（必填） |
| `wlt customer delete-list` | 批量删除客户 | `--ids`（必填，逗号分隔） |

## 供应商管理 (`wlt supplier`)

与客户命令结构相同，将 `customer` 替换为 `supplier`。**注意**：供应商没有 `credit` 子命令组。

## 客户/供应商发票 (`wlt customer|supplier invoice`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `invoice list` | 分页查询发票 | `--partner-id`, `--name`, `--page-no`, `--page-size` |
| `invoice get --id <N>` | 获取发票详情 | `--id`（必填） |
| `invoice create --data '<json>'` | 创建发票 | `--data`（必填） |
| `invoice update --data '<json>'` | 更新发票 | `--data`（必填） |
| `invoice delete --id <N>` | 删除发票 | `--id`（必填） |
| `invoice list-by-partner` | 按合作伙伴查询 | `--partner-id`（必填） |
| `invoice delete-list` | 批量删除 | `--ids`（必填，逗号分隔） |

## 客户/供应商结算 (`wlt customer|supplier settlement`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `settlement list` | 分页查询结算 | `--partner-id`, `--name`, `--status`, `--page-no`, `--page-size` |
| `settlement get --id <N>` | 获取结算详情 | `--id`（必填） |
| `settlement create --data '<json>'` | 创建结算 | `--data`（必填） |
| `settlement update --data '<json>'` | 更新结算 | `--data`（必填） |
| `settlement delete --id <N>` | 删除结算 | `--id`（必填） |
| `settlement page-count` | 统计结算数量 | `--partner-id`, `--name`, `--status` |
| `settlement list-by-partner` | 按合作伙伴查询 | `--partner-id`（必填） |
| `settlement delete-list` | 批量删除 | `--ids`（必填，逗号分隔） |

## 客户信用额度 (`wlt customer credit`) — 仅客户

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `credit list` | 分页查询信用额度 | `--customer-id`, `--status`, `--page-no`, `--page-size` |
| `credit get --id <N>` | 获取信用额度详情 | `--id`（必填） |
| `credit create --data '<json>'` | 创建信用额度 | `--data`（必填） |
| `credit update --data '<json>'` | 更新信用额度 | `--data`（必填） |
| `credit delete --id <N>` | 删除信用额度 | `--id`（必填） |
| `credit page-count` | 统计数量 | `--customer-id`, `--status` |
| `credit cancel --id <N>` | 取消信用额度 | `--id`（必填） |
| `credit valid-credit` | 查询有效额度 | `--customer-id`（必填） |
| `credit delete-batch` | 批量删除 | `--ids`（必填，逗号分隔） |
