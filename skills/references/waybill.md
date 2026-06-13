# 运单管理 (waybill)

## 运单源 (`wlt waybill source`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt waybill source list` | 分页查询运单 | `--waybill-no`, `--license-plate`, `--status`, `--warehouse-id`, `--customer-id`, `--supplier-id`, `--page-no`, `--page-size` |
| `wlt waybill source page-count` | 统计运单数量 | `--waybill-no`, `--license-plate`, `--status`, `--warehouse-id`, `--customer-id`, `--supplier-id` |
| `wlt waybill source get --id <N>` | 获取运单详情 | `--id`（必填） |
| `wlt waybill source get-events --id <N>` | 获取运单事件 | `--id`（必填） |
| `wlt waybill source get-order-info --id <N>` | 获取订单信息 | `--id`（必填） |
| `wlt waybill source create --data '<json>'` | 创建运单 | `--data`（必填） |
| `wlt waybill source update --data '<json>'` | 更新运单 | `--data`（必填） |
| `wlt waybill source delete --id <N>` | 删除运单 | `--id`（必填） |
| `wlt waybill source delete-list` | 批量删除运单 | `--ids`（必填，逗号分隔） |
| `wlt waybill source load --data '<json>'` | 装车 | `--data`（必填） |
| `wlt waybill source unload --data '<json>'` | 卸车 | `--data`（必填） |
| `wlt waybill source sign --data '<json>'` | 签收 | `--data`（必填） |
| `wlt waybill source sign-batch --data '<json>'` | 批量签收（数据） | `--data`（必填） |
| `wlt waybill source batch-sign --ids <id1,id2>` | 批量签收（ID列表） | `--ids`（必填，逗号分隔） |
| `wlt waybill source import --file <path>` | 导入 Excel | `--file`（必填，Excel 文件路径） |
| `wlt waybill source export` | 导出 Excel | `--waybill-no`, `--license-plate`, `--status`, `--warehouse-id`, `--customer-id`, `--supplier-id` |

## 推送配置 (`wlt waybill push-config`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt waybill push-config get` | 获取推送配置 | 无 |
| `wlt waybill push-config update --data '<json>'` | 更新推送配置 | `--data`（必填） |
| `wlt waybill push-config generate-secret-key` | 生成密钥 | 无 |

## 运单生命周期

运单状态流转：创建 → 装车(load) → 发车 → 卸车(unload) → 签收(sign)

```bash
# 创建运单
wlt waybill source create --data '{"waybillNo":"WB-001",...}'
# 装车
wlt waybill source load --data '{"id":1,...}'
# 卸车
wlt waybill source unload --data '{"id":1,...}'
# 签收
wlt waybill source sign --data '{"id":1,...}'
# 或批量签收
wlt waybill source batch-sign --ids 1,2,3
```
