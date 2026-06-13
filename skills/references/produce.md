# 生产管理 (produce)

## 生产单 (`wlt produce main`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt produce main list` | 分页查询生产单 | `--produce-no`, `--status`, `--warehouse-id`, `--product-id`, `--page-no`, `--page-size` |
| `wlt produce main page-count` | 统计生产单数量 | `--produce-no`, `--status`, `--warehouse-id`, `--product-id` |
| `wlt produce main get --id <N>` | 获取生产单详情 | `--id`（必填） |
| `wlt produce main create --data '<json>'` | 创建生产单 | `--data`（必填） |
| `wlt produce main update --data '<json>'` | 更新生产单 | `--data`（必填） |
| `wlt produce main delete --id <N>` | 删除生产单 | `--id`（必填） |
| `wlt produce main update-status --data '<json>'` | 更新生产单状态 | `--data`（必填） |
| `wlt produce main export` | 导出 Excel | `--produce-no`, `--status`, `--warehouse-id`, `--product-id` |
| `wlt produce main quality-page` | 质检分页 | `--produce-id`, `--page-no`, `--page-size` |
| `wlt produce main quality-count` | 质检统计 | `--produce-id` |
| `wlt produce main quality-export` | 质检导出 | `--produce-id` |

## 生产计划 (`wlt produce plan`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt produce plan list` | 分页查询计划 | `--plan-no`, `--status`, `--warehouse-id`, `--product-id`, `--page-no`, `--page-size` |
| `wlt produce plan page-count` | 统计计划数量 | `--plan-no`, `--status`, `--warehouse-id`, `--product-id` |
| `wlt produce plan simple-list` | 计划精简列表 | `--warehouse-id` |
| `wlt produce plan get --id <N>` | 获取计划详情 | `--id`（必填） |
| `wlt produce plan create --data '<json>'` | 创建计划 | `--data`（必填） |
| `wlt produce plan update --data '<json>'` | 更新计划 | `--data`（必填） |
| `wlt produce plan delete --id <N>` | 删除计划 | `--id`（必填） |
| `wlt produce plan update-status --data '<json>'` | 更新计划状态 | `--data`（必填） |
| `wlt produce plan export` | 导出 Excel | `--plan-no`, `--status`, `--warehouse-id`, `--product-id` |
