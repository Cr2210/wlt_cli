# API 约定与通用模式

## 后端 API 规范

### 基础路径

- **SIT 环境**: `https://erpsit.api.w-lian.com/admin-api`
- **生产环境**: `https://erpapi.w-lian.com/admin-api`
- **本地开发**: `http://localhost:50000/admin-api`（dev/local profile）

### 统一响应格式 `CommonResult<T>`

```json
{
  "code": 0,
  "data": T,
  "msg": ""
}
```

- `code: 0` 表示成功
- 非 0 表示业务错误，`msg` 中包含错误描述

### 认证（无状态）

- 请求头：`Authorization: Bearer {token}` —— 由 `--token <accessToken>` 提供，CLI 自动加 `Bearer ` 前缀
- CLI 不保存登录态、不自动刷新；token 过期需调用方重新获取并通过 `--token` 传入

### 多租户

- 请求头：`tenant-id`（由 `--tenant-id` 提供，**必填**）、`enterprise-type`（来自 profile，可选）
- `tenant-id` 不再从配置文件读取，每次调用由 flag 提供

### 分页

- 查询参数：`pageNo=1&pageSize=20`
- 响应格式：`{ "list": [...], "total": N }`
- CLI 默认 `page_no=1`, `page_size=20`

### 删除

- 单个删除：`DELETE /delete?id=N` 或 `DELETE /delete --data '{"id":N}'`
- 批量删除：`DELETE /delete?ids=1,2,3` 或 `DELETE /delete-list --data '{"ids":[1,2,3]}'`

## CRUD 模式

CLI 使用三种 CRUD 模式：

### 1. 标准 CRUD（通过 CRUDConfig）

```
list          → GET  /page
get           → GET  /get?id=N
create        → POST /create
update        → PUT  /update
delete        → DELETE /delete
update-status → PUT /update-status
simple-list   → GET /simple-list（可选）
page-count    → GET /page-count（可选）
```

### 2. Stock 子命令 CRUD（newCRUDSubCmd）

```
list          → GET  /page
get           → GET /get?id=N
create        → POST /create
update        → PUT /update
delete        → DELETE /delete
update-status → PUT /update-status
```

### 3. 自定义命令

每个模块可能有额外的自定义命令，如：
- `export` → 导出 Excel
- `import` → 导入 Excel
- `summary` → 汇总统计
- `cancel` / `reopen` / `complete` → 生命周期操作
- `link-waybill` / `unlink-waybill` → 关联操作

## 通用 API 调用

当快捷命令不覆盖所需操作时：

```bash
# GET 请求
wlt api GET /erp/custom-endpoint

# POST 请求
wlt api POST /erp/custom-endpoint --data '{"key":"value"}'

# 带查询参数
wlt api GET /erp/custom-endpoint --params '{"key":"value"}'

# 预览请求（不发送）
wlt api GET /erp/custom-endpoint --dry-run
```

## 输出协议

- **stdout**: JSON 数据（成功时）
- **stderr**: JSON 错误（失败时），格式：`{ "type": "...", "message": "...", "hint": "..." }`
- **退出码**:
  - `0` 成功 / `1` 通用错误 / `2` 配置错误 / `6` 网络错误
  - `4` 参数错误 —— **缺少 `--token`/`--tenant-id` 也归此类**
  - `5` API 错误 —— 服务端 token 失效(`code:401 账号未登录`)当前也表现为 `5`,需重新获取 token 并通过 `--token` 传入
  - `3` 鉴权(保留常量,当前未使用;统一错误分类为后续优化)
