# internal 层测试与验证约定

> 本层是纯逻辑包（HTTP 客户端、配置、参数收集、信封），**可以也应该**做单元测试——
> 命令层做不到的部分（`RunE` 依赖全局状态）在这里补齐。

---

## 1. 本层测试文件

| 文件 | 覆盖 | 风格 |
|------|------|------|
| `internal/cmdutil/helpers_test.go` | `flagToParamKey` kebab→camel（表驱动 9 例） | 表驱动 `cases := []struct{ in, want string }` |
| `internal/cmdutil/timeflags_test.go` | `AddTimeRangeFlags` / `CollectTimeRangeFlags` 往返；空 end 不产 key | 构造 `&cobra.Command{}` + `Flags().Set` |
| `internal/cmdutil/crud_test.go` | `DocumentCmds` 命令结构、`page` 别名、flag 注册、TimeKey 空时不注册时间 flag | 命令树遍历 + 辅助 `mustHaveSub` / `subNamed` |
| `internal/output/output_test.go` | stdout 两种信封形状、stderr 形状、`WriteRaw` | `bytes.Buffer` 断言子串 |
| `internal/apierr/apierr_test.go` | `FromCommonResult` 三态、`ExitType` 全表 | 表驱动 |
| `internal/client/client_test.go` | GET/POST/DELETE 成功、`code!=0` 报错、`DoRaw`、buildRequest | `httptest.NewServer` |
| `internal/config/config_test.go` | 默认配置、保存重载、`SetProfileOverride`、非法 `SetActive` | 临时目录 + `NewManagerWithPath` |

---

## 2. 风格约定

- **包内测试**（`package <same>`），不是 `_test` 外部包 —— 需要访问包私有符号
  （`flagToParamKey`、`collectDocumentTimeRange`、`saveLocked` 等）。
- 函数名 `Test<Subject><Aspect>`。
- 表驱动优先（`helpers_test.go`、`apierr_test.go` 是样例）。
- HTTP 相关用 `httptest.NewServer`（`internal/client/client_test.go:11-24` 的 `setupTestClient`），
  **不 mock `http.Client`**。
- `client` 由 `RequestContext` 直构，无需 auth manager（无状态鉴权的直接体现）。
- 配置类测试用 `NewManagerWithPath(tempPath)`，不要碰真实的 `~/.wlt/config.yaml`。

---

## 3. 验证命令

```bash
go build ./... && go vet ./... && go test ./...
```

| 目标 | 命令 | 锚点 |
|------|------|------|
| 全量测试 | `make test` → `go test ./...` | `Makefile:10-11` |
| 构建 | `make build` | `Makefile:7-8` |

---

## 4. 改本层代码时的额外检查

- [ ] 改了 `internal/client` 的请求/响应契约 → 跑全量 `go test ./...`（`cmd/**` 依赖它）。
- [ ] 改了 `internal/output` 信封 → 同步更新 `internal/output/output_test.go` 断言。
- [ ] 改了 exit code 常量或 `ExitType` → 确认没有把 code 3/6 变成产出路径（保留值语义）；
      同步 `apierr_test.go` 的 `ExitType` 全表。
- [ ] 改了 `internal/config` 的 `Profile` → 同步 `UpdateProfileField` switch、`defaultConfig()`、
      以及 `cmd/config/` 的 show/set 输出（**不要**重新引入凭证字段）。
- [ ] 改了 `internal/cmdutil` 的工厂签名 → 这是 cli 层的公共契约，
      同步 [../cli/command-factories.md](../cli/command-factories.md) 并跑 `cmd/**` 的 flag 装配测试。