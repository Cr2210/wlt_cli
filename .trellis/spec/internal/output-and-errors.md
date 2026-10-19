# 输出与错误（`internal/output` + `internal/apierr`）

> 本文件讲**实现级**：信封函数、`ExitError` 结构、退出码常量与类型映射的来源。
> 命令层「怎么选 `ParsePagedJSON` 还是 `OutputJSON`」以及**退出码完整表**见
> [../cli/output-protocol.md](../cli/output-protocol.md)（避免两处重复）。

---

## 1. `internal/output` —— 信封

### 1.1 `ExitError`

```go
// internal/output/output.go:13-19
type ExitError struct {
	Code    int
	Type    string
	Message string
	Hint    string
	Detail  map[string]any
}
```

构造函数：

| 函数 | 行 | 说明 |
|------|----|------|
| `NewExitError(code int, msg string, hint string) *ExitError` | `:26` | `Type` 自动由 `apierr.ExitType(code)` 派生 |
| `NewExitErrorWithDetail(code, msg, hint, detail)` | `:36` | 带 `detail` map |

→ 调用方**不需要**手写 `Type`；只要给对 `code`，`type` 字符串就正确。

### 1.2 写出函数

| 函数 | 行 | 目标流 | 产出 |
|------|----|--------|------|
| `WriteJSON(out, data)` | `:60` | stdout | `{"data": <data>}` |
| `WritePagedJSON(out, list, total, pageNo, pageSize)` | `:67` | stdout | `{"data": <list>, "meta": {"page_no": N, "page_size": N, "total": N}}` |
| `WriteError(errOut, exitErr)` | `:81` | **stderr** | `{"error": {"type","code","message","hint?","detail?"}}` |
| `WriteRaw(out, data)` | `:96` | stdout | 原样字节 |

### 1.3 单测

`internal/output/output_test.go`：

- `TestWriteJSON` `:10-27` 断言 `data` 字段；
- `TestWritePagedJSON` `:29-47` 断言 `meta.total`；
- `TestWriteError` `:49-62` 用 `NewExitError(3, ...)` 断言输出含 `"type":"authentication"` 与 `"code":3`
  —— 这是**该退出码在全仓的唯一使用处**（测试），不是产出路径；
- `TestWriteRaw` `:64-70`。

---

## 2. `internal/apierr` —— 后端错误与退出码

### 2.1 `APIError`

```go
// internal/apierr/apierr.go:9-13
type APIError struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	HTTPStatus int    `json:"-"`
}
```

- `Error()` = `fmt.Sprintf("API error (http=%d, code=%d): %s", e.HTTPStatus, e.Code, e.Msg)`（`:17-19`）。
- `FromCommonResult(httpStatus int, body []byte) *APIError`（`:28-37`）：
  - 非法 JSON → `New(httpStatus, -1, "failed to parse response: ...")`；
  - `code == 0` → `nil`；
  - 否则 `New(httpStatus, code, msg)`。
- 测试：`internal/apierr/apierr_test.go`（成功 `:7-13`、`code:401` `:15-30`、非法 JSON `:32-41`、`ExitType` 全表 `:43-62`）。

### 2.2 退出码常量与类型映射

```go
// internal/apierr/apierr.go:43-49
const (
	ExitGeneral        = 1
	ExitConfig         = 2
	ExitAuthentication = 3
	ExitValidation     = 4
	ExitAPI            = 5
	ExitNetwork        = 6
)

// ExitType(code) 映射为字符串（:53-68）
// 2→config, 3→authentication, 4→validation, 5→api_error, 6→network, default→general
```

- `ExitType(99)` → `"general"`（default 分支，测试 `apierr_test.go:43-62` 覆盖）。
- ⚠️ **`ExitAuthentication = 3` 与 `ExitNetwork = 6` 有定义、有类型映射，但全仓无产出路径**。
  这是保留值；不要把退出码 6 与网络错误挂钩描述（详见下节与 cli 层的完整表）。

---

## 3. ⚠️ 已知文档漂移（本 spec 的纠正对象）

| 漂移 | 现状 |
|------|------|
| `skills/SKILL.md` 与 `skills/references/api-conventions.md` 把退出码 6 描述为 network 类（并将网络错误列为它的含义、附排障流程） | **错误**。全仓 `grep -rho "NewExitError([0-9]"` 无 `(6)`；网络错误经 `client.do` 返回普通 error，被命令层包成 `NewExitError(5, ...)`。`ExitNetwork` 仅有定义 + `ExitType` 映射。 |
| 同一批文档把已删除的 legacy 子命令工厂列为 CRUD 模式之一（标识符已随文件删除） | **已删除**：该文件与标识符在 `internal/cmdutil/` 下均不存在，全仓零命中。 |
| 「分页子命令全仓统一叫 `list`」 | **不成立**：`finance`(6)/`produce`(2)/`quality`(1)/`waybill`(1) 仍以 `page` 为主名。 |

### 实际产出分布（实测）

```bash
grep -rho "NewExitError([0-9]" --include=*.go . | sort | uniq -c
#  (1)×1  (2)×6  (4)×102  (5)×319   ← 无 (3) 于生产代码、无 (6)
```

- 生产代码实际产出：**1 / 2 / 4 / 5**。
- code 3 仅出现在测试（`internal/output/output_test.go:51`）。

> `skills/` 的修正属独立任务（本轮只记录，不改 `skills/`）。

---

## 4. 命令层封装（`internal/cmdutil/output.go`）

| 函数 | 行 | 转发到 |
|------|----|--------|
| `OutputJSON(data any) error` | `:12` | `output.WriteJSON(os.Stdout, data)` |
| `OutputPagedJSON(list, total, pageNo, pageSize) error` | `:17` | `output.WritePagedJSON(...)` |
| `OutputRaw(data []byte)` | `:22` | `output.WriteRaw(...)` |
| `ParsePagedJSON(data json.RawMessage, pageNo, pageSize int) error` | `:29` | 解析 `{list,total}` → `OutputPagedJSON` |

`ParsePagedJSON` 的实现要点（`Total` 用 `json.Number` 以容忍字符串型 total 等）见
[../cli/output-protocol.md](../cli/output-protocol.md) §3。

---

## 5. 退出码写出的唯一总入口

`cmd/root.go:67-85` 的 `Execute()`：

```go
err := rootCmd.Execute()
if err != nil {
	var exitErr *output.ExitError
	if as, ok := err.(*output.ExitError); ok {
		exitErr = as
	} else {
		exitErr = &output.ExitError{Code: 1, Type: "general", Message: err.Error()}
	}
	_ = output.WriteError(os.Stderr, exitErr)
	os.Exit(exitErr.Code)
}
return nil
```

- 非 `*ExitError` 的 error → 兜底 `Code: 1`、`Type: "general"`。
- `main.go` 里的 `os.Exit(1)` 实际不可达（`os.Exit` 已在 `Execute()` 内发生）。