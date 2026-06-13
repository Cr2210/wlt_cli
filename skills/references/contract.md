# 合同管理 (contract)

## 长期合同 (`wlt contract`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt contract list` | 分页查询合同 | `--name`, `--status`, `--code`, `--partner-id`, `--page-no`, `--page-size` |
| `wlt contract get --id <N>` | 获取合同详情 | `--id`（必填）, `--code` |
| `wlt contract create --data '<json>'` | 创建合同 | `--data`（必填） |
| `wlt contract update --data '<json>'` | 更新合同 | `--data`（必填） |
| `wlt contract delete --ids <id1,id2>` | 删除合同 | `--ids`（必填，逗号分隔） |
| `wlt contract update-status` | 更新合同状态 | `--id`（必填）, `--status`（必填） |
| `wlt contract page-count` | 统计合同数量 | `--name`, `--status`, `--code`, `--partner-id` |

## 供货合同 (`wlt contract provision`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt contract provision list` | 分页查询 | `--name`, `--status`, `--no`, `--partner-id`, `--page-no`, `--page-size` |
| `wlt contract provision get --id <N>` | 获取详情 | `--id`（必填）, `--no` |
| `wlt contract provision create --data '<json>'` | 创建 | `--data`（必填） |
| `wlt contract provision update --data '<json>'` | 更新 | `--data`（必填） |
| `wlt contract provision delete --id <N>` | 删除 | `--id`（必填） |
| `wlt contract provision update-status` | 更新状态 | `--id`（必填）, `--status`（必填） |
| `wlt contract provision page-count` | 统计数量 | `--name`, `--status`, `--no`, `--partner-id` |
| `wlt contract provision delete-batch` | 批量删除 | `--ids`（必填，逗号分隔） |
| `wlt contract provision from-long` | 从长期合同生成 | `--long-contract-id`（必填） |

## 服务合同 (`wlt contract service`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt contract service list` | 分页查询 | `--name`, `--status`, `--code`, `--partner-id`, `--page-no`, `--page-size` |
| `wlt contract service get --id <N>` | 获取详情 | `--id`（必填）, `--code` |
| `wlt contract service create --data '<json>'` | 创建 | `--data`（必填） |
| `wlt contract service update --data '<json>'` | 更新 | `--data`（必填） |
| `wlt contract service delete --ids <id1,id2>` | 删除 | `--ids`（必填，逗号分隔） |
| `wlt contract service update-status` | 更新状态 | `--id`（必填）, `--status`（必填） |
| `wlt contract service page-count` | 统计数量 | `--name`, `--status`, `--code`, `--partner-id` |

## 运输合同 (`wlt contract transport`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt contract transport list` | 分页查询 | `--name`, `--status`, `--code`, `--partner-id`, `--page-no`, `--page-size` |
| `wlt contract transport get --id <N>` | 获取详情 | `--id`（必填） |
| `wlt contract transport create --data '<json>'` | 创建 | `--data`（必填） |
| `wlt contract transport update --data '<json>'` | 更新 | `--data`（必填） |
| `wlt contract transport delete --ids <id1,id2>` | 删除 | `--ids`（必填，逗号分隔） |
| `wlt contract transport update-status` | 更新状态 | `--id`（必填）, `--status`（必填） |
| `wlt contract transport page-count` | 统计数量 | `--name`, `--status`, `--code`, `--partner-id` |
