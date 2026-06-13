# 系统管理 (system)

## 用户管理 (`wlt system user`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt system user list` | 分页查询用户 | `--username`, `--mobile`, `--status`, `--dept-id`, `--page-no`, `--page-size` |
| `wlt system user simple-list` | 用户精简列表 | 无 |
| `wlt system user get --id <N>` | 获取用户详情 | `--id`（必填） |
| `wlt system user create --data '<json>'` | 创建用户 | `--data`（必填） |
| `wlt system user update --data '<json>'` | 更新用户 | `--data`（必填） |
| `wlt system user delete --id <N>` | 删除用户 | `--id`（必填） |
| `wlt system user update-password --data '<json>'` | 修改密码 | `--data`（必填） |
| `wlt system user update-status --data '<json>'` | 更新用户状态 | `--data`（必填） |
| `wlt system user export` | 导出 Excel | 无 |
| `wlt system user get-import-template` | 下载导入模板 | 无 |
| `wlt system user import --data '<json>'` | 导入用户 | `--data`（必填） |

## 部门管理 (`wlt system dept`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt system dept list` | 查询部门（非分页） | `--name`, `--status` |
| `wlt system dept get --id <N>` | 获取部门详情 | `--id`（必填） |
| `wlt system dept create --data '<json>'` | 创建部门 | `--data`（必填） |
| `wlt system dept update --data '<json>'` | 更新部门 | `--data`（必填） |
| `wlt system dept delete --id <N>` | 删除部门 | `--id`（必填） |

## 角色管理 (`wlt system role`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt system role list` | 查询角色（非分页） | `--name`, `--status` |
| `wlt system role get --id <N>` | 获取角色详情 | `--id`（必填） |
| `wlt system role create --data '<json>'` | 创建角色 | `--data`（必填） |
| `wlt system role update --data '<json>'` | 更新角色 | `--data`（必填） |
| `wlt system role delete --id <N>` | 删除角色 | `--id`（必填） |

## 菜单管理 (`wlt system menu`)

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt system menu list` | 查询菜单（非分页） | `--name`, `--status` |
| `wlt system menu get --id <N>` | 获取菜单详情 | `--id`（必填） |
| `wlt system menu create --data '<json>'` | 创建菜单 | `--data`（必填） |
| `wlt system menu update --data '<json>'` | 更新菜单 | `--data`（必填） |
| `wlt system menu delete --id <N>` | 删除菜单 | `--id`（必填） |

## 字典管理 (`wlt system dict`)

### 字典类型

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt system dict type-list` | 查询字典类型列表 | 无 |
| `wlt system dict type-get --id <N>` | 获取字典类型详情 | `--id`（必填） |
| `wlt system dict type-create --data '<json>'` | 创建字典类型 | `--data`（必填） |
| `wlt system dict type-update --data '<json>'` | 更新字典类型 | `--data`（必填） |
| `wlt system dict type-delete --id <N>` | 删除字典类型 | `--id`（必填） |

### 字典数据

| 命令 | 说明 | 关键参数 |
|------|------|---------|
| `wlt system dict data-list` | 分页查询字典数据 | `--dict-type`, `--status`, `--page-no`, `--page-size` |
| `wlt system dict data-get --id <N>` | 获取字典数据详情 | `--id`（必填） |
| `wlt system dict data-create --data '<json>'` | 创建字典数据 | `--data`（必填） |
| `wlt system dict data-update --data '<json>'` | 更新字典数据 | `--data`（必填） |
| `wlt system dict data-delete --id <N>` | 删除字典数据 | `--id`（必填） |
