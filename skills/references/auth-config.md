# 鉴权与配置

## 鉴权（无状态）

`wlt` 采用**无状态鉴权**：CLI 不保存登录态，也没有 `login`/`logout` 命令。每条业务命令都由调用方（AI agent）在命令行直接传入鉴权信息。

### 必填 flag（业务命令）

| flag | 对应请求头 | 说明 |
|------|-----------|------|
| `--token <accessToken>` | `Authorization: Bearer <accessToken>` | 后端登录返回的 accessToken；CLI 自动加 `Bearer ` 前缀 |
| `--tenant-id <租户ID>` | `tenant-id` | 租户 ID |

- 缺失任意一个 → 退出码 **4（validation）**，提示「缺少必填鉴权参数」。
- token 失效（服务端返回 `code:401 账号未登录`）→ 退出码 **5（api_error）**，message 含 `code=401`；需重新获取 token 并通过 `--token` 传入。
- `wlt version` / `wlt config *` / `wlt completion` / `wlt --help` **不需要**鉴权 flag。

### 可选 flag

| flag | 说明 |
|------|------|
| `--profile sit\|prod` | 选择 profile（默认 `sit`），提供 `base_url`/`api_prefix` |
| `--base-url <url>` | 可选，临时覆盖 profile 的 `base_url` |
| `--quiet` | 静默模式，仅输出数据 |

### 示例

```bash
# 查询客户列表（必带 --token 与 --tenant-id）
wlt customer list --token <accessToken> --tenant-id 1

# 分页
wlt customer list --token <accessToken> --tenant-id 1 --page-no 2 --page-size 50

# 临时指向生产环境
wlt customer list --token <accessToken> --tenant-id 1 --profile prod

# 临时覆盖 base_url
wlt customer list --token <accessToken> --tenant-id 1 --base-url https://erpapi.w-lian.com
```

## 配置命令

| 命令 | 说明 | 参数 |
|------|------|------|
| `wlt config init` | 交互式初始化配置（名称 / API 地址 / API 前缀） | 无 |
| `wlt config show` | 显示当前配置 | 无 |
| `wlt config set <key> <value>` | 设置配置项 | key, value |

`config set` 的 key 格式为 `<profile>.<field>`，可设字段：`base_url`、`api_prefix`、`enterprise_type`。

### 配置文件

配置文件位置：`~/.wlt/config.yaml`。**只保存连接信息，不保存任何凭证。**

```yaml
active: sit
profiles:
  sit:
    base_url: https://erpsit.api.w-lian.com
    api_prefix: /admin-api
    enterprise_type: ""
  prod:
    base_url: https://erpapi.w-lian.com
    api_prefix: /admin-api
    enterprise_type: ""
```

### 环境切换

```bash
# 默认使用 sit 环境
wlt stock warehouse list --token <accessToken> --tenant-id 1

# 切换到生产环境
wlt stock warehouse list --token <accessToken> --tenant-id 1 --profile prod

# 静默模式（只输出数据）
wlt stock warehouse list --token <accessToken> --tenant-id 1 --quiet
```

## 版本信息

```bash
wlt version    # 无需鉴权
```
