# Bootstrap 任务：填充项目开发约定（已据实改写）

> **本文件已由任务 `10-16-superpowers-to-trellis`（2026-10-16）按仓库实际完成情况改写。**
> 原 `trellis init` 生成的 fullstack 模板（backend/ + frontend/ 两层 + "All documentation should be
> written in English"）与本项目现实不符：本项目是 **Go 编写的 CLI**，无前端、无数据库 ORM。
> 模板中的 backend/frontend 占位层已被删除，改为 `cli/` + `internal/` 两层。
> 语言以**中文**为准（与 `CLAUDE.md`、`skills/`、`docs/`、commit message 一致）。

**你（AI）是本任务的执行者，开发者不读本文件。**

`.trellis/spec/` 是 `trellis-implement` / `trellis-check` 子代理的唯一上下文来源：
每个任务通过 `implement.jsonl` / `check.jsonl` 声明要加载哪些 spec 文件，
平台 hook 会自动把它们连同 `prd.md` 注入子代理提示词。
**spec 留空 → 子代理写出通用代码；spec 真实 → 子代理匹配团队实际模式。**

---

## 状态（按实际完成情况）

- [x] 填充 cli 层约定（命令分包/工厂选型/flag 映射/输出协议/后端同步/测试）
- [x] 填充 internal 层约定（鉴权与客户端/配置/输出与错误/测试）
- [x] 加入真实代码示例与 `path:line` 锚点（非泛化最佳实践）
- [x] 删除与本项目无关的 backend/frontend 占位层
- [x] 明确记录仓库中**真实存在的不一致**（见下方「如实记录的既有不一致」）
- [ ] （可选）后续任务：收敛 `.trellis/spec/` 与 `CLAUDE.md` / `skills/` 的重复内容

---

## Spec 文件清单（cli 层）

| 文件 | 记录什么 | 状态 |
|------|----------|------|
| `.trellis/spec/cli/index.md` | 层索引 + Pre-Development Checklist + Quality Check | 已填充 |
| `.trellis/spec/cli/directory-structure.md` | 分包、`Register(parent)`、`init()` 注册、新增业务域步骤、文件命名、`cmd/partner/` | 已填充 |
| `.trellis/spec/cli/command-factories.md` | `DocumentCmds` / `AddCRUDToParent` / `NewCRUDGroup` / `Crud*Cmd` / 手写的选型判据与真实调用点 | 已填充 |
| `.trellis/spec/cli/flags-and-params.md` | flag 命名、kebab→camel、时间数组、分页、`--headers`、`--data`、Put/PutQuery | 已填充 |
| `.trellis/spec/cli/output-protocol.md` | stdout/stderr 信封、退出码 0-6、`ParsePagedJSON` vs `OutputJSON` | 已填充 |
| `.trellis/spec/cli/backend-sync-playbook.md` | 后端基线差异分析、A/B/排除分类、薄透传原则、挂载决策 | 已填充 |
| `.trellis/spec/cli/quality-guidelines.md` | 测试清单/风格、Makefile 验证命令、命令树核对 | 已填充 |

## Spec 文件清单（internal 层）

| 文件 | 记录什么 | 状态 |
|------|----------|------|
| `.trellis/spec/internal/index.md` | 层索引 + Checklist + Quality Check + 层间职责边界 | 已填充 |
| `.trellis/spec/internal/auth-and-client.md` | 无状态鉴权契约、`RequestContext`、请求头、重试、`CommonResult` 解析 | 已填充 |
| `.trellis/spec/internal/config.md` | `Profile` 结构、`~/.wlt/config.yaml`、`ActiveProfile()`、`UpdateProfileField` | 已填充 |
| `.trellis/spec/internal/output-and-errors.md` | `output` 信封与 `ExitError`、`apierr` 退出码映射来源 | 已填充 |
| `.trellis/spec/internal/quality-guidelines.md` | 内部包测试约定与验证命令 | 已填充 |

## 思维指南（保留，未修改）

`.trellis/spec/guides/` 下 3 个文件（`index.md`、`code-reuse-thinking-guide.md`、
`cross-layer-thinking-guide.md`）为框架无关的通用思维指南，**保持不变**。

---

## 本章程固化时的两个原则

### 原则 1：从既有约定文档导入优先

本项目的既有约定文档是 `CLAUDE.md`、`skills/`、以及历史设计文档
（原 `docs/superpowers/` 下 8 篇，已归档为 Trellis 归档任务并清理）。

### 原则 2：Document reality, not ideals

**写代码实际做了什么，而不是应该做什么。** 已知技术债如实记录，改进是独立话题。

### 如实记录的既有不一致（不粉饰）

1. 分页子命令主名**并非**全仓统一 `list`：`finance`(6)/`produce`(2)/`quality`(1)/`waybill`(1)
   仍以 `page` 为主名；`page` 仅在 stock 四个单据上是 `list` 的别名。
2. `DocumentCmds` vs `AddCRUDToParent` **无强制规则**，域间并存且不一致
   （`product`/`stock_warehouse` 用后者，`product_unit`/`product_metrics` 用 `NewCRUDGroup`）。
3. `--headers` 只部分命令支持：`settlement main export` **无**；logistics / 销售核算 / 运费申请有；
   `DocumentCmds` 的 `page-count` 刻意不收集。
4. export 类命令输出**不是 JSON 信封**（23 处 `fmt.Println`），是待 `-o` 方案统一的已知例外。
5. 时间范围有两套并存实现：通用 `AddTimeRangeFlags/CollectTimeRangeFlags` 与 `crud.go` 内 6 组域内 helper；
   `stock query/record` 的 `in-time`/`create-time` 是**单值**，与 `stock in/out` 的 `{timeKey}[0]/[1]` 数组契约不同。

### 已纠正的两处文档漂移（不得回退）

1. **退出码 6 `network` 从未产出** —— 常量 0–6 均已定义，但生产代码实际产出仅 **1/2/4/5**；
   3（authentication）与 6（network）为保留值。网络错误与 token 失效都表现为退出码 5。
2. **legacy 工厂已删除** —— `internal/cmdutil/crud_legacy.go` 与 `NewCRUDSubCmd` 全仓零命中。

---

## 运行时说明（当开发者问「为什么需要 spec」时分享）

- 每个 AI 编码任务会派两个子代理：`trellis-implement`（写代码）与 `trellis-check`（验质量）。
- 每个任务有 `implement.jsonl` / `check.jsonl`，列出要加载哪些 spec 文件。
- 平台 hook 会自动把这些 spec 与任务的 `prd.md` 注入子代理提示词。
- 唯一事实来源是 `.trellis/spec/`，所以现在填好，长期受益。

---

## 完成

当开发者确认上述清单项已用真实示例（而非占位）完成，运行：

```bash
python ./.trellis/scripts/task.py finish
python ./.trellis/scripts/task.py archive 00-bootstrap-guidelines
```

归档后，新加入的开发者会收到 `00-join-<slug>` 入驻任务而不是本 bootstrap 任务。

---

## Spec 契约修正规则（原模板的 Pre-Modification Rule 摘要）

**改任何值之前先搜索**：

```bash
grep -r "value_to_change" .
```

改 spec 前先确认该约定在代码里仍成立（锚点 `path:line` 仍然对得上）；
spec 与代码冲突时，**以代码为准并修正 spec**。