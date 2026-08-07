# wlt 库存命令优化设计

日期:2026-08-07 ｜ 状态:已确认,待实现

## 背景

库存相关命令分布在 `cmd/stock/`、`cmd/purchase/`、`cmd/sale/` 三个包,实现风格三套并存:

- **采购入库 / 销售出库**:用现代工厂 `cmdutil.SalePurchaseCmds`,每个文件仅 ~18 行,自动生成 `list/page-count/get/create/update/delete/update-status`。
- **采购退货 / 销售退货**:用 legacy 工厂 `cmdutil.NewCRUDSubCmd`(全仓仅剩这两个调用方),list 的筛选 flag 是硬编码的 7 个通用字段,无法表达退货单业务字段。
- **stock 的 8 个文件(warehouse/query/record/in/out/move/check)全部完全手写**,无一使用任何工厂。其中 `stock in/out/move/check` 四个单据命令结构和代码几乎逐行相同(都是同样 7 个子命令 + 各自 20 个左右 flag),每个约 290 行,合计 ~1160 行重复样板。

由此衍生若干真实问题:

1. **命名割裂**:工厂生成的分页子命令叫 `list`,stock 手写的叫 `page` —— `wlt purchase in list` 但 `wlt stock in page`。
2. **`--headers` 自定义导出缺失**:只有 stock 手写命令支持,采购入库/销售出库(工厂生成)没有。
3. **分页解析不健壮**:`stock in/out/move/check/warehouse`、`stock query ledger` 还用手写 `int64 total` 解析(遇字符串型 `total` 报错),而同包的 `query page`、`record page` 已改用健壮的 `cmdutil.ParsePagedJSON`。
4. **flag 列表重复**:stock in/out/move/check 的 page 与 page-count 各抄一份 20+ flag(收集列表 + 注册列表);`stock query` 的 page/page-count/ledger/ledger-count 共享几乎相同的 10 字段却复制 4 份;`stock record` 同理。
5. **文档失真**:`CLAUDE.md` 称 stock 用 legacy 工厂,实际早全手写。

## 目标

- 把全仓单据命令(采购入库/销售出库/其他入库/其他出库/调拨/盘点/退货)归一到**一套通用工厂**,消除"工厂 / legacy / 手写"三套风格。
- 库存域命令命名、筛选 flag、输出行为对外一致。
- 削减 ~1000 行重复样板,删除 legacy 工厂。
- 修正失真文档。

## 非目标

- 不改变后端 API 契约(仅调整 CLI 层 flag 与传参映射)。
- 不重构 `cmd/report/`、`cmd/stats/` 下的库存报表/统计命令(它们属于报表模块,不在本次 8 项库存操作范围)。
- 不改变鉴权、输出协议(stdout JSON / stderr 错误 / 退出码 0-6)。

## 设计

### 1. 通用单据工厂 `DocumentConfig`

把 `SalePurchaseConfig` 重命名为 `DocumentConfig`、`SalePurchaseCmds` → `DocumentCmds`(采购/销售仅 2 处调用点机械改名),并扩展能力:

```go
// DocumentConfig 描述一个单据 CRUD 子模块(采购入库/销售出库/其他入库/出库/调拨/盘点/退货)。
type DocumentConfig struct {
    Name      string     // 子命令名:in / out / move / check / return …
    APIPath   string     // 后端根路径:/erp/stock-in
    Label     string     // 中文标签:入库单
    Filters   []FlagSpec // 模块特有筛选 flag(不含分页/时间/headers 这些通用项)
    TimeKey   string     // 时间字段:空=无时间筛选;非空=注册 --start-time/--end-time,传 {TimeKey}[0]/{TimeKey}[1] 数组
    Headers   bool       // 是否注册 --headers 自定义导出表头
    PageAlias bool       // 是否给 list 注册 page 别名(向后兼容)
}

// DocumentCmds 为一组单据子模块生成标准 CRUD 子命令,统一注册到 parent。
// 每个子模块自动拥有:list(+可选 page 别名) / page-count / get / create / update / delete / update-status。
func DocumentCmds(parent *cobra.Command, cfgs ...DocumentConfig)
```

**list 命令的 flag 集合** = `--page-no`/`--page-size` + `Filters` + (TimeKey 非空时)`--start-time`/`--end-time` + (Headers 时)`--headers`。**无隐式 base flags** —— 所有筛选字段(含原采购/销售的 `warehouse-id`/`product-id`/`product-name`/`no`/`status`/`type`,它们原先由 `addSalePurchaseBaseFlags` 隐式注入)一律通过 `Filters` 显式配置。
**list 传参**:`pageNo`/`pageSize` + `Filters` 收集(camelCase) + 时间范围收集(转 `{TimeKey}[0]/{TimeKey}[1]`) + `headers` 收集。
**list 输出**:统一 `cmdutil.ParsePagedJSON`(健壮的 `json.Number` total)。

> 实现时把 `collectSalePurchaseTime`/`addSalePurchaseTimeFlags` 重命名为通用名(`collectDocumentTimeRange`/`addDocumentTimeRangeFlags`),并删除已无用的 `addSalePurchaseBaseFlags`(其 6 个字段并入各单据 `Filters`)。

`page-count` 共享 list 的全部筛选 flag(含时间、headers 与否同 list),只是去掉 `page-no`/`page-size`,调用 `{APIPath}/page-count`,输出 `OutputJSON`。

`get/create/update/delete/update-status` 直接复用现有 `CrudGetCmd`/`CrudCreateCmd`/`CrudUpdateCmd`/`CrudDeleteCmd(apiPath, label, false)`/`CrudUpdateStatusCmd`(stock 单据均用 `ids` 批量删除,故 `single=false`)。

> 已确认:stock in/out/move/check 的后端 `page` 接口时间参数认**数组范围**(`inTime[0]/inTime[1]` 等,同采购入库),故 `TimeKey` 非空即数组模式已覆盖全部场景,无需单值模式(YAGNI;若未来出现单值时间端点再扩展)。

### 2. 命令归并映射

| 命令 | 现状 | 归并后配置 | 预计行数 |
|---|---|---|---|
| `purchase in` | `SalePurchaseCmds` | `DocumentCmds`{APIPath:/erp/purchase-in, TimeKey:inTime, Headers, Filters:[warehouse-id,product-id,product-name,no,status,type,supplier-id,metrics-name]} | ~18 |
| `sale out` | `SalePurchaseCmds` | `DocumentCmds`{APIPath:/erp/sale-out, TimeKey:outTime, Headers, Filters:[warehouse-id,product-id,product-name,no,status,type,customer-id,batch-no]} | ~18 |
| `stock in` | 手写 ~290 行 | `DocumentCmds`{APIPath:/erp/stock-in, TimeKey:inTime, Headers, PageAlias, Filters:[全部业务筛选字段,见下]} | ~20 |
| `stock out` | 手写 ~290 行 | `DocumentCmds`{APIPath:/erp/stock-out, TimeKey:outTime, Headers, PageAlias, Filters:[与 in 对称,supplier→customer]} | ~20 |
| `stock move` | 手写 ~290 行 | `DocumentCmds`{APIPath:/erp/stock-move, TimeKey:moveTime, Headers, PageAlias, Filters:[含 from/to-warehouse-id、updater 等]} | ~20 |
| `stock check` | 手写 ~290 行 | `DocumentCmds`{APIPath:/erp/stock-check, TimeKey:checkTime, Headers, PageAlias, Filters:[含 warehouse-id 等]} | ~20 |
| `stock warehouse` | 手写 ~207 行 | `AddCRUDToParent` + `CrudSimpleListCmd`(仓库有 simple-list,无 page-count) | ~15 |
| `purchase return` | `NewCRUDSubCmd` | `DocumentCmds`{APIPath:/erp/purchase-return, TimeKey:"", Filters:[实现时确认退货业务字段]} | ~10 |
| `sale return` | `NewCRUDSubCmd` | `DocumentCmds`{APIPath:/erp/sale-return, TimeKey:"", Filters:[实现时确认]} | ~10 |

stock in/out/move/check 的 `PageAlias=true`(保留 `page` 兼容旧调用);采购/销售/退货 `PageAlias=false`(它们本来就是 `list`)。

stock in 的 `Filters` 含 21 个业务筛选字段(no/supplier-id/supplier-name/status/remark/creator/product-id/product-name/warehouse-id/warehouse-name/metrics-name/creator-name/user-id/receive-address/send-address/batch-no/create-time/updater-name/update-time/custom-order/keyword;`in-time` 抽出作 `TimeKey:inTime`)声明一次,工厂自动在 list 与 page-count 注册并收集 —— 消除当前"收集列表 + 注册列表复制两份"的重复。stock out/move/check 同理(out 把 supplier→customer、`in-time`→`out-time`;move 含 `from-warehouse-id`/`to-warehouse-id`/`updater`;check 含 `warehouse-id`,各自字段实现时按现状逐字保留)。

### 3. 命名统一 page → list

stock 的分页子命令由 `page` 改名为 `list`(与采购/销售一致)。`page` 通过 cobra `Aliases` 保留为兼容别名(`wlt stock in page` 仍可用,不破坏旧脚本/AI 记忆)。仅 stock 单据设 `PageAlias=true`。

### 4. `--headers` 导出补齐

`purchase in` / `sale out` 现缺失的 `--headers`,通过 `Headers:true` 补齐,与 stock 单据对齐。

### 5. stock query / record 重构(stock 包内)

query(7 子命令)/ record(6 子命令)含非 CRUD 业务端点,不套 `DocumentCmds`,在 stock 包内消除 flag 重复:

- **共享 FlagSpec**(声明一次,多处复用):
  - `queryPageFilters`:page/page-count/ledger/ledger-count 共享(product-id/warehouse-id/metrics-name/batch-no/send-address/plan-no/supplier-id/supplier-name/is-detail/count-positive/keyword)。
  - `recordPageFilters`:page/page-count 共享(product-id/category-id/warehouse-id/biz-type/biz-no/create-time/in-time/metrics-name/product-name/batch-no/keyword)。
  - `recordDimensionFilters`:record-page/total-cost 共享(product-id/warehouse-id/metrics-name/batch-no/supplier-id/source-supplier-id/out-count)。
- **本地辅助函数**:
  - `newStockPagedCmd(use, apiPath, endpoint, short, filters, withHeaders) *cobra.Command` —— 可配端点(`page`/`batch-detail`/`record-page`)的分页命令,输出 `ParsePagedJSON`。
  - `newStockCountCmd(apiPath, endpoint, short, filters) *cobra.Command` —— 可配端点(`page-count`/`detail-count`/`total-cost`)的计数命令,输出 `OutputJSON`。
- **保留手写**(真实业务逻辑):`query get`(id/product-id/warehouse-id 三选一)、`query count`(合并 get-count + get-warehouse-count)、`query stock-record-count`(维度键)、`record get`、`record count`(无参)。
- query/record 的分页命令同样加 `page` 别名保持兼容。
- query/record 的时间相关 flag(`in-time`/`create-time`)保持现有单值语义不动(它们走 `/erp/stock-record` 端点,与 stock-in 的数组范围是不同契约;本次零风险,不改其传参)。

### 6. 附带清理

- 手写的 `int64 total` 分页解析(stock in/out/move/check/warehouse、query ledger)全部替换为 `cmdutil.ParsePagedJSON`。
- 退货迁离后删除 `internal/cmdutil/crud_legacy.go`(224 行,全仓再无调用方)。
- 修正 `CLAUDE.md` 关于"stock 用 legacy 工厂 / `NewCRUDSubCmd` 库存子域使用"的失真描述。
- 更新 `skills/` 下相关 Markdown:`page→list` 别名说明、`purchase in`/`sale out` 新增 `--headers`。

## 关键决策记录

1. **方案 A(扩展为通用工厂)** 而非新建库存专用工厂(B)或最小重构(C):唯一能达成"全仓一套工厂、消除三套风格"目标;一次性投入、长期收益最大。
2. **时间字段=数组范围**:已确认 stock in/out/move/check 后端认 `inTime[0]/inTime[1]` 等,与采购入库一致,故 flag 统一为 `--start-time/--end-time`、`TimeKey` 非空即数组模式,无需单值模式。
3. **`page` 保留为别名**:命名统一为 `list`,但用 cobra `Aliases` 保留 `page` 兼容,避免破坏现有 AI 调用记忆/脚本。
4. **query/record 不套工厂**:其 count/ledger/record-page/total-cost 等非 CRUD 端点有真实业务逻辑,套工厂会过度抽象;改为 FlagSpec 复用 + 本地辅助函数,消除 flag 重复即可。

## 测试策略

- **回归(高优先级)**:工厂改造触及采购/销售,验证 `purchase in` / `sale out` 的 `list`/`get`/`create`/`update`/`delete`/`update-status` 行为与改造前一致(flag 集合、传参、输出)。
- **等价性**:stock in/out/move/check 用工厂后,筛选 flag、API 路径、分页输出与原手写实现等价(差异仅:`page` 子命令改名 `list` 并保留别名、新增 `--headers`)。
- **新增单测**:为 `DocumentCmds` 的 list 传参构造(TimeKey 数组映射、Headers、Filters 收集)写表格化单元测试。
- `make test` 全部通过;`make build` 成功。

## 迁移与文档

- `skills/` 文档与 CHANGELOG 标注:`stock in/out/move/check/query/record` 的分页子命令 `page` 现以 `list` 为主名(`page` 仍可用作别名);`purchase in`/`sale out` 新增 `--headers`。
- 本次为 CLI 行为小幅演进(命名 + headers),无配置迁移。

## 风险

- **工厂改造波及采购/销售**:通过回归测试 + `page`/`page-count` 等价性比对控制。
- **退货业务字段不确定**:当前 legacy 给退货的 flag 是硬编码通用字段,迁移到 `DocumentCmds` 时需确认退货实际筛选字段(实现时核对后端/前端);若无法确认,先沿用通用字段,不臆造。
- **`page` 别名兼容**:cobra `Aliases` 已验证可在 help 中提示且不破坏调用;若发现别名在 `wlt stock in page --help` 链路有异常,改为注册隐藏的 `page` 子命令兜底。
