# 运单管理 (waybill)

> 后端运单服务现仅一组路径 `/erp/waybill/*`。`get` / `page` / `page-count` 是只读查询命令,顶层平铺。`load` 是装货写入命令（危险写操作,需要 token 鉴权）。

## 命令一览 (`wlt waybill *`)

| 命令 | 后端 | 说明 | 关键参数 |
|------|------|------|---------|
| `wlt waybill get` | `GET /erp/waybill/get` | 运单详情（`WaybillSourceDetailVO`,含 metricsList/关联订单/合同/图片 URL/结算标志） | `--id`（必填） |
| `wlt waybill page` | `GET /erp/waybill/page` | 分页查询（18 个筛选字段） | 见本页筛选字段表 + `--page-no` `--page-size` |
| `wlt waybill page-count` | `GET /erp/waybill/page-count` | 统计数量（同 page 的筛选字段,无分页参数） | 同 page 的筛选字段 |
| `wlt waybill load` | `POST /erp/waybill/load` | 运单装货（写操作,要求 UN_LOAD） | `--waybill-id`（必填）、`--load-time`（必填）、`--load-weight`（选填）、`--loaded-img-url`（选填） |
| `wlt waybill unload` | `POST /erp/waybill/unload` | 运单卸货（写操作,要求 ON_LOAD） | `--waybill-id`（必填）、`--unload-time`（必填）、`--unload-weight`（必填）、`--unloaded-img-url`（选填） |
| `wlt waybill sign-batch` | `POST /erp/waybill/sign-batch` | 批量运单签收（写操作） | `--waybill-id`（必填,逗号分隔）、`--sign-time`（必填）;或 `--data '[...]'` 透传 JSON 数组 |
| `wlt waybill push-config get` | — | 推送配置 | 无 |
| `wlt waybill push-config update` | — | 更新推送配置 | `--data`（必填） |
| `wlt waybill push-config generate-secret-key` | — | 生成推送密钥 | 无 |

> ⚠️ `waybill load` 是写入操作（类似 `stock in update-status`、`sale out update-status`），会变更运单状态。
>
> **业务规则:只有 `UN_LOAD`(待发货)状态的运单才允许装货**。`load` 命令在发 POST 前会自动调用 `GET /erp/waybill/get` 自检:状态 OK 放行;状态不符(如 `ON_LOAD`/`DELIVERED`/`SIGNED_FOR` 等)直接报错退出(退出码 4)。可用 `--force` 跳过自检,慎用。

## 运单状态枚举

| 状态值 | 含义 | 可否装货 |
|--------|------|---------|
| `FINISHED` | 已完成 | ✗ |
| `UNDEFINED` | 未定义 | ✗ |
| `UN_LOAD` | 待发货 | ✅ 唯一可进入装货的状态 |
| `ON_LOAD` | 运输中 | ✅ 唯一可进入卸货的状态 |
| `DELIVERED` | 已送达 | ✗ |
| `SIGNED_FOR` | 已签收 | ✗ |

## 筛选字段（18 个,list / page-count 共用）

| 字段 | CLI flag | 说明 | 示例值 |
|------|----------|------|-------|
| `waybillNo` | `--waybill-no` | 运单号 | `YD20260701000013` |
| `carNumber` | `--car-number` | 车牌号 | `皖A12345` |
| `orderType` | `--order-type` | 订单类型 | `SALE_OUT` / `PURCHASE_IN` |
| `addressName` | `--address-name` | 起始地/目的地（模糊） | `起始地、目的地` |
| `status` | `--status` | 状态 | （运单生命周期值） |
| `mediumName` | `--medium-name` | 货物名称（模糊） | `货物名称` |
| `metricsName` | `--metrics-name` | 规格指标 | `规格指标` |
| `realLoadDate` | `--real-load-date-start` / `--real-load-date-end` | 实际装车时间范围 | `2026-06-30 00:00:00` |
| `realUnloadDate` | `--real-unload-date-start` / `--real-unload-date-end` | 实际卸车时间范围 | `2026-07-01 00:00:00` |
| `capacityName` | `--capacity-name` | 承运商 | `承运商` |
| `userName` | `--user-name` | 业务负责人 | `业务负责人` |
| `projectName` | `--project-name` | 项目名称 | `项目名称` |
| `createTime` | `--create-time-start` / `--create-time-end` | 创建时间范围 | `2026-06-30 00:00:00` |
| `inputType` | `--input-type` | 录入方式 | `ACQUIRE` / `MANUAL` |
| `dataSource` | `--data-source` | 数据来源 | `WLCLOUDS` |
| `outWaybillNo` | `--out-waybill-no` | 外部数据编号 | `外部数据编号` |

> 日期范围 flag: `--{range}-start` / `--{range}-end`，range 取值 `real-load-date`、`real-unload-date`、`create-time`。后端以 `key[0]`/`key[1]` 数组形式接收 (`foo[0]=...&foo[1]=...`)。时间格式统一 `YYYY-MM-DD HH:mm:ss`。

## 调用样例

```bash
# 详情（按 ID）
wlt waybill get --id 17386 --token <t> --tenant-id <T>

# 按运单号查
wlt waybill page --waybill-no YD20260701000013 --token <t> --tenant-id <T>

# 按订单类型 + 装车时间
wlt waybill page --order-type SALE_OUT \
  --real-load-date-start "2026-06-30 00:00:00" --real-load-date-end "2026-08-04 23:59:59" \
  --token <t> --tenant-id <T>

# 多字段联合 + 分页
wlt waybill page --medium-name 货物名称 --order-type SALE_OUT \
  --create-time-start "2026-06-30 00:00:00" --create-time-end "2026-07-29 23:59:59" \
  --page-no 1 --page-size 10 --token <t> --tenant-id <T>

# 计数（同 page 的筛选字段，不传分页参数）
wlt waybill page-count --order-type SALE_OUT --token <t> --tenant-id <T>

# 装货（写入操作,必填 waybill-id + load-time；负载可选）
wlt waybill load --waybill-id 1001 \
  --load-time "2026-07-02 08:00:00" \
  --load-weight 25.6 \
  --loaded-img-url "https://oss.wlclouds.com/images/load_20250405.jpg" \
  --token <t> --tenant-id <T>

# 也可用 --data 透传 JSON（互斥于上方 flag）
wlt waybill load --data '{"waybillId":1001,"loadTime":"2026-07-02 08:00:00"}' --token <t> --tenant-id <T>

# 卸货（写入操作,要求 ON_LOAD 运输中；必填 waybill-id + unload-time + unload-weight）
wlt waybill unload --waybill-id 1001 \
  --unload-time "2026-07-02 18:00:00" \
  --unload-weight 25.4 \
  --unloaded-img-url "https://oss.wlclouds.com/images/unload_20250405.jpg" \
  --token <t> --tenant-id <T>

# 也可用 --data 透传 JSON
wlt waybill unload --data '{"waybillId":1001,"unloadTime":"2026-07-02 18:00:00","unloadWeight":25.4}' --token <t> --tenant-id <T>

# 批量签收(单条便捷:waybill-id + sign-time)
wlt waybill sign-batch --waybill-id 1001 --sign-time "2026-07-02 20:00:00" --token <t> --tenant-id <T>

# 批量签收(多条,逗号分隔 ID;同一签收时间)
wlt waybill sign-batch --waybill-id 1001,1002,1003 --sign-time "2026-07-02 20:00:00" --token <t> --tenant-id <T>

# 批量签收(--data 透传 JSON 数组,每条可独立 signTime)
wlt waybill sign-batch --data '[{"waybillId":1001,"signTime":"2026-07-02 20:00:00"},{"waybillId":1002,"signTime":"2026-07-02 20:05:00"}]' --token <t> --tenant-id <T>
```

## 运单生命周期

运单状态流转：创建 → 装车 → 发车 → 卸车 → 签收。可用状态值: `UNDEFINED` / `ON_LOAD` / `UN_LOAD` / `DELIVERED` / `SIGNED_FOR` / `FINISHED`（详见 `get` 响应 `status` 字段）。
