# 配置与 profile（`internal/config`）

> 本包只管**连接信息**：后端地址与 API 前缀。凭证（token / tenant-id）**不落在配置里**，
> 每次调用经 flag 传入。见 [auth-and-client.md](./auth-and-client.md)。

---

## 1. 数据结构

```go
// internal/config/config.go:23-27
type Profile struct {
	BaseURL        string `yaml:"base_url"`
	APIPrefix      string `yaml:"api_prefix"`
	EnterpriseType string `yaml:"enterprise_type"`
}
```

- **无凭证字段**。历史上曾有 `AccessToken` / `RefreshToken` / `ExpiresTime` / `TenantID`，
  已在无状态鉴权重构（2026-06-26）中删除。
- 旧 config 中若残留这些键，YAML 反序列化时会**静默忽略**（不报错）——破坏性变更的兼容处理。

---

## 2. 配置文件位置与权限

| 项 | 值 | 锚点 |
|----|----|------|
| 默认路径 | `~/.wlt/config.yaml` | `internal/config/config.go:39-45`（`os.UserHomeDir()`） |
| 目录权限 | `0700`（`os.MkdirAll(dir, 0700)`） | `internal/config/config.go:88` |
| 文件权限 | `0600`（`os.WriteFile(..., 0600)`） | `internal/config/config.go:95` |
| 写入入口 | `Manager.Save()` → `saveLocked()` | `internal/config/config.go:80-95` |
| 测试用构造 | `NewManagerWithPath(p)` | `internal/config/config.go:48` |

---

## 3. 默认配置

`defaultConfig()`（`internal/config/config.go:181-195`）：

| profile | `base_url` | `api_prefix` |
|---------|-----------|--------------|
| `sit`（默认 active） | `https://erpsit.api.w-lian.com` | `/admin-api` |
| `prod` | `https://erpapi.w-lian.com` | `/weiliantong/admin-api` |

- 缺配置文件时 `NewManager()` 直接持有 `defaultConfig()`，因此 **sit 开箱即用**。
- `enterprise_type` 两个默认 profile 都未设置 → `buildHeaders` 不发该头。

---

## 4. 激活 profile 的解析

`ActiveProfile()`（`internal/config/config.go:100-113`）：

```go
name := m.profileOverride          // --profile flag（经 SetProfileOverride 注入）
if name == "" {
	name = m.cfg.Active            // 配置文件里的 active
}
p, ok := m.cfg.Profiles[name]
if !ok {
	return nil, fmt.Errorf("profile %q not found", name)
}
```

- **优先级**：`--profile` flag（非空）> 配置文件 `active`。
- `--profile` 默认值 `"sit"`（`cmd/root.go:88`）——注意这意味着**默认永远覆盖**配置文件里的 `active`；
  想用配置里的 active，需显式传入对应名字。
- 找不到 profile 时返回 error → 上层 `EnsureClient` 转成 `NewExitError(2, "当前环境配置不存在: ...", ...)`
  （`internal/cmdutil/client.go:48`）。

---

## 5. 字段更新

`UpdateProfileField(profileName, field, value)`（`internal/config/config.go:142-162`）：

```go
switch field {
case "base_url":        p.BaseURL = value
case "api_prefix":      p.APIPrefix = value
case "enterprise_type": p.EnterpriseType = value
default:                return fmt.Errorf("unknown field: %s", field)
}
```

- **只支持这 3 个字段**。`tenant_id` 分支已在无状态鉴权重构中删除——
  不要再往这里加租户/凭证字段。
- 更新后立即 `saveLocked()` 落盘。
- 调用方：`cmd/config/`（`config set`）。

---

## 6. 相关命令

| 命令 | 作用 |
|------|------|
| `wlt config init` | 初始化配置文件（写默认 sit/prod） |
| `wlt config show` | 展示当前配置（**不**含 token / tenant_id / expires_time —— 那些字段已不存在） |
| `wlt config set` | 更新 profile 字段（走 `UpdateProfileField`） |

`cmd/config/` 下命令**不调用 `EnsureClient()`**，因此不需要 `--token`/`--tenant-id`。

---

## 7. 测试

`internal/config/config_test.go:9-88`：

- 默认配置创建（`active == sit`，含 `sit`/`prod` 两个 profile）；
- 保存后重载一致；
- `SetProfileOverride` 生效；
- 非法 `SetActive` 报错。