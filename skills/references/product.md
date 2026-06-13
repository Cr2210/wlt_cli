# 产品管理 (product)

## 产品 CRUD (`wlt product`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt product list` | 分页查询产品 | `--name`, `--category-id`, `--status`, `--bar-code`, `--spec`, `--page-no`, `--page-size` |
| `wlt product get --id <N>` | 获取产品详情 | `--id`（必填） |
| `wlt product create --data '<json>'` | 创建产品 | `--data`（必填） |
| `wlt product update --data '<json>'` | 更新产品 | `--data`（必填） |
| `wlt product delete --id <N>` | 删除产品 | `--id`（必填，单个删除） |
| `wlt product update-status` | 更新产品状态 | `--id`, `--status` |
| `wlt product simple-list` | 产品精简列表 | 无 |
| `wlt product get-metrics --id <N>` | 获取产品指标 | `--id`（必填） |
| `wlt product get-history-metrics --id <N>` | 获取历史指标 | `--id`（必填）, `--name` |

## 产品单位 (`wlt product unit`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt product unit list` | 分页查询单位 | `--name`, `--status`, `--page-no`, `--page-size` |
| `wlt product unit get --id <N>` | 获取单位详情 | `--id`（必填） |
| `wlt product unit create --data '<json>'` | 创建单位 | `--data`（必填） |
| `wlt product unit update --data '<json>'` | 更新单位 | `--data`（必填） |
| `wlt product unit delete --id <N>` | 删除单位 | `--id`（必填） |
| `wlt product unit update-status` | 更新单位状态 | `--id`, `--status` |
| `wlt product unit simple-list` | 单位精简列表 | 无 |

## 产品计量 (`wlt product metrics`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt product metrics list` | 分页查询计量 | `--name`, `--status`, `--page-no`, `--page-size` |
| `wlt product metrics get --id <N>` | 获取计量详情 | `--id`（必填） |
| `wlt product metrics create --data '<json>'` | 创建计量 | `--data`（必填） |
| `wlt product metrics update --data '<json>'` | 更新计量 | `--data`（必填） |
| `wlt product metrics delete --id <N>` | 删除计量 | `--id`（必填） |
| `wlt product metrics update-status` | 更新计量状态 | `--id`, `--status` |
| `wlt product metrics simple-list` | 计量精简列表 | `--name` |
| `wlt product metrics order-item-list` | 订单项列表 | `--id`（必填，订单项 ID） |
| `wlt product metrics item-list` | 计量项列表 | `--metric-id`（必填） |
| `wlt product metrics add-metrics-items` | 添加计量项 | `--data`（必填） |

## 产品分类 (`wlt product category`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt product category list` | 查询分类（非分页） | `--name`, `--status`, `--parent-id` |
| `wlt product category get --id <N>` | 获取分类详情 | `--id`（必填） |
| `wlt product category create --data '<json>'` | 创建分类 | `--data`（必填） |
| `wlt product category update --data '<json>'` | 更新分类 | `--data`（必填） |
| `wlt product category delete --id <N>` | 删除分类 | `--id`（必填） |
| `wlt product category update-status` | 更新分类状态 | `--id`, `--status` |
| `wlt product category simple-list` | 分类精简列表 | 无 |
