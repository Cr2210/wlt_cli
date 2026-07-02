# wlt — 维链通 ERP 命令行工具

AI Agent 原生设计的 ERP CLI 工具。

## 安装

### 从源码构建

```bash
# 需要 Go 1.22+
git clone <repo>/weiliantong-cli.git
cd weiliantong-cli
make build
```

### npm 安装（需先发布）

```bash
npm install -g @weiliantong/cli
```

## 快速开始

```bash
# 初始化配置
wlt config init

# 登录
wlt auth login

# 查看状态
wlt auth status

# 查询仓库
wlt stock warehouse list

# 查询库存
wlt stock query page --warehouse-id 1

# 通用 API 调用
wlt api GET /erp/warehouse/simple-list
```

## 配置

配置文件位于 `~/.wlt/config.yaml`，支持多环境：

```bash
wlt config show              # 查看配置
wlt config set sit.base_url https://erpsit.api.w-lian.com
wlt --profile prod auth login  # 使用生产环境
```

## 库存模块命令

| 子域 | 命令 | 操作 |
|------|------|------|
| 仓库 | `wlt stock warehouse` | list, get, create, update, delete, update-status, simple-list |
| 库存 | `wlt stock query` | get, page, page-count, count, ledger, ledger-count |
| 入库 | `wlt stock in` | page, page-count, get, create, update, delete, update-status |
| 出库 | `wlt stock out` | 同入库 |
| 调拨 | `wlt stock move` | 同入库 |
| 盘点 | `wlt stock check` | 同入库 |
| 明细 | `wlt stock record` | page, page-count, get, count, record-page, total-cost |

## 输出协议

- **stdout**：JSON 格式数据
- **stderr**：JSON 格式错误 + 进度信息
- **退出码**：0 成功 / 1 通用 / 2 配置 / 3 认证 / 4 验证 / 5 API / 6 网络

## AI Agent 集成

详见 `skills/wlt.md`（主入口）和 `skills/references/` 目录下的各模块参考文档，包含完整的命令参考、工作流示例和错误处理指南。

## 开发

```bash
make build    # 构建
make test     # 测试
make install  # 安装
```
