# 销售采购 (sale / purchase)

销售和采购模块均使用标准 CRUD 子命令模式（list / get / create / update / delete / update-status）。

## 销售出库 (`wlt sale out`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt sale out list` | 分页查询销售出库单 | `--page-no`, `--page-size`, `--no`, `--status` |
| `wlt sale out get --id <N>` | 获取出库单详情 | `--id`（必填） |
| `wlt sale out create --data '<json>'` | 创建销售出库单 | `--data`（必填） |
| `wlt sale out update --data '<json>'` | 更新销售出库单 | `--data`（必填） |
| `wlt sale out delete --ids <id1,id2>` | 删除销售出库单 | `--ids`（必填） |
| `wlt sale out update-status --data '<json>'` | 更新出库单状态 | `--data`（必填，含 id 和 status） |

**API 路径**：`/erp/sale-out`

## 销售退货 (`wlt sale return`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt sale return list` | 分页查询销售退货单 | `--page-no`, `--page-size`, `--no`, `--status` |
| `wlt sale return get --id <N>` | 获取退货单详情 | `--id`（必填） |
| `wlt sale return create --data '<json>'` | 创建销售退货单 | `--data`（必填） |
| `wlt sale return update --data '<json>'` | 更新销售退货单 | `--data`（必填） |
| `wlt sale return delete --ids <id1,id2>` | 删除销售退货单 | `--ids`（必填） |
| `wlt sale return update-status --data '<json>'` | 更新退货单状态 | `--data`（必填，含 id 和 status） |

**API 路径**：`/erp/sale-return`

## 采购入库 (`wlt purchase in`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt purchase in list` | 分页查询采购入库单 | `--page-no`, `--page-size`, `--no`, `--status` |
| `wlt purchase in get --id <N>` | 获取入库单详情 | `--id`（必填） |
| `wlt purchase in create --data '<json>'` | 创建采购入库单 | `--data`（必填） |
| `wlt purchase in update --data '<json>'` | 更新采购入库单 | `--data`（必填） |
| `wlt purchase in delete --ids <id1,id2>` | 删除采购入库单 | `--ids`（必填） |
| `wlt purchase in update-status --data '<json>'` | 更新入库单状态 | `--data`（必填，含 id 和 status） |

**API 路径**：`/erp/purchase-in`

## 采购退货 (`wlt purchase return`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt purchase return list` | 分页查询采购退货单 | `--page-no`, `--page-size`, `--no`, `--status` |
| `wlt purchase return get --id <N>` | 获取退货单详情 | `--id`（必填） |
| `wlt purchase return create --data '<json>'` | 创建采购退货单 | `--data`（必填） |
| `wlt purchase return update --data '<json>'` | 更新采购退货单 | `--data`（必填） |
| `wlt purchase return delete --ids <id1,id2>` | 删除采购退货单 | `--ids`（必填） |
| `wlt purchase return update-status --data '<json>'` | 更新退货单状态 | `--data`（必填，含 id 和 status） |

**API 路径**：`/erp/purchase-return`
