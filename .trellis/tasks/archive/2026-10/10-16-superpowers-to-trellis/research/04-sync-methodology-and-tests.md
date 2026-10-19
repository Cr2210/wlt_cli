<!-- 本文件是 `conventions-extraction` 的拆分件之一（按研究产物约定：一主题一文件）。
     初衷：单文件 57KB 超过 context_injection.max_file_bytes(32768) 会被截断。
     全量导航见同目录 `00-index.md`。文中形如「§C-1」「§4.5」的交叉引用按该导航定位。 -->

## 6. 后端基线同步方法论

来源：`docs/superpowers/specs/2026-09-28-backend-sync-0923-design.md`（+ `docs/superpowers/plans/2026-09-28-backend-sync-0923.md`）。

### 6.1 「CLI 是薄透传层」原则

- `sync-design:7`：差异分析结论「**3 个全新模块 + 3 处已有域增量**，其余为后端内部口径调整或响应字段透传（CLI 是薄透传层，无需改动）」。
- `sync-plan:7`：`Architecture: 纯薄透传层：每个子命令构造 query/body 调用后端 REST 端点，输出遵循 stdout 数据 JSON / stderr 错误 JSON / 退出码 0-6 协议。`
- 推论（代码可直接验证）：提交字段（如 `invoiceFlag`/`isPrepaid`）经 `--data` JSON 直传天然支持，**无需加 flag**；只有 `PageReqVO` 的新**筛选**字段需要加 flag（`sync-design:18`、后端仓库 memory）。
- `sync-design:27-31` 非目标：不封装组合工作流（保持单接口单命令）；不改 `DocumentCmds`/`AddCRUDToParent`/`CrudXxxCmd` 工厂签名，「不足之处用域内自定义命令补齐」。

### 6.2 差异分析方法（可复现步骤）

1. **定位基线两端**：CLI 侧上次同步 commit ↔ 后端仓库目标 commit。
2. **git log 差异分析**（memory `wlt-cli-backend-repo`）：`新 Controller / SaveReqVO 新字段 / PageReqVO 新筛选` 三类是 CLI 需要跟进的信号。
3. **按域统计端点数**。踩坑记录（memory，务必写进 spec）：「注意后端一个域可能有多个 Controller（如 `/erp/settlement` 与 `/erp/logistics-settlement`），grep 端点清单时分开看，避免把两个文件的输出混在一起数错端点」。初次分析曾因此把物流结算误判为 14 端点，纠正为 13（`sync-design:108`）。
4. **A/B/排除 三分分类**，逐类处置。

### 6.3 A / B / 排除 分类法（`sync-design:9-18`）

| 类别 | 定义 | 0923 实例 | 处置 |
|---|---|---|---|
| **A** | 后端全新模块（前缀 + 端点数） | A1 `/erp/logistics-settlement`（13 端点，08-28）；A2 `/erp/sale-accounting`（11，09-15）；A3 `/erp/finance-freight-apply`（11，09-17） | 新增命令组，挂到既有域父命令 |
| **B** | 已有域增量（新端点 / 新 flag / 新提交字段） | B1 结算单忽略剩余 `/erp/settlement/ignore-remainder` + `invoiceFlag`；B2 预付关联双业务类型 `bizType`；B3 收付款预付标识 `isPrepaid` 查询筛选 | 改既有命令（加子命令 / 加 flag） |
| **排除** | 后端内部逻辑或响应透传 | 发票核销/开票申请卡口过滤、出入库详情磅单回填、驾驶舱/资金占用/物流预付统计口径、WebSocket 鉴权、通知公告评论回复（CLI 未接 notice）、扫码接单/费用下发 | **CLI 无改动** |

`sync-design:18` 原文：「`invoiceFlag`/`isPrepaid` 的提交字段经 `--data` JSON 直传天然支持，仅在 skills 文档补充说明。」

### 6.4 挂载决策模式（`sync-design:33-42`）

新模块挂载需显式决策并落表：模块 → 挂载路径 → 文件。0923 结果：A1→`wlt settlement logistics <sub>`（`cmd/settlement/logistics.go`）、A2→`wlt sale accounting <sub>`（`cmd/sale/accounting.go`）、A3→`wlt finance freight-apply <sub>`（`cmd/finance/freight_apply.go`）；B1→`wlt settlement main ignore-remainder`（与 main CRUD 平级）、B2→`wlt order prepayment-relation list/create` 加 `--biz-type`、B3→`wlt finance` 收付款 list 加 `--is-prepaid`。
挂载进既有域时**不改 `cmd/root.go`**（`sync-plan:7`：「新子命令组通过 `init()` 注册到既有域父命令变量，不动 root.go」）。

### 6.5 实施顺序与提交粒度（`sync-design:114-118`）

① B 类增量（小，先行）→ ② A1 → ③ A2 → ④ A3 → ⑤ 文档收尾与 backlog 更新。「按子项目分 4+1 次 commit」，conventional commits + 中文描述。

### 6.6 重做同步所需的工具/配置事实

| 事实 | 值 | 出处 |
|---|---|---|
| 后端仓库位置 | `D:\Projects\weiliantong\weiliantong-api`（**只读参考，勿改**） | `sync-plan:22`；memory `wlt-cli-backend-repo` |
| 后端分支线 | `dev-1.7` → `ge-dev-1.8`（云效 CodeUp，MR 开发流）；0930 同步时后端分支 `ge-dev-1.8` | `sync-plan:22`；memory |
| CLI 侧上次同步基线 | 2026-08-21（commit `6318e49`） | `sync-design:7` |
| 后端目标基线 | 2026-09-23（commit `2e52721b`） | `sync-design:7`；memory |
| CLI 侧完成 commit | `398d6ee`（2026-09-28） | memory |
| Swagger 导出（Round 1 用） | `D:\Download\erp.md`（510 端点：GET 303 / POST 65 / PUT 92 / DELETE 50）；参数表在 `**请求参数**` 段 | `backlog:7`、`backlog:83-84` |
| grep 端点参数的示例命令 | `grep -nA40 "接口地址**:\`/admin-api/erp/waybill/page\`" D:/Download/erp.md` | `backlog:84` |
| 差异分析入口 | 后端 git log（新 Controller / SaveReqVO 新字段 / PageReqVO 新筛选） | memory |
| 已知递延 | export 命令 `fmt.Println` 非 JSON 输出是仓库级既有模式（20+ 处），待专项清理 | memory；本文 §5.1 实测 21 处 |

---

## 7. 测试与验证约定

### 7.1 测试文件清单与检查内容

| 文件 | 检查什么 | 关键锚点 |
|---|---|---|
| `internal/cmdutil/helpers_test.go` | `flagToParamKey` kebab→camel 映射（表驱动 9 例） | `:5-22` |
| `internal/cmdutil/crud_test.go` | `DocumentCmds` 命令结构（子命令齐全、`page` 别名、flag 注册）+ `collectDocumentTimeRange` 往返/空值 | `:23-47`、`:49-87`、`:89-98`；辅助 `mustHaveSub` `:10-21`、`subNamed` `:101-108` |
| `internal/cmdutil/timeflags_test.go` | `AddTimeRangeFlags`/`CollectTimeRangeFlags` 往返 + 空 end 不产 key | `:9-26` |
| `internal/output/output_test.go` | stdout 信封形状（`data`、`meta.total`）、stderr 形状（`type`/`code`）、`WriteRaw` | `:10-70` |
| `internal/apierr/apierr_test.go` | `FromCommonResult` 成功/`code:401`/非法 JSON（code=-1）、`ExitType` 全表（含 `99→general`） | `:8-70` |
| `internal/client/client_test.go` | httptest 驱动 GET/POST/DELETE 成功、`code!=0` 报错、`DoRaw` 原样返回；client 由 `RequestContext` 直构（无 auth manager） | `:26-128`；`setupTestClient` `:11-24` |
| `internal/config/config_test.go` | 默认配置创建（active=sit + sit/prod）、保存重载、`SetProfileOverride`、非法 `SetActive` | `:9-96` |
| `cmd/stock/stock_doc_test.go` | 包内辅助 `newStockPagedCmd`/`newStockCountCmd` 的 flag 集合与「count 不应有 page/headers」 | `:9-33` |
| `cmd/order/order_prepayment_test.go` | flag 装配（list/create/update-amount/available-*）+ 子命令存在性（5 个） | `:5-69` |
| `cmd/settlement/settlement_main_test.go` | `newSettlementMainIgnoreRemainderCmd` 的 id/type/reason | `:5-12` |
| `cmd/settlement/logistics_test.go` | 13 子命令存在性 + list/ignore-remainder flag | `:5-51` |
| `cmd/sale/accounting_test.go` | 11 子命令存在性 + list/waybill-list flag | `:5-43` |
| `cmd/finance/freight_apply_test.go` | 11 子命令存在性 + list/available-settlement flag | `:5-43` |
| `cmd/finance/finance_receipt_payment_test.go` | `newFinanceReceiptPaymentPageCmd` 的 `is-prepaid`/`page-no`/`page-size` | `:5-12` |

### 7.2 命名模式与测试风格

- 文件名：`<被测文件>_test.go`；包内测试（`package <same>`，非 `_test` 外部包），以便访问包私有构造函数（如 `newStockPagedCmd`、`newOrderPrepayRelationListCmd`）。
- 函数名：`Test<Subject><Aspect>`，如 `TestFlagToParamKey`、`TestDocumentCmdsStructure`、`TestSaleAccountingCmdSubcommands`、`TestLogisticsSettlementListCmdFlags`。
- 两大类断言模式（`sync-plan:20` 明确「测试风格对齐 `cmd/order/order_prepayment_test.go`：flag 存在性 + 子命令存在性（不 mock HTTP）」）：
  1. **flag 装配**：`for _, f := range []string{...} { if c.Flags().Lookup(f) == nil { t.Errorf(...) } }`
  2. **子命令存在性**：遍历 `parent.Commands()` 与 `want` map 比对
  - **反向断言**也出现：`stock_doc_test.go:28-32` 断言 count 命令**不应**有 `page-no/page-size/headers`；`crud_test.go:89-98` 断言 `TimeKey` 为空时不注册 `start-time`。
- **不 mock HTTP**（命令层）；HTTP 层用 `httptest.NewServer`（`client_test.go:11-24`）。
- 已知限制（`stock-plan:18`）：「RunE 依赖全局 `APIClient`（见 `internal/cmdutil/client.go:39`），无法纯单测；工厂层逻辑（时间转换、命令结构）用单元测试守护，迁移类任务靠编译 + `make test` + `--help` 结构核对验证等价性」。

### 7.3 Makefile 验证命令

| 目标 | 命令 | 锚点 |
|---|---|---|
| `make build` | `go build -ldflags "-X github.com/weiliantong/cli/internal/build.Version=$(VERSION)" -o wlt .` | `Makefile:7-8`；`VERSION` 由 `git describe --tags --always --dirty` 注入（`:4`） |
| `make test` | `go test ./...` | `Makefile:10-11` |
| `make install` | `build` 后 `cp wlt /usr/local/bin/` | `Makefile:13-14` |
| `make clean` | 删 `wlt`/`wlt.exe`/`dist` + `go clean` | `Makefile:16-19` |
| `make snapshot` / `make release` | `goreleaser release --snapshot --clean` / `goreleaser release --clean` | `Makefile:25-31` / `:34-40` |
| `make cross-compile` | windows/linux/macOS amd64+arm64 裸二进制 | `Makefile:43-51` |

计划中反复出现的验证序列（`sync-plan:1449`、`stock-plan:17`）：`go build ./... && go vet ./... && go test ./...`，外加 `wlt <group> --help` 命令树核对（`sync-design:112`）。CLAUDE.md 的 `make build/test/install/clean` 与实际一致。

---

