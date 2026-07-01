# 合同管理 (contract)

## 总览

合同模块覆盖 4 大类、7 个子类合同，统一通过 `type` 字段在 3 个后端端点之间区分：

| 子命令 | 后端 type | HTTP 路径 |
|---|---|---|
| `purchase-long-cooperate` | `PURCHASE_LONG_COOPERATE` | `/erp/contract` |
| `sale-contract` | `SALE_CONTRACT` | `/erp/service-contract` |
| `sale-long-cooperate` | `SALE_LONG_COOPERATE` | `/erp/contract` |
| `transport` | `TRANSPORT` | `/erp/transport-contract` |
| `transport-long` | `TRANSPORT_LONG` | `/erp/contract` |
| `service-contract` | `SERVICE` | `/erp/provision-contract` |
| `service-long` | `SERVICE_LONG` | `/erp/provision-contract` |

每个子类合同的标准子命令（7 个）：`list` / `page-count` / `get` / `create` / `update` / `delete` / `update-status`。

## 共用的筛选与分页标志

所有 list / page-count 命令均支持以下筛选参数：

| 标志 | 说明 |
|---|---|
| `--keyword` | 关键字搜索（合同编号等） |
| `--enterprise-id` | 企业 ID |
| `--order-start` | 下单日期起始（如 `2026-07-06 00:00:00`） |
| `--order-end` | 下单日期结束（如 `2026-08-03 23:59:59`） |
| `--end-start` | 合同到期起始（如 `2026-07-13 00:00:00`） |
| `--end-end` | 合同到期结束（如 `2026-08-03 23:59:59`） |
| `--page-no` | 页码（默认 1） |
| `--page-size` | 每页数量（默认 20） |

后端会把 `order-start/order-end/end-start/end-end` 转成 `orderDate[0]/[1]` 与 `endTime[0]/[1]` 数组参数。

## 采购长协

```bash
# 分页查询
wlt contract purchase-long-cooperate list \
  --keyword XY20260403000001 \
  --enterprise-id 2001489305039032322 \
  --order-start "2026-07-07 00:00:00" --order-end "2026-08-13 23:59:59" \
  --end-start "2026-07-13 00:00:00" --end-end "2026-08-13 23:59:59" \
  --page-no 1 --page-size 10

# 统计数量
wlt contract purchase-long-cooperate page-count --keyword XY20260403000001

# CRUD
wlt contract purchase-long-cooperate get --id <N>
wlt contract purchase-long-cooperate create --data '<json>'
wlt contract purchase-long-cooperate update --data '<json>'
wlt contract purchase-long-cooperate delete --ids <id1,id2>
wlt contract purchase-long-cooperate update-status --data '{"id":1,"status":2}'
```

## 销售合同 / 销售长协

```bash
# 销售合同
wlt contract sale-contract list --keyword HT20260401000001 --order-start "..." --order-end "..."
wlt contract sale-contract get --id <N>

# 销售长协
wlt contract sale-long-cooperate list --keyword XY20260401000001 --order-start "..." --order-end "..."
wlt contract sale-long-cooperate page-count --keyword XY20260401000001
wlt contract sale-long-cooperate create --data '<json>'
wlt contract sale-long-cooperate delete --ids <id1,id2>
wlt contract sale-long-cooperate update-status --data '{"id":1,"status":2}'
```

## 运输合同 / 运输长协

```bash
# 运输合同
wlt contract transport list --keyword HT20260319000001 --order-start "..." --order-end "..."
wlt contract transport get --id <N>
wlt contract transport page-count --keyword HT20260319000001

# 运输长协
wlt contract transport-long list --keyword XY20260515000002 --order-start "..." --order-end "..."
wlt contract transport-long page-count --keyword XY20260515000002
```

## 服务合同 / 服务长协

```bash
# 服务合同（后端端点同 /erp/provision-contract?type=SERVICE）
wlt contract service-contract list --keyword HT20260402000001 --order-start "..." --order-end "..."
wlt contract service-contract get --id <N>

# 服务长协
wlt contract service-long list --keyword XY20260409000001 --order-start "..." --order-end "..."
wlt contract service-long page-count --keyword XY20260409000001
```
