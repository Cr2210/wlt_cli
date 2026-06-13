# 认证与配置

## 认证命令

| 命令 | 说明 | 参数 |
|------|------|------|
| `wlt auth login` | 交互式登录（表单） | 无（交互式输入用户名、密码、企业名称） |
| `wlt auth logout` | 登出 | 无 |
| `wlt auth status` | 查看认证状态 | 无 |

### 认证状态检查

```bash
wlt auth status
# 成功输出:
# { "status": "logged_in", "tenantId": "...", "expiresTime": "..." }
```

### 登录流程

```bash
wlt auth login
# 交互式表单：
# ? 用户名: admin
# ? 密码: ****
# ? 企业名称: 维链通
```

## 配置命令

| 命令 | 说明 | 参数 |
|------|------|------|
| `wlt config init` | 交互式初始化配置 | 无 |
| `wlt config show` | 显示当前配置 | 无 |
| `wlt config set <key> <value>` | 设置配置项 | key, value |

### 配置文件

配置文件位置：`~/.wlt/config.yaml`

```yaml
profiles:
  sit:
    base_url: https://erpsit.api.w-lian.com
    api_prefix: /admin-api
    tenant_id: "1"
    enterprise_type: "1"
    access_token: "..."
    refresh_token: "..."
    expires_time: "2024-01-01T00:00:00Z"
  prod:
    base_url: https://erpapi.w-lian.com
    api_prefix: /admin-api
    tenant_id: "1"
    enterprise_type: "1"
```

### 环境切换

```bash
# 默认使用 sit 环境
wlt stock warehouse list

# 切换到生产环境
wlt stock warehouse list --profile prod

# 静默模式（只输出数据）
wlt stock warehouse list --quiet
```

## 版本信息

```bash
wlt version
```
