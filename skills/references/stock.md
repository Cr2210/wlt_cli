# 库存管理 (stock)

## 仓库管理 (`wlt stock warehouse`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt stock warehouse list` | 分页查询仓库 | `--page-no`, `--page-size`, `--name`, `--status` |
| `wlt stock warehouse get --id <N>` | 获取仓库详情 | `--id`（必填） |
| `wlt stock warehouse simple-list` | 仓库精简列表 | 无 |
| `wlt stock warehouse create --data '<json>'` | 创建仓库 | `--data`（必填） |
| `wlt stock warehouse update --data '<json>'` | 更新仓库 | `--data`（必填） |
| `wlt stock warehouse delete --ids <id1,id2>` | 删除仓库 | `--ids`（必填，逗号分隔） |
| `wlt stock warehouse update-status --data '<json>'` | 更新仓库状态 | `--data`（必填，含 id 和 status） |

## 库存查询 (`wlt stock query`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt stock query get` | 获取单个库存 | `--id` 或 `--product-id` + `--warehouse-id` |
| `wlt stock query list` | 分页查询库存 | `--page-no`, `--page-size`, `--warehouse-id`, `--product-name` |
| `wlt stock query count` | 统计库存数量 | `--product-id`（**必填**，后端强制）, `--warehouse-id`（可选）, `--metric-name`（可选） |
| `wlt stock query batch-detail` | 获取批次明细 | `--page-no`, `--page-size`, `--product-id`, `--warehouse-id` |

## 入库单 (`wlt stock in`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt stock in list` | 分页查询入库单 | `--page-no`, `--page-size`, `--no`, `--status` |
| `wlt stock in get --id <N>` | 获取入库单详情 | `--id`（必填） |
| `wlt stock in create --data '<json>'` | 创建入库单 | `--data`（必填） |
| `wlt stock in update --data '<json>'` | 更新入库单 | `--data`（必填） |
| `wlt stock in delete --ids <id1,id2>` | 删除入库单 | `--ids`（必填） |
| `wlt stock in update-status --data '<json>'` | 更新入库单状态 | `--data`（必填，含 id 和 status） |

## 出库单 (`wlt stock out`)

与入库单结构相同，将 `in` 替换为 `out`。

## 调拨单 (`wlt stock move`)

与入库单结构相同，将 `in` 替换为 `move`。

## 盘点单 (`wlt stock check`)

与入库单结构相同，将 `in` 替换为 `check`。

## 库存明细 (`wlt stock record`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt stock record list` | 分页查询明细 | `--page-no`, `--page-size`, `--product-id`, `--warehouse-id` |
| `wlt stock record get --id <N>` | 获取明细详情 | `--id`（必填） |
| `wlt stock record count` | 获取明细总数量 | ⚠️ SIT 后端异常，暂不可用 |

## 常见工作流

### 查询某仓库库存

```bash
wlt stock warehouse simple-list                    # 获取仓库列表
wlt stock query list --warehouse-id <ID>           # 查询库存
```

### 创建入库单

```bash
wlt stock warehouse simple-list                    # 获取仓库 ID
wlt product simple-list                            # 获取产品 ID
wlt stock in create --data '{"warehouseId":1,"items":[{"productId":1,"count":10}]}'
```

### 审核出入库单

```bash
wlt stock in list --status 0                       # 查询待审核
wlt stock in update-status --data '{"id":1,"status":2}' # 审核通过
```
