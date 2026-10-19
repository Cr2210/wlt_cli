# A 类 JSON 查询端点盘点（后端 ge-dev-1.8 ↔ CLI 现状）

> `$C` = 后端 `wlclouds-module-erp/src/main/java/com/wlclouds/module/erp/controller/admin/`。
> 结论三态：**待做**（后端有、CLI 无）/ **已完成**（CLI 已实现，至多待实测）/ **漂移**（CLI 调的端点后端已不存在）。

## 盘点总表

| # | 模块 | 端点 | 后端锚点 | CLI 现状 | 结论 |
|---|------|------|---------|---------|------|
| 1 | invoice-item | `/erp/invoice-item/list` | `$C/invoice/InvoiceItemController.java:31`（仅此 1 个方法） | `cmd/` 零命中 | **待做**（单端点，建议挂 `invoice` 域作子命令 `item list`，或独立小组） |
| 2 | plan-prepayment-relation | `/available-initials` `/available-payments` `/page` | `$C/finance/ErpOrderPrepaymentRelationController.java`（`@RequestMapping("/erp/order-prepayment-relation")`） | 已接入 `wlt order prepayment-relation`（2026-08-21，另含 create/update-amount） | **已完成** |
| 3 | position | `/erp/position/get-latest-waybill-follow` | `$C/position/ErpPositionController.java:30`「运单车辆最新轨迹」 | `cmd/` 零命中（无 position 目录） | **待做**（新模块，2 端点） |
| 4 | position | `/erp/position/get-history-waybill-follow` | `$C/position/ErpPositionController.java:59`「运单车辆历史轨迹」 | 同上 | **待做** |
| 5-10 | screen | `/erp/screen/line` `/line-price` `/purchase-enterprise` `/purchase-in-price` `/purchase-in-product` `/settle-amount` | `$C/screen/ErpScreenController.java:96 :90 :103 :70 :76 :83` | `cmd/screen/screen.go` 仅有 purchase-sale/stock-count/amount-used/project-count（`:39 :53 :73 :93`） | **待做 6 个**（后端 10 端点全存在，CLI 已做 4） |
| 11 | system/account | `/list` | `$C/system/ErpSystemAccountController.java:47`「关联账号列表」 | `cmd/` 中 `system/account`、`weight-account` 零命中 | **待做**（新模块） |
| 12 | system/account | `/getCarDetail` `/getCarrierDetail` `/getDriverDetail` `/getEnterpriseInfo` / `/getLatestWaybillFollow` `/getWaybill` `/getWaybillFollow` `/getWaybillSourceInfo` | 同文件 `:92 :98 :86 :34 :73 :80 :67 :61` | 零命中 | **待做 8 个**（prd 所列 9 个全部核实存在） |
| 13 | system/account（prd 未列） | `/sendSms` `/create` `/delete` | 同文件 `:28 :40 :54` | 零命中 | **范围决策**：prd 只列查询，写/短信命令是否纳入 |
| 14 | system/weight-account | `/erp/system/weight-account/get` | `$C/system/ErpSystemWeightAccountController.java:24`「获取已绑定的信息」 | 零命中 | **待做** |
| 15 | system/weight-account（prd 未列） | `/save` | 同文件 `:31`「新建或更新已绑定」 | 零命中 | **范围决策**（写操作） |
| 16 | waybill | `/erp/waybill/get` | `$C/waybill/WaybillSourceController.java:92` | `cmd/waybill/waybill_get.go:22` | **已完成** |
| 17 | waybill | `/erp/waybill/page` `/page-count` | 同文件 `:132 :126` | `cmd/waybill/waybill_page.go:101 :121` | **已完成** |
| 18 | waybill | `/erp/waybill/get-events` | 同文件 `:98`「运单节点事件列表」 | `cmd/waybill/` 无此命令 | **待做** |
| 19 | waybill | `/erp/waybill/get-order-info` | 同文件 `:104`「运单关联的订单信息」 | 同上 | **待做** |
| 20 | service-contract | `/page` `/page-count` | `$C/contract/ErpServiceContractController.java:61 :67` | `cmd/contract/contract_main.go:24-25` 已配 `APIPath:"/erp/service-contract"` | **已完成**（待实测；详见 02 号文件） |
| 21 | transport-contract | `/page` `/page-count` | `$C/contract/ErpTransportContractController.java:65 :71` | `cmd/contract/contract_main.go:39-40` 已配 `APIPath:"/erp/transport-contract"` | **已完成**（待实测） |
| 22 | finance-record | `/erp/finance-record/init` | `$C/finance/ErpFinanceRecordController.java:25-28`：**无参、返回 `CommonResult<Boolean>`** | 零命中 | **移出 A 类**：实为一次性初始化写操作，非查询；若要暴露建议按「维护命令」单独处理 |

## waybill 模块全景（顺带盘点，供 design 参考）

后端 `WaybillSourceController`（`@RequestMapping("/erp/waybill")`）共 24 个方法；CLI 已做：
get / page / page-count / load（`waybill_load.go`）/ sign-batch（`waybill_sign_batch.go`）/ push-config 组。
未做：get-events、get-order-info（本轮 A 类）、create `:78`、update `:85`、delete `:110`、delete-list `:118`、
unload `:146`、sign `:153`、batch-sign `:176`、batch-unlink-stock `:169`、calc-settlement `:71`、
get-settlement `:64`、update-settlement `:57`、export-excel `:183`、import `:197`、import-update-freight `:324`。
（这些不属于本轮 A 类范围，仅记录。）

## 端点漂移（新发现）

| CLI 命令 | 调用端点 | 后端现状 |
|---|---|---|
| `finance invoice page-count`（`cmd/finance/finance_invoice.go:108-119`） | `/erp/finance-invoice/page-count` | **不存在**。`FinanceInvoiceController.java` 方法集仅 create/update/update-status/delete/page/summary/get/export-excel（`:33-78`），全后端 grep 无 `finance-invoice/page-count`。6 月 erp.md 有此端点、后端其后删除。处理建议：实测确认 404 后删除该子命令（或留待与待决策 2 一并处理） |

## 汇总

- **已完成**（无需开发，至多待 token 实测）：#2、#16、#17、#20、#21 ≈ 9 个端点
- **待做**（后端有、CLI 无）：#1、#3-#12、#14、#18、#19 ≈ 19 个端点（position 2 + screen 6 + system/account 9 + weight-account 1 + waybill 2 + invoice-item 1，其中部分依赖范围决策）
- **移出/降级**：#22 finance-record/init（写操作）
- **漂移**：finance-invoice/page-count 1 处
