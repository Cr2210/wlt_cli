# 技术设计：superpowers 残留清理与 spec/文档迁移

## 边界

本任务只触碰**文档与任务元数据**，不触碰任何 Go 源码：

| 区域 | 处理 |
|------|------|
| `.trellis/spec/**` | 改造结构 + 写入真实约定 |
| `.trellis/tasks/**` | 新增 5 个任务目录（4 归档 + 1 开放） |
| `.trellis/.template-hashes.json` | 移除 52 个 `.cursor/*` 键（R9） |
| `docs/**`、`.superpowers/**`、`.flow/**`、`.cursor/**` | 删除 |
| `cmd/**`、`internal/**`、`CLAUDE.md`、`skills/**`、`Makefile`、`.goreleaser.yml`、`AGENTS.md` | **不动**（R8/R9） |

---

## 1. spec 目标结构

依据 `.claude/skills/trellis-meta/references/local-architecture/spec-system.md`：
单仓模式下 `.trellis/spec/<layer>/`，每层一个 `index.md`（含 **Pre-Development Checklist** 与
**Quality Check** 两段），细则放同级文件。`guides/` 为跨层思维指南，**保留不动**。

```
.trellis/spec/
├── cli/                              命令层（cmd/ 下各业务域包）
│   ├── index.md                      入口：checklist + quality check + 索引
│   ├── directory-structure.md        分包、Register/init 注册、新增业务域步骤
│   ├── command-factories.md          工厂选型决策（DocumentCmds/AddCRUDToParent/Crud*/手写）
│   ├── flags-and-params.md           flag 命名、kebab→camel、分页、时间数组、--headers/--data
│   ├── output-protocol.md            stdout/stderr 分离、退出码 0-6、ParsePagedJSON vs OutputJSON
│   ├── backend-sync-playbook.md      后端基线差异分析＋新增端点对接流程
│   └── quality-guidelines.md         测试类型（flag 装配/表驱动）、构建验证、文档同步
├── internal/                         共享内部包层（internal/）
│   ├── index.md
│   ├── auth-and-client.md            无状态鉴权契约、RequestContext、多租户头、重试
│   ├── config.md                     profile 结构、~/.wlt/config.yaml、UpdateProfileField
│   ├── output-and-errors.md          output 包、apierr、退出码映射来源
│   └── quality-guidelines.md
└── guides/                           保留（已有 3 个文件，内容为框架无关的思维指南）
```

删除：`.trellis/spec/backend/`（5 个 To fill 占位）、`.trellis/spec/frontend/`（6 个前端专用占位）。
理由：本项目是 Go CLI，无 ORM/数据库迁移（`database-guidelines` 无对象）、无前端
（`hook-guidelines`/`state-management`/`component-guidelines` 无对象）。

**为何 `cli` + `internal` 而不是别的命名**：与 `cmd/`、`internal/` 两个真实顶层源码目录一一对应，
命名即可推断归属；且 `packages_context.py:30-40` 的 `_scan_spec_layers()` 直接扫描 `spec/` 子目录
（排除 `guides`）作为 layer，改完即被 `get_context.py --mode packages` 识别，**无需改脚本或 config.yaml**。
保持单仓模式，不引入 `packages:` 配置（那会要求 `spec/<pkg>/<layer>/` 再深一层，收益不足）。

---

## 2. spec 内容来源映射

每条规则必须能落到真实 `path:line`。来源分两类：**历史设计文档中的决策**（知识来源）与
**代码本身**（事实权威）。

| spec 文件 | 主要来源 |
|-----------|---------|
| `cli/directory-structure.md` | `CLAUDE.md` 项目结构/架构约定；`cmd/*/` 真实包；`cmd/root.go` init 注册 |
| `cli/command-factories.md` | stock-cmd 设计文档（工厂演进与选型决策、page→list 别名、`--headers` 补齐）；`internal/cmdutil/crud.go` |
| `cli/flags-and-params.md` | query-endpoints backlog（kebab→camel 修复）；stock-cmd 设计（TimeKey 数组、flag 复用）；backup-sync 设计（`{timeKey}[0]/[1]`）；`internal/cmdutil/helpers.go`、`timeflags.go` |
| `cli/output-protocol.md` | 无状态鉴权设计（退出码 4/5/3 的边界）；prepayment 设计（裸数组响应不可用 ParsePagedJSON）；`internal/cmdutil/output.go`、`internal/output/` |
| `cli/backend-sync-playbook.md` | 0923 同步设计（差异分析法、A/B/排除项分类、薄透传原则、挂载决策）；backlog（从 Swagger 取参数的方法） |
| `cli/quality-guidelines.md` | 各设计中反复出现的验证要求；仓库真实 `_test.go` 文件；`Makefile` |
| `internal/auth-and-client.md` | 无状态鉴权设计（完整契约表 + 校验边界）；`internal/cmdutil/client.go`、`internal/client/client.go` |
| `internal/config.md` | 无状态鉴权设计（Profile 字段删除）；`internal/config/config.go` |
| `internal/output-and-errors.md` | 无状态鉴权设计（错误分类说明）；`internal/output/`、`internal/apierr/` |
| 两个 `quality-guidelines.md` | 各自层的测试/验证约定 |

> 具体 file:line 锚点由 `research/` 下的研究产物提供并在实施时逐条核对。
> 研究产物按主题拆为 `00-index.md` + `01-auth-contract.md` … `06-backlog-facts.md`：
> 原始单文件 57 KB 超过 `context_injection.max_file_bytes`（32768），注入子代理时会被**截断**，
> 而截断的约定清单比不注入更危险（子代理会照半截事实写 spec）。拆分后各片 ≤12 KB，
> 全量注入无截断；副带好处是未来任务能按主题只注入需要的分片。

### 2.1 研究已证实必须纠正的漂移（写入 spec 时的硬约束）

研究在代码中核实了 10 条"文档说一套、代码做另一套"的漂移（见 `research/05-spec-layout-and-drift.md` §C）。
其中 3 条会**直接毒化** spec，必须显式纠正而不得照抄：

1. **退出码 6 `network` 从未产出** —— 全仓无 `NewExitError(6)`；网络错误实际落 exit 5。
   `skills/SKILL.md:318,324` 与 `skills/references/api-conventions.md:105` 的"6 = 网络错误"是失真。
   spec 须写"0–6 常量已定义，**实际产出仅 1/2/4/5**；3、6 为保留值"。
   （`internal/apierr/apierr.go:43-68` 定义，`cmd/root.go:74-79` 消费。）
2. **`skills/references/api-conventions.md:62-71` 仍在描述已删除的 legacy CRUD**
   （`newCRUDSubCmd`，列为"三种 CRUD 模式"之一），而 `internal/cmdutil/crud_legacy.go` 已删除、
   全仓零命中。spec 只写 `DocumentCmds`/`AddCRUDToParent`/`NewCRUDGroup`/`Crud*Cmd`/手写五类。
3. **"分页子命令全仓统一叫 `list`"不成立** —— `finance`(6)/`produce`(2)/`quality`(1)/`waybill`(1)
   仍以 `page` 为主名；`list`+`page` 别名只覆盖 `DocumentCmds` 族与 4 个 stock 单据。

另有 7 条属"现实不整齐但非错误"，spec 须如实记录（见 prd R6b）。
注意 `CLAUDE.md` 的历史漂移（曾称 stock 用 legacy 工厂）**已自行修复**（`CLAUDE.md:64,79-80` 现正确），
不要再当作现存问题去"修正"。

---

## 3. 历史归档映射

**4 个归档任务**（每个一个目录，`archive/{YYYY-MM}/`）：

| # | 归档任务 slug | prd.md | design.md | implement.md |
|---|---|---|---|---|
| 1 | `stateless-auth` | 新写：目标/约束/验收（从原文档 Goal 段提取） | 无状态鉴权设计原文 | 从原文档「代码改动/文档改动」段生成（原文即实施清单） |
| 2 | `stock-cmd-optimization` | 新写：目标/非目标 | stock-cmd 设计原文 | stock-cmd 计划原文 |
| 3 | `order-prepayment-relation` | 新写：目标/非目标 | prepayment 设计原文 | prepayment 计划原文 |
| 4 | `backend-sync-0923` | 新写：目标/非目标 | 0923 设计原文 | 0923 计划原文 |

拆分原则：**原文内容逐字保留，不重写**，仅按 trellis 三件套归位；每个文件头部加一段
溯源说明（原始路径 + 交付日期 + 对应 commit + 说明本文档是历史归档而非待办）。

**1 个开放任务**：`query-endpoints-round2`（status 保持 `planning`）
- `prd.md`：第 2 轮待对接端点清单（逐字保留）、export-excel 的 `-o` 下载方案诉求、
  2 处待决策、以及"第 1 轮已完成"记录与未完成的实测项作为背景。
- 不提供 design.md / implement.md（本轮不进入实施）。

---

## 4. 已核实的 Trellis 机制约束

实施时必须遵守以下已从脚本源码核实的行为，否则会踩坑：

| 约束 | 证据 | 应对 |
|------|------|------|
| `task.py create` 会**自动抢占当前会话的 active task 指针** | `task.py create` 输出 "Activated task for this session" | 建这 5 个任务一律加 `--no-start`，避免覆盖本任务的 active 指针 |
| 归档目录的 `{YYYY-MM}` 取自**归档时刻**，非交付日期 | `task_utils.py:110-114` `datetime.now().strftime("%Y-%m")` | 4 个历史任务都落在 `archive/2026-10/`；历史日期写进 `task.json` 与 prd |
| `rename` 拒绝以日期前缀开头的新 slug | `task_store.py:804` | 无法得到 `2026-08-07-*` 形式的目录名；接受 `MM-DD-<slug>` 形式，历史日期靠内部字段承载 |
| 归档前校验 branch 元数据；仓库有 remote 且 `base_branch=main`、`branch=null` 会被拒 | `task_store.py:1183-1225`（docstring：有 base_branch + 有 remote ⇒ 视为 PR-backed）+ `git remote` 有 2 个 remote | 归档时加 `--skip-branch-validation` |
| `archive` 会设置 `status=completed`、`completedAt=今天`，**默认自动 git 提交** | `task_store.py:1350+`；`--no-commit` 开关 | 用 `--no-commit`；是否提交留待用户在收尾时决定（见下） |
| `_scan_spec_layers` 排除 `guides`、只认子目录 | `packages_context.py:40` | `guides/` 不会被当成业务层，符合预期 |
| `task.py create` 的 `--description` 必填非空 | `task.py create --help` | 5 个任务都要给 description |

**提交策略**：本任务不自动提交（`--no-commit`），全部改动完成后向用户汇报 diff，
由用户决定是否提交以及提交粒度。理由：本任务改动量大（删除已跟踪文件 + 新增大量
`.trellis/` 内容），且 `.goreleaser.yml:92` 的 changelog 过滤器会排除 `^docs:` 前缀的 commit，
提交信息措辞需要用户确认。

### 4.1 移除 Cursor 平台（R9）——为什么必须动哈希清单

`trellis platforms` 的输出**不来自目录探测**。`@mindfoldhq/trellis` 的

```js
// dist/configurators/index.js:103-133
export function getConfiguredPlatforms(cwd) {
  const hashes = loadHashes(cwd);          // 读 .trellis/.template-hashes.json
  for (const id of PLATFORM_IDS) {
    const configDir = AI_TOOLS[id].configDir;   // cursor → ".cursor"
    const hasTrackedTemplate = [...templates.keys()].some(p =>
      (p === configDir || p.startsWith(`${configDir}/`)) && hashes[p] !== undefined);
    if (hasTrackedTemplate) platforms.add(id);
  }
```

即：**某平台"已配置" ⟺ 哈希清单里存在该平台 `configDir` 下的条目**。
`commands/update.js:435` 又正是遍历这个结果去管理平台文件：

```js
for (const platformId of getConfiguredPlatforms(cwd)) { ... collectPlatformTemplates(platformId) ... }
```

推论（这是 R9 的整个设计依据）：
- **只删 `.cursor/` 不删哈希** ⇒ `platforms` 仍报 Cursor、`update` 仍会重新生成 52 个文件，
  用户"我不需要 Cursor"的意图不成立 → 假删除。
- **先删哈希再删目录**（或反之，两者都做）⇒ `getConfiguredPlatforms` 只剩 `claude-code`，
  `update` 不再管理 Cursor → 真删除。
- trellis 自带 `utils/manifest-prune.js` 的 `pruneOrphanManifestKeys`（被 `update.js:1663` 使用）
  证明"清理清单中的悬空键"是该工具认可的常规操作，不是野路子。
- 未跟踪性：`.cursor/` 与 `.template-hashes.json` 都**未纳入 git**（`git status` 显示为 `??`），
  因此删除**无法经 git 恢复**；恢复手段是重跑 `trellis init --cursor`。删除前需确认这一点被用户知晓。
- `AGENTS.md` 不在任何平台的 `configDir` 下（哈希键就是字面 `AGENTS.md`），
  故它与平台判定无关，保留不影响 R9 目标。

---

## 5. 兼容性与影响面

- **spec 层改名会打断既有引用**：`.trellis/tasks/00-bootstrap-guidelines/task.json` 的
  `relatedFiles` 现引用 `.trellis/spec/backend/`、`.trellis/spec/frontend/`；其 `prd.md` 也按
  backend/frontend 列出待填文件。改造后必须同步改写，否则留下指向不存在路径的悬空引用。
  处理：改写 `relatedFiles` 为 `cli/`、`internal/`，并在 prd 中据实改写校验项，
  最后按 trellis 自身指引归档该 bootstrap 任务（其交付物即本次 spec 填充）。
- **git 影响**：删除 8 个已跟踪文件（可从 `git log`/`git show` 恢复）；新增
  `.trellis/spec/**`、`.trellis/tasks/**` 均为未跟踪文件（`.trellis/.gitignore` 未排除 spec/skills，
  仅排除运行时会话文件——实施时核对一遍）。
- **与 `CLAUDE.md` / `skills/` 的重复**：用户已决定本轮不动。spec 中的鉴权/API 约定会与
  `skills/references/api-conventions.md`、`auth-config.md` 部分重叠。这是**有意接受**的：
  `skills/` 面向 CLI 最终使用者（Agent 调用者），`.trellis/spec/` 面向改代码的 implement/check
  子代理，受众不同；收敛留待后续单独任务。
- **`docs/` 空目录**：删除 `docs/superpowers/` 后 `docs/` 内无内容，一并移除（已确认无代码、
  Makefile、goreleaser 引用它；`.goreleaser.yml` 只在打包 `skills/`，`^docs:` 是 commit 前缀过滤器）。

---

## 6. 取舍

| 取舍 | 选定 | 被放弃的替代方案及原因 |
|------|------|---------------------|
| spec 层命名 | `cli` + `internal` | 保留 `backend`（语义失真，项目无后端代码）；新建三层细分（过度设计） |
| spec 语言 | 中文 | 英文（模板默认；与全仓中文现实冲突，且 bootstrap prd 要求"记录现实"） |
| 原文处置 | 逐字保留进归档任务 | 重写摘要（丢失设计推理过程与锚点）；直接删除只留 git（检索需绕道 git） |
| 归档粒度 | 每交付一任务 | 打包单任务（丢单条可检索性） |
| 提交 | 不自动提交，交用户决定 | 随 `task.py` 自动提交（4+ 个隐式 commit，粒度用户未确认） |

---

## 7. 回滚

- 已跟踪内容：`git checkout -- docs/` 即恢复 8 篇原文。
- 新增内容：删除 `.trellis/spec/cli/`、`.trellis/spec/internal/`、`archive/` 下 4 个目录、
  开放任务目录即可；`backend/`/`frontend/` 内容可从 `git`（若曾提交）或 `trellis init` 重新生成。
- 无数据/配置迁移，无代码影响，回滚成本低。