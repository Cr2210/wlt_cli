# 生产管理 (produce)

## 生产任务 (`wlt produce main`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt produce main page` | 分页查询生产任务 | `--page-no`, `--page-size`, `--no`, `--project-name`, `--status`, `--produce-time`, `--batch-no`, `--metrics-name`, `--warehouse-id`, `--warehouse-name`, `--product-id`, `--product-name`, `--plan-no`, `--plan-name`, `--user-id`, `--user-name`, `--order-id`, `--order-no`, `--remark`, `--creator`, `--creator-name`, `--create-time`, `--update-time`, `--updater-name`, `--custom-order`, `--keyword`, `--headers` |
| `wlt produce main page-count` | 按筛选统计生产任务数量 | 同 `main page`（去 `--headers`） |
| `wlt produce main get --id <N>` | 获取生产任务详情 | `--id`（必填） |
| `wlt produce main create --data '<json>'` | 创建生产任务 | `--data`（必填） |
| `wlt produce main update --data '<json>'` | 更新生产任务 | `--data`（必填） |
| `wlt produce main delete --id <N>` | 删除生产任务 | `--id`（必填） |
| `wlt produce main update-status --data '<json>'` | 更新生产任务状态 | `--data`（必填） |
| `wlt produce main export` | 导出生产任务 Excel | 同 `main page` |
| `wlt produce main quality-page` | 分页查询生产质检 | `--page-no`, `--page-size`, `--no`, `--produce-time`, `--product-name`, `--metrics-name`, `--inspection-nos-str`, `--inspection-result`, `--warehouse-name`, `--plan-name`, `--user-name`, `--project-name`, `--plan-no`, `--order-no`, `--batch-no`, `--remark`, `--status`, `--creator`, `--create-time`, `--updater`, `--update-time`, `--keyword`, `--headers` |
| `wlt produce main quality-count` | 按筛选统计生产质检数量 | 同 `quality-page`（去 `--headers`） |
| `wlt produce main quality-export` | 导出生产质检 Excel | 同 `quality-page` |

## 生产方案 (`wlt produce plan`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt produce plan page` | 分页查询生产方案 | `--page-no`, `--page-size`, `--plan-no`, `--plan-name`, `--status`, `--product-id`, `--product-name`, `--metrics-name`, `--plan-time`, `--customer-id`, `--customer-name`, `--remark`, `--creator`, `--creator-name`, `--create-time`, `--update-time`, `--updater-name`, `--custom-order`, `--keyword`, `--headers` |
| `wlt produce plan page-count` | 按筛选统计生产方案数量 | 同 `plan page`（去 `--headers`） |
| `wlt produce plan simple-list` | 方案精简列表 | `--warehouse-id` |
| `wlt produce plan get --id <N>` | 获取方案详情 | `--id`（必填） |
| `wlt produce plan create --data '<json>'` | 创建方案 | `--data`（必填） |
| `wlt produce plan update --data '<json>'` | 更新方案 | `--data`（必填） |
| `wlt produce plan delete --id <N>` | 删除方案 | `--id`（必填） |
| `wlt produce plan update-status --data '<json>'` | 更新方案状态 | `--data`（必填） |
| `wlt produce plan export` | 导出生产方案 Excel | 同 `plan page` |
