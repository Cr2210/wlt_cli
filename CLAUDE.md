# CLAUDE.md — weiliantong-cli

## 项目概述

weiliantong-cli（`wlt`）是维链通 ERP 系统的命令行工具，AI Agent 原生设计。

## 技术栈

- Go 1.26+ / Cobra CLI 框架
- 关键依赖：`github.com/spf13/cobra`, `github.com/tidwall/gjson`, `gopkg.in/yaml.v3`, `github.com/charmbracelet/huh`

## 开发命令

```bash
make build          # 构建 wlt 二进制
make test           # 运行所有测试
make install        # 安装到 /usr/local/bin
make clean          # 清理构建产物
```

## 项目结构

```
cmd/                    # Cobra 命令层，每个业务域一个独立 Go 包
├── root.go             # 根命令 + 全局 PersistentPreRunE（配置加载）
├── api.go              # wlt api — 通用 HTTP 调用
├── version.go          # wlt version
├── completion.go       # Shell 自动补全
├── help.go             # 中文帮助
├── config/             # 配置命令（init/show/set）
├── system/             # 系统管理（user/dept/role/menu/dict）
├── stock/              # 库存管理（warehouse/query/in/out/move/check/record）
├── product/            # 产品管理（product/unit/category/metrics）
├── contract/           # 合同管理（main/provision/service/transport）
├── customer/           # 客户管理（复用 partner 包）
├── supplier/           # 供应商管理（复用 partner 包）
├── partner/            # 客户/供应商共享逻辑（CRUD/invoice/settlement/credit）
├── purchase/           # 采购管理（in/return）
├── sale/               # 销售管理（out/return）
├── report/             # 报表管理（stock/purchase/sale）
├── stats/              # 数据总览（overview/finance/purchase/sale/produce/stock）
├── finance/            # 财务管理（account/payment/receipt/refund/settlement/write-off/invoice-apply）
├── waybill/            # 运单管理（source/push-config）
├── order/              # 订单管理（main/plan）
├── produce/            # 生产管理（main/plan）
├── quality/            # 质检管理（inspection/weight）
├── weight/             # 称重管理（waybill）
├── settlement/         # 结算管理（main）
├── invoice/            # 发票管理（main）
├── homepage/           # 首页数据总览（dashboard1-6/inventory-backlog/product-ranking）
├── operate_log/        # 操作日志查询
├── job_trigger/        # 定时任务触发
├── profit/             # 利润管理（profit-event/profit-calculation）
├── data_sync/          # 数据同步管理
└── screen/             # 大屏数据管理
internal/
├── apierr/             # 结构化错误类型
├── build/              # 版本信息（ldflags 注入）
├── client/             # HTTP 客户端（认证头 + 多租户 + 重试 + CommonResult 解析，无状态）
├── cmdutil/            # 命令共享工具
│   ├── client.go       #   EnsureClient(), InitManagers(), GetClient(), SetAuthFlags()
│   ├── helpers.go      #   ParseJSONData, CollectStringFlag/Flags, CollectIntFlags
│   ├── output.go       #   OutputJSON, OutputPagedJSON, OutputRaw, ParsePagedJSON
│   ├── crud.go         #   CRUDConfig, NewCRUDGroup, AddCRUDToParent + 工厂函数
│   ├── crud_legacy.go  #   NewCRUDSubCmd（库存/采购/销售子域通用 CRUD）
│   └── stats.go        #   NewStatsGetCmd, CollectTimeRangeFlags, AddStatsFlags
├── config/             # 配置文件管理（~/.wlt/config.yaml）
└── output/             # JSON 输出格式化（stdout/stderr 分离）
skills/                 # AI Agent Skills Markdown
```

## 架构约定

- **按功能分包**：每个业务域是独立的 Go 包（`cmd/<domain>/`），暴露 `Register(parent *cobra.Command)` 函数
- **无状态鉴权**：CLI 不保存登录态、无 `auth` 命令；每条业务命令通过全局 flag `--token`/`--tenant-id` 传入鉴权，`cmdutil.EnsureClient()` 从 profile（base_url/api_prefix）+ flag 组装 `client.RequestContext` 并校验必填
- **命令注册**：子命令在各自文件中通过 `init()` 注册到父命令变量，`root.go` 调用每个包的 `Register()`
- **输出协议**：stdout 数据 JSON / stderr 错误 JSON / 退出码 0-6
- **CRUD 模式**：
  - 现代：`cmdutil.AddCRUDToParent(parent, CRUDConfig{...})` — 大多数域使用
  - 遗留：`cmdutil.NewCRUDSubCmd(name, apiPath, label)` — 库存/采购/销售子域使用
- **共享逻辑**：`cmd/partner/` 包封装了客户/供应商共享的 CRUD、开票、结算、授信逻辑

## 后端 API 约定

- 基础路径：`{base_url}{api_prefix}`（如 `https://erpsit.api.w-lian.com/admin-api`），由 `--profile` 提供，可用 `--base-url` 覆盖
- 统一响应：`CommonResult<T>` — `{ "code": 0, "data": T, "msg": "" }`
- 认证头：`Authorization: Bearer {token}` —— 由必填 flag `--token <accessToken>` 提供（CLI 自动加 `Bearer ` 前缀）
- 多租户头：`tenant-id`（由必填 flag `--tenant-id` 提供）、`enterprise-type`（来自 profile，可选）
- 分页：`GET /page?pageNo=1&pageSize=20` → `{ "list": [...], "total": N }`
- 删除：`DELETE /delete?ids=1,2,3`

## 新增业务域指南

添加新业务域只需：

1. 创建 `cmd/<domain>/` 目录
2. 定义父命令变量（如 `var xxxCmd = &cobra.Command{...}`）
3. 添加 `Register(parent *cobra.Command)` 函数
4. 可选：使用 `cmdutil.AddCRUDToParent(xxxCmd, cmdutil.CRUDConfig{...})` 快速生成 CRUD
5. 在 `cmd/root.go` 的 `init()` 中添加 `<domain>.Register(rootCmd)`

## 发布流程（云效 CodeUp）

### 前置准备

一次性安装 GoReleaser：
```bash
scoop install goreleaser
```

### 发布步骤

1. **打 tag 并推送到远程**
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. **本地构建所有平台产物**
   ```bash
   cd weiliantong-cli
   make snapshot    # 测试构建（不需要 tag）
   # 或
   make release     # 正式构建（需要 tag）
   ```

3. **构建产物位置**
   
   构建完成后，所有产物在 `dist/` 目录下：
   - `wlt-linux-amd64.tar.gz`
   - `wlt-linux-arm64.tar.gz`
   - `wlt-darwin-amd64.tar.gz`
   - `wlt-darwin-arm64.tar.gz`
   - `wlt-windows-amd64.zip`
   - `checksums.txt`
   - `wlt-skills.zip`

4. **上传到云效**
   
   手动上传到云效 CodeUp：
   - 进入代码库 → 标签 → 发行版
   - 点击"新建发行版"
   - 上传 `dist/` 目录下的所有文件

### 注意事项

install 脚本中的 `<owner>/weiliantong-cli` 需要替换为实际的下载地址，取决于团队内部分发方式：
- 直接从云效 Release 页面下载
- 通过内部文件服务器分发
