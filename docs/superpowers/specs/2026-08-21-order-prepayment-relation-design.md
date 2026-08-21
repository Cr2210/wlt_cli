# wlt 订单预付关联命令设计

日期:2026-08-21 ｜ 状态:已确认,待实现

## 背景

后端新增订单预付关联模块(`/erp/order-prepayment-relation`),用于把订单与预付款来源(付款单 / 供应商期初)建立关联并维护关联金额。共 5 个接口,CLI 尚无对应命令:

| 接口 | 方法 | 说明 |
|---|---|---|
| `/erp/order-prepayment-relation/update-amount` | PUT | 更新关联金额(id + relationAmount) |
| `/erp/order-prepayment-relation/create` | POST | 创建关联(orderId/relationType/relationId/relationAmount) |
| `/erp/order-prepayment-relation/page` | GET | 按订单分页查询关联记录(orderId 必填) |
| `/erp/order-prepayment-relation/available-payments` | GET | 按供应商查询可用付款单(supplierId 必填) |
| `/erp/order-prepayment-relation/available-initials` | GET | 按供应商查询可用期初(supplierId 必填) |

`relationType` 枚举:`PAYMENT`(付款单)/ `SUPPLIER`(供应商期初)。

## 目标

- 在 order 域下新增 `wlt order prepayment-relation` 子命令组,覆盖上述 5 个接口。
- 命令命名、flag 风格、输出协议与全仓既有约定一致。
- 同步更新 skills 文档与 CLAUDE.md。

## 非目标

- 不封装组合工作流(查询可用余额 → 创建关联 → 改金额)为一条命令,保持单接口单命令。
- 不改动 `DocumentCmds` / `AddCRUDToParent` 等通用工厂(本模块 5 个接口形状各异,不适合套用标准 CRUD 工厂)。

## 设计

### 1. 归属与命名

挂在 order 域:`wlt order prepayment-relation`,新文件 `cmd/order/order_prepayment.go`(package order,经 `init()` 注册到 `orderCmd`)。理由:后端路径即 order-prepayment-relation,`page` 接口 `orderId` 必填,语义是"订单的预付关联";order 域已有 main/plan,加第三个子模块顺理成章;与 finance 域 `prepayment-apply`(对应 `/erp/finance-prepayment-apply`)命名规则一致。

### 2. 子命令与 flag

```
wlt order prepayment-relation
├── list                 分页查询预付关联（GET /page）
│     --order-id（必填，int64） --relation-type PAYMENT|SUPPLIER --headers --page-no --page-size
├── create               创建预付关联（POST /create）→ 输出新关联 ID
│     --order-id --relation-type --relation-id --relation-amount（均必填；金额单位元，float64）
├── update-amount        更新关联金额（PUT /update-amount）
│     --id --relation-amount（均必填）
├── available-payments   查询可用的付款单列表（GET /available-payments）
│     --supplier-id（必填，int64） --order-id --biz-date-start --biz-date-end --no --headers --page-no --page-size
└── available-initials   查询可用的供应商期初列表（GET /available-initials，与上共享工厂函数）
```

- **create / update-amount 用离散 flag** 而非 `--data` JSON:字段少且固定(4 个 / 2 个),`relationType` 枚举写入 flag 帮助文本,Agent 免拼 JSON;有 `data_sync resend --id`(map 请求体)先例。请求体以 `map[string]any` 构造,字段名转 camelCase(orderId/relationType/relationId/relationAmount)。
- **available-payments / available-initials 共享一个工厂函数** `newPrepayAvailableCmd(path, label)`:二者 flag 与响应结构完全一致,仅端点与标签不同。
- `--order-id` 在 available-* 上是可选参数(排除已关联的订单)。

### 3. 输出处理

- `page` 接口响应 `data` 为**裸数组**(无 `{list,total}` 包裹,见接口文档响应示例)→ 直接 `cmdutil.OutputJSON(resp.Data)`,`ParsePagedJSON` 会因数组解进 struct 而报错,不可用。若实测响应为包裹结构,切换为 `ParsePagedJSON`(一行改动)。
- `available-*` 返回标准 PageResult(`{list,total,statistics}`)→ 走 `cmdutil.ParsePagedJSON` 统一分页信封(`{data, meta:{page_no,page_size,total}}`);`statistics` 目前为空对象,被丢弃可接受。
- 错误输出遵循全仓约定:API 失败退出码 5,JSON 解析失败退出码 4,错误走 stderr 的 `output.NewExitError`。

### 4. 文档同步

- `skills/references/order.md`:新增"订单预付关联"章节(命令表 + 示例)。
- `skills/SKILL.md`:order 域行与触发词补"预付关联"。
- `CLAUDE.md`:项目结构中 order 目录注释补 prepayment-relation。

## 测试与验证

- 仿 `cmd/stock/stock_doc_test.go` 在 `cmd/order/` 加 flag 装配测试(各子命令的必填 flag、分页 flag、available-* 的筛选 flag 是否注册齐全)。
- `go build ./...` + `go test ./...` + `wlt order prepayment-relation --help` 冒烟。
